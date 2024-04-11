package database

import (
	"context"
	"log"

	"github.com/cnstr/arbiter/v2/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateClient() *pgxpool.Pool {
	connectionUrl := utils.LoadEnvOrFatal("POSTGRES_URL")
	conn, err := pgxpool.New(context.Background(), connectionUrl)
	if err != nil {
		log.Println("[database] Failed to connect to the database:", err)
		return nil
	}

	return conn
}

func CloseClient(conn *pgxpool.Pool) {
	conn.Close()
}
