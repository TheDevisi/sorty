//go:build linux
// +build linux

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"sorty/logger"

	"github.com/rs/zerolog"
)

// Structure of config.json
type Config struct {
	LogLevel     int16               `json:"log_level"`
	WatchFolder  string              `json:"watch_folder"`
	MonitorFiles map[string][]string `json:"monitor_files`
}

var log *zerolog.Logger

// logger init
func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

func CheckIfConfigExists() {
	userName, _ := user.Current()
	var systemConfigFolder string = fmt.Sprintf("%v/.config/sorty/config.json", userName.HomeDir)

	_, err := os.Stat(systemConfigFolder) // path to config. HARDCODED >:)
	if errors.Is(err, os.ErrNotExist) {
		createConfig()
	}
}

func createConfig() {
	configFile, systemConfigFolder := generateConfig()

	data, err := json.MarshalIndent(configFile, "", " ")
	if err != nil {
		// ...
	}
	os.WriteFile(systemConfigFolder, data, 0644)

}

func generateConfig() (Config, string) {

	userName, _ := user.Current()
	os.Mkdir(fmt.Sprintf("%v/.config/sorty", userName.HomeDir), 0755)

	var systemConfigFolder string = fmt.Sprintf("%v/.config/sorty/config.json", userName.HomeDir)
	log.Info().Msg("Can't find config. Creating a new one")
	_, err := os.Create(systemConfigFolder)
	if err != nil {
		//...
	}
	downloadsPath := fmt.Sprintf("/home/%v/Downloads", userName.Username)

	imgsFolder := fmt.Sprintf("/home/%v/Pictures", userName.Username)
	docsFolder := fmt.Sprintf("/home/%v/Documents", userName.Username)
	muscFolder := fmt.Sprintf("/home/%v/Music", userName.Username)
	vidsFolder := fmt.Sprintf("/home/%v/Videos", userName.Username)

	defaultMonitorFolders := map[string][]string{
		imgsFolder: {".png", ".jpeg", ".jpg", ".webp"},
		docsFolder: {".pdf", ".doc", ".docx"},
		muscFolder: {".mp3", ".wav", ".waw", ".test"},
		vidsFolder: {".mp4", ".mov", ".mkv"},
	}
	var configFile Config = Config{
		LogLevel:     1,
		WatchFolder:  downloadsPath,
		MonitorFiles: defaultMonitorFolders,
	}
	return configFile, systemConfigFolder
}
