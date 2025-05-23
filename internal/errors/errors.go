package errors

import (
	"log"
	// "errors"
)

func ErrorsHandler(err error, level string) {
	log.Println(err)
}
