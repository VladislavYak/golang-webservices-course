package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgPool() (*pgxpool.Pool, error) {
	connStr := "postgres://postgres:love@localhost:5432/golang?sslmode=disable"

	Pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	if err := Pool.Ping(context.Background()); err != nil {
		fmt.Println("i was here NewPgPool")
		fmt.Println("errr", err)
		return nil, err
	}
	fmt.Println("connected to pg")

	return Pool, nil

}
