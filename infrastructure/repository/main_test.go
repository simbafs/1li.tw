package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"1litw/sqlc"

	_ "modernc.org/sqlite"
)

var (
	testQueries *sqlc.Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	// Use in-memory SQLite database for testing
	testDB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// Read and execute the schema to create tables
	schema, err := os.ReadFile("../../sql/schema.sql")
	if err != nil {
		log.Fatalf("could not read schema.sql: %v", err)
	}

	_, err = testDB.ExecContext(context.Background(), string(schema))
	if err != nil {
		log.Fatalf("could not apply schema: %v", err)
	}

	testQueries = sqlc.New(testDB)

	// Run the tests
	code := m.Run()

	// Clean up
	testDB.Close()
	os.Exit(code)
}
