// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_add is a generated GoMock package.
package mock_add

import (
	context "context"
	reflect "reflect"

	postgres "github.com/frutonanny/wallet-service/internal/postgres"
	add "github.com/frutonanny/wallet-service/internal/services/add"
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

// Add mocks base method.
func (m *MockRepository) Add(ctx context.Context, walletID, cash int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, walletID, cash)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockRepositoryMockRecorder) Add(ctx, walletID, cash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockRepository)(nil).Add), ctx, walletID, cash)
}

// AddTransaction mocks base method.
func (m *MockRepository) AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTransaction", ctx, walletID, action, payload, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTransaction indicates an expected call of AddTransaction.
func (mr *MockRepositoryMockRecorder) AddTransaction(ctx, walletID, action, payload, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTransaction", reflect.TypeOf((*MockRepository)(nil).AddTransaction), ctx, walletID, action, payload, amount)
}

// CreateWallet mocks base method.
func (m *MockRepository) CreateWallet(ctx context.Context, userID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWallet", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWallet indicates an expected call of CreateWallet.
func (mr *MockRepositoryMockRecorder) CreateWallet(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWallet", reflect.TypeOf((*MockRepository)(nil).CreateWallet), ctx, userID)
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

// MockRepoBuilder is a mock of RepoBuilder interface.
type MockRepoBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockRepoBuilderMockRecorder
}

// MockRepoBuilderMockRecorder is the mock recorder for MockRepoBuilder.
type MockRepoBuilderMockRecorder struct {
	mock *MockRepoBuilder
}

// NewMockRepoBuilder creates a new mock instance.
func NewMockRepoBuilder(ctrl *gomock.Controller) *MockRepoBuilder {
	mock := &MockRepoBuilder{ctrl: ctrl}
	mock.recorder = &MockRepoBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepoBuilder) EXPECT() *MockRepoBuilderMockRecorder {
	return m.recorder
}

// NewRepository mocks base method.
func (m *MockRepoBuilder) NewRepository(db postgres.Database) add.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRepository", db)
	ret0, _ := ret[0].(add.Repository)
	return ret0
}

// NewRepository indicates an expected call of NewRepository.
func (mr *MockRepoBuilderMockRecorder) NewRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRepository", reflect.TypeOf((*MockRepoBuilder)(nil).NewRepository), db)
}
