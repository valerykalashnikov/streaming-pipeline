package db

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetDBConnection() (*sql.DB, error) {
	buffer := bytes.NewBufferString("sslmode=disable")

	user := os.Getenv("POSTGRES_ENV_USERNAME")
	if user == "" {
		user = os.Getenv("DB_USER")
	}
	if user == "" {
		return nil, errors.New("DB user is expected to be passed as DB_USER environment variable")
	}
	buffer.WriteString(fmt.Sprintf(" user=%s", user))

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, errors.New("DB name is expected to be passed as DB_NAME environment variable")
	}
	buffer.WriteString(fmt.Sprintf(" dbname=%s", dbName))

	password := os.Getenv("POSTGRES_ENV_PASSWORD")
	if password == "" {
		password = os.Getenv("DB_PASSWORD")
	}
	if password != "" {
		buffer.WriteString(fmt.Sprintf(" password='%s'", password))
	}

	host := os.Getenv("DB_HOST")

	if host != "" {
		buffer.WriteString(fmt.Sprintf(" host=%s", host))
	}

	port := os.Getenv("DB_PORT")
	if port != "" {
		buffer.WriteString(fmt.Sprintf(" port=%s", port))
	}

	return sql.Open("postgres", buffer.String())
}
