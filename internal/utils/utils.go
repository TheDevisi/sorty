// some utils for any other pkgs.
package utils

import (
	"encoding/json"
	"fmt"
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

// PlaceToMove moves a file to a destination directory based on its extension and the current operating system's configuration.
// 
// The function reads a configuration file that maps file extensions to destination directories. It determines the appropriate config file path depending on the OS, loads the configuration, and moves the file if a matching extension is found. If errors occur during user lookup, config loading, or file moving, they are handled and logged with appropriate severity. The function returns after the first successful move attempt or upon encountering a fatal error.
func PlaceToMove(fileExt string, filePath string) {
	OS := GetOperatingSystem()

	userInfo, err := user.Current()
	if err != nil {
		errors.ErrorsHandler(err, "FATAL")
		return
	}

	var configPath string
	if OS == "windows" {
		programFiles := os.Getenv("ProgramFiles")
		if programFiles == "" {
			errors.ErrorsHandler(fmt.Errorf("ProgramFiles env variable not found"), "FATAL")
			return
		}
		configPath = filepath.Join(programFiles, "Sorty", "config.json")
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
