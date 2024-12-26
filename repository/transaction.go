package repository

import (
    "github.com/jmoiron/sqlx"
)

// handleTransaction ensures proper handling of sqlx transactions.
func handleTransaction(tx *sqlx.Tx) func() {
    return func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }
}
