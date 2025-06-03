//go:build linux
// +build linux

package utils

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sorty/logger"
	"sorty/pkg/settings"

	"github.com/getlantern/systray"
	"github.com/rs/zerolog"
)

var log *zerolog.Logger

//go:embed imgs/logo.ico
var iconData []byte

// setup logging
func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

/*
this function checks if service already installed.
if not then creating systemd service
*/
func EnableAutoStart() error {
	log.Debug().Msg("Checking if Sorty service exists in /usr/local/bin")

	_, err := os.Stat("/usr/local/bin/sorty")
	if errors.Is(err, os.ErrNotExist) {
		log.Info().Msg("Sorty service not found, proceeding with installation")
		fmt.Println("Seems like there's no sorty service. Run sorty using sudo for autostartup!")

		user, err := user.Current()
		if err != nil {
			log.Error().Err(err).Msg("failed to get current user")
			return fmt.Errorf("failed to get current user: %v", err)
		}
		log.Debug().Str("username", user.Username).Msg("got current user")

		path, err := os.Executable()
		if err != nil {
			log.Error().Err(err).Msg("failed to get executable path")
			return fmt.Errorf("failed to get executable path: %v", err)
		}
		log.Debug().Str("path", path).Msg("got executable path")

		if err := os.Rename(path, "/usr/local/bin/sorty"); err != nil {
			log.Error().Err(err).Str("from", path).Str("to", "/usr/local/bin/sorty").Msg("failed to move executable")
			return fmt.Errorf("failed to move executable: %v", err)
		}
		log.Info().Msg("moved executable to /usr/local/bin/sorty")
		// configuration for service
		systemdConfig := fmt.Sprintf(`[Unit]
Description=Sorty
After=network.target

[Service]
ExecStart=/usr/local/bin/sorty
Restart=on-failure
User=%v

[Install]
WantedBy=multi-user.target`, user.Username)

		if err := os.WriteFile("/etc/systemd/system/sorty.service", []byte(systemdConfig), 0644); err != nil {
			log.Error().Err(err).Msg("failed to create service file")
			return fmt.Errorf("failed to create service file: %v", err)
		}
		cmd := exec.Command("bash", "sudo systemctl enable sorty")
		cmd.Start()
		log.Info().Msg("created systemd service file successfully")

		return nil
	}
	log.Warn().Msg("sorty is already installed in /usr/local/bin")
	return fmt.Errorf("sorty is already installed in /usr/local/bin")
}

func InitTray() {
	log.Debug().Msg("Initializing system tray (Linux version)")
	systray.Run(tray, exitTray)
}

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
