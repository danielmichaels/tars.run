package main

import (
	"github.com/danielmichaels/shortlink-go/cmd/api/config"
	"github.com/danielmichaels/shortlink-go/internal/logger"
	"os"
)

type application struct {
	config *config.Conf
	logger *logger.Logger
}

func main() {
	cfg := config.AppConfig()
	log := logger.New(os.Stdout, logger.LevelInfo)
	app := &application{
		config: cfg,
		logger: log,
	}

	err := app.serve()
	if err != nil {
		log.PrintFatal(err, nil)
	}
}
