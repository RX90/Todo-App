package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RX90/Todo-App/server/internal/service"
	mock_service "github.com/RX90/Todo-App/server/internal/service/mocks"
	"github.com/RX90/Todo-App/server/internal/todo"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestList_createList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId string, list todo.List)

	testTable := []struct {
		name                 string
		inputBody            string
		inputList            todo.List
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid input",
			inputBody: `{"title":"New Title List"}`,
			inputList: todo.List{
				Title: "New Title List",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId string, list todo.List) {
				s.EXPECT().Create(userId, list).Return("2", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"list_id":"2"}`,
		},
		{
			name:                 "No title",
			inputBody:            `{"title":""}`,
			mockBehavior:         func(s *mock_service.MockTodoList, userId string, list todo.List) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"list title is empty"}`,
		},
		{
			name:                 "Long title",
			inputBody:            `{"title":"Ultimate Productivity Checklist 24/7"}`,
			mockBehavior:         func(s *mock_service.MockTodoList, userId string, list todo.List) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"list title exceeds 32 characters"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"title":"Жизнь/Работа"}`,
			inputList: todo.List{
				Title: "Жизнь/Работа",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId string, list todo.List) {
				s.EXPECT().Create(userId, list).Return("", errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't create list: service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var userId = "1"

			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, userId, testCase.inputList)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/create-list", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.createList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create-list", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestList_getAllLists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid input",
			mockBehavior: func(s *mock_service.MockTodoList, userId string) {
				s.EXPECT().GetAll(userId).Return([]todo.List{{Id: "1", Title: "Homework"}, {Id: "5", Title: "Health"}}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":"1","title":"Homework"},{"id":"5","title":"Health"}]`,
		},
		{
			name: "Service error",
			mockBehavior: func(s *mock_service.MockTodoList, userId string) {
				s.EXPECT().GetAll(userId).Return(nil, errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't get all lists: service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var userId = "1"

			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, userId)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/get-all-lists", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.getAllLists)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get-all-lists", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestList_updateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId string, list todo.List)

	testTable := []struct {
		name                 string
		inputBody            string
		inputList            todo.List
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid input",
			inputBody: `{"title":"New Unique Title"}`,
			inputList: todo.List{
				Title: "New Unique Title",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId string, list todo.List) {
				s.EXPECT().Update(userId, listId, list).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "No title",
			inputBody:            `{"title":""}`,
			mockBehavior:         func(s *mock_service.MockTodoList, userId, listId string, list todo.List) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"list title is empty"}`,
		},
		{
			name:                 "Long title",
			inputBody:            `{"title":"Ultimate Productivity Checklist 24/7"}`,
			mockBehavior:         func(s *mock_service.MockTodoList, userId, listId string, list todo.List) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"list title exceeds 32 characters"}`,
		},
		{
			name:      "Service error",
			inputBody: `{"title":"Not new title"}`,
			inputList: todo.List{
				Title: "Not new title",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId string, list todo.List) {
				s.EXPECT().Update(userId, listId, list).Return(errors.New("not unique title"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't update list: not unique title"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				userId = "1"
				listId = "2"
			)

			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, userId, listId, testCase.inputList)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/update-list/:listId", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.updateList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/update-list/%s", listId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestList_deleteList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid input",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId string) {
				s.EXPECT().Delete(userId, listId).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name: "Service error",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId string) {
				s.EXPECT().Delete(userId, listId).Return(errors.New("task does not exist"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't delete list: task does not exist"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				userId = "1"
				listId = "2"
			)

			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(todoList, userId, listId)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/delete-list/:listId", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.deleteList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/delete-list/%s", listId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
