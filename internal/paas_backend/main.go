// Main backend logic
package paas_backend

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/repofetch"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/secretsprovider"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"

	_ "github.com/ThomasRubini/cloud-paas/internal/paas_backend/docs"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func printBuildInfo() {
	var rev string
	var time string
	var modified bool
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				rev = setting.Value[0:7]
			} else if setting.Key == "vcs.time" {
				time = setting.Value
			} else if setting.Key == "vcs.modified" {
				modified = setting.Value == "true"
			}
		}
	}
	if rev != "" || time != "" {
		fmt.Printf("Program built at %v from commit %v (dirty=%v)\n", time, rev, modified)
	} else {
		fmt.Printf("Could not retrieve build info\n")
	}
}

func connectToDB() (*gorm.DB, error) {
	c := config.Get()
	logrus.Debug("Connecting to database..")
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("%s password=%s", c.DB_URL, c.DB_PASSWORD)), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	logrus.Debug("Connected to database !")

	db = db.Set("gorm:auto_preload", true)

	return db, nil
}

func MigrateModels(db *gorm.DB) error {
	logrus.Debug("Running database migrations..")
	models := []interface{}{models.DBApplication{}, models.DBEnvironment{}}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to run database migrations for model %v: %w", model, err)
		}
	}
	return nil
}

func setupLogging() {
	var cliVerbose bool
	var cliTrace bool
	flag.BoolVar(&cliVerbose, "v", false, "enable verbose logging")
	flag.BoolVar(&cliTrace, "vv", false, "enable trace logging")
	flag.Parse()

	if cliTrace {
		logrus.SetLevel(logrus.TraceLevel)
	} else if cliVerbose || config.Get().VERBOSE {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Debug("Verbose logging enabled")
	logrus.Trace("Trace logging enabled")
}

func getSecretsProvider() secretsprovider.Helper {
	c := config.Get()
	impl := c.SECRETS_IMPL
	if impl == "file" {
		if c.SECRETS_IMPL_FILE == "" {
			panic("file secrets backend chosen but SECRETS_IMPL_FILE environment variable not set")
		}
		return secretsprovider.Helper{Core: secretsprovider.FromFile(c.SECRETS_IMPL_FILE)}
	} else if impl == "vault" {
		panic("TODO")
	} else {
		panic("Currently supported secrets backends: [file, vault]")
	}
}

func constructState(conf config.Config) (*utils.State, error) {
	// Connect to DB
	db, err := connectToDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := MigrateModels(db); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Get docker client
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Test registry connection
	_, err = dockerClient.RegistryLogin(context.Background(), registry.AuthConfig{
		Username:      conf.REGISTRY_USER,
		Password:      conf.REGISTRY_PASSWORD,
		ServerAddress: conf.REGISTRY_REPO_URI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to login to registry: %w", err)
	}

	// Get helm client
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), "default", "memory", log.Printf); err != nil {
		return nil, fmt.Errorf("error initializing config: %w", err)
	}
	// Test it
	if err = actionConfig.KubeClient.IsReachable(); err != nil {
		return nil, fmt.Errorf("failed to connect to kubernetes cluster: %w", err)
	}

	// Construct state
	return &utils.State{
		Db:              db,
		DockerClient:    dockerClient,
		HelmConfig:      actionConfig,
		SecretsProvider: getSecretsProvider(),
	}, nil
}

func Entrypoint() {
	printBuildInfo()
	config.Init()
	setupLogging()

	// Setup state (note: we assign to the global variable here)
	state, err := constructState(config.Get())
	if err != nil {
		logrus.Fatalf("Failed to construct state: %v", err)
	}
	utils.SetState(*state)

	// Setup web server
	g := SetupWebServer(*state)

	// init crontab for fetching repos
	if config.Get().REPO_FETCH_ENABLE {
		logrus.Info("Starting repository fetch crontab")
		repofetch.Init(config.Get().REPO_FETCH_PERIOD_SECS)
	} else {
		logrus.Info("Repository fetch crontab disabled")
	}

	logrus.Info("Setup finished successfully. Starting to serve incoming requests.")
	// Launch web server. Function will never return
	launchWebServer(g)
}
