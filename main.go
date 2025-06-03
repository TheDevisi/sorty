package main

import (
	"sorty/config"
	"sorty/internal/utils"
	"sorty/logger"
	"sorty/pkg/settings"
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
	go watcher.WatchDirectory()

	log.Info().Msg("Starting system tray initialization")
	go utils.InitTray()
	log.Info().Msg("Initializing directory watcher")

	settings.RunSettingsWindow()
}
