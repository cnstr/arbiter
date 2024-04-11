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
	serverHost := loadEnvOrFatal("TYPESENSE_HOST")
	privateKey := loadEnvOrFatal("TYPESENSE_PRIVATE_API_KEY")

	client := typesense.NewClient(
		typesense.WithServer(serverHost),
		typesense.WithAPIKey(privateKey),
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

	repositoryPresent, packagePresent := false, false
	for _, collection := range collections {
		if collection.Name == "repositories" {
			log.Println("[typesense] Found existing repositories collection")
			repositoryPresent = true
			continue
		}

		if collection.Name == "packages" {
			log.Println("[typesense] Found existing packages collection")
			packagePresent = true
			continue
		}
	}

	if packagePresent && repositoryPresent {
		return true
	}

	if !repositoryPresent {
		log.Println("[typesense] Creating repositories collection")
		_, err = client.Collections().Create(
			context.Background(),
			RepositorySchema(),
		)

		if err != nil {
			log.Println("[typesense] Error creating repositories collection:", err)
			return false
		}
	}

	if !packagePresent {
		log.Println("[typesense] Creating packages collection")
		_, err = client.Collections().Create(
			context.Background(),
			PackageSchema(),
		)

		if err != nil {
			log.Println("[typesense] Error creating packages collection:", err)
			return false
		}
	}

	return true
}
