package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrTaskNotFound = errors.New("task not found")

type DB struct {
	pool *pgxpool.Pool
}

type Item struct {
	ID     int
	Task   string
	Status string
}

func New(dbUrl string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &DB{pool: pool}, nil
}

func (db *DB) InsertItem(ctx context.Context, item Item) error {
	query := `INSERT INTO todo_items(task, status) VALUES ($1, $2)`
	_, err := db.pool.Exec(ctx, query, item.Task, item.Status)
	return err
}

func (db *DB) UpdateItemStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE todo_items SET status = $1 WHERE id = $2`
	result, err := db.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found with id: %d", id)
	}
	return nil
}

func (db *DB) GetAllItems(ctx context.Context) ([]Item, error) {
	query := `SELECT id, task, status from todo_items`
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Task, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (db *DB) DeleteItem(ctx context.Context, id int) error {
	query := `DELETE FROM todo_items WHERE id = $1`
	result, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found with id: %d", id)
	}
	return nil
}

func (db *DB) Close() {
	db.pool.Close()
}
