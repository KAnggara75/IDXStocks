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

	// 1. Ensure schema exists
	_, err := Pool.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS idxstock")
	if err != nil {
		logrus.Fatalf("Failed to create schema idxstock: %v", err)
	}

	// 2. Create migrations tracking table
	_, err = Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS idxstock.schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT now()
		)
	`)
	if err != nil {
		logrus.Fatalf("Failed to create migration tracking table: %v", err)
	}

	// 3. Handle transition from old skip logic:
	// If stocks table exists but migrations haven't been tracked, mark 000 and 001 as applied
	var stocksExists bool
	err = Pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'idxstock' AND table_name = 'stocks')").Scan(&stocksExists)
	if err != nil {
		logrus.Debugf("Note: could not check if stocks table exists: %v", err)
	}

	if stocksExists {
		_, err = Pool.Exec(ctx, "INSERT INTO idxstock.schema_migrations (filename) VALUES ('migrations/000_board_type.sql'), ('migrations/001_create_stocks_table.sql') ON CONFLICT DO NOTHING")
		if err != nil {
			logrus.Debugf("Note: could not seed initial migrations: %v", err)
		}
	}

	// 4. Get migration files
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		logrus.Fatalf("Failed to read migrations directory: %v", err)
	}

	sort.Strings(files)

	for _, file := range files {
		// Check if already applied
		var applied bool
		err = Pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM idxstock.schema_migrations WHERE filename = $1)", file).Scan(&applied)
		if err != nil {
			logrus.Errorf("Failed to check migration status for %s: %v", file, err)
			continue
		}

		if applied {
			continue
		}

		logrus.Infof("Running migration: %s", file)
		cleanFile := filepath.Clean(file)
		// #nosec G304
		content, err := os.ReadFile(cleanFile)
		if err != nil {
			logrus.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		_, err = Pool.Exec(ctx, string(content))
		if err != nil {
			logrus.Fatalf("Failed to execute migration %s: %v", file, err)
		}

		// Mark as applied
		_, err = Pool.Exec(ctx, "INSERT INTO idxstock.schema_migrations (filename) VALUES ($1)", file)
		if err != nil {
			logrus.Fatalf("Failed to record migration %s: %v", file, err)
		}
	}

	logrus.Info("Migrations synchronization completed")
}
