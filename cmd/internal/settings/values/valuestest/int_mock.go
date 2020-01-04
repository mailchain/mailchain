// Code generated by MockGen. DO NOT EDIT.
// Source: int.go

// Package valuestest is a generated GoMock package.
package valuestest

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	output "github.com/mailchain/mailchain/cmd/internal/settings/output"
)

// MockInt is a mock of Int interface
type MockInt struct {
	ctrl     *gomock.Controller
	recorder *MockIntMockRecorder
}

// MockIntMockRecorder is the mock recorder for MockInt
type MockIntMockRecorder struct {
	mock *MockInt
}

// NewMockInt creates a new mock instance
func NewMockInt(ctrl *gomock.Controller) *MockInt {
	mock := &MockInt{ctrl: ctrl}
	mock.recorder = &MockIntMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInt) EXPECT() *MockIntMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockInt) Get() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(int)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockIntMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInt)(nil).Get))
}

// Set mocks base method
func (m *MockInt) Set(v int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", v)
}

// Set indicates an expected call of Set
func (mr *MockIntMockRecorder) Set(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockInt)(nil).Set), v)
}

// Attribute mocks base method
func (m *MockInt) Attribute() output.Attribute {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attribute")
	ret0, _ := ret[0].(output.Attribute)
	return ret0
}

// Attribute indicates an expected call of Attribute
func (mr *MockIntMockRecorder) Attribute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attribute", reflect.TypeOf((*MockInt)(nil).Attribute))
}
