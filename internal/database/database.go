package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var Pool *pgxpool.Pool

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logrus.Fatal("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logrus.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logrus.Fatalf("Unable to create connection pool: %v", err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		logrus.Fatalf("Unable to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
}
