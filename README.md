# 🚀 Sorty

*Your Smart File Organization Assistant*

[![Go Version](https://img.shields.io/badge/Go-1.24.3-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)](https://www.microsoft.com/)
[![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)](https://www.linux.org/)

![Sorty Demo](demo/demo.mp4)


## 🌟 Overview

Sorty is your intelligent file organization companion that automatically keeps your downloads folder tidy. It watches for new files and instantly moves them to the appropriate location based on their type - be it documents, pictures, music, or videos.

## ✨ Key Features

### 🔄 Real-time File Management
- **Instant Detection**: Powered by [`fsnotify`](https://pkg.go.dev/github.com/fsnotify/fsnotify) for real-time file monitoring
- **Smart Categorization**: Automatically sorts files based on their extensions
- **Zero Configuration**: Works out of the box with sensible defaults

### 🛡️ Platform Support
- **Windows Integration**
  - System tray icon for easy quit
  - Registry-based auto-start (in progress)
- **Linux Support**
  - Systemd service integration
  - XDG compliant configuration
  - Desktop environment integration

### 🔧 Advanced Features
- **Configurable Rules**: Customize file sorting patterns in config file
- **Logging System**: Detailed activity logs using [`zerolog`](https://pkg.go.dev/github.com/rs/zerolog)
- **Error Recovery**: Robust error handling and recovery mechanisms

## 🚀 Quick Start

### Prerequisites
- Go 1.24.3 or higher
- Windows 10/11 or Linux (kernel 5.x recommended)

### Installation

1. **Get the Code**
```bash
git clone https://github.com/TheDevisi/sorty.git
cd sorty
```

2. **Build the Project**
```bash
# For Linux
go mod tidy
go build -o sorty

# For Windows
go mod tidy
GOOS=windows go build -ldflags="-H windowsgui" -o sorty.exe
```

3. **Run Sorty**
```bash
# Linux
./sorty # if it first startup, run with sudo to enable auto-start
sudo systemctl enable sorty
# Windows
sorty.exe
```

## ⚙️ Configuration

Sorty creates its configuration automatically on first run:

- **Windows**: `C:\Users\<username>\AppData\Local\sorty\config.json`
- **Linux**: `~/.config/sorty/config.json`

### Sample Configuration (Linux)
```json
{
 "log_level": 1,
 "watch_folder": "/home/<username>/Downloads",
 "monitor_files": {
  "/home/<username>/Documents": [
   ".pdf",
   ".doc",
   ".docx"
  ],
  "/home/<username>/Music": [
   ".mp3",
   ".wav",
   ".waw",
   ".test"
  ],
  "/home/<username>/Pictures": [
   ".png",
   ".jpeg",
   ".jpg",
   ".webp"
  ],
  "/home/<username>/Videos": [
   ".mp4",
   ".mov",
   ".mkv"
  ]
    // add your own path(s)&extension(s) here if you want
 }
}
```

### Sample Configuration (Windows)

```  json
{
 "log_level": 1,
 "watch_folder": "C:\\Users\\<username>\\Downloads",
 "monitor_files": {
  "C:\\Users\\<username>\\Dildo": [
   ".pdf",
   ".doc",
   ".docx"
  ],
  "C:\\Users\\<username>\\Music": [
   ".mp3",
   ".wav",
   ".waw",
   ".test"
  ],
  "C:\\Users\\<username>\\Pictures": [
   ".png",
   ".jpeg",
   ".jpg",
   ".webp"
  ],
  "C:\\Users\\<username>\\Videos": [
   ".mp4",
   ".mov",
   ".mkv"
  ]
    // add your own path(s)&extension(s) here if you want
 }
}
```
## 📁 Project Structure

```
sorty/
├── 📂 config/           # Configuration management
├── 📂 internal/         # Private application code
├── 📂 logger/           # Logging setup
├── 📂 pkg/             # Public libraries
│   ├── 📂 utils/       # Utility functions
│   └── 📂 watcher/     # File system watcher
└── 📄 main.go          # Application entry point
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/okak-feature`)
3. Commit your changes (`git commit -m 'Add random feature'`)
4. Push to the branch (`git push origin feature/okak-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [fsnotify](https://github.com/fsnotify/fsnotify) for file system notifications
- [zerolog](https://github.com/rs/zerolog) for structured logging
- [systray](https://github.com/getlantern/systray) for system tray support

---

Made by [TheDevisi](https://thedevisi.com) with ♥️
