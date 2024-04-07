package utils

import (
	"log"
	"runtime/debug"
)

func GetRuntimeVersion() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Println("[arbiter] Could not pull build information")
		return
	}

	log.Println("[arbiter] Version:", buildInfo.Main.Version)
}
