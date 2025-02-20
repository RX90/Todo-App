package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuth_isUsernameTaken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newAuthDB(sqlxDB)

	testTable := []struct {
		name       string
		username   string
		mockFunc   func(mock sqlmock.Sqlmock, username string)
		wantExists bool
		wantErr    bool
	}{
		{
			name:     "Username exists",
			username: "Test_Username",
			mockFunc: func(mock sqlmock.Sqlmock, username string) {
				query := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:     "Username does not exist",
			username: "New_Username",
			mockFunc: func(mock sqlmock.Sqlmock, username string) {
				query := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name:     "DB error",
			username: "New_Username",
			mockFunc: func(mock sqlmock.Sqlmock, username string) {
				query := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`)
				mock.ExpectQuery(query).WithArgs(username).WillReturnError(errors.New("db error"))
			},
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.username)

			exists, err := repo.isUsernameTaken(testCase.username)

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

func TestAuth_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newAuthDB(sqlxDB)

	testTable := []struct {
		name     string
		user     todo.User
		mockFunc func(mock sqlmock.Sqlmock, username, password string)
		wantErr  bool
	}{
		{
			name: "Successful user creation",
			user: todo.User{
				Username: "New_Username",
				Password: "hashed_passw0rd",
			},
			mockFunc: func(mock sqlmock.Sqlmock, username, password string) {
				query := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)

				query = regexp.QuoteMeta(`INSERT INTO users (username, password_hash) values ($1, $2)`)
				mock.ExpectExec(query).WithArgs(username, password).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Username is already taken",
			user: todo.User{
				Username: "New_Username",
				Password: "hashed_passw0rd",
			},
			mockFunc: func(mock sqlmock.Sqlmock, username, password string) {
				queryExists := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(queryExists).WithArgs(username).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "DB error on username check",
			user: todo.User{
				Username: "New_Username",
				Password: "hashed_passw0rd",
			},
			mockFunc: func(mock sqlmock.Sqlmock, username, password string) {
				queryExists := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`)
				mock.ExpectQuery(queryExists).WithArgs(username).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "DB error on user insertion",
			user: todo.User{
				Username: "New_Username",
				Password: "hashed_passw0rd",
			},
			mockFunc: func(mock sqlmock.Sqlmock, username, password string) {
				queryExists := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))`)
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(queryExists).WithArgs(username).WillReturnRows(rows)

				queryInsert := regexp.QuoteMeta(`INSERT INTO users (username, password_hash) values ($1, $2)`)
				mock.ExpectExec(queryInsert).WithArgs(username, password).WillReturnError(errors.New("insert error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.user.Username, testCase.user.Password)

			err := repo.CreateUser(testCase.user)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuth_GetUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newAuthDB(sqlxDB)

	testTable := []struct {
		name       string
		user       todo.User
		mockFunc   func(mock sqlmock.Sqlmock, username, password string)
		expectedID string
		wantErr    bool
	}{
		{
			name: "Successful user retrieval",
			user: todo.User{
				Username: "New_Username",
				Password: "hashed_passw0rd",
			},
			mockFunc: func(mock sqlmock.Sqlmock, username, password string) {
				query := regexp.QuoteMeta("SELECT id FROM users WHERE username = $1 AND password_hash = $2")
				rows := sqlmock.NewRows([]string{"id"}).AddRow("21")
				mock.ExpectQuery(query).WithArgs(username, password).WillReturnRows(rows)
			},
			expectedID: "21",
			wantErr:    false,
		},
		{
			name: "User not found",
			user: todo.User{
				Username: "Non-existing_Username",
				Password: "hashed_passw0rd",
			},
			mockFunc: func(mock sqlmock.Sqlmock, username, password string) {
				query := regexp.QuoteMeta("SELECT id FROM users WHERE username = $1 AND password_hash = $2")
				mock.ExpectQuery(query).WithArgs(username, password).WillReturnError(sql.ErrNoRows)
			},
			expectedID: "",
			wantErr:    true,
		},
		{
			name: "DB error",
			user: todo.User{
				Username: "New_Username",
				Password: "hashed_passw0rd",
			},
			mockFunc: func(mock sqlmock.Sqlmock, username, password string) {
				query := regexp.QuoteMeta("SELECT id FROM users WHERE username = $1 AND password_hash = $2")
				mock.ExpectQuery(query).WithArgs(username, password).WillReturnError(errors.New("db error"))
			},
			expectedID: "",
			wantErr:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.user.Username, testCase.user.Password)

			id, err := repo.GetUserId(testCase.user)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedID, id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuth_NewRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newAuthDB(sqlxDB)

	testTable := []struct {
		name      string
		token     string
		userId    string
		expiresAt time.Time
		mockFunc  func(mock sqlmock.Sqlmock, token, userId string, expiresAt time.Time)
		wantErr   bool
	}{
		{
			name:      "Successful insertion",
			token:     "new_refresh_token",
			userId:    "1",
			expiresAt: time.Now().Add(30 * 24 * time.Hour),
			mockFunc: func(mock sqlmock.Sqlmock, token, userId string, expiresAt time.Time) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT ut.token_id
					FROM users_tokens ut
					INNER JOIN users t ON ut.token_id = t.id
					WHERE ut.user_id = $1`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(sql.ErrNoRows)

				query = regexp.QuoteMeta("INSERT INTO tokens (refresh_token, expires_at) values ($1, $2) RETURNING id")
				rows := sqlmock.NewRows([]string{"id"}).AddRow("12")
				mock.ExpectQuery(query).WithArgs(token, expiresAt).WillReturnRows(rows)

				query = regexp.QuoteMeta("INSERT INTO users_tokens (user_id, token_id) values ($1, $2)")
				mock.ExpectExec(query).WithArgs(userId, "12").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:      "Successful update",
			token:     "new_refresh_token",
			userId:    "1",
			expiresAt: time.Now().Add(30 * 24 * time.Hour),
			mockFunc: func(mock sqlmock.Sqlmock, token, userId string, expiresAt time.Time) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT ut.token_id
					FROM users_tokens ut
					INNER JOIN users t ON ut.token_id = t.id
					WHERE ut.user_id = $1`,
				)
				rows := sqlmock.NewRows([]string{"token_id"}).AddRow("12")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta("UPDATE tokens SET refresh_token = $1, expires_at = $2 WHERE id = $3")
				mock.ExpectExec(query).WithArgs(token, expiresAt, "12").WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:      "DB error on token check",
			token:     "new_refresh_token",
			userId:    "1",
			expiresAt: time.Now().Add(30 * 24 * time.Hour),
			mockFunc: func(mock sqlmock.Sqlmock, token, userId string, expiresAt time.Time) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT ut.token_id
					FROM users_tokens ut
					INNER JOIN users t ON ut.token_id = t.id
					WHERE ut.user_id = $1`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:      "DB error on token insertion (tokens table)",
			token:     "new_refresh_token",
			userId:    "1",
			expiresAt: time.Now().Add(30 * 24 * time.Hour),
			mockFunc: func(mock sqlmock.Sqlmock, token, userId string, expiresAt time.Time) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT ut.token_id
					FROM users_tokens ut
					INNER JOIN users t ON ut.token_id = t.id
					WHERE ut.user_id = $1`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(sql.ErrNoRows)

				query = regexp.QuoteMeta("INSERT INTO tokens (refresh_token, expires_at) values ($1, $2) RETURNING id")
				mock.ExpectQuery(query).WithArgs(token, expiresAt).WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:      "DB error on token insert (users_tokens table)",
			token:     "new_refresh_token",
			userId:    "1",
			expiresAt: time.Now().Add(30 * 24 * time.Hour),
			mockFunc: func(mock sqlmock.Sqlmock, token, userId string, expiresAt time.Time) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT ut.token_id
					FROM users_tokens ut
					INNER JOIN users t ON ut.token_id = t.id
					WHERE ut.user_id = $1`,
				)
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(sql.ErrNoRows)

				query = regexp.QuoteMeta("INSERT INTO tokens (refresh_token, expires_at) values ($1, $2) RETURNING id")
				rows := sqlmock.NewRows([]string{"id"}).AddRow("12")
				mock.ExpectQuery(query).WithArgs(token, expiresAt).WillReturnRows(rows)

				query = regexp.QuoteMeta("INSERT INTO users_tokens (user_id, token_id) values ($1, $2)")
				mock.ExpectExec(query).WithArgs(userId, "12").WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:      "DB error on token update",
			token:     "new_refresh_token",
			userId:    "1",
			expiresAt: time.Now().Add(30 * 24 * time.Hour),
			mockFunc: func(mock sqlmock.Sqlmock, token, userId string, expiresAt time.Time) {
				mock.ExpectBegin()

				query := regexp.QuoteMeta(`
					SELECT ut.token_id
					FROM users_tokens ut
					INNER JOIN users t ON ut.token_id = t.id
					WHERE ut.user_id = $1`,
				)
				rows := sqlmock.NewRows([]string{"token_id"}).AddRow("12")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta("UPDATE tokens SET refresh_token = $1, expires_at = $2 WHERE id = $3")
				mock.ExpectExec(query).WithArgs(token, expiresAt, "12").WillReturnError(errors.New("update error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.token, testCase.userId, testCase.expiresAt)

			err := repo.NewRefreshToken(testCase.token, testCase.userId, testCase.expiresAt)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuth_CheckRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newAuthDB(sqlxDB)

	testTable := []struct {
		name         string
		userId       string
		refreshToken string
		mockFunc     func(mock sqlmock.Sqlmock, userId, refreshToken string)
		wantErr      bool
	}{
		{
			name:         "Valid token",
			userId:       "1",
			refreshToken: "refresh_token",
			mockFunc: func(mock sqlmock.Sqlmock, userId, refreshToken string) {
				query := regexp.QuoteMeta("SELECT ut.token_id FROM users_tokens ut WHERE ut.user_id = $1")
				rows := sqlmock.NewRows([]string{"token_id"}).AddRow("123")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta("SELECT t.refresh_token, t.expires_at FROM tokens t WHERE t.id = $1")
				rows = sqlmock.NewRows([]string{"refresh_token", "expires_at"}).AddRow(refreshToken, time.Now().Add(30*24*time.Hour))
				mock.ExpectQuery(query).WithArgs("123").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:         "Token id is missing",
			userId:       "1",
			refreshToken: "refresh_token",
			mockFunc: func(mock sqlmock.Sqlmock, userId, refreshToken string) {
				query := regexp.QuoteMeta("SELECT ut.token_id FROM users_tokens ut WHERE ut.user_id = $1")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnError(errors.New("no token_id in db"))
			},
			wantErr: true,
		},
		{
			name:         "Token is missing",
			userId:       "1",
			refreshToken: "refresh_token",
			mockFunc: func(mock sqlmock.Sqlmock, userId, refreshToken string) {
				query := regexp.QuoteMeta("SELECT ut.token_id FROM users_tokens ut WHERE ut.user_id = $1")
				rows := sqlmock.NewRows([]string{"token_id"}).AddRow("123")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta("SELECT t.refresh_token, t.expires_at FROM tokens t WHERE t.id = $1")
				mock.ExpectQuery(query).WithArgs("123").WillReturnError(errors.New("no token in db"))
			},
			wantErr: true,
		},
		{
			name:         "Tokens are different",
			userId:       "1",
			refreshToken: "refresh_token",
			mockFunc: func(mock sqlmock.Sqlmock, userId, refreshToken string) {
				query := regexp.QuoteMeta("SELECT ut.token_id FROM users_tokens ut WHERE ut.user_id = $1")
				rows := sqlmock.NewRows([]string{"token_id"}).AddRow("123")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta("SELECT t.refresh_token, t.expires_at FROM tokens t WHERE t.id = $1")
				rows = sqlmock.NewRows([]string{"refresh_token", "expires_at"}).AddRow("another_refresh_token", time.Now().Add(30*24*time.Hour))
				mock.ExpectQuery(query).WithArgs("123").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:         "Token has expired",
			userId:       "1",
			refreshToken: "refresh_token",
			mockFunc: func(mock sqlmock.Sqlmock, userId, refreshToken string) {
				query := regexp.QuoteMeta("SELECT ut.token_id FROM users_tokens ut WHERE ut.user_id = $1")
				rows := sqlmock.NewRows([]string{"token_id"}).AddRow("123")
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				query = regexp.QuoteMeta("SELECT t.refresh_token, t.expires_at FROM tokens t WHERE t.id = $1")
				rows = sqlmock.NewRows([]string{"refresh_token", "expires_at"}).AddRow(refreshToken, time.Now().Add(-24*time.Hour))
				mock.ExpectQuery(query).WithArgs("123").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId, testCase.refreshToken)

			err := repo.CheckRefreshToken(testCase.userId, testCase.refreshToken)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuth_DeleteRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := newAuthDB(sqlxDB)

	testTable := []struct {
		name         string
		userId       string
		mockFunc     func(mock sqlmock.Sqlmock, userId string)
		wantErr      bool
	}{
		{
			name:         "Successful deletion",
			userId:       "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta("DELETE FROM tokens t USING users_tokens ut WHERE t.id = ut.token_id AND ut.user_id = $1")
				mock.ExpectExec(query).WithArgs(userId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:         "DB error on deletion",
			userId:       "1",
			mockFunc: func(mock sqlmock.Sqlmock, userId string) {
				query := regexp.QuoteMeta("DELETE FROM tokens t USING users_tokens ut WHERE t.id = ut.token_id AND ut.user_id = $1")
				mock.ExpectExec(query).WithArgs(userId).WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockFunc(mock, testCase.userId)

			err := repo.DeleteRefreshToken(testCase.userId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
