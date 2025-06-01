//go:build windows
// +build windows

package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sorty/logger"
	"syscall"

	"github.com/rs/zerolog"
	"golang.org/x/sys/windows"
)

// Structure of config.json
type Config struct {
	LogLevel     int16               `json:"log_level"`
	WatchFolder  string              `json:"watch_folder"`
	MonitorFiles map[string][]string `json:"monitor_files"`
}

var configPath string
var log *zerolog.Logger

func init() {
	localAppData := os.Getenv("LOCALAPPDATA")
	configPath = filepath.Join(localAppData, "sorty", "config.json")
	os.MkdirAll(filepath.Dir(configPath), 0755)

	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

var (
	modShell32               = syscall.NewLazyDLL("shell32.dll")
	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")

	// FOLDERID_Documents: {FDD39AD0-238F-46AF-ADB4-6C85480369C7}
	folderDocuments = windows.KNOWNFOLDERID{Data1: 0xFDD39AD0, Data2: 0x238F, Data3: 0x46AF, Data4: [8]byte{0xAD, 0xB4, 0x6C, 0x85, 0x48, 0x03, 0x69, 0xC7}}

	// FOLDERID_Pictures: {33E28130-4E1E-4676-835A-98395C3BC3BB}
	folderPictures = windows.KNOWNFOLDERID{Data1: 0x33E28130, Data2: 0x4E1E, Data3: 0x4676, Data4: [8]byte{0x83, 0x5A, 0x98, 0x39, 0x5C, 0x3B, 0xC3, 0xBB}}

	// FOLDERID_Videos: {18989B1D-99B5-455B-841C-AB7C74E4DDFC}
	folderVideos = windows.KNOWNFOLDERID{Data1: 0x18989B1D, Data2: 0x99B5, Data3: 0x455B, Data4: [8]byte{0x84, 0x1C, 0xAB, 0x7C, 0x74, 0xE4, 0xDD, 0xFC}}

	// FOLDERID_Music: {4BD8D571-6D19-48D3-BE97-422220080E43}
	folderMusic = windows.KNOWNFOLDERID{Data1: 0x4BD8D571, Data2: 0x6D19, Data3: 0x48D3, Data4: [8]byte{0xBE, 0x97, 0x42, 0x22, 0x20, 0x08, 0x0E, 0x43}}

	// FOLDERID_Downloads: {374DE290-123F-4565-9164-39C4925E467B}
	folderDownloads = windows.KNOWNFOLDERID{Data1: 0x374DE290, Data2: 0x123F, Data3: 0x4565, Data4: [8]byte{0x91, 0x64, 0x39, 0xC4, 0x92, 0x5E, 0x46, 0x7B}}
)

// logger init
func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

func CheckIfConfigExists() {

	_, err := os.Stat(configPath)
	if errors.Is(err, os.ErrNotExist) {
		createConfig()
	}
}

func createConfig() {
	configFile := generateConfig()

	data, err := json.MarshalIndent(configFile, "", " ")
	if err != nil {
		// ...
	}
	os.WriteFile(configPath, data, 0644)

}

func generateConfig() Config {
	log.Info().Msg("Can't find config. Creating a new one")
	_, err := os.Create(configPath)
	if err != nil {
		//...
	}

	picturesPath, _ := windows.KnownFolderPath(&folderPictures, 0)
	documentsPath, _ := windows.KnownFolderPath(&folderDocuments, 0)
	musicPath, _ := windows.KnownFolderPath(&folderMusic, 0)
	videosPath, _ := windows.KnownFolderPath(&folderVideos, 0)

	defaultMonitorFolders := map[string][]string{
		picturesPath:  {".png", ".jpeg", ".jpg", ".webp"},
		documentsPath: {".pdf", ".doc", ".docx"},
		musicPath:     {".mp3", ".wav", ".waw", ".test"},
		videosPath:    {".mp4", ".mov", ".mkv"},
	}
	downloadsPath, _ := windows.KnownFolderPath(&folderDownloads, 0)
	var configFile Config = Config{
		LogLevel:     1,
		WatchFolder:  downloadsPath,
		MonitorFiles: defaultMonitorFolders,
	}
	return configFile
}
