package database

import _ "embed"
import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations.sql
var MigrationSql string

func RunMigrations(Client *pgxpool.Pool) {
	statements := strings.Split(MigrationSql, ";")

	for i, statement := range statements {
		rows, err := Client.Query(context.Background(), statement)
		rows.Close()

		if err != nil {
			log.Fatal("[database] Failed to run migration[", i, "]:", err)
		}
	}

	log.Println("[database] Finished running migrations")
}
