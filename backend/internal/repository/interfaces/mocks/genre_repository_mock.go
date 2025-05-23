// Code generated by MockGen. DO NOT EDIT.
// Source: src/src/internal/repository/interfaces/genre_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "music-service/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockGenreRepository is a mock of GenreRepository interface.
type MockGenreRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGenreRepositoryMockRecorder
}

// MockGenreRepositoryMockRecorder is the mock recorder for MockGenreRepository.
type MockGenreRepositoryMockRecorder struct {
	mock *MockGenreRepository
}

// NewMockGenreRepository creates a new mock instance.
func NewMockGenreRepository(ctrl *gomock.Controller) *MockGenreRepository {
	mock := &MockGenreRepository{ctrl: ctrl}
	mock.recorder = &MockGenreRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenreRepository) EXPECT() *MockGenreRepositoryMockRecorder {
	return m.recorder
}

// AddGenreToTrack mocks base method.
func (m *MockGenreRepository) AddGenreToTrack(trackID, genreID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddGenreToTrack", trackID, genreID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddGenreToTrack indicates an expected call of AddGenreToTrack.
func (mr *MockGenreRepositoryMockRecorder) AddGenreToTrack(trackID, genreID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGenreToTrack", reflect.TypeOf((*MockGenreRepository)(nil).AddGenreToTrack), trackID, genreID)
}

// Delete mocks base method.
func (m *MockGenreRepository) Delete(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockGenreRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockGenreRepository)(nil).Delete), id)
}

// FindByID mocks base method.
func (m *MockGenreRepository) FindByID(id uuid.UUID) (*models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockGenreRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockGenreRepository)(nil).FindByID), id)
}

// GetGenresForTrack mocks base method.
func (m *MockGenreRepository) GetGenresForTrack(trackID uuid.UUID) ([]*models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenresForTrack", trackID)
	ret0, _ := ret[0].([]*models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenresForTrack indicates an expected call of GetGenresForTrack.
func (mr *MockGenreRepositoryMockRecorder) GetGenresForTrack(trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenresForTrack", reflect.TypeOf((*MockGenreRepository)(nil).GetGenresForTrack), trackID)
}

// ListAll mocks base method.
func (m *MockGenreRepository) ListAll() ([]*models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAll")
	ret0, _ := ret[0].([]*models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAll indicates an expected call of ListAll.
func (mr *MockGenreRepositoryMockRecorder) ListAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAll", reflect.TypeOf((*MockGenreRepository)(nil).ListAll))
}

// RemoveGenreFromTrack mocks base method.
func (m *MockGenreRepository) RemoveGenreFromTrack(trackID, genreID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveGenreFromTrack", trackID, genreID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveGenreFromTrack indicates an expected call of RemoveGenreFromTrack.
func (mr *MockGenreRepositoryMockRecorder) RemoveGenreFromTrack(trackID, genreID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveGenreFromTrack", reflect.TypeOf((*MockGenreRepository)(nil).RemoveGenreFromTrack), trackID, genreID)
}

// Save mocks base method.
func (m *MockGenreRepository) Save(genre *models.Genre) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", genre)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockGenreRepositoryMockRecorder) Save(genre interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockGenreRepository)(nil).Save), genre)
}
