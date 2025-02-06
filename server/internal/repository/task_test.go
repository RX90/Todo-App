package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTask_isTitleExistsInTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newTaskDB(sqlxDB)

	testTable := []struct {
		name       string
		listId     string
		title      string
		mockFunc   func(mock sqlmock.Sqlmock, listId, title string)
		wantExists bool
		wantErr    bool
	}{
		{
			name:   "Title exists",
			listId: "2",
			title:  "Not new title",
			mockFunc: func(mock sqlmock.Sqlmock, listId, title string) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WithArgs(listId, title).WillReturnRows(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:   "Title does not exist",
			listId: "2",
			title:  "New title",
			mockFunc: func(mock sqlmock.Sqlmock, listId, title string) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(listId, title).WillReturnRows(rows)
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name:   "DB error",
			listId: "2",
			title:  "Not new title",
			mockFunc: func(mock sqlmock.Sqlmock, listId, title string) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				mock.ExpectQuery(query).WithArgs(listId, title).WillReturnError(errors.New("db error"))
			},
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.listId, testCase.title)

			exists, err := repo.isTitleExistsInTasks(testCase.listId, testCase.title)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.wantExists, exists)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTask_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newTaskDB(sqlxDB)

	testTable := []struct {
		name     string
		listId   string
		task     todo.Task
		mockFunc func(mock sqlmock.Sqlmock, listId string, task todo.Task)
		wantErr  bool
	}{
		{
			name:   "Successful task creation",
			listId: "2",
			task: todo.Task{
				Title: "Need to buy milk (Нужно купить молоко)🥛",
			},
			mockFunc: func(mock sqlmock.Sqlmock, listId string, task todo.Task) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(listId, task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO tasks (title) VALUES ($1) RETURNING id`)
				rows = sqlmock.NewRows([]string{"id"}).AddRow("6")
				mock.ExpectQuery(query).WithArgs(task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO lists_tasks (list_id, task_id) VALUES ($1, $2)`)
				mock.ExpectExec(query).WithArgs(listId, "6").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:   "Task is already exists",
			listId: "2",
			task: todo.Task{
				Title: "Need to buy milk (Нужно купить молоко)🥛",
			},
			mockFunc: func(mock sqlmock.Sqlmock, listId string, task todo.Task) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WithArgs(listId, task.Title).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "DB error on title check",
			listId: "2",
			task: todo.Task{
				Title: "Need to buy milk (Нужно купить молоко)🥛",
			},
			mockFunc: func(mock sqlmock.Sqlmock, listId string, task todo.Task) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				mock.ExpectQuery(query).WithArgs(listId, task.Title).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "DB error on task insertion",
			listId: "2",
			task: todo.Task{
				Title: "Need to buy milk (Нужно купить молоко)🥛",
			},
			mockFunc: func(mock sqlmock.Sqlmock, listId string, task todo.Task) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(listId, task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO tasks (title) VALUES ($1) RETURNING id`)
				mock.ExpectQuery(query).WithArgs(task.Title).WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:   "DB error on lists_tasks insertion",
			listId: "2",
			task: todo.Task{
				Title: "Need to buy milk (Нужно купить молоко)🥛",
			},
			mockFunc: func(mock sqlmock.Sqlmock, listId string, task todo.Task) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(listId, task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO tasks (title) VALUES ($1) RETURNING id`)
				rows = sqlmock.NewRows([]string{"id"}).AddRow("6")
				mock.ExpectQuery(query).WithArgs(task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO lists_tasks (list_id, task_id) VALUES ($1, $2)`)
				mock.ExpectExec(query).WithArgs(listId, "6").WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.listId, testCase.task)

			taskId, err := repo.Create(testCase.listId, testCase.task)

			if testCase.wantErr {
				assert.Error(t, err)
				assert.Empty(t, taskId)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, taskId)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskDB_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newTaskDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		listId   string
		mockFunc func(mock sqlmock.Sqlmock, userId, listId string)
		want     []todo.Task
		wantErr  bool
	}{
		{
			name:   "Successful tasks fetch",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string) {
				query := regexp.QuoteMeta(`
					SELECT t.id, t.title, t.done
					FROM tasks t
					INNER JOIN lists_tasks lt ON lt.task_id = t.id
					INNER JOIN users_lists ul ON ul.list_id = lt.list_id
					WHERE lt.list_id = $1 AND ul.user_id = $2
					ORDER BY t.id`,
				)
				rows := sqlmock.NewRows([]string{"id", "title", "done"}).
					AddRow("1", "Need to finish the Todo App", false).
					AddRow("2", "Need to find a new project", false)
				mock.ExpectQuery(query).WithArgs(listId, userId).WillReturnRows(rows)
			},
			want: []todo.Task{
				{
					Id:    "1",
					Title: "Need to finish the Todo App",
					Done:  false,
				},
				{
					Id:    "2",
					Title: "Need to find a new project",
					Done:  false,
				},
			},
			wantErr: false,
		},
		{
			name:   "Successful nil fetch",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string) {
				query := regexp.QuoteMeta(`
					SELECT t.id, t.title, t.done
					FROM tasks t
					INNER JOIN lists_tasks lt ON lt.task_id = t.id
					INNER JOIN users_lists ul ON ul.list_id = lt.list_id
					WHERE lt.list_id = $1 AND ul.user_id = $2
					ORDER BY t.id`,
				)
				mock.ExpectQuery(query).WithArgs(listId, userId).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "done"}))
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:   "DB error",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string) {
				query := regexp.QuoteMeta(`
					SELECT t.id, t.title, t.done
					FROM tasks t
					INNER JOIN lists_tasks lt ON lt.task_id = t.id
					INNER JOIN users_lists ul ON ul.list_id = lt.list_id
					WHERE lt.list_id = $1 AND ul.user_id = $2
					ORDER BY t.id`,
				)
				mock.ExpectQuery(query).WithArgs(listId, userId).WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.listId)

			tasks, err := repo.GetAll(testCase.userId, testCase.listId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, tasks)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func toPointer[T any](v T) *T {
	return &v
}

func TestTask_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newTaskDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		listId   string
		taskId   string
		task     todo.UpdateTaskInput
		mockFunc func(mock sqlmock.Sqlmock, userId, listId, taskId string, task todo.UpdateTaskInput)
		wantErr  bool
	}{
		{
			name:   "Successful task update 1",
			userId: "1",
			listId: "2",
			taskId: "3",
			task: todo.UpdateTaskInput{
				Title: toPointer("New Task Title"),
				Done:  toPointer(true),
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string, task todo.UpdateTaskInput) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(listId, *task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					UPDATE tasks t
					SET title = $1, done = $2
					FROM lists_tasks lt
					INNER JOIN users_lists ul ON lt.list_id = ul.list_id
					WHERE t.id = lt.task_id AND ul.user_id = $3 AND lt.list_id = $4 AND t.id = $5`,
				)
				mock.ExpectExec(query).WithArgs(*task.Title, *task.Done, userId, listId, taskId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:   "Successful task update 2",
			userId: "1",
			listId: "2",
			taskId: "3",
			task: todo.UpdateTaskInput{
				Title: toPointer("New Task Title"),
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string, task todo.UpdateTaskInput) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(listId, *task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					UPDATE tasks t
					SET title = $1
					FROM lists_tasks lt
					INNER JOIN users_lists ul ON lt.list_id = ul.list_id
					WHERE t.id = lt.task_id AND ul.user_id = $2 AND lt.list_id = $3 AND t.id = $4`,
				)
				mock.ExpectExec(query).WithArgs(*task.Title, userId, listId, taskId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:   "Successful task update 3",
			userId: "1",
			listId: "2",
			taskId: "3",
			task: todo.UpdateTaskInput{
				Done: toPointer(true),
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string, task todo.UpdateTaskInput) {
				query := regexp.QuoteMeta(`
					UPDATE tasks t
					SET done = $1
					FROM lists_tasks lt
					INNER JOIN users_lists ul ON lt.list_id = ul.list_id
					WHERE t.id = lt.task_id AND ul.user_id = $2 AND lt.list_id = $3 AND t.id = $4`,
				)
				mock.ExpectExec(query).WithArgs(*task.Done, userId, listId, taskId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:   "Task is already exists",
			userId: "1",
			listId: "2",
			taskId: "3",
			task: todo.UpdateTaskInput{
				Title: toPointer("Not Unique Task Title"),
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string, task todo.UpdateTaskInput) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WithArgs(listId, *task.Title).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "DB error on title check",
			userId: "1",
			listId: "2",
			taskId: "3",
			task: todo.UpdateTaskInput{
				Title: toPointer("Not Unique Task Title"),
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string, task todo.UpdateTaskInput) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				mock.ExpectQuery(query).WithArgs(listId, *task.Title).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "DB error on task update",
			userId: "1",
			listId: "2",
			taskId: "3",
			task: todo.UpdateTaskInput{
				Title: toPointer("New Task Title"),
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string, task todo.UpdateTaskInput) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
					SELECT 1
					FROM tasks t
					INNER JOIN lists_tasks lt ON t.id = lt.task_id
					WHERE lt.list_id = $1 AND LOWER(t.title) = LOWER($2))`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(listId, *task.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					UPDATE tasks t
					SET title = $1
					FROM lists_tasks lt
					INNER JOIN users_lists ul ON lt.list_id = ul.list_id
					WHERE t.id = lt.task_id AND ul.user_id = $2 AND lt.list_id = $3 AND t.id = $4`,
				)
				mock.ExpectExec(query).WithArgs(*task.Title, userId, listId, taskId).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.listId, testCase.taskId, testCase.task)

			err := repo.Update(testCase.userId, testCase.listId, testCase.taskId, testCase.task)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTask_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newTaskDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		listId   string
		taskId   string
		mockFunc func(mock sqlmock.Sqlmock, userId, listId, taskId string)
		wantErr  bool
	}{
		{
			name:   "Successful task deletion",
			userId: "1",
			listId: "2",
			taskId: "3",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string) {
				query := regexp.QuoteMeta(`
					DELETE FROM tasks t
					USING lists_tasks lt
					INNER JOIN users_lists ul ON lt.list_id = ul.list_id
					WHERE t.id = lt.task_id AND ul.user_id = $1 AND lt.list_id = $2 AND t.id = $3`,
				)
				mock.ExpectExec(query).WithArgs(userId, listId, taskId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:   "Task does not exist / Task does not belong to user",
			userId: "1",
			listId: "2",
			taskId: "3",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string) {
				query := regexp.QuoteMeta(`
					DELETE FROM tasks t
					USING lists_tasks lt
					INNER JOIN users_lists ul ON lt.list_id = ul.list_id
					WHERE t.id = lt.task_id AND ul.user_id = $1 AND lt.list_id = $2 AND t.id = $3`,
				)
				mock.ExpectExec(query).WithArgs(userId, listId, taskId).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
		{
			name:   "DB error",
			userId: "1",
			listId: "2",
			taskId: "3",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId, taskId string) {
				query := regexp.QuoteMeta(`
					DELETE FROM tasks t
					USING lists_tasks lt
					INNER JOIN users_lists ul ON lt.list_id = ul.list_id
					WHERE t.id = lt.task_id AND ul.user_id = $1 AND lt.list_id = $2 AND t.id = $3`,
				)
				mock.ExpectExec(query).WithArgs(userId, listId, taskId).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.listId, testCase.taskId)

			err := repo.Delete(testCase.userId, testCase.listId, testCase.taskId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
