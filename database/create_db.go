package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riskiapl/fiber-app/config"
)

func CreateDatabase() {
	// Koneksi ke postgres default database
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres",
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", "password"),
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "5432"),
	)

	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close()

	// Cek apakah database sudah ada
	var exists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		"cashier_assistant_fiber").Scan(&exists)

	if err != nil {
		log.Fatalf("Error checking database existence: %v", err)
	}

	// Jika database belum ada, buat baru
	if !exists {
		_, err = conn.Exec(context.Background(),
			"CREATE DATABASE cashier_assistant_fiber")

		if err != nil {
			log.Fatalf("Error creating database: %v", err)
		}

		log.Println("Database cashier_assistant_fiber created successfully!")
	} else {
		log.Println("Database cashier_assistant_fiber already exists")
	}
}
