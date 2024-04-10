package typesense

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/typesense/typesense-go/typesense"
)

func loadEnvOrFatal(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("[typesense] %s is required", key)
	}

	return value
}

func CreateClient() *typesense.Client {
	client := typesense.NewClient(
		typesense.WithServer("http://localhost:7700"),
		typesense.WithAPIKey("typesense"),
		typesense.WithConnectionTimeout(5*time.Second),
	)

	return client
}

func EnsureCollections(client *typesense.Client) bool {
	collections, err := client.Collections().Retrieve(context.Background())
	if err != nil {
		log.Println("[typesense] Error retrieving collections:", err)
		return false
	}

	repositoryPresent, packagePresent := false, true
	for _, collection := range collections {
		if collection.Name == "repositories" {
			log.Println("[typesense] Found existing repositories collection")
			repositoryPresent = true
		}

		// if collection.Name == "packages" {
		// 	packagePresent = true
		// }
	}

	if packagePresent && repositoryPresent {
		return true
	}

	log.Println("[typesense] Creating repositories collection")
	_, err = client.Collections().Create(
		context.Background(),
		RepositorySchema(),
	)

	if err != nil {
		log.Println("[typesense] Error creating repositories collection:", err)
		return false
	}

	return true
}
