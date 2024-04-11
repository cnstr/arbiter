package main

import (
	"github.com/cnstr/arbiter/v2/internal/arbiter"
	"github.com/cnstr/arbiter/v2/internal/database"
	"github.com/cnstr/arbiter/v2/internal/typesense"
	"github.com/cnstr/arbiter/v2/internal/utils"
)

func main() {
	utils.GetRuntimeVersion()
	indexClient := typesense.CreateClient()
	typesense.EnsureApiKey(indexClient)
	typesense.EnsureCollections(indexClient)

	databaseClient := database.CreateClient()
	database.RunMigrations(databaseClient)
	database.CloseClient(databaseClient)

	arbiter.StartServer()
}
