// Main backend logic
package backend

import (
	"cloud-paas/internal/backend/config"
	"cloud-paas/internal/backend/endpoints"
	"cloud-paas/internal/backend/models"
	"cloud-paas/internal/backend/repofetch"
	"cloud-paas/internal/backend/state"
	"flag"
	"fmt"

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

	models := []interface{}{models.DBProject{}}

	logrus.Debug("Running database migrations..")
	for _, model := range models {
		if db.AutoMigrate(model) != nil {
			return nil, fmt.Errorf("failed to run database migrations for model %v", model)
		}
	}

	return db, nil
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

func Entrypoint() {
	config.Init()
	setupLogging()

	// Connect to DB
	db, err := connectToDB()
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup state
	state.Set(state.T{
		Db: db,
	})

	// Setup web server
	g := setupWebServer()
	endpoints.Init(g.Group("/api/v1"))

	// init crontab for fetching repos
	repofetch.Init(config.Get().REPO_FETCH_PERIOD_SECS)

	// Launch web server. Function will never return
	launchWebServer(g)
}
