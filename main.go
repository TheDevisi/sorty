package main

import (
	_ "embed"
	"sorty/logger"
	"sorty/pkg/utils"
	"sorty/pkg/watcher"

	"github.com/rs/zerolog"
)

var log *zerolog.Logger

func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

func main() {

	utils.EnableAutoStart()
	watcher.WatchDirectory()

}
