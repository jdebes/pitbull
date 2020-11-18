package db

import (
	"context"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type pitbullDBKey string

var (
	noDBError = errors.New("context does not contain DB")

	key = pitbullDBKey("pitbulldb")
)

func NewDB() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", "user:password@tcp(127.0.0.1:3306)/pitbull")
}

func WithDB(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, key, db)
}

func DB(ctx context.Context) (*sqlx.DB, error) {
	if db, ok := ctx.Value(key).(*sqlx.DB); ok {
		return db, nil
	}

	return nil, noDBError
}
