// This package watches directory 24/7, downloads for ex.

package watcher

import (
	"fmt"
	"log"
	"os/user"
	"sorty/internal/errors"
	"sorty/pkg/utils"

	"github.com/fsnotify/fsnotify"
)

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
				log.Println("event:", event)
				utils.MoveFile(event.Name)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

// now handle write events
