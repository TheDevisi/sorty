//go:build linux
// +build linux

package utils

import (
	"os"
	"path/filepath"
)

// возвращает путь к домашней папке + подкаталог
func GetDefaultFolder(name string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch name {
	case "Documents":
		return filepath.Join(home, "Documents")
	case "Pictures":
		return filepath.Join(home, "Pictures")
	case "Music":
		return filepath.Join(home, "Music")
	case "Videos":
		return filepath.Join(home, "Videos")
	default:
		return home
	}
}
