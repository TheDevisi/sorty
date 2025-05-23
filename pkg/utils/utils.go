// some utils for any other pkgs.
package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sorty/internal/errors"
	"time"
)

type FileExtStruct struct {
	Images    []string
	Videos    []string
	Documents []string
	Music     []string
}

type FolderStruct struct {
	LinuxImages      string
	LinuxMusic       string
	LinuxDocuments   string
	LinuxVideos      string
	WindowsImages    string
	WindowsMusic     string
	WindowsDocuments string
	WindowsVideos    string
}

var FileExtArr FileExtStruct = FileExtStruct{
	Images:    []string{".png", ".wepb", ".jpg", ".jpeg"},
	Videos:    []string{".mp4", ".mov", ".mkv"},
	Documents: []string{".pdf"},         // TODO: add any other later  :3
	Music:     []string{".mp3", ".wav"}, // answer is below ^^
}

func GetOperatingSystem() string {
	OS := runtime.GOOS
	return OS
}

// README: not used for now. will be removed
// func getUserName() string {
// 	username, err := user.Current()
// 	if err != nil {
// 		errors.ErrorsHandler(err, "FATAL")
// 	}

// 	return username.Username
// }

func MoveFile(filePath string) {
	canMove := checkForTempFile(filePath)
	if canMove {
		fileExt := filepath.Ext(filePath)
		PlaceToMove(fileExt, filePath)
	}
}

// this function checking if current file is temp file. yes ? - wait. no? - move.
func checkForTempFile(filename string) bool {
	var tmpFileExts []string = []string{".part", ".crdownload", ".tmp", ".temp", ".dmp", ".adadownload", ".cache", ".partial", ""}
	for _, tmpFileExt := range tmpFileExts {
		if path.Ext(filename) == tmpFileExt {
			time.Sleep(1 * time.Second)
			continue
		}
		fileStat, err := os.Stat(filename)
		if err != nil {
			errors.ErrorsHandler(err, "FATAL")
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

func PlaceToMove(fileExt string, filePath string) {
	for _, ext := range FileExtArr.Documents {
		if ext == fileExt {
			filename := filepath.Base(filePath)
			fileFolder := GetDefaultFolder("Documents")
			if fileFolder == "" {
				errors.ErrorsHandler(
					fmt.Errorf("GetDefaultFolder returned empty path for Documents"), "ERROR")
				return
			}
			newPath := filepath.Join(fileFolder, filename)
			oldPath, _ := filepath.Abs(filePath)
			if err := os.Rename(oldPath, newPath); err != nil {
				errors.ErrorsHandler(err, "ERROR")
			}
			return
		}
	}
	for _, ext := range FileExtArr.Images {
		if ext == fileExt {
			filename := filepath.Base(filePath)
			fileFolder := GetDefaultFolder("Pictures")
			newPath := filepath.Join(fileFolder, filename)
			oldPath, _ := filepath.Abs(filePath)
			if err := os.Rename(oldPath, newPath); err != nil {
				errors.ErrorsHandler(err, "ERROR")
			}
			return
		}
	}
	for _, ext := range FileExtArr.Videos {
		if ext == fileExt {
			filename := filepath.Base(filePath) // исправлено!
			fileFolder := GetDefaultFolder("Videos")
			newPath := filepath.Join(fileFolder, filename)
			oldPath, _ := filepath.Abs(filePath)
			if err := os.Rename(oldPath, newPath); err != nil {
				errors.ErrorsHandler(err, "ERROR")
			}
			return
		}
	}
	for _, ext := range FileExtArr.Music {
		if ext == fileExt {
			filename := filepath.Base(filePath) // исправлено!
			fileFolder := GetDefaultFolder("Music")
			newPath := filepath.Join(fileFolder, filename)
			oldPath, _ := filepath.Abs(filePath)
			if err := os.Rename(oldPath, newPath); err != nil {
				errors.ErrorsHandler(err, "ERROR")
			}
			return
		}
	}
}
