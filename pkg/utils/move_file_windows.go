//go:build windows
// +build windows

package utils

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modShell32               = syscall.NewLazyDLL("shell32.dll")
	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")

	// FOLDERID_Documents: {FDD39AD0-238F-46AF-ADB4-6C85480369C7}
	folderDocuments = windows.GUID{Data1: 0xFDD39AD0, Data2: 0x238F, Data3: 0x46AF, Data4: [8]byte{0xAD, 0xB4, 0x6C, 0x85, 0x48, 0x03, 0x69, 0xC7}}

	// FOLDERID_Pictures: {33E28130-4E1E-4676-835A-98395C3BC3BB}
	folderPictures = windows.GUID{Data1: 0x33E28130, Data2: 0x4E1E, Data3: 0x4676, Data4: [8]byte{0x83, 0x5A, 0x98, 0x39, 0x5C, 0x3B, 0xC3, 0xBB}}

	// FOLDERID_Videos: {18989B1D-99B5-455B-841C-AB7C74E4DDFC}
	folderVideos = windows.GUID{Data1: 0x18989B1D, Data2: 0x99B5, Data3: 0x455B, Data4: [8]byte{0x84, 0x1C, 0xAB, 0x7C, 0x74, 0xE4, 0xDD, 0xFC}}

	// FOLDERID_Music: {4BD8D571-6D19-48D3-BE97-422220080E43}
	folderMusic = windows.GUID{Data1: 0x4BD8D571, Data2: 0x6D19, Data3: 0x48D3, Data4: [8]byte{0xBE, 0x97, 0x42, 0x22, 0x20, 0x08, 0x0E, 0x43}}
)

func GetDefaultFolder(name string) string {
	var folderID *windows.GUID

	switch name {
	case "Documents":
		folderID = &folderDocuments
	case "Pictures":
		folderID = &folderPictures

	case "Music":
		folderID = &folderMusic

	case "Videos":
		folderID = &folderVideos
	default:
		return ""
	}

	var path *uint16
	hr, _, _ := procSHGetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(folderID)),
		uintptr(0),
		uintptr(0),
		uintptr(unsafe.Pointer(&path)),
	)
	if hr != 0 {
		return ""
	}
	defer windows.CoTaskMemFree(unsafe.Pointer(path))
	return windows.UTF16PtrToString(path)
}
