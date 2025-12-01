package utils

import (
	"log"
)

func LogError(errMessage string) {
	log.Printf("ERROR: %s", errMessage)
}