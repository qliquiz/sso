package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sso/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var configPath, migrationsPath, migrationsTable, command string

	flag.StringVar(&configPath, "config-path", "", "Path to config file")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to migration files")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Path to migration table")
	flag.StringVar(&command, "command", "up", "Migration command: 'up' or 'down'")
	flag.Parse()

	if configPath == "" || migrationsPath == "" {
		log.Fatal("error: migrations-path is a required flag")
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var cfg config.Config
	if err = yaml.Unmarshal(configData, &cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&x-migrations-table=%s",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DBName,
		cfg.Storage.SSLMode,
		cfg.Storage.MigTable,
	)

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	var migrateErr error
	switch command {
	case "up":
		migrateErr = m.Up()
		fmt.Println("Applying 'up' migrations...")
	case "down":
		migrateErr = m.Steps(-1) // rollbacks the last migration
		fmt.Println("Applying 'down' migration...")
	default:
		log.Fatalf("unknown command: %s. Use 'up' or 'down'", command)
	}

	if migrateErr != nil {
		if errors.Is(migrateErr, migrate.ErrNoChange) {
			fmt.Println("✅ No changes to apply.")
			return
		}
		log.Fatalf("migration failed: %v", migrateErr)
	}

	fmt.Println("✅ Migrations applied successfully!")
}
