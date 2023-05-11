package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// For Docker deployment set docker to true
const docker = true
const dbPath = "/database/social-network.db?_foreign_keys=on"

var db *sql.DB

func ConnectDB() {
	var wd string
	var err error
	if docker {
		wd = "/app/back-end"
	} else {
		wd, err = os.Getwd()
		if err != nil {
			log.Fatalf("Error getting working directory: %v", err)
		}
	}

	db, err = sql.Open("sqlite3", wd+dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = MigrateDB(db, wd+"/database/migrations")
	if err != nil {
		log.Fatalf("Error applying database migrations: %v", err)
	}

	log.Println("Database connection established")
}

func MigrateDB(db *sql.DB, migrationsDir string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize sqlite3 driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply database migrations: %w", err)
	}

	return nil
}

func CloseDB() {
	db.Close()
	log.Println("Database connection closed")
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Prepare(query string) (*sql.Stmt, error) {
	return db.Prepare(query)
}
