package typesense

import (
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func RepositorySchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: "repositories",
		Fields: []api.Field{
			{
				Name:     "id",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.False(),
			},
			{
				// Legacy field
				Name:     "slug",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.False(),
			},
			{
				Name:     "name",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.True(),
			},
			{
				Name:     "description",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.True(),
			},
			{
				Name:     "aliases",
				Type:     "string[]",
				Facet:    pointer.False(),
				Optional: pointer.False(),
			},
			{
				Name:     "uri",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.False(),
			},
			{
				// Legacy field
				Name:     "tier",
				Type:     "int32",
				Facet:    pointer.True(),
				Optional: pointer.False(),
			},
			{
				Name:     "quality",
				Type:     "int32",
				Facet:    pointer.True(),
				Optional: pointer.False(),
			},
			{
				Name:     "bootstrap",
				Type:     "bool",
				Facet:    pointer.True(),
				Optional: pointer.False(),
			},
			{
				// Legacy field
				Name:     "isBootstrap",
				Type:     "bool",
				Facet:    pointer.True(),
				Optional: pointer.False(),
			},
		},
	}
}
