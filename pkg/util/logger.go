package util

import "log"

func LogInfo(message string) {
	log.Println("[INFO]", message)
}

func LogError(message string, err error) {
	log.Fatalf("[ERROR]: %s - %v", message, err)
}
