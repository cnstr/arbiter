package main

import (
	"github.com/cnstr/arbiter/v2/internal/arbiter"
	"github.com/cnstr/arbiter/v2/internal/typesense"
	"github.com/cnstr/arbiter/v2/internal/utils"
)

func main() {
	utils.GetRuntimeVersion()
	client := typesense.CreateClient()
	typesense.EnsureApiKey(client)
	typesense.EnsureCollections(client)
	arbiter.StartServer()
}
