package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const connectionStringFmt = `postgres://%s:%s@%s:%s`

// CreatePostgresConnection creates a new PostgreSQL database connection
func CreatePostgresConnection(dbHost, dbPort, dbName, dbUser, dbPassword, sslmode string) (*sql.DB, error) {
	connectionString := fmt.Sprintf(connectionStringFmt, dbUser, dbPassword, dbHost, dbPort)
	if len(dbName) > 0 {
		connectionString += "/" + dbName
	}

	connectionString = fmt.Sprintf("%s?sslmode=%s", connectionString, sslmode)

	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to server=%s user=%s sslmode=%s, err:%w",
			dbHost,
			dbUser,
			sslmode,
			err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping postgres db server=%s user=%s sslmode=%s, err:%w",
			dbHost,
			dbUser,
			sslmode,
			err)
	}

	return sqlDB, err
}
