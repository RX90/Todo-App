// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	todo "github.com/RX90/Todo-App/server/internal/todo"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// CheckRefreshToken mocks base method.
func (m *MockAuthorization) CheckRefreshToken(userId, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRefreshToken", userId, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckRefreshToken indicates an expected call of CheckRefreshToken.
func (mr *MockAuthorizationMockRecorder) CheckRefreshToken(userId, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRefreshToken", reflect.TypeOf((*MockAuthorization)(nil).CheckRefreshToken), userId, refreshToken)
}

// CreateUser mocks base method.
func (m *MockAuthorization) CreateUser(user todo.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorization)(nil).CreateUser), user)
}

// DeleteRefreshToken mocks base method.
func (m *MockAuthorization) DeleteRefreshToken(userId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRefreshToken", userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRefreshToken indicates an expected call of DeleteRefreshToken.
func (mr *MockAuthorizationMockRecorder) DeleteRefreshToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRefreshToken", reflect.TypeOf((*MockAuthorization)(nil).DeleteRefreshToken), userId)
}

// GetUserId mocks base method.
func (m *MockAuthorization) GetUserId(user todo.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserId", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserId indicates an expected call of GetUserId.
func (mr *MockAuthorizationMockRecorder) GetUserId(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserId", reflect.TypeOf((*MockAuthorization)(nil).GetUserId), user)
}

// NewAccessToken mocks base method.
func (m *MockAuthorization) NewAccessToken(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAccessToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewAccessToken indicates an expected call of NewAccessToken.
func (mr *MockAuthorizationMockRecorder) NewAccessToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAccessToken", reflect.TypeOf((*MockAuthorization)(nil).NewAccessToken), userId)
}

// NewRefreshToken mocks base method.
func (m *MockAuthorization) NewRefreshToken(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRefreshToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRefreshToken indicates an expected call of NewRefreshToken.
func (mr *MockAuthorizationMockRecorder) NewRefreshToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRefreshToken", reflect.TypeOf((*MockAuthorization)(nil).NewRefreshToken), userId)
}

// ParseAccessToken mocks base method.
func (m *MockAuthorization) ParseAccessToken(token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseAccessToken", token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseAccessToken indicates an expected call of ParseAccessToken.
func (mr *MockAuthorizationMockRecorder) ParseAccessToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseAccessToken", reflect.TypeOf((*MockAuthorization)(nil).ParseAccessToken), token)
}

// MockTodoList is a mock of TodoList interface.
type MockTodoList struct {
	ctrl     *gomock.Controller
	recorder *MockTodoListMockRecorder
}

// MockTodoListMockRecorder is the mock recorder for MockTodoList.
type MockTodoListMockRecorder struct {
	mock *MockTodoList
}

// NewMockTodoList creates a new mock instance.
func NewMockTodoList(ctrl *gomock.Controller) *MockTodoList {
	mock := &MockTodoList{ctrl: ctrl}
	mock.recorder = &MockTodoListMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoList) EXPECT() *MockTodoListMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodoList) Create(userId string, list todo.List) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userId, list)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTodoListMockRecorder) Create(userId, list interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodoList)(nil).Create), userId, list)
}

// Delete mocks base method.
func (m *MockTodoList) Delete(userId, listId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userId, listId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTodoListMockRecorder) Delete(userId, listId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTodoList)(nil).Delete), userId, listId)
}

// GetAll mocks base method.
func (m *MockTodoList) GetAll(userId string) ([]todo.List, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", userId)
	ret0, _ := ret[0].([]todo.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTodoListMockRecorder) GetAll(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTodoList)(nil).GetAll), userId)
}

// Update mocks base method.
func (m *MockTodoList) Update(userId, listId string, list todo.List) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", userId, listId, list)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTodoListMockRecorder) Update(userId, listId, list interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTodoList)(nil).Update), userId, listId, list)
}

// MockTodoTask is a mock of TodoTask interface.
type MockTodoTask struct {
	ctrl     *gomock.Controller
	recorder *MockTodoTaskMockRecorder
}

// MockTodoTaskMockRecorder is the mock recorder for MockTodoTask.
type MockTodoTaskMockRecorder struct {
	mock *MockTodoTask
}

// NewMockTodoTask creates a new mock instance.
func NewMockTodoTask(ctrl *gomock.Controller) *MockTodoTask {
	mock := &MockTodoTask{ctrl: ctrl}
	mock.recorder = &MockTodoTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoTask) EXPECT() *MockTodoTaskMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodoTask) Create(userId, listId string, task todo.Task) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userId, listId, task)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTodoTaskMockRecorder) Create(userId, listId, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodoTask)(nil).Create), userId, listId, task)
}

// Delete mocks base method.
func (m *MockTodoTask) Delete(userId, listId, taskId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userId, listId, taskId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTodoTaskMockRecorder) Delete(userId, listId, taskId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTodoTask)(nil).Delete), userId, listId, taskId)
}

// GetAll mocks base method.
func (m *MockTodoTask) GetAll(userId, listId string) ([]todo.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", userId, listId)
	ret0, _ := ret[0].([]todo.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTodoTaskMockRecorder) GetAll(userId, listId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTodoTask)(nil).GetAll), userId, listId)
}

// Update mocks base method.
func (m *MockTodoTask) Update(userId, listId, taskId string, task todo.UpdateTaskInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", userId, listId, taskId, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTodoTaskMockRecorder) Update(userId, listId, taskId, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTodoTask)(nil).Update), userId, listId, taskId, task)
}
