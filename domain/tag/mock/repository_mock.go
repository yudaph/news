// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package tag_mock is a generated GoMock package.
package tag_mock

import (
	context "context"
	entities "news/domain/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateTag mocks base method.
func (m *MockRepository) CreateTag(ctx context.Context, tag *entities.Tag) (*entities.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTag", ctx, tag)
	ret0, _ := ret[0].(*entities.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTag indicates an expected call of CreateTag.
func (mr *MockRepositoryMockRecorder) CreateTag(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTag", reflect.TypeOf((*MockRepository)(nil).CreateTag), ctx, tag)
}

// DeleteTag mocks base method.
func (m *MockRepository) DeleteTag(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTag", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockRepositoryMockRecorder) DeleteTag(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockRepository)(nil).DeleteTag), ctx, id)
}

// GetAllTag mocks base method.
func (m *MockRepository) GetAllTag(ctx context.Context) (*entities.Tags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTag", ctx)
	ret0, _ := ret[0].(*entities.Tags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTag indicates an expected call of GetAllTag.
func (mr *MockRepositoryMockRecorder) GetAllTag(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTag", reflect.TypeOf((*MockRepository)(nil).GetAllTag), ctx)
}

// GetTagByIds mocks base method.
func (m *MockRepository) GetTagByIds(ctx context.Context, id []string) (*entities.Tags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagByIds", ctx, id)
	ret0, _ := ret[0].(*entities.Tags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagByIds indicates an expected call of GetTagByIds.
func (mr *MockRepositoryMockRecorder) GetTagByIds(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagByIds", reflect.TypeOf((*MockRepository)(nil).GetTagByIds), ctx, id)
}

// GetTagLike mocks base method.
func (m *MockRepository) GetTagLike(ctx context.Context, like string) (*entities.Tags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagLike", ctx, like)
	ret0, _ := ret[0].(*entities.Tags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagLike indicates an expected call of GetTagLike.
func (mr *MockRepositoryMockRecorder) GetTagLike(ctx, like interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagLike", reflect.TypeOf((*MockRepository)(nil).GetTagLike), ctx, like)
}

// UpdateTag mocks base method.
func (m *MockRepository) UpdateTag(ctx context.Context, tag *entities.Tag) (*entities.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTag", ctx, tag)
	ret0, _ := ret[0].(*entities.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTag indicates an expected call of UpdateTag.
func (mr *MockRepositoryMockRecorder) UpdateTag(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTag", reflect.TypeOf((*MockRepository)(nil).UpdateTag), ctx, tag)
}