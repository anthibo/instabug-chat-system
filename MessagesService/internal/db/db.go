package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func NewMySQLConn(dsn string) (*sql.DB, error) {
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(10)
	dbConn.SetConnMaxLifetime(2 * time.Minute)
	log.Println("Connected Successfully to MySQL")
	return dbConn, nil
}
