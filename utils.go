package main

import (
	"context"
	"database/sql"
)

func transact(db *sql.DB, ctx context.Context, f func(tx *sql.Tx, ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := f(tx, ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	err = tx.Commit()
	return err
}
