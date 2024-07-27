// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/api_favorites_service.go
//
// Generated by this command:
//
//	mockgen -source=./internal/services/api_favorites_service.go -package=controllers
//

// Package controllers is a generated GoMock package.
package controllers

import (
	context "context"
	reflect "reflect"

	models "github.com/vskurikhin/gofavorites/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockApiFavoritesService is a mock of ApiFavoritesService interface.
type MockApiFavoritesService struct {
	ctrl     *gomock.Controller
	recorder *MockApiFavoritesServiceMockRecorder
}

// MockApiFavoritesServiceMockRecorder is the mock recorder for MockApiFavoritesService.
type MockApiFavoritesServiceMockRecorder struct {
	mock *MockApiFavoritesService
}

// NewMockApiFavoritesService creates a new mock instance.
func NewMockApiFavoritesService(ctrl *gomock.Controller) *MockApiFavoritesService {
	mock := &MockApiFavoritesService{ctrl: ctrl}
	mock.recorder = &MockApiFavoritesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApiFavoritesService) EXPECT() *MockApiFavoritesServiceMockRecorder {
	return m.recorder
}

// ApiFavoritesGet mocks base method.
func (m *MockApiFavoritesService) ApiFavoritesGet(ctx context.Context, model models.Favorites) (models.Favorites, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApiFavoritesGet", ctx, model)
	ret0, _ := ret[0].(models.Favorites)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApiFavoritesGet indicates an expected call of ApiFavoritesGet.
func (mr *MockApiFavoritesServiceMockRecorder) ApiFavoritesGet(ctx, model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApiFavoritesGet", reflect.TypeOf((*MockApiFavoritesService)(nil).ApiFavoritesGet), ctx, model)
}

// ApiFavoritesGetForUser mocks base method.
func (m *MockApiFavoritesService) ApiFavoritesGetForUser(ctx context.Context, model models.User) ([]models.Favorites, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApiFavoritesGetForUser", ctx, model)
	ret0, _ := ret[0].([]models.Favorites)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApiFavoritesGetForUser indicates an expected call of ApiFavoritesGetForUser.
func (mr *MockApiFavoritesServiceMockRecorder) ApiFavoritesGetForUser(ctx, model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApiFavoritesGetForUser", reflect.TypeOf((*MockApiFavoritesService)(nil).ApiFavoritesGetForUser), ctx, model)
}

// ApiFavoritesSet mocks base method.
func (m *MockApiFavoritesService) ApiFavoritesSet(ctx context.Context, model models.Favorites) (models.Favorites, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApiFavoritesSet", ctx, model)
	ret0, _ := ret[0].(models.Favorites)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApiFavoritesSet indicates an expected call of ApiFavoritesSet.
func (mr *MockApiFavoritesServiceMockRecorder) ApiFavoritesSet(ctx, model any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApiFavoritesSet", reflect.TypeOf((*MockApiFavoritesService)(nil).ApiFavoritesSet), ctx, model)
}
