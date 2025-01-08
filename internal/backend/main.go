package backend

import (
	"cloud-paas/internal/backend/config"
	"cloud-paas/internal/backend/endpoints"
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDB() *gorm.DB {
	c := config.Get()
	logrus.Debug("Connecting to database..")
	db, err := gorm.Open(postgres.Open(c.DB_URL), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %v", err))
	}
	logrus.Debug("Connected to database !")

	return db
}

func setupLogging() {
	var cliVerbose bool
	flag.BoolVar(&cliVerbose, "v", false, "enable verbose logging")
	flag.Parse()

	if cliVerbose || config.Get().VERBOSE {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Debug("Verbose logging enabled")
}

func Entrypoint() {
	config.Init()
	setupLogging()

	SetState(BackendState{
		db: connectToDB(),
	})

	g := setupWebServer()
	endpoints.Init(g.Group("/api/v1"))
	launchWebServer(g)

	fmt.Println("Backend main")
}
