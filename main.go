package main

import (
	"sorty/config"
	"sorty/internal/utils"
	"sorty/logger"
	"sorty/pkg/watcher"

	"github.com/rs/zerolog"
)

var log *zerolog.Logger

// init initializes the global logger for the Sorty application.
func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
	log.Info().Msg("Initializing Sorty application")
}

// main initializes the Sorty application, verifies configuration, sets up auto-start, starts directory watching in the background, and initializes the system tray.
func main() {
	log.Info().Msg("Checking configuration file")
	config.CheckIfConfigExists()

	log.Info().Msg("Setting up auto startup capability if it's not set")
	if err := utils.EnableAutoStart(); err != nil {
		log.Error().Err(err).Msg("Failed to set up autostart")
	}
	go watcher.WatchDirectory()

	log.Info().Msg("Starting system tray initialization")
	utils.InitTray()
	log.Info().Msg("Initializing directory watcher")
}
