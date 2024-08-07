// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/sync_util_service.go
//
// Generated by this command:
//
//	mockgen -source=./internal/services/sync_util_service.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"

	entity "github.com/vskurikhin/gofavorites/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockSyncUtilService is a mock of SyncUtilService interface.
type MockSyncUtilService struct {
	ctrl     *gomock.Controller
	recorder *MockSyncUtilServiceMockRecorder
}

// MockSyncUtilServiceMockRecorder is the mock recorder for MockSyncUtilService.
type MockSyncUtilServiceMockRecorder struct {
	mock *MockSyncUtilService
}

// NewMockSyncUtilService creates a new mock instance.
func NewMockSyncUtilService(ctrl *gomock.Controller) *MockSyncUtilService {
	mock := &MockSyncUtilService{ctrl: ctrl}
	mock.recorder = &MockSyncUtilServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncUtilService) EXPECT() *MockSyncUtilServiceMockRecorder {
	return m.recorder
}

// Sync mocks base method.
func (m *MockSyncUtilService) Sync(ctx context.Context, mongodbFavorites, pgDBFavorites []entity.Favorites) ([]entity.Favorites, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", ctx, mongodbFavorites, pgDBFavorites)
	ret0, _ := ret[0].([]entity.Favorites)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockSyncUtilServiceMockRecorder) Sync(ctx, mongodbFavorites, pgDBFavorites any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockSyncUtilService)(nil).Sync), ctx, mongodbFavorites, pgDBFavorites)
}
