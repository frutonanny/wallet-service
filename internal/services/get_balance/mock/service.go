// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_get_balance is a generated GoMock package.
package mock_get_balance

import (
	context "context"
	reflect "reflect"

	postgres "github.com/frutonanny/wallet-service/internal/postgres"
	get_balance "github.com/frutonanny/wallet-service/internal/services/get_balance"
	gomock "github.com/golang/mock/gomock"
)

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *Mocklogger) Error(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", msg)
}

// Error indicates an expected call of Error.
func (mr *MockloggerMockRecorder) Error(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mocklogger)(nil).Error), msg)
}

// Info mocks base method.
func (m *Mocklogger) Info(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", msg)
}

// Info indicates an expected call of Info.
func (mr *MockloggerMockRecorder) Info(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklogger)(nil).Info), msg)
}

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

// ExistWallet mocks base method.
func (m *MockRepository) ExistWallet(ctx context.Context, userID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistWallet", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExistWallet indicates an expected call of ExistWallet.
func (mr *MockRepositoryMockRecorder) ExistWallet(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistWallet", reflect.TypeOf((*MockRepository)(nil).ExistWallet), ctx, userID)
}

// GetBalance mocks base method.
func (m *MockRepository) GetBalance(ctx context.Context, walletID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, walletID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockRepositoryMockRecorder) GetBalance(ctx, walletID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockRepository)(nil).GetBalance), ctx, walletID)
}

// Mockdependencies is a mock of dependencies interface.
type Mockdependencies struct {
	ctrl     *gomock.Controller
	recorder *MockdependenciesMockRecorder
}

// MockdependenciesMockRecorder is the mock recorder for Mockdependencies.
type MockdependenciesMockRecorder struct {
	mock *Mockdependencies
}

// NewMockdependencies creates a new mock instance.
func NewMockdependencies(ctrl *gomock.Controller) *Mockdependencies {
	mock := &Mockdependencies{ctrl: ctrl}
	mock.recorder = &MockdependenciesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdependencies) EXPECT() *MockdependenciesMockRecorder {
	return m.recorder
}

// NewRepository mocks base method.
func (m *Mockdependencies) NewRepository(db postgres.Database) get_balance.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRepository", db)
	ret0, _ := ret[0].(get_balance.Repository)
	return ret0
}

// NewRepository indicates an expected call of NewRepository.
func (mr *MockdependenciesMockRecorder) NewRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRepository", reflect.TypeOf((*Mockdependencies)(nil).NewRepository), db)
}
