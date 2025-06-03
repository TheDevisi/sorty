package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sorty/config"
	"sorty/logger"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog"
)

var log *zerolog.Logger
var settingsChan chan bool

func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
	settingsChan = make(chan bool)
}

func ShowSettingsWindow() {
	settingsChan <- true
}

func RunSettingsWindow() {
	for {
		select {
		case <-settingsChan:
			showSettingsWindow()
		}
	}
}

func showSettingsWindow() {
	myApp := app.New()
	window := myApp.NewWindow("Sorty Settings")

	watchFolderEntry := widget.NewEntry()
	watchFolderEntry.SetPlaceHolder("Watch Folder Path")

	configData := loadConfig()
	if configData != nil {
		watchFolderEntry.SetText(configData.WatchFolder)
	}

	mappingsList := widget.NewList(
		func() int {
			if configData == nil {
				return 0
			}
			return len(configData.MonitorFiles)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Folder: "),
				widget.NewLabel("Extensions: "),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if configData == nil {
				return
			}
			var folder string
			var extensions []string
			i := 0
			for f, exts := range configData.MonitorFiles {
				if i == int(id) {
					folder = f
					extensions = exts
					break
				}
				i++
			}
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText("Folder: " + folder)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("Extensions: " + strings.Join(extensions, ", "))
		},
	)

	addMappingBtn := widget.NewButton("Add Mapping", func() {
		folderEntry := widget.NewEntry()
		folderEntry.SetPlaceHolder("Target Folder Path")
		extEntry := widget.NewEntry()
		extEntry.SetPlaceHolder("Extensions (comma-separated, e.g., .pdf,.doc,.docx)")

		dialog.ShowCustom("Add File Mapping", "Cancel", container.NewVBox(
			widget.NewLabel("Target Folder:"),
			folderEntry,
			widget.NewLabel("File Extensions:"),
			extEntry,
			widget.NewButton("Add", func() {
				if configData != nil && folderEntry.Text != "" && extEntry.Text != "" {
					extensions := strings.Split(extEntry.Text, ",")
					for i, ext := range extensions {
						extensions[i] = strings.TrimSpace(ext)
					}
					configData.MonitorFiles[folderEntry.Text] = extensions
					mappingsList.Refresh()
				}
			}),
		), window)
	})

	saveButton := widget.NewButton("Save", func() {
		if configData != nil {
			configData.WatchFolder = watchFolderEntry.Text
			saveConfig(configData)
		}
		window.Close()
	})

	form := container.NewVBox(
		widget.NewLabel("Settings"),
		widget.NewLabel("Watch Folder:"),
		watchFolderEntry,
		widget.NewLabel("File Extension Mappings:"),
		mappingsList,
		addMappingBtn,
		saveButton,
	)

	window.SetContent(form)
	window.Resize(fyne.NewSize(600, 400))
	window.ShowAndRun()
}

func loadConfig() *config.Config {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user home directory")
		return nil
	}

	var configPath string
	if os.Getenv("OS") == "Windows_NT" {
		localAppData := os.Getenv("LOCALAPPDATA")
		configPath = filepath.Join(localAppData, "sorty", "config.json")
	} else {
		configPath = filepath.Join(userHome, ".config", "sorty", "config.json")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read config file")
		return nil
	}

	var configData config.Config
	if err := json.Unmarshal(data, &configData); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal config")
		return nil
	}

	return &configData
}

func saveConfig(configData *config.Config) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user home directory")
		return
	}

	var configPath string
	if os.Getenv("OS") == "Windows_NT" {
		localAppData := os.Getenv("LOCALAPPDATA")
		configPath = filepath.Join(localAppData, "sorty", "config.json")
	} else {
		configPath = filepath.Join(userHome, ".config", "sorty", "config.json")
	}

	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal config")
		return
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		log.Error().Err(err).Msg("Failed to write config file")
		return
	}

	log.Info().Msg("Settings saved successfully")
}
