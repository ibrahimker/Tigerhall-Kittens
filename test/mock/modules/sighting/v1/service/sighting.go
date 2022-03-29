// Code generated by MockGen. DO NOT EDIT.
// Source: ./modules/sighting/v1/service/sighting.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/ibrahimker/tigerhall-kittens/modules/sighting/v1/entity"
)

// MockTigerSighting is a mock of TigerSighting interface.
type MockTigerSighting struct {
	ctrl     *gomock.Controller
	recorder *MockTigerSightingMockRecorder
}

// MockTigerSightingMockRecorder is the mock recorder for MockTigerSighting.
type MockTigerSightingMockRecorder struct {
	mock *MockTigerSighting
}

// NewMockTigerSighting creates a new mock instance.
func NewMockTigerSighting(ctrl *gomock.Controller) *MockTigerSighting {
	mock := &MockTigerSighting{ctrl: ctrl}
	mock.recorder = &MockTigerSightingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTigerSighting) EXPECT() *MockTigerSightingMockRecorder {
	return m.recorder
}

// CreateSighting mocks base method.
func (m *MockTigerSighting) CreateSighting(ctx context.Context, tigerID int32, sighting *entity.Sighting) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSighting", ctx, tigerID, sighting)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSighting indicates an expected call of CreateSighting.
func (mr *MockTigerSightingMockRecorder) CreateSighting(ctx, tigerID, sighting interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSighting", reflect.TypeOf((*MockTigerSighting)(nil).CreateSighting), ctx, tigerID, sighting)
}

// CreateTiger mocks base method.
func (m *MockTigerSighting) CreateTiger(ctx context.Context, tiger *entity.Tiger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTiger", ctx, tiger)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTiger indicates an expected call of CreateTiger.
func (mr *MockTigerSightingMockRecorder) CreateTiger(ctx, tiger interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTiger", reflect.TypeOf((*MockTigerSighting)(nil).CreateTiger), ctx, tiger)
}

// GetSightingsByTigerID mocks base method.
func (m *MockTigerSighting) GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSightingsByTigerID", ctx, tigerID)
	ret0, _ := ret[0].([]*entity.Sighting)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSightingsByTigerID indicates an expected call of GetSightingsByTigerID.
func (mr *MockTigerSightingMockRecorder) GetSightingsByTigerID(ctx, tigerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSightingsByTigerID", reflect.TypeOf((*MockTigerSighting)(nil).GetSightingsByTigerID), ctx, tigerID)
}

// GetTigers mocks base method.
func (m *MockTigerSighting) GetTigers(ctx context.Context) ([]*entity.Tiger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTigers", ctx)
	ret0, _ := ret[0].([]*entity.Tiger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTigers indicates an expected call of GetTigers.
func (mr *MockTigerSightingMockRecorder) GetTigers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTigers", reflect.TypeOf((*MockTigerSighting)(nil).GetTigers), ctx)
}

// MockTigerSightingRepository is a mock of TigerSightingRepository interface.
type MockTigerSightingRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTigerSightingRepositoryMockRecorder
}

// MockTigerSightingRepositoryMockRecorder is the mock recorder for MockTigerSightingRepository.
type MockTigerSightingRepositoryMockRecorder struct {
	mock *MockTigerSightingRepository
}

// NewMockTigerSightingRepository creates a new mock instance.
func NewMockTigerSightingRepository(ctrl *gomock.Controller) *MockTigerSightingRepository {
	mock := &MockTigerSightingRepository{ctrl: ctrl}
	mock.recorder = &MockTigerSightingRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTigerSightingRepository) EXPECT() *MockTigerSightingRepositoryMockRecorder {
	return m.recorder
}

// CreateSighting mocks base method.
func (m *MockTigerSightingRepository) CreateSighting(ctx context.Context, tigerID int32, sighting *entity.Sighting) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSighting", ctx, tigerID, sighting)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSighting indicates an expected call of CreateSighting.
func (mr *MockTigerSightingRepositoryMockRecorder) CreateSighting(ctx, tigerID, sighting interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSighting", reflect.TypeOf((*MockTigerSightingRepository)(nil).CreateSighting), ctx, tigerID, sighting)
}

// CreateTiger mocks base method.
func (m *MockTigerSightingRepository) CreateTiger(ctx context.Context, tiger *entity.Tiger) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTiger", ctx, tiger)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTiger indicates an expected call of CreateTiger.
func (mr *MockTigerSightingRepositoryMockRecorder) CreateTiger(ctx, tiger interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTiger", reflect.TypeOf((*MockTigerSightingRepository)(nil).CreateTiger), ctx, tiger)
}

// GetSightingsByTigerID mocks base method.
func (m *MockTigerSightingRepository) GetSightingsByTigerID(ctx context.Context, tigerID int32) ([]*entity.Sighting, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSightingsByTigerID", ctx, tigerID)
	ret0, _ := ret[0].([]*entity.Sighting)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSightingsByTigerID indicates an expected call of GetSightingsByTigerID.
func (mr *MockTigerSightingRepositoryMockRecorder) GetSightingsByTigerID(ctx, tigerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSightingsByTigerID", reflect.TypeOf((*MockTigerSightingRepository)(nil).GetSightingsByTigerID), ctx, tigerID)
}

// GetTigers mocks base method.
func (m *MockTigerSightingRepository) GetTigers(ctx context.Context) ([]*entity.Tiger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTigers", ctx)
	ret0, _ := ret[0].([]*entity.Tiger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTigers indicates an expected call of GetTigers.
func (mr *MockTigerSightingRepositoryMockRecorder) GetTigers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTigers", reflect.TypeOf((*MockTigerSightingRepository)(nil).GetTigers), ctx)
}
