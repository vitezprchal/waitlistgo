package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func InitDatabase() Repository {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS subscribers (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT
		)
	`)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		os.Exit(1)
	}

	return NewSubscriberRepository(conn)
}
