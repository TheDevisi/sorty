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
	log.Info().Msg("Checking configuration file")
	config.CheckIfConfigExists()

	log.Info().Msg("Setting up auto startup capability if it's not set")
	// if err := utils.EnableAutoStart(); err != nil {
	// 	log.Error().Err(err).Msg("Failed to set up autostart")
	// }
	go watcher.WatchDirectory()

	log.Info().Msg("Starting system tray initialization")
	utils.InitTray()
	log.Info().Msg("Initializing directory watcher")
}
