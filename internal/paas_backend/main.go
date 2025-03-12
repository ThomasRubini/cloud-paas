// Main backend logic
package paas_backend

import (
	"flag"
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/repofetch"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/secretsprovider"
	"github.com/ThomasRubini/cloud-paas/internal/utils"

	_ "github.com/ThomasRubini/cloud-paas/internal/paas_backend/docs"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDB() (*gorm.DB, error) {
	c := config.Get()
	logrus.Debug("Connecting to database..")
	db, err := gorm.Open(postgres.Open(c.DB_URL), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	logrus.Debug("Connected to database !")

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

func Entrypoint() {
	config.Init()
	setupLogging()

	// Connect to DB
	db, err := connectToDB()
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	if err := MigrateModels(db); err != nil {
		logrus.Fatalf("Failed to run database migrations: %v", err)
	}

	// Setup state (note: we assign to the global variable here)
	state := utils.State{
		Db:              db,
		SecretsProvider: getSecretsProvider(),
	}
	utils.SetState(state)

	// Setup web server
	g := SetupWebServer(state)

	// init crontab for fetching repos
	if config.Get().REPO_FETCH_ENABLE {
		logrus.Info("Starting repository fetch crontab")
		repofetch.Init(config.Get().REPO_FETCH_PERIOD_SECS)
	} else {
		logrus.Info("Repository fetch crontab disabled")
	}

	// Launch web server. Function will never return
	launchWebServer(g)
}
