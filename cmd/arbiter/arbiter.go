package main

import (
	"github.com/cnstr/arbiter/v2/internal/arbiter"
	"github.com/cnstr/arbiter/v2/internal/utils"
)

func main() {
	utils.GetRuntimeVersion()
	arbiter.StartServer()
}
