package repository

import (
	"fmt"
	"strings"

	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/jmoiron/sqlx"
)

type TaskDB struct {
	db *sqlx.DB
}

func newTaskDB(db *sqlx.DB) *TaskDB {
	return &TaskDB{db: db}
}

func (r *TaskDB) isTitleExistsInTasks(listId, title string) (bool, error) {
	var exists bool

	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1
			FROM %s t
			INNER JOIN %s lt ON t.id = lt.task_id
			WHERE lt.list_id = ? AND LOWER(t.title) = LOWER(?)
		)`, tasksTable, listsTasksTable,
	)
	err := r.db.QueryRow(query, listId, title).Scan(&exists)

	return exists, err
}

func (r *TaskDB) countUserTasks(listId string) (int, error) {
	var count int

	query := fmt.Sprintf(`
		SELECT COUNT(*) AS tasks_count
		FROM %s
		WHERE list_id = ?`, listsTasksTable,
	)
	err := r.db.QueryRow(query, listId).Scan(&count)

	return count, err
}

func (r *TaskDB) Create(listId string, task todo.Task) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	exists, err := r.isTitleExistsInTasks(listId, task.Title)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("task '%s' is already exists", task.Title)
	}

	count, err := r.countUserTasks(listId)
	if err != nil {
		return "", err
	}
	if count >= 50 {
		return "", fmt.Errorf("reached the limit of tasks")
	}

	var taskId string

	query := fmt.Sprintf("INSERT INTO %s (title) VALUES (?) RETURNING id", tasksTable)
	if err = tx.QueryRow(query, task.Title).Scan(&taskId); err != nil {
		return "", err
	}

	query = fmt.Sprintf("INSERT INTO %s (list_id, task_id) VALUES (?, ?)", listsTasksTable)
	_, err = tx.Exec(query, listId, taskId)
	if err != nil {
		return "", err
	}

	return taskId, tx.Commit()
}

func (r *TaskDB) GetAll(userId, listId string) ([]todo.Task, error) {
	var tasks []todo.Task

	query := fmt.Sprintf(`
		SELECT t.id, t.title, t.done
		FROM %s t
		INNER JOIN %s lt ON lt.task_id = t.id
		INNER JOIN %s ul ON ul.list_id = lt.list_id
		WHERE lt.list_id = ? AND ul.user_id = ?
		ORDER BY t.id`,
		tasksTable, listsTasksTable, usersListsTable,
	)
	err := r.db.Select(&tasks, query, listId, userId)

	return tasks, err
}

func (r *TaskDB) Update(userId, listId, taskId string, task todo.UpdateTaskInput) error {
	setValues := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if task.Title != nil {
		exists, err := r.isTitleExistsInTasks(listId, *task.Title)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("task '%s' is already exists", *task.Title)
		}

		setValues = append(setValues, "title = ?")
		args = append(args, *task.Title)
	}

	if task.Done != nil {
		setValues = append(setValues, "done = ?")
		args = append(args, *task.Done)
	}

	if len(setValues) == 0 {
		return nil
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE id = ?
		AND id IN (
			SELECT t.id
			FROM %s t
			INNER JOIN %s lt ON t.id = lt.task_id
			INNER JOIN %s ul ON lt.list_id = ul.list_id
			WHERE ul.user_id = ? AND lt.list_id = ?
		)`,
		tasksTable, setQuery, tasksTable, listsTasksTable, usersListsTable,
	)

	args = append(args, taskId, userId, listId)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TaskDB) Delete(userId, listId, taskId string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = ?
		AND EXISTS (
			SELECT 1
			FROM %s lt
			INNER JOIN %s ul ON lt.list_id = ul.list_id
			WHERE lt.task_id = %s.id AND lt.list_id = ? AND ul.user_id = ?
		)`,
		tasksTable, listsTasksTable, usersListsTable, tasksTable,
	)
	_, err := r.db.Exec(query, taskId, listId, userId)
	return err
}
