package repository

import (
	"github.com/jmoiron/sqlx"
)

// Use Sqlx
type repositoryDB struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) repositoryDB {
	return repositoryDB{db: db}
}
