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

	// 1. Check if the table idxstock.stocks already exists
	var tableExists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE  table_schema = 'idxstock'
			AND    table_name   = 'stocks'
		);
	`
	err := Pool.QueryRow(ctx, checkQuery).Scan(&tableExists)
	if err != nil {
		logrus.Errorf("Failed to check if table exists: %v", err)
	}

	if tableExists {
		logrus.Info("Table idxstock.stocks already exists, skipping migrations")
		return
	}

	// 2. Ensure schema exists
	_, err = Pool.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS idxstock")
	if err != nil {
		logrus.Fatalf("Failed to create schema idxstock: %v", err)
	}

	// 3. Get migration files
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
