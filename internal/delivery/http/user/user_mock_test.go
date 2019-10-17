// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package user is a generated GoMock package.
package user

import (
	context "context"
	user "github.com/essajiwa/teratur/internal/service/user"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIUserSvc is a mock of IUserSvc interface
type MockIUserSvc struct {
	ctrl     *gomock.Controller
	recorder *MockIUserSvcMockRecorder
}

// MockIUserSvcMockRecorder is the mock recorder for MockIUserSvc
type MockIUserSvcMockRecorder struct {
	mock *MockIUserSvc
}

// NewMockIUserSvc creates a new mock instance
func NewMockIUserSvc(ctrl *gomock.Controller) *MockIUserSvc {
	mock := &MockIUserSvc{ctrl: ctrl}
	mock.recorder = &MockIUserSvcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUserSvc) EXPECT() *MockIUserSvcMockRecorder {
	return m.recorder
}

// GetUserByID mocks base method
func (m *MockIUserSvc) GetUserByID(ctx context.Context, userID int64) (user.User, error) {
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID
func (mr *MockIUserSvcMockRecorder) GetUserByID(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockIUserSvc)(nil).GetUserByID), ctx, userID)
}
