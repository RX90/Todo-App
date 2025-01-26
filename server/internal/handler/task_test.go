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

func TestTask_createTask(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoTask, userId, listId string, task todo.Task)

	testTable := []struct {
		name                 string
		listId               string
		inputBody            string
		inputTask            todo.Task
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid input",
			listId:    "2",
			inputBody: `{"title":"New task"}`,
			inputTask: todo.Task{
				Title: "New task",
			},
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId string, task todo.Task) {
				s.EXPECT().Create(userId, listId, task).Return("3", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"task_id":"3"}`,
		},
		{
			name:                 "Invalid listId",
			listId:               `{"listId":"2"}`,
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId string, task todo.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't get list id: strconv.Atoi: parsing \"{\\\"listId\\\":\\\"2\\\"}\": invalid syntax"}`,
		},
		{
			name:                 "No title",
			listId:               "2",
			inputBody:            `{"title":""}`,
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId string, task todo.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"empty input"}`,
		},
		{
			name:                 "Invalid title",
			listId:               "2",
			inputBody:            `{"title":"New taskkkkkkkkkkkkkkkkkkkkkkkkkkkkkk"}`,
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId string, task todo.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"input exceeds 32 characters"}`,
		},
		{
			name:      "Service error",
			listId:    "2",
			inputBody: `{"title":"Какая-то задача"}`,
			inputTask: todo.Task{
				Title: "Какая-то задача",
			},
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId string, task todo.Task) {
				s.EXPECT().Create(userId, listId, task).Return("", errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't create task: service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var userId = "1"

			c := gomock.NewController(t)
			defer c.Finish()

			todoTask := mock_service.NewMockTodoTask(c)
			testCase.mockBehavior(todoTask, userId, testCase.listId, testCase.inputTask)

			services := &service.Service{TodoTask: todoTask}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/lists/:listId/create-task", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.createTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/lists/%s/create-task", testCase.listId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestTask_getAllTasks(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoTask, userId, listId string)

	testTable := []struct {
		name                 string
		listId               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Valid input",
			listId: "2",
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId string) {
				s.EXPECT().GetAll(userId, listId).Return([]todo.Task{{Id: "1", Title: "Task №1", Done: true}, {Id: "5", Title: "Task №2", Done: false}}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":"1","title":"Task №1","done":true},{"id":"5","title":"Task №2","done":false}]`,
		},
		{
			name:                 "Invalid listId",
			listId:               `{"listId":"2"}`,
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't get list id: strconv.Atoi: parsing \"{\\\"listId\\\":\\\"2\\\"}\": invalid syntax"}`,
		},
		{
			name:   "Service error",
			listId: "2",
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId string) {
				s.EXPECT().GetAll(userId, listId).Return(nil, errors.New("service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't get all tasks: service error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var userId = "1"

			c := gomock.NewController(t)
			defer c.Finish()

			todoTask := mock_service.NewMockTodoTask(c)
			testCase.mockBehavior(todoTask, userId, testCase.listId)

			services := &service.Service{TodoTask: todoTask}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/lists/:listId/get-all-tasks", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.getAllTasks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/lists/%s/get-all-tasks", testCase.listId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func toPointer[T any](v T) *T {
	return &v
}

func TestTask_updateTask(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput)

	testTable := []struct {
		name                 string
		listId               string
		taskId               string
		inputBody            string
		inputTask            todo.UpdateTaskInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid input 1",
			listId:    "2",
			taskId:    "3",
			inputBody: `{"title":"New Task Title"}`,
			inputTask: todo.UpdateTaskInput{
				Title: toPointer("New Task Title"),
			},
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput) {
				s.EXPECT().Update(userId, listId, taskId, task).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:      "Valid input 2",
			listId:    "2",
			taskId:    "3",
			inputBody: `{"done":true}`,
			inputTask: todo.UpdateTaskInput{
				Done: toPointer(true),
			},
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput) {
				s.EXPECT().Update(userId, listId, taskId, task).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:      "Valid input 3",
			listId:    "2",
			taskId:    "3",
			inputBody: `{"title":"New Task Title", "done":true}`,
			inputTask: todo.UpdateTaskInput{
				Title: toPointer("New Task Title"),
				Done:  toPointer(true),
			},
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput) {
				s.EXPECT().Update(userId, listId, taskId, task).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:      "No values",
			listId:    "2",
			taskId:    "3",
			inputBody: `{"title": ""}`,
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"empty input"}`,
		},
		{
			name:                 "Invalid listId",
			listId:               `{"listId":"2"}`,
			taskId:               "3",
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't get list id: strconv.Atoi: parsing \"{\\\"listId\\\":\\\"2\\\"}\": invalid syntax"}`,
		},
		{
			name:                 "Invalid taskId",
			listId:               "2",
			taskId:               `{"taskId":"3"}`,
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't get task id: strconv.Atoi: parsing \"{\\\"taskId\\\":\\\"3\\\"}\": invalid syntax"}`,
		},
		{
			name:                 "No values",
			listId:               "2",
			taskId:               "3",
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId, taskId string, task todo.UpdateTaskInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't bind JSON: EOF"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var userId = "1"

			c := gomock.NewController(t)
			defer c.Finish()

			todoTask := mock_service.NewMockTodoTask(c)
			testCase.mockBehavior(todoTask, userId, testCase.listId, testCase.taskId, testCase.inputTask)

			services := &service.Service{TodoTask: todoTask}
			handler := NewHandler(services)

			r := gin.New()
			r.PUT("/lists/:listId/update-task/:taskId", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.updateTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/lists/%s/update-task/%s", testCase.listId, testCase.taskId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestTask_deleteTask(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoTask, userId, listId, taskId string)

	testTable := []struct {
		name                 string
		listId               string
		taskId               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Valid input",
			listId: "2",
			taskId: "3",
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId, taskId string) {
				s.EXPECT().Delete(userId, listId, taskId).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "Invalid listId",
			listId:               `{"listId":"2"}`,
			taskId:               "3",
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId, taskId string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't get list id: strconv.Atoi: parsing \"{\\\"listId\\\":\\\"2\\\"}\": invalid syntax"}`,
		},
		{
			name:                 "Invalid taskId",
			listId:               "2",
			taskId:               `{"taskId":"3"}`,
			mockBehavior:         func(s *mock_service.MockTodoTask, userId, listId, taskId string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"err":"can't get task id: strconv.Atoi: parsing \"{\\\"taskId\\\":\\\"3\\\"}\": invalid syntax"}`,
		},
		{
			name:   "Service error",
			listId: "2",
			taskId: "3",
			mockBehavior: func(s *mock_service.MockTodoTask, userId, listId, taskId string) {
				s.EXPECT().Delete(userId, listId, taskId).Return(errors.New("task does not exist"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"err":"can't delete task: task does not exist"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var userId = "1"

			c := gomock.NewController(t)
			defer c.Finish()

			todoTask := mock_service.NewMockTodoTask(c)
			testCase.mockBehavior(todoTask, userId, testCase.listId, testCase.taskId)

			services := &service.Service{TodoTask: todoTask}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/lists/:listId/delete-task/:taskId", func(ctx *gin.Context) {
				ctx.Set(userCtx, userId)
				ctx.Next()
			}, handler.deleteTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/lists/%s/delete-task/%s", testCase.listId, testCase.taskId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
