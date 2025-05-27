// This package watches directory 24/7, downloads for ex.

package watcher

import (
	"fmt"
	"os/user"
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

// TODO replace with new config-based logic
// based on operating system returning "downloads path"
func downloadsPath(os string) string {
	if os == "linux" {
		userName, _ := user.Current()

		return fmt.Sprintf("%v/Downloads", userName.HomeDir)
	} else if os == "windows" {
		userName, err := user.Current()
		if err != nil {
			errors.ErrorsHandler(err, "")
		}
		return fmt.Sprintf("%v\\Downloads", userName.HomeDir) // idk if this works
	}
	return "~/Downloads" // Default to unix-style path
}

func WatchDirectory() {
	folderPath := downloadsPath(utils.GetOperatingSystem())
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
				log.Info().Msgf("Event: %v", event)
				utils.MoveFile(event.Name)
				if event.Has(fsnotify.Write) {
					log.Info().Msgf("Modified file: %v", event.Name)
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
