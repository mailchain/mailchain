// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package processortest is a generated GoMock package.
package processortest

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	processor "github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	reflect "reflect"
)

// MockTransaction is a mock of Transaction interface
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockTransaction) Run(ctx context.Context, protocol, network string, tx interface{}, txOpts processor.TransactionOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx, protocol, network, tx, txOpts)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
func (mr *MockTransactionMockRecorder) Run(ctx, protocol, network, tx, txOpts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockTransaction)(nil).Run), ctx, protocol, network, tx, txOpts)
}

// MockTransactionOptions is a mock of TransactionOptions interface
type MockTransactionOptions struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionOptionsMockRecorder
}

// MockTransactionOptionsMockRecorder is the mock recorder for MockTransactionOptions
type MockTransactionOptionsMockRecorder struct {
	mock *MockTransactionOptions
}

// NewMockTransactionOptions creates a new mock instance
func NewMockTransactionOptions(ctrl *gomock.Controller) *MockTransactionOptions {
	mock := &MockTransactionOptions{ctrl: ctrl}
	mock.recorder = &MockTransactionOptionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTransactionOptions) EXPECT() *MockTransactionOptionsMockRecorder {
	return m.recorder
}
