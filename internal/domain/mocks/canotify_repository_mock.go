// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases (interfaces: CanNotifyRepository)
//
// Generated by this command:
//
//	mockgen -destination=../mocks/canotify_repository_mock.go -package=mocks github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases CanNotifyRepository
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	entities "github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockCanNotifyRepository is a mock of CanNotifyRepository interface.
type MockCanNotifyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCanNotifyRepositoryMockRecorder
}

// MockCanNotifyRepositoryMockRecorder is the mock recorder for MockCanNotifyRepository.
type MockCanNotifyRepositoryMockRecorder struct {
	mock *MockCanNotifyRepository
}

// NewMockCanNotifyRepository creates a new mock instance.
func NewMockCanNotifyRepository(ctrl *gomock.Controller) *MockCanNotifyRepository {
	mock := &MockCanNotifyRepository{ctrl: ctrl}
	mock.recorder = &MockCanNotifyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCanNotifyRepository) EXPECT() *MockCanNotifyRepositoryMockRecorder {
	return m.recorder
}

// CheckSentNotifications mocks base method.
func (m *MockCanNotifyRepository) CheckSentNotifications(arg0 entities.Identifier, arg1 string) (entities.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSentNotifications", arg0, arg1)
	ret0, _ := ret[0].(entities.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckSentNotifications indicates an expected call of CheckSentNotifications.
func (mr *MockCanNotifyRepositoryMockRecorder) CheckSentNotifications(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSentNotifications", reflect.TypeOf((*MockCanNotifyRepository)(nil).CheckSentNotifications), arg0, arg1)
}

// RegisterNotification mocks base method.
func (m *MockCanNotifyRepository) RegisterNotification(arg0 entities.Notification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterNotification", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterNotification indicates an expected call of RegisterNotification.
func (mr *MockCanNotifyRepositoryMockRecorder) RegisterNotification(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterNotification", reflect.TypeOf((*MockCanNotifyRepository)(nil).RegisterNotification), arg0)
}
