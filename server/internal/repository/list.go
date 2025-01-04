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

func (r *ListDB) isTitleExistsInLists(userId, title string) (bool, error) {
	var exists bool

	query := fmt.Sprintf(`
		SELECT EXISTS(
		SELECT 1
		FROM %s l
		INNER JOIN %s ul ON l.id = ul.list_id
		WHERE ul.user_id = $1 AND LOWER(l.title) = LOWER($2))`,
		listsTable, usersListsTable,
	)
	err := r.db.QueryRow(query, userId, title).Scan(&exists)

	return exists, err
}

func (r *ListDB) Create(userId string, list todo.List) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	exists, err := r.isTitleExistsInLists(userId, list.Title)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("list '%s' is already exists", list.Title)
	}

	var listId string

	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", listsTable)
	if err := tx.QueryRow(query, list.Title).Scan(&listId); err != nil {
		tx.Rollback()
		return "", err
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(query, userId, listId)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	return listId, tx.Commit()
}

func (r *ListDB) GetAll(userId string) ([]todo.List, error) {
	var lists []todo.List

	query := fmt.Sprintf(`
		SELECT l.id, l.title
		FROM %s l
		INNER JOIN %s ul ON l.id = ul.list_id
		WHERE ul.user_id = $1
		ORDER BY l.id`,
		listsTable, usersListsTable,
	)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *ListDB) GetById(userId, listId string) (todo.List, error) {
	var list todo.List

	query := fmt.Sprintf(`
		SELECT l.id, l.title
		FROM %s l
		INNER JOIN %s ul on l.id = ul.list_id
		WHERE ul.user_id = $1 AND ul.list_id = $2`,
		listsTable, usersListsTable,
	)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *ListDB) Update(userId, listId string, list todo.List) error {
	exists, err := r.isTitleExistsInLists(userId, list.Title)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("list '%s' is already exists", list.Title)
	}

	query := fmt.Sprintf(`
		UPDATE %s l
		SET title = $1
		FROM %s ul
		WHERE l.id = ul.list_id AND ul.list_id = $2 AND ul.user_id = $3`,
		listsTable, usersListsTable,
	)
	_, err = r.db.Exec(query, list.Title, listId, userId)

	return err
}

func (r *ListDB) Delete(userId, listId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Deleting all tasks from list
	query := fmt.Sprintf(`
		DELETE FROM %s t
		USING %s lt
		INNER JOIN %s ul ON lt.list_id = ul.list_id
		WHERE t.id = lt.task_id AND ul.user_id = $1 AND lt.list_id = $2`,
		tasksTable, listsTasksTable, usersListsTable,
	)

	_, err = r.db.Exec(query, userId, listId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Deleting list
	query = fmt.Sprintf(`
		DELETE FROM %s l
		USING %s ul
		WHERE l.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`,
		listsTable, usersListsTable,
	)

	_, err = r.db.Exec(query, userId, listId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
