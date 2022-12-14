// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/swordhealth-api/internal/service (interfaces: AuthService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// DecodeBasicAuth mocks base method.
func (m *MockAuthService) DecodeBasicAuth(arg0 context.Context, arg1 string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeBasicAuth", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DecodeBasicAuth indicates an expected call of DecodeBasicAuth.
func (mr *MockAuthServiceMockRecorder) DecodeBasicAuth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeBasicAuth", reflect.TypeOf((*MockAuthService)(nil).DecodeBasicAuth), arg0, arg1)
}
