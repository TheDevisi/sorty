//go:build linux
// +build linux

package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
)

func EnableAutoStart() error {
	_, err := os.Stat("/usr/local/bin/sorty")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Seems like there's no sorty service. Run sorty using sudo for autostartup!")
		user, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %v", err)
		}

		path, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %v", err)
		}

		if err := os.Rename(path, "/usr/local/bin/sorty"); err != nil {
			return fmt.Errorf("failed to move executable: %v", err)
		}

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
			return fmt.Errorf("failed to create service file: %v", err)
		}

		return nil
	}
	return fmt.Errorf("sorty is already installed in /usr/local/bin")
}
