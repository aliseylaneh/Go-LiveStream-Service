package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectToPostgres establishes a connection to a PostgreSQL database with the provided connection parameters.
// It returns a pointer to the database connection if successful, or an error if the connection fails.
func ConnectToPostgres(dbHost string, dbPort int, dbName string, dbUsername string, dbPassword string) (*sql.DB, error) {
	// Create the connection string using the provided parameters.
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUsername, dbPassword)

	for {
		// Open a database connection with the PostgreSQL driver.
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			// If there's an error while opening the connection, return the error.
			return nil, err
		}

		// Ping the database to verify the connection.
		err = db.Ping()
		if err == nil {
			// If the connection is successful, return the database connection.
			return db, nil
		} else {
			// If there's an error during the ping, return the error.
			return nil, err
		}
	}
}
