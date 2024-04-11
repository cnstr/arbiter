package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PruneStaleData(Client *pgxpool.Pool, Ids []string) {
	rows, err := Client.Query(
		context.Background(),
		"UPDATE repository SET visible = false WHERE id != ALL($1)",
		Ids,
	)

	rows.Close()

	if err != nil {
		log.Println("[database] Failed to prune stale repositories:", err)
	}

	rows, err = Client.Query(
		context.Background(),
		"UPDATE package SET visible = false WHERE repository_id != ALL($1)",
		Ids,
	)

	rows.Close()

	if err != nil {
		log.Println("[database] Failed to prune stale packages:", err)
	}

	log.Println("[database] Finished pruning stale data")
	return
}
