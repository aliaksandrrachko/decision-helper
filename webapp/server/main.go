package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
	"io.github.aliaksandrrachko/decision-helper/webapp/decision"
	"io.github.aliaksandrrachko/decision-helper/webapp/pkg/config"
	"io.github.aliaksandrrachko/decision-helper/webapp/pkg/dbcontext"
	"io.github.aliaksandrrachko/decision-helper/webapp/pkg/ginlogrus"
)

func main() {
	// TODO: think change to singleton
	var logger = &logrus.Logger{
		Out: os.Stderr,
		// %d{HH:mm:ss.SSS} %-5level {%thread} [%logger{20}] : %msg%n
		Formatter: &logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.DebugLevel,
	}

	cfgPath, err := ParseFlags()
	if err != nil {
		logger.Fatal(err)
	}
	config, err := config.NewConfig(cfgPath)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := sql.Open(config.DataSource.DriverName, config.DataSource.Url)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	logger.Info("Data source connected")

	buildHandler(*logger, db, config)
}

func buildHandler(logger logrus.Logger, db *sql.DB, cfg *config.Config) http.Handler {
	router := gin.New()
	router.Use(ginlogrus.Logger(&logger))
	router.Use(gin.Recovery())

	dbContext := dbcontext.New(db, logger)

	decision.RegisterHandlers(
		router.Group("/decision-helper/api/v1"),
		decision.NewService(decision.NewRepository(dbContext, logger), dbContext, logger),
		logger,
	)

	router.Run(
		fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.Port),
	)
	router.SetTrustedProxies(cfg.Server.TrustedProxies)

	return router
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "application", "../../application.yaml", "path to config file")
	flag.Parse()
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
}
