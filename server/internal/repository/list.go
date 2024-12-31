package repository

import (
	"fmt"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/jmoiron/sqlx"
)

type ListDB struct {
	db *sqlx.DB
}

func newListDB(db *sqlx.DB) *ListDB {
	return &ListDB{db: db}
}

func (r *ListDB) Create(userId string, list todo.List) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var listId string

	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", listsTable)
	if err := tx.QueryRow(query, list.Title).Scan(&listId); err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(query, userId, listId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
