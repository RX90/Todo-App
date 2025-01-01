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

func (r *ListDB) isTitleExists(userId, title string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s l INNER JOIN %s ul ON l.id = ul.list_id WHERE ul.user_id = $1 AND LOWER(l.title) = LOWER($2))", listsTable, usersListsTable)
	err := r.db.QueryRow(query, userId, title).Scan(&exists)
	return exists, err
}

func (r *ListDB) Create(userId string, list todo.List) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	exists, err := r.isTitleExists(userId, list.Title)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("list '%s' is already exists", list.Title)
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

func (r *ListDB) GetAll(userId string) ([]todo.List, error) {
	var lists []todo.List

	query := fmt.Sprintf("SELECT l.id, l.title FROM %s l INNER JOIN %s ul on l.id = ul.list_id WHERE ul.user_id = $1", listsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}
