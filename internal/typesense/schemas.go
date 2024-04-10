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

func PackageSchema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name: "packages",
		Fields: []api.Field{
			{
				Name:     "id",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.False(),
			},
			{
				// Legacy field
				Name:     "package",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.False(),
			},
			{
				Name:     "package_id",
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
				Name:     "author",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.True(),
			},
			{
				Name:     "maintainer",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.True(),
			},
			{
				Name:     "section",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.True(),
			},
			{
				Name:     "version",
				Type:     "string",
				Facet:    pointer.False(),
				Optional: pointer.False(),
			},
			{
				// Legacy field
				Name:     "repositoryTier",
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
		},
	}
}
