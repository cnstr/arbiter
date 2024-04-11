package database

import (
	"context"
	"log"

	"github.com/cnstr/arbiter/v2/internal/utils"
	"github.com/jackc/pgx/v5"
)

func CreateClient() *pgx.Conn {
	connectionUrl := utils.LoadEnvOrFatal("POSTGRES_URL")
	conn, err := pgx.Connect(context.Background(), connectionUrl)
	if err != nil {
		log.Println("[database] Failed to connect to the database:", err)
		return nil
	}

	return conn
}

func CloseClient(conn *pgx.Conn) {
	conn.Close(context.Background())
}
