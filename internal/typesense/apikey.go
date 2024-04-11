package typesense

import (
	"context"
	"log"
	"strings"

	"github.com/cnstr/arbiter/v2/internal/utils"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func EnsureApiKey(client *typesense.Client) bool {
	desiredApiKey := utils.LoadEnvOrFatal("TYPESENSE_PUBLIC_API_KEY")
	keys, err := client.Keys().Retrieve(context.Background())
	if err != nil {
		log.Println("[typesense] Error retrieving API keys:", err)
		return false
	}

	for _, key := range keys {
		if strings.HasPrefix(desiredApiKey, *key.ValuePrefix) {
			log.Println("[typesense] Found existing public API key")
			return true
		}
	}

	_, err = client.Keys().Create(context.Background(), &api.ApiKeySchema{
		Description: "Public API key for search",
		Actions:     []string{"documents:search"},
		Collections: []string{"*"},
		Value:       &desiredApiKey,
	})

	if err != nil {
		log.Println("[typesense] Error creating public API key:", err)
		return false
	}

	log.Println("[typesense] Created public API search key")
	return true
}
