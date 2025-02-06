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
		WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
		tasksTable, listsTasksTable,
	)
	err := r.db.QueryRow(query, listId, title).Scan(&exists)

	return exists, err
}

func (r *TaskDB) Create(listId string, task todo.Task) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	exists, err := r.isTitleExistsInTasks(listId, task.Title)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("task '%s' is already exists", task.Title)
	}

	var taskId string

	query := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", tasksTable)
	if err = tx.QueryRow(query, task.Title).Scan(&taskId); err != nil {
		tx.Rollback()
		return "", err
	}

	query = fmt.Sprintf("INSERT INTO %s (list_id, task_id) VALUES ($1, $2)", listsTasksTable)
	_, err = tx.Exec(query, listId, taskId)
	if err != nil {
		tx.Rollback()
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
		WHERE lt.list_id = $1 AND ul.user_id = $2
		ORDER BY t.id`,
		tasksTable, listsTasksTable, usersListsTable,
	)
	err := r.db.Select(&tasks, query, listId, userId)

	return tasks, err
}

func (r *TaskDB) Update(userId, listId, taskId string, task todo.UpdateTaskInput) error {
	setValues := make([]string, 0, 2)
	args := make([]any, 0, 2)
	argId := 1

	if task.Title != nil {
		exists, err := r.isTitleExistsInTasks(listId, *task.Title)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("task '%s' is already exists", *task.Title)
		}

		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *task.Title)
		argId++
	}

	if task.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *task.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
		UPDATE %s t
		SET %s
		FROM %s lt
		INNER JOIN %s ul ON lt.list_id = ul.list_id
		WHERE t.id = lt.task_id AND ul.user_id = $%d AND lt.list_id = $%d AND t.id = $%d`,
		tasksTable, setQuery, listsTasksTable, usersListsTable, argId, argId+1, argId+2,
	)
	args = append(args, userId, listId, taskId)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r TaskDB) Delete(userId, listId, taskId string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s t
		USING %s lt
		INNER JOIN %s ul ON lt.list_id = ul.list_id
		WHERE t.id = lt.task_id AND ul.user_id = $1 AND lt.list_id = $2 AND t.id = $3`,
		tasksTable, listsTasksTable, usersListsTable,
	)
	_, err := r.db.Exec(query, userId, listId, taskId)

	return err
}
