// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/swordhealth-api/internal/repository (interfaces: HealthRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHealthRepository is a mock of HealthRepository interface.
type MockHealthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockHealthRepositoryMockRecorder
}

// MockHealthRepositoryMockRecorder is the mock recorder for MockHealthRepository.
type MockHealthRepositoryMockRecorder struct {
	mock *MockHealthRepository
}

// NewMockHealthRepository creates a new mock instance.
func NewMockHealthRepository(ctrl *gomock.Controller) *MockHealthRepository {
	mock := &MockHealthRepository{ctrl: ctrl}
	mock.recorder = &MockHealthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthRepository) EXPECT() *MockHealthRepositoryMockRecorder {
	return m.recorder
}

// Health mocks base method.
func (m *MockHealthRepository) Health(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Health indicates an expected call of Health.
func (mr *MockHealthRepositoryMockRecorder) Health(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockHealthRepository)(nil).Health), arg0)
}
