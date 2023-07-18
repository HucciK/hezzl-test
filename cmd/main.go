package main

import (
	"hezzl/config"
	"hezzl/internal/app"
	"log"
	"os"
	"path"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(path.Join(cfg.LogFilePath, "logfile"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", 0)

	app := app.New(cfg, logger)
	if err := app.Run(); err != nil {
		logger.Printf("crictical error: %v", err)
		panic(err)
	}
}
