package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// ExecTx executes a function within a database transaction.
func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// create queries with transaction
	q := New(tx)
	// execute function
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			// Return both errors if both roll back and function return an error
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
