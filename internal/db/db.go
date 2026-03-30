package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(
		context.Background(),
		DATABASE_URL)

	if err != nil {
		return nil, err
	}
	
	return conn, nil
}
