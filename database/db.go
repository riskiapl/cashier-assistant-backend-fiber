package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riskiapl/fiber-app/config"
)

var DB *pgxpool.Pool

func ConnectDB() {
	config.LoadEnv()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.GetEnv("DB_USER", "fiber_user"),
		config.GetEnv("DB_PASSWORD", "fiber_pass"),
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "5432"),
		config.GetEnv("DB_NAME", "fiber_db"),
	)

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to the database!")
}
