package logger

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

type LogConfig struct {
	Level  int16
	Format string
}

func getInt(value string, defaultValue int16) (int16, error) {
	data := os.Getenv(value)
	if data == "" {
		return defaultValue, nil
	}
	strConvInt, err := strconv.Atoi(data)
	if err != nil {
		return defaultValue, err
	}

	if strConvInt < 0 && strConvInt > 65535 {
		return defaultValue, nil
	}

	intFromEnv := int16(strConvInt)

	return intFromEnv, nil

}

func getStr(value, defaultValue string) (string, error) {
	data := os.Getenv(value)
	if data == "" {
		return defaultValue, nil
	}

	return data, nil
}

func NewLogConfig() *LogConfig {
	level, err := getInt("LOG_LEVEL", 0)
	if err != nil {
		level = 0 // fallback to default if error occurs
	}
	format, err := getStr("LOG_FORMAT", "json")
	if err != nil {
		format = "json" // fallback to default if error occurs
	}
	return &LogConfig{
		Level:  level,
		Format: format,
	}
}

func NewLogger(config *LogConfig) *zerolog.Logger {
	var logger zerolog.Logger
	zerolog.SetGlobalLevel(zerolog.Level(config.Level))
	if config.Format == "json" {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		return &logger
	} else {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		return &logger
	}

}
