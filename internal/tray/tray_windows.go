//go:build windows
// +build windows

package tray

import (
	_ "embed"
	"os"

	"sorty/pkg/settings"

	"fyne.io/systray"
	"github.com/rs/zerolog/log"
)

//go:embed imgs/logo.ico
var iconData []byte

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
