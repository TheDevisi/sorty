//go:build linux
// +build linux

package utils

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/emersion/go-autostart"
)

const appName = "Sorty"

func EnsureInstalled() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get home dir")
	}

	target := filepath.Join(home, ".local/bin/sorty")

	self, err := os.Executable()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get executable path")
	}
	self, _ = filepath.EvalSymlinks(self)

	if self != target {
		log.Info().Msg("installing to " + target)

		if err := installSelf(self, target); err != nil {
			log.Fatal().Err(err).Msg("install failed")
		}

		if err := enableAutostart(target); err != nil {
			log.Fatal().Err(err).Msg("autostart failed")
		}

		log.Info().Msg("restarting from installed location")
		exec.Command(target).Start()
		os.Exit(0)
	}
	if err := installDesktopEntry(target); err != nil {
		log.Error().Err(err).Msg("failed to install desktop entry")
	}
	if err := enableAutostart(target); err != nil {
		log.Error().Err(err).Msg("autostart check failed")
	}
}

func installSelf(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return os.Chmod(dst, 0755)
}

func installDesktopEntry(execPath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".local/share/applications")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	desktop := `[Desktop Entry]
Type=Application
Name=Sorty
Comment=Sorty application
Exec=` + execPath + `
Icon=utilities-terminal
Terminal=false
Categories=Utility;
`

	return os.WriteFile(
		filepath.Join(dir, "Sorty.desktop"),
		[]byte(desktop),
		0644,
	)
}

func enableAutostart(execPath string) error {
	app := &autostart.App{
		Name:        appName,
		DisplayName: appName,
		Exec:        []string{execPath},
	}

	if app.IsEnabled() {
		log.Debug().Msg("autostart already enabled")
		return nil
	}

	log.Info().Msg("enabling autostart")
	return app.Enable()
}
