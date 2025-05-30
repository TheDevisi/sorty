// This package watches directory 24/7, downloads for ex.

package watcher

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"sorty/config"
	"sorty/internal/errors"
	"sorty/internal/utils"
	"sorty/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

var log *zerolog.Logger

func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

func monitorFolder() string {
	var fileData *config.Config
	OS := utils.GetOperatingSystem()
	userInfo, err := user.Current()
	if err != nil {
		errors.ErrorsHandler(err, "FATAL")
	}
	userHome := userInfo.HomeDir
	if OS != "windows" {
		file, err := os.ReadFile(fmt.Sprintf("%v/.config/sorty/config.json", userHome))
		if err != nil {
			errors.ErrorsHandler(err, "FATAL")
		}

		json.Unmarshal(file, &fileData)
	} else {
		file, err := os.ReadFile("C\\Program Files\\sorty\\config.json")
		if err != nil {
			errors.ErrorsHandler(err, "FATAL")
		}
		json.Unmarshal(file, &fileData)

	}
	return fileData.WatchFolder

}
// WatchDirectory monitors a configured directory for file system events and processes file changes as they occur.
//
// Sets up a file system watcher on the directory specified in the configuration, handling file events by invoking file processing logic and logging errors as needed. The function blocks indefinitely to keep the watcher active.
func WatchDirectory() {
	folderPath := monitorFolder()
	// now let's setup fsnotify...
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// log error with stack trace before handling errors
		log.Error().Stack().Err(err).Msg("Unable to create fsnotify watcher")
		errors.ErrorsHandler(err, "")
	}

	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Printf("Event: %v", event)
				utils.MoveFile(event.Name)
				if event.Has(fsnotify.Write) {
					fmt.Printf("Modified file: %v", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error().Stack().Err(err).Msg("Watcher error")
			}
		}
	}()

	// Add a path.
	err = watcher.Add(folderPath)
	if err != nil {
		log.Error().Stack().Str("path", folderPath).Err(err).Msg("Failed to add directory watch")
	} else {
		log.Info().Str("path", folderPath).Msg("Starting directory watch")
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
