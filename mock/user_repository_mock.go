// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/swordhealth-api/internal/repository (interfaces: UserRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/viniosilva/swordhealth-api/internal/dto"
	model "github.com/viniosilva/swordhealth-api/internal/model"
	repository "github.com/viniosilva/swordhealth-api/internal/repository"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(arg0 context.Context, arg1 dto.CreateUserDto) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), arg0, arg1)
}

// GetUserByID mocks base method.
func (m *MockUserRepository) GetUserByID(arg0 context.Context, arg1 int) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepositoryMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), arg0, arg1)
}

// GetUserByUsernameAndPassword mocks base method.
func (m *MockUserRepository) GetUserByUsernameAndPassword(arg0 context.Context, arg1, arg2 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsernameAndPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsernameAndPassword indicates an expected call of GetUserByUsernameAndPassword.
func (mr *MockUserRepositoryMockRecorder) GetUserByUsernameAndPassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsernameAndPassword", reflect.TypeOf((*MockUserRepository)(nil).GetUserByUsernameAndPassword), arg0, arg1, arg2)
}

// ListUsers mocks base method.
func (m *MockUserRepository) ListUsers(arg0 context.Context, arg1, arg2 int, arg3 ...repository.WhereOpt) ([]model.User, int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListUsers", varargs...)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockUserRepositoryMockRecorder) ListUsers(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockUserRepository)(nil).ListUsers), varargs...)
}
