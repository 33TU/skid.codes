package database

import (
	"backend/config"
	"context"
	"errors"
	"log"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

// init initializes DB connection pool.
func init() {
	var err error

	connString, ok := config.Get("DB_URL")
	if !ok {
		log.Fatalln("DB_URI env not found.")
	}

	if pool, err = pgxpool.New(context.Background(), connString); err != nil {
		log.Fatalln(err)
	}
}

// Pool returns internal connection pool.
func Pool() *pgxpool.Pool {
	return pool
}

// Select using internal pool. Scans into dst type.
func Select[T any](timeout time.Duration, query string, args ...interface{}) ([]*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dst := []*T{}
	err := pgxscan.Select(ctx, pool, &dst, query, args...)

	return dst, err
}

// Select one row using internal pool. Scans into dst type.
func SelectOne[T any](timeout time.Duration, query string, args ...interface{}) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dstArr := []T{}

	if err := pgxscan.Select(ctx, pool, &dstArr, query, args...); err != nil {
		return nil, err
	}

	if len(dstArr) == 0 {
		return nil, pgx.ErrNoRows
	}

	return &dstArr[0], nil
}

// QueryRow fetches single row. If the query selects no rows, will return ErrNoRows.
func QueryRow[T any](timeout time.Duration, query string, args ...interface{}) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	arr := [1]T{}
	row := pool.QueryRow(ctx, query, args)

	err := row.Scan(&arr)
	if err == nil {
		return &arr[0], err
	}

	return nil, err
}

// GetError tries to get pgconn.PgError from err
func GetError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	return pgErr
}

// ErrorCodeEquals tries to check if error is pg.error and matches code.
func ErrorCodeEquals(err error, code string) bool {
	if pgErr := GetError(err); pgErr != nil {
		return pgErr.Code == code
	}

	return false
}
