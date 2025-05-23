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
	// check if its document
	for _, ext := range FileExtArr.Documents {
		if ext != fileExt {
			continue
		} else {
			filename := path.Base(filePath)
			newPath := fmt.Sprintf("/home/thedevisi/Documents/%v", filename)
			oldPath, _ := filepath.Abs(filePath)
			os.Rename(oldPath, newPath)
		}
	}
	for _, ext := range FileExtArr.Images {
		if ext != fileExt {
			continue
		} else {
			filename := path.Base(filePath)
			newPath := fmt.Sprintf("/home/thedevisi/Pictures/%v", filename)
			oldPath, _ := filepath.Abs(filePath)
			os.Rename(oldPath, newPath)
		}
	}

	for _, ext := range FileExtArr.Videos {
		if ext != fileExt {
			continue
		} else {
			filename := path.Base(filePath)
			newPath := fmt.Sprintf("/home/thedevisi/Videos/%v", filename)
			oldPath, _ := filepath.Abs(filePath)
			os.Rename(oldPath, newPath)
		}
	}

	for _, ext := range FileExtArr.Music {
		if ext != fileExt {
			continue
		} else {
			filename := path.Base(filePath)
			newPath := fmt.Sprintf("/home/thedevisi/Music/%v", filename)
			oldPath, _ := filepath.Abs(filePath)
			os.Rename(oldPath, newPath)
		}
	}

}
