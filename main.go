package main

import (
	"sorty/config"
	"sorty/internal/utils"
	"sorty/logger"
	"sorty/pkg/watcher"

	"github.com/rs/zerolog"
)

var log *zerolog.Logger

func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
	log.Info().Msg("Initializing Sorty application")
}

func main() {
	log.Info().Msg("Starting system tray initialization")
	go utils.InitTray()

	log.Info().Msg("Checking configuration file")
	config.CheckIfConfigExists()

	log.Info().Msg("Setting up auto startup capability if it's not it")
	if err := utils.EnableAutoStart(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to enable auto-start")
	}

	log.Info().Msg("Initializing directory watcher")
	watcher.WatchDirectory()

}
