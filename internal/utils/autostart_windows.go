//go:build windows
// +build windows

package utils

import (
	_ "embed"
	"os"
	"path/filepath"
	"sorty/logger"
	"sorty/pkg/settings"

	"github.com/getlantern/systray"
	"github.com/rs/zerolog"
	"golang.org/x/sys/windows/registry"
)

var log *zerolog.Logger

//go:embed imgs/logo.ico
var iconData []byte

func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

// FIXME: fix autostart on win11

/*
getting path of executable and putting this into registry key

	(HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run)
*/
func EnableAutoStart() error {
	log.Debug().Msg("Starting autostart setup")

	InitTray()
	p, err := os.Executable()
	if err != nil {
		log.Error().Err(err).Msg("failed to get executable path")
		return err
	}
	exePath := filepath.Clean(p)
	log.Debug().Str("path", exePath).Msg("got executable path")
	// creating registry key for automatic startup
	k, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create/open registry key")
		return err
	}
	defer k.Close()

	if err := k.SetStringValue("Sorty", exePath); err != nil {
		log.Error().Err(err).Str("path", exePath).Msg("failed to set registry value")
		return err
	}

	log.Info().Msg("Successfully set up autostart in Windows registry")
	return nil
}

// init sorty tray on Windows
func InitTray() {
	log.Debug().Msg("Initializing system tray")
	systray.Run(tray, exitTray)
}

// setting up the tray
func tray() {
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetTitle("Sorty")
	systray.SetTooltip("IM WORKING YAY")

	mSettings := systray.AddMenuItem("Settings", "Open settings window")
	mQuit := systray.AddMenuItem("Quit", "dude, name says for itself. just close the tray")
	log.Info().Msg("System tray initialized successfully")

	go func() {
		for {
			select {
			case <-mSettings.ClickedCh:
				log.Info().Msg("Settings requested through system tray")
				go settings.ShowSettingsWindow()
			case <-mQuit.ClickedCh:
				log.Info().Msg("Quit requested through system tray")
				systray.Quit()
			}
		}
	}()
}

// function to exit the tray
func exitTray() {
	log.Info().Msg("Application shutting down")
	os.Exit(0)
}
