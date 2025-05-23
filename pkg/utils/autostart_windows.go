//go:build windows
// +build windows

package utils

import (
	_ "embed"
	"os"
	"path/filepath"
	"sorty/logger"

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

func EnableAutoStart() error {
	systray.Run(tray, exitTray)
	p, err := os.Executable()
	if err != nil {
		return err
	}
	exePath := filepath.Clean(p)

	k, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer k.Close()

	return k.SetStringValue("Sorty", exePath)
}

func tray() {
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetTitle("Sorty")
	systray.SetTooltip("IM WORKING YAY")

	mQuit := systray.AddMenuItem("Quit", "dude, name says for itself. just close the tray")

	log.Info().Msg("App started in the background")

	// Обработка кнопки выхода
	go func() {
		<-mQuit.ClickedCh
		log.Info().Msg("Tray close.")
		systray.Quit()
	}()
}

func exitTray() {
	os.Exit(0)
}
