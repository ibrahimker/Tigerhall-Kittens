// Code generated by MockGen. DO NOT EDIT.
// Source: ./modules/sighting/v1/repository/postgres/sighting.go

// Package mock_postgres is a generated GoMock package.
package mock_postgres

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgconn "github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

// MockPgxPoolIface is a mock of PgxPoolIface interface.
type MockPgxPoolIface struct {
	ctrl     *gomock.Controller
	recorder *MockPgxPoolIfaceMockRecorder
}

// MockPgxPoolIfaceMockRecorder is the mock recorder for MockPgxPoolIface.
type MockPgxPoolIfaceMockRecorder struct {
	mock *MockPgxPoolIface
}

// NewMockPgxPoolIface creates a new mock instance.
func NewMockPgxPoolIface(ctrl *gomock.Controller) *MockPgxPoolIface {
	mock := &MockPgxPoolIface{ctrl: ctrl}
	mock.recorder = &MockPgxPoolIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPgxPoolIface) EXPECT() *MockPgxPoolIfaceMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockPgxPoolIface) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range arguments {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockPgxPoolIfaceMockRecorder) Exec(ctx, sql interface{}, arguments ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, arguments...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockPgxPoolIface)(nil).Exec), varargs...)
}

// Ping mocks base method.
func (m *MockPgxPoolIface) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPgxPoolIfaceMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPgxPoolIface)(nil).Ping), ctx)
}

// Query mocks base method.
func (m *MockPgxPoolIface) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockPgxPoolIfaceMockRecorder) Query(ctx, sql interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockPgxPoolIface)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockPgxPoolIface) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockPgxPoolIfaceMockRecorder) QueryRow(ctx, sql interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockPgxPoolIface)(nil).QueryRow), varargs...)
}
