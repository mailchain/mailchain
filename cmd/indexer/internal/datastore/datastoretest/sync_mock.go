// Code generated by MockGen. DO NOT EDIT.
// Source: sync.go

// Package datastoretest is a generated GoMock package.
package datastoretest

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSyncStore is a mock of SyncStore interface
type MockSyncStore struct {
	ctrl     *gomock.Controller
	recorder *MockSyncStoreMockRecorder
}

// MockSyncStoreMockRecorder is the mock recorder for MockSyncStore
type MockSyncStoreMockRecorder struct {
	mock *MockSyncStore
}

// NewMockSyncStore creates a new mock instance
func NewMockSyncStore(ctrl *gomock.Controller) *MockSyncStore {
	mock := &MockSyncStore{ctrl: ctrl}
	mock.recorder = &MockSyncStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSyncStore) EXPECT() *MockSyncStoreMockRecorder {
	return m.recorder
}

// GetBlockNumber mocks base method
func (m *MockSyncStore) GetBlockNumber(ctx context.Context, protocol, network string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockNumber", ctx, protocol, network)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockNumber indicates an expected call of GetBlockNumber
func (mr *MockSyncStoreMockRecorder) GetBlockNumber(ctx, protocol, network interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockNumber", reflect.TypeOf((*MockSyncStore)(nil).GetBlockNumber), ctx, protocol, network)
}

// PutBlockNumber mocks base method
func (m *MockSyncStore) PutBlockNumber(ctx context.Context, protocol, network string, blockNo uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutBlockNumber", ctx, protocol, network, blockNo)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutBlockNumber indicates an expected call of PutBlockNumber
func (mr *MockSyncStoreMockRecorder) PutBlockNumber(ctx, protocol, network, blockNo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutBlockNumber", reflect.TypeOf((*MockSyncStore)(nil).PutBlockNumber), ctx, protocol, network, blockNo)
}
