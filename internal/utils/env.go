package utils

import (
	"log"
	"os"
)

func LoadEnvOrFatal(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("[typesense] %s is required", key)
	}

	return value
}
