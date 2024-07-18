package repositories

// query helper with logging features for time and query sql
import (
	"context"
	"database/sql"
	"log"
	"time"
)

// I created DBQueryer and DBCommander to add support for *sql.Tx and *sql.DB for sql command and query helpers
type DBQueryer interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type DBCommander interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func singleRowQueryWrapper(ctx context.Context, q DBQueryer, query string, args ...interface{}) *sql.Row {
	startTime := time.Now()
	row := q.QueryRowContext(ctx, query, args...)
	duration := time.Since(startTime)
	log.Printf("Executing Query: %s with parameters: %v", query, args)
	log.Printf("Query executed in %s", duration)
	return row
}

func commandQueryWrapper(ctx context.Context, cm DBCommander, query string, params ...interface{}) (sql.Result, error) {
	startTime := time.Now()
	result, err := cm.ExecContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime)
	log.Printf("Executing Query: %s with parameters: %v", query, params)
	log.Printf("Query executed in %s", duration)
	return result, err
}
