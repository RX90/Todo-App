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
			WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
		)`,
		listsTable, usersListsTable,
	)
	err := r.db.QueryRow(query, userId, title).Scan(&exists)

	return exists, err
}

func (r *ListDB) countUserLists(userId string) (int, error) {
	var count int

	query := fmt.Sprintf(`
		SELECT COUNT(*) AS lists_count
		FROM %s
		WHERE user_id = ?`,
		usersListsTable,
	)
	err := r.db.QueryRow(query, userId).Scan(&count)

	return count, err
}

func (r *ListDB) Create(userId string, list todo.List) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	exists, err := r.isTitleExistsInLists(userId, list.Title)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("list '%s' is already exists", list.Title)
	}

	count, err := r.countUserLists(userId)
	if err != nil {
		return "", err
	}
	if count >= 10 {
		return "", fmt.Errorf("reached the limit of lists")
	}

	var listId string

	query := fmt.Sprintf("INSERT INTO %s (title) VALUES (?) RETURNING id", listsTable)
	if err := tx.QueryRow(query, list.Title).Scan(&listId); err != nil {
		return "", err
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES (?, ?)", usersListsTable)
	_, err = tx.Exec(query, userId, listId)
	if err != nil {
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
		WHERE ul.user_id = ?
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
		WHERE ul.user_id = ? AND ul.list_id = ?`,
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
		UPDATE %s
		SET title = ?
		WHERE id = (
			SELECT list_id
			FROM %s
			WHERE list_id = ? AND user_id = ?
		)`,
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
	defer tx.Rollback()

	// Deleting all tasks from list
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id IN (
			SELECT task_id
			FROM %s lt
			WHERE lt.list_id = ? AND EXISTS (
				SELECT 1
				FROM %s ul
				WHERE ul.user_id = ? AND ul.list_id = lt.list_id
			)
		)`,
		tasksTable, listsTasksTable, usersListsTable,
	)

	_, err = r.db.Exec(query, listId, userId)
	if err != nil {
		return err
	}

	// Deleting list
	query = fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = ? AND EXISTS (
			SELECT 1
			FROM %s ul
			WHERE ul.user_id = ? AND ul.list_id = id
		)`,
		listsTable, usersListsTable,
	)
	_, err = r.db.Exec(query, listId, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
