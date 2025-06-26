package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestList_isTitleExistsInLists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newListDB(sqlxDB)

	testTable := []struct {
		name       string
		userId     string
		title      string
		mockFunc   func(mock sqlmock.Sqlmock, userId, title string)
		wantExists bool
		wantErr    bool
	}{
		{
			name:   "Title exists",
			userId: "1",
			title:  "Not new title",
			mockFunc: func(mock sqlmock.Sqlmock, userId, title string) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WithArgs(userId, title).WillReturnRows(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:   "Title does not exist",
			userId: "1",
			title:  "New title",
			mockFunc: func(mock sqlmock.Sqlmock, userId, title string) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, title).WillReturnRows(rows)
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name:   "DB error",
			userId: "1",
			title:  "Not new title",
			mockFunc: func(mock sqlmock.Sqlmock, userId, title string) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				mock.ExpectQuery(query).WithArgs(userId, title).WillReturnError(errors.New("db error"))
			},
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.title)

			exists, err := repo.isTitleExistsInLists(testCase.userId, testCase.title)

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

func TestList_countUserLists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newListDB(sqlxDB)

	testTable := []struct {
		name          string
		userId        string
		mockFunc      func(mock sqlmock.Sqlmock, userId string)
		expectedCount int
		wantErr       bool
	}{
		{
			name:   "Didn't reach the limit",
			userId: "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				rows := sqlmock.NewRows([]string{"lists_count"}).AddRow(14)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)
			},
			expectedCount: 14,
			wantErr:       false,
		},
		{
			name:   "Reached the limit",
			userId: "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				rows := sqlmock.NewRows([]string{"lists_count"}).AddRow(20)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)
			},
			expectedCount: 20,
			wantErr:       false,
		},
		{
			name:   "DB error",
			userId: "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(errors.New("db error"))
			},
			expectedCount: 0,
			wantErr:       true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId)

			count, err := repo.countUserLists(testCase.userId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.expectedCount, count)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestList_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newListDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		list     todo.List
		mockFunc func(mock sqlmock.Sqlmock, userId string, list todo.List)
		wantErr  bool
	}{
		{
			name:   "Successful list creation",
			userId: "1",
			list: todo.List{
				Title: "Shopping list",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId string, list todo.List) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				rows = sqlmock.NewRows([]string{"lists_count"}).AddRow(5)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO lists (title) VALUES (?) RETURNING id`)
				rows = sqlmock.NewRows([]string{"id"}).AddRow("5")
				mock.ExpectQuery(query).WithArgs(list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO users_lists (user_id, list_id) VALUES (?, ?)`)
				mock.ExpectExec(query).WithArgs(userId, "5").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:   "List is already exists",
			userId: "1",
			list: todo.List{
				Title: "Shopping list",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId string, list todo.List) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "DB error on title check",
			userId: "1",
			list: todo.List{
				Title: "Shopping list",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId string, list todo.List) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "Reached lists limit",
			userId: "1",
			list: todo.List{
				Title: "Shopping list",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId string, list todo.List) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				rows = sqlmock.NewRows([]string{"lists_count"}).AddRow(10)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "DB error on limit check",
			userId: "1",
			list: todo.List{
				Title: "Shopping list",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId string, list todo.List) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "DB error on list insertion",
			userId: "1",
			list: todo.List{
				Title: "Shopping list",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId string, list todo.List) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				rows = sqlmock.NewRows([]string{"lists_count"}).AddRow(5)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO lists (title) VALUES (?) RETURNING id`)
				mock.ExpectQuery(query).WithArgs(list.Title).WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:   "DB error on users_lists insertion",
			userId: "1",
			list: todo.List{
				Title: "Shopping list",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId string, list todo.List) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					SELECT COUNT(*) AS lists_count
					FROM users_lists
					WHERE user_id = ?`,
				)
				rows = sqlmock.NewRows([]string{"lists_count"}).AddRow(5)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO lists (title) VALUES (?) RETURNING id`)
				rows = sqlmock.NewRows([]string{"id"}).AddRow("5")
				mock.ExpectQuery(query).WithArgs(list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO users_lists (user_id, list_id) VALUES (?, ?)`)
				mock.ExpectExec(query).WithArgs(userId, "5").WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.list)

			listId, err := repo.Create(testCase.userId, testCase.list)

			if testCase.wantErr {
				assert.Error(t, err)
				assert.Empty(t, listId)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, listId)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListDB_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newListDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		mockFunc func(mock sqlmock.Sqlmock, userId string)
		want     []todo.List
		wantErr  bool
	}{
		{
			name:   "Successful lists fetch",
			userId: "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta(`
					SELECT l.id, l.title
					FROM lists l
					INNER JOIN users_lists ul ON l.id = ul.list_id
					WHERE ul.user_id = ?
					ORDER BY l.id`,
				)
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow("1", "Shopping list").
					AddRow("2", "Work list")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)
			},
			want: []todo.List{
				{
					Id:    "1",
					Title: "Shopping list",
				},
				{
					Id:    "2",
					Title: "Work list",
				},
			},
			wantErr: false,
		},
		{
			name:   "Successful nil fetch",
			userId: "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta(`
					SELECT l.id, l.title
					FROM lists l
					INNER JOIN users_lists ul ON l.id = ul.list_id
					WHERE ul.user_id = ?
					ORDER BY l.id`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:   "DB error",
			userId: "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta(`
					SELECT l.id, l.title
					FROM lists l
					INNER JOIN users_lists ul ON l.id = ul.list_id
					WHERE ul.user_id = ?
					ORDER BY l.id`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId)

			lists, err := repo.GetAll(testCase.userId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, lists)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListDB_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newListDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		listId   string
		mockFunc func(mock sqlmock.Sqlmock, userId, listId string)
		want     todo.List
		wantErr  bool
	}{
		{
			name:   "Successful list fetch",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string) {
				query := regexp.QuoteMeta(`
					SELECT l.id, l.title
					FROM lists l
					INNER JOIN users_lists ul on l.id = ul.list_id
					WHERE ul.user_id = ? AND ul.list_id = ?`,
				)
				rows := sqlmock.NewRows([]string{"id", "title"}).AddRow("2", "Shopping list")
				mock.ExpectQuery(query).WithArgs(userId, listId).WillReturnRows(rows)
			},
			want:    todo.List{Id: "2", Title: "Shopping list"},
			wantErr: false,
		},
		{
			name:   "List not found",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string) {
				query := regexp.QuoteMeta(`
					SELECT l.id, l.title
					FROM lists l
					INNER JOIN users_lists ul on l.id = ul.list_id
					WHERE ul.user_id = ? AND ul.list_id = ?`,
				)
				mock.ExpectQuery(query).WithArgs(userId, listId).WillReturnError(sql.ErrNoRows)
			},
			want:    todo.List{},
			wantErr: true,
		},
		{
			name:   "DB error",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string) {
				query := regexp.QuoteMeta(`
					SELECT l.id, l.title
					FROM lists l
					INNER JOIN users_lists ul on l.id = ul.list_id
					WHERE ul.user_id = ? AND ul.list_id = ?`,
				)
				mock.ExpectQuery(query).WithArgs(userId, listId).WillReturnError(errors.New("db error"))
			},
			want:    todo.List{},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.listId)

			list, err := repo.GetById(testCase.userId, testCase.listId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, list)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestList_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newListDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		listId   string
		list     todo.List
		mockFunc func(mock sqlmock.Sqlmock, userId, listId string, list todo.List)
		wantErr  bool
	}{
		{
			name:   "Successful list update",
			userId: "1",
			listId: "2",
			list: todo.List{
				Title: "New unique list title",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string, list todo.List) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					UPDATE lists
					SET title = ?
					WHERE id = (
						SELECT list_id
						FROM users_lists
						WHERE list_id = ? AND user_id = ?
					)`,
				)
				mock.ExpectExec(query).WithArgs(list.Title, listId, userId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:   "List is already exists",
			userId: "1",
			listId: "2",
			list: todo.List{
				Title: "Not new title",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string, list todo.List) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "DB error on title check",
			userId: "1",
			listId: "2",
			list: todo.List{
				Title: "Not new title",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string, list todo.List) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "DB error on list update",
			userId: "1",
			listId: "2",
			list: todo.List{
				Title: "New list title",
			},
			mockFunc: func(mock sqlmock.Sqlmock, userId, listId string, list todo.List) {
				query := regexp.QuoteMeta(`
					SELECT EXISTS(
						SELECT 1
						FROM lists l
						INNER JOIN users_lists ul ON l.id = ul.list_id
						WHERE ul.user_id = ? AND LOWER(l.title) = LOWER(?)
					)`,
				)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(userId, list.Title).WillReturnRows(rows)

				query = regexp.QuoteMeta(`
					UPDATE lists
					SET title = ?
					WHERE id = (
						SELECT list_id
						FROM users_lists
						WHERE list_id = ? AND user_id = ?
					)`,
				)
				mock.ExpectExec(query).WithArgs(list.Title, listId, userId).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.listId, testCase.list)

			err := repo.Update(testCase.userId, testCase.listId, testCase.list)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestList_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newListDB(sqlxDB)

	testTable := []struct {
		name     string
		userId   string
		listId   string
		mockFunc func(mock sqlmock.Sqlmock, userId, listId string)
		wantErr  bool
	}{
		{
			name:   "Successful list deletion",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, listId, userId string) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					DELETE FROM tasks
					WHERE id IN (
						SELECT task_id
						FROM lists_tasks lt
						WHERE lt.list_id = ? AND EXISTS (
							SELECT 1
							FROM users_lists ul
							WHERE ul.user_id = ? AND ul.list_id = lt.list_id
						)
					)`,
				)
				mock.ExpectExec(query).WithArgs(listId, userId).WillReturnResult(sqlmock.NewResult(0, 3))

				query = regexp.QuoteMeta(`
					DELETE FROM lists
					WHERE id = ? AND EXISTS (
						SELECT 1
						FROM users_lists ul
						WHERE ul.user_id = ? AND ul.list_id = id
					)`,
				)
				mock.ExpectExec(query).WithArgs(listId, userId).WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:   "List does not exist / List does not belong to user",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, listId, userId string) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					DELETE FROM tasks
					WHERE id IN (
						SELECT task_id
						FROM lists_tasks lt
						WHERE lt.list_id = ? AND EXISTS (
							SELECT 1
							FROM users_lists ul
							WHERE ul.user_id = ? AND ul.list_id = lt.list_id
						)
					)`,
				)
				mock.ExpectExec(query).WithArgs(listId, userId).WillReturnResult(sqlmock.NewResult(0, 0))

				query = regexp.QuoteMeta(`
					DELETE FROM lists
					WHERE id = ? AND EXISTS (
						SELECT 1
						FROM users_lists ul
						WHERE ul.user_id = ? AND ul.list_id = id
					)`,
				)
				mock.ExpectExec(query).WithArgs(listId, userId).WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:   "DB error on tasks deletion",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, listId, userId string) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					DELETE FROM tasks
					WHERE id IN (
						SELECT task_id
						FROM lists_tasks lt
						WHERE lt.list_id = ? AND EXISTS (
							SELECT 1
							FROM users_lists ul
							WHERE ul.user_id = ? AND ul.list_id = lt.list_id
						)
					)`,
				)
				mock.ExpectExec(query).WithArgs(listId, userId).WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:   "DB error on list deletion",
			userId: "1",
			listId: "2",
			mockFunc: func(mock sqlmock.Sqlmock, listId, userId string) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					DELETE FROM tasks
					WHERE id IN (
						SELECT task_id
						FROM lists_tasks lt
						WHERE lt.list_id = ? AND EXISTS (
							SELECT 1
							FROM users_lists ul
							WHERE ul.user_id = ? AND ul.list_id = lt.list_id
						)
					)`,
				)
				mock.ExpectExec(query).WithArgs(listId, userId).WillReturnResult(sqlmock.NewResult(0, 3))

				query = regexp.QuoteMeta(`
					DELETE FROM lists
					WHERE id = ? AND EXISTS (
						SELECT 1
						FROM users_lists ul
						WHERE ul.user_id = ? AND ul.list_id = id
					)`,
				)
				mock.ExpectExec(query).WithArgs(listId, userId).WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.listId)

			err := repo.Delete(testCase.listId, testCase.userId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
