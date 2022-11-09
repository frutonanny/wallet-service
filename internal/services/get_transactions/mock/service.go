// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_get_transactions is a generated GoMock package.
package mock_get_transactions

import (
	context "context"
	reflect "reflect"

	postgres "github.com/frutonanny/wallet-service/internal/postgres"
	transaction "github.com/frutonanny/wallet-service/internal/repositories/transaction"
	get_transactions "github.com/frutonanny/wallet-service/internal/services/get_transactions"
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

// MockWalletRepository is a mock of WalletRepository interface.
type MockWalletRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWalletRepositoryMockRecorder
}

// MockWalletRepositoryMockRecorder is the mock recorder for MockWalletRepository.
type MockWalletRepositoryMockRecorder struct {
	mock *MockWalletRepository
}

// NewMockWalletRepository creates a new mock instance.
func NewMockWalletRepository(ctrl *gomock.Controller) *MockWalletRepository {
	mock := &MockWalletRepository{ctrl: ctrl}
	mock.recorder = &MockWalletRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletRepository) EXPECT() *MockWalletRepositoryMockRecorder {
	return m.recorder
}

// ExistWallet mocks base method.
func (m *MockWalletRepository) ExistWallet(ctx context.Context, userID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistWallet", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExistWallet indicates an expected call of ExistWallet.
func (mr *MockWalletRepositoryMockRecorder) ExistWallet(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistWallet", reflect.TypeOf((*MockWalletRepository)(nil).ExistWallet), ctx, userID)
}

// MockTransactionRepository is a mock of TransactionRepository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// GetTransactions mocks base method.
func (m *MockTransactionRepository) GetTransactions(ctx context.Context, walletID, limit, offset int64, sortBy transaction.SortBy, direction transaction.Direction) ([]transaction.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions", ctx, walletID, limit, offset, sortBy, direction)
	ret0, _ := ret[0].([]transaction.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions.
func (mr *MockTransactionRepositoryMockRecorder) GetTransactions(ctx, walletID, limit, offset, sortBy, direction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockTransactionRepository)(nil).GetTransactions), ctx, walletID, limit, offset, sortBy, direction)
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

// NewTransactionRepository mocks base method.
func (m *Mockdependencies) NewTransactionRepository(db postgres.Database) get_transactions.TransactionRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTransactionRepository", db)
	ret0, _ := ret[0].(get_transactions.TransactionRepository)
	return ret0
}

// NewTransactionRepository indicates an expected call of NewTransactionRepository.
func (mr *MockdependenciesMockRecorder) NewTransactionRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTransactionRepository", reflect.TypeOf((*Mockdependencies)(nil).NewTransactionRepository), db)
}

// NewWalletRepository mocks base method.
func (m *Mockdependencies) NewWalletRepository(db postgres.Database) get_transactions.WalletRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWalletRepository", db)
	ret0, _ := ret[0].(get_transactions.WalletRepository)
	return ret0
}

// NewWalletRepository indicates an expected call of NewWalletRepository.
func (mr *MockdependenciesMockRecorder) NewWalletRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWalletRepository", reflect.TypeOf((*Mockdependencies)(nil).NewWalletRepository), db)
}