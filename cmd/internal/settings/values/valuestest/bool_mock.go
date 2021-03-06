// Code generated by MockGen. DO NOT EDIT.
// Source: bool.go

// Package valuestest is a generated GoMock package.
package valuestest

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	output "github.com/mailchain/mailchain/cmd/internal/settings/output"
)

// MockBool is a mock of Bool interface.
type MockBool struct {
	ctrl     *gomock.Controller
	recorder *MockBoolMockRecorder
}

// MockBoolMockRecorder is the mock recorder for MockBool.
type MockBoolMockRecorder struct {
	mock *MockBool
}

// NewMockBool creates a new mock instance.
func NewMockBool(ctrl *gomock.Controller) *MockBool {
	mock := &MockBool{ctrl: ctrl}
	mock.recorder = &MockBoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBool) EXPECT() *MockBoolMockRecorder {
	return m.recorder
}

// Attribute mocks base method.
func (m *MockBool) Attribute() output.Attribute {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attribute")
	ret0, _ := ret[0].(output.Attribute)
	return ret0
}

// Attribute indicates an expected call of Attribute.
func (mr *MockBoolMockRecorder) Attribute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attribute", reflect.TypeOf((*MockBool)(nil).Attribute))
}

// Get mocks base method.
func (m *MockBool) Get() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockBoolMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBool)(nil).Get))
}

// Set mocks base method.
func (m *MockBool) Set(v bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", v)
}

// Set indicates an expected call of Set.
func (mr *MockBoolMockRecorder) Set(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockBool)(nil).Set), v)
}
