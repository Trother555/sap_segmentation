package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sap_segmentation/internal/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Db interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	GetDb() *sqlx.DB
}

type DbImpl struct {
	Db *sqlx.DB
}

func New(cfg *config.Config) Db {
	constring := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)
	db, err := sqlx.Connect("pgx", constring)
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}
	return &DbImpl{
		Db: db,
	}
}

func (db *DbImpl) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return db.Db.NamedExecContext(ctx, query, arg)
}

func (db *DbImpl) GetDb() *sqlx.DB {
	return db.Db
}
