// some utils for any other pkgs.
package utils

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"sorty/config"
	"sorty/internal/errors"
	"time"
)

func GetOperatingSystem() string {
	OS := runtime.GOOS
	return OS
}

// function to move file. getting info (true or false) to move from checkFromTempFile
func MoveFile(filePath string) {
	canMove := checkForTmpFile(filePath)
	if canMove {
		fileExt := filepath.Ext(filePath)
		PlaceToMove(fileExt, filePath)
	}
}

// this function checking if current file is temp file. yes ? - wait. no? - move.
func checkForTmpFile(filename string) bool {
	var tmpFileExts []string = []string{".part", ".crdownload", ".tmp", ".temp", ".dmp", ".adadownload", ".cache", ".partial"}
	for _, tmpFileExt := range tmpFileExts {
		if path.Ext(filename) == tmpFileExt {
			time.Sleep(1 * time.Second)
			continue
		}
		fileStat, err := os.Stat(filename)
		if err != nil {
			errors.ErrorsHandler(err, "WARN")
			return false
		}

		size1 := fileStat.Size()
		time.Sleep(2 * time.Second)
		size2 := fileStat.Size()
		if size1 != size2 {
			continue
		}
		return true
	}
	return false
}

// Searching place to move based on operating system and moving
func PlaceToMove(fileExt string, filePath string) {
	OS := GetOperatingSystem()

	userInfo, err := user.Current()
	if err != nil {
		errors.ErrorsHandler(err, "FATAL")
		return
	}

	var configPath string
	if OS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		configPath = filepath.Join(localAppData, "Sorty", "config.json")
	} else {
		configPath = filepath.Join(userInfo.HomeDir, ".config", "sorty", "config.json")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		errors.ErrorsHandler(err, "FATAL")
		return
	}

	var unmarshaledData config.Config
	err = json.Unmarshal(data, &unmarshaledData)
	if err != nil {
		errors.ErrorsHandler(err, "FATAL")
		return
	}

	for path, exts := range unmarshaledData.MonitorFiles {
		for _, ext := range exts {
			if ext == fileExt {
				newPath := filepath.Join(path, filepath.Base(filePath))
				log.Debug().Str("constructedPath", newPath).Msg("Constructed path for file move")
				err := os.Rename(filePath, newPath)
				if err != nil {
					errors.ErrorsHandler(err, "WARN")
				} else {
					log.Info().Str("source", filePath).Str("destination", newPath).Msg("File moved successfully")
				}
				return
			}
		}
	}
}
