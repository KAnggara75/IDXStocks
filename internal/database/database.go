package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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

	// Run migrations
	Migrate()
}

func Migrate() {
	ctx := context.Background()

	// Ensure schema exists
	_, err := Pool.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS idxstock")
	if err != nil {
		logrus.Fatalf("Failed to create schema idxstock: %v", err)
	}

	// Get migration files
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		logrus.Fatalf("Failed to read migrations directory: %v", err)
	}

	sort.Strings(files)

	for _, file := range files {
		logrus.Infof("Running migration: %s", file)
		content, err := os.ReadFile(file)
		if err != nil {
			logrus.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		_, err = Pool.Exec(ctx, string(content))
		if err != nil {
			logrus.Fatalf("Failed to execute migration %s: %v", file, err)
		}
	}

	logrus.Info("Migrations completed successfully")
}
