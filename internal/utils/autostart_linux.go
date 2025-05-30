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

	"github.com/getlantern/systray"
	"github.com/rs/zerolog"
)

var log *zerolog.Logger

//go:embed imgs/logo.ico
var iconData []byte

// init initializes the package-level logger with a new logging configuration.
func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

// EnableAutoStart installs and enables the Sorty service as a systemd service if it is not already installed.
//
// If the Sorty executable is not found in /usr/local/bin, this function moves the current executable there,
// creates a systemd service file to run it as the current user, and enables the service to start on boot.
// Returns an error if the service is already installed or if any installation step fails.
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

// InitTray initializes and starts the system tray icon with menu options for the Linux version.
func InitTray() {
	log.Debug().Msg("Initializing system tray (Linux version)")
	systray.Run(tray, exitTray)
}

// tray configures the system tray icon, title, tooltip, and adds a "Quit" menu item with its event handler.
func tray() {
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetTitle("Sorty")
	systray.SetTooltip("IM WORKING YAY")

	mQuit := systray.AddMenuItem("Quit", "dude, name says for itself. just close the tray")
	log.Info().Msg("System tray initialized successfully")

	go func() {
		<-mQuit.ClickedCh
		log.Info().Msg("Quit requested through system tray")
		systray.Quit()
	}()
}

// exitTray handles cleanup and terminates the application when the system tray is closed.
func exitTray() {
	log.Info().Msg("Application shutting down")
	os.Exit(0)
}
