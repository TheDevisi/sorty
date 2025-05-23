package errors

import (
	"sorty/logger"

	"github.com/rs/zerolog"
)

var log *zerolog.Logger

func init() {
	config := logger.NewLogConfig()
	log = logger.NewLogger(config)
}

var ErrNoData string = "seems like there's no setup value. using a default one if it provided"

func ErrorsHandler(err error, level string) {
	switch level {
	case "FATAL":
		log.Fatal().Err(err).Msg("Fatal error occurred")
	case "ERROR":
		log.Error().Err(err).Msg("Error occurred")
	default:
		log.Warn().Err(err).Msg("Warning occurred")
	}
}
