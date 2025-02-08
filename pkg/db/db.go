package dbpkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Configurable PostgreSQL credentials (using environment variables)
var (
	dbUser     = getEnv("PG_USER", "postgres")
	dbPassword = getEnv("PG_PASSWORD", "secret")
	dbHost     = getEnv("PG_HOST", "localhost")
	dbPort     = getEnv("PG_PORT", "5432")
	dbName     = getEnv("PG_DB", "mydatabase")
)

type Dbconn struct {
	Db *pgxpool.Pool
}

// Default DSN (for connecting to 'postgres' to check database existence)
var defaultDSN = fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
	dbUser, dbPassword, dbHost, dbPort)

// Target DSN (actual application database)
var targetDSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	dbUser, dbPassword, dbHost, dbPort, dbName)

// InitDB initializes the database connection and checks if the DB exists.
func InitDB() (*Dbconn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to PostgreSQL to check if the database exists
	conn, err := pgxpool.New(ctx, defaultDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Check if database exists
	exists, err := databaseExists(ctx, &Dbconn{Db: conn}, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %v", err)
	}

	// Create database if it doesn't exist
	if !exists {
		log.Printf("Database %s does not exist. Creating it...", dbName)
		if err := createDatabase(ctx, &Dbconn{Db: conn}, dbName); err != nil {
			return nil, fmt.Errorf("failed to create database: %v", err)
		}
		log.Println("Database created successfully!")
	}

	// Now, connect to the actual database
	pool, err := pgxpool.New(ctx, targetDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to target database: %v", err)
	}

	log.Println("Connected to PostgreSQL database successfully.")
	return &Dbconn{Db: pool}, nil
}

// databaseExists checks if a database already exists.
func databaseExists(ctx context.Context, dbconn *Dbconn, name string) (bool, error) {
	var exists bool
	err := dbconn.Db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", name).Scan(&exists)
	return exists, err
}

// createDatabase creates a new database.
func createDatabase(ctx context.Context, dbconn *Dbconn, name string) error {
	_, err := dbconn.Db.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s;", name))
	return err
}

// GetDB returns the database connection from Dbconn.
func (dbconn *Dbconn) GetDB() *pgxpool.Pool {
	return dbconn.Db
}

// getEnv fetches an environment variable or returns a default value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
