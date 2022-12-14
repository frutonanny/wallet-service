// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_reserve is a generated GoMock package.
package mock_reserve

import (
	context "context"
	reflect "reflect"

	postgres "github.com/frutonanny/wallet-service/internal/postgres"
	reserve "github.com/frutonanny/wallet-service/internal/services/reserve"
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

// Reserve mocks base method.
func (m *MockWalletRepository) Reserve(ctx context.Context, walletID, cash int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reserve", ctx, walletID, cash)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reserve indicates an expected call of Reserve.
func (mr *MockWalletRepositoryMockRecorder) Reserve(ctx, walletID, cash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reserve", reflect.TypeOf((*MockWalletRepository)(nil).Reserve), ctx, walletID, cash)
}

// MockOrderRepository is a mock of OrderRepository interface.
type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
}

// MockOrderRepositoryMockRecorder is the mock recorder for MockOrderRepository.
type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

// NewMockOrderRepository creates a new mock instance.
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

// AddOrderTransactions mocks base method.
func (m *MockOrderRepository) AddOrderTransactions(ctx context.Context, orderID int64, nameType string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrderTransactions", ctx, orderID, nameType)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOrderTransactions indicates an expected call of AddOrderTransactions.
func (mr *MockOrderRepositoryMockRecorder) AddOrderTransactions(ctx, orderID, nameType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrderTransactions", reflect.TypeOf((*MockOrderRepository)(nil).AddOrderTransactions), ctx, orderID, nameType)
}

// CreateOrder mocks base method.
func (m *MockOrderRepository) CreateOrder(ctx context.Context, walletID, externalID, serviceID, amount int64, status string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, walletID, externalID, serviceID, amount, status)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderRepositoryMockRecorder) CreateOrder(ctx, walletID, externalID, serviceID, amount, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderRepository)(nil).CreateOrder), ctx, walletID, externalID, serviceID, amount, status)
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

// AddTransaction mocks base method.
func (m *MockTransactionRepository) AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTransaction", ctx, walletID, action, payload, amount)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTransaction indicates an expected call of AddTransaction.
func (mr *MockTransactionRepositoryMockRecorder) AddTransaction(ctx, walletID, action, payload, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).AddTransaction), ctx, walletID, action, payload, amount)
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

// NewOrderRepository mocks base method.
func (m *Mockdependencies) NewOrderRepository(db postgres.Database) reserve.OrderRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewOrderRepository", db)
	ret0, _ := ret[0].(reserve.OrderRepository)
	return ret0
}

// NewOrderRepository indicates an expected call of NewOrderRepository.
func (mr *MockdependenciesMockRecorder) NewOrderRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewOrderRepository", reflect.TypeOf((*Mockdependencies)(nil).NewOrderRepository), db)
}

// NewTransactionRepository mocks base method.
func (m *Mockdependencies) NewTransactionRepository(db postgres.Database) reserve.TransactionRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTransactionRepository", db)
	ret0, _ := ret[0].(reserve.TransactionRepository)
	return ret0
}

// NewTransactionRepository indicates an expected call of NewTransactionRepository.
func (mr *MockdependenciesMockRecorder) NewTransactionRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTransactionRepository", reflect.TypeOf((*Mockdependencies)(nil).NewTransactionRepository), db)
}

// NewWalletRepository mocks base method.
func (m *Mockdependencies) NewWalletRepository(db postgres.Database) reserve.WalletRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWalletRepository", db)
	ret0, _ := ret[0].(reserve.WalletRepository)
	return ret0
}

// NewWalletRepository indicates an expected call of NewWalletRepository.
func (mr *MockdependenciesMockRecorder) NewWalletRepository(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWalletRepository", reflect.TypeOf((*Mockdependencies)(nil).NewWalletRepository), db)
}
