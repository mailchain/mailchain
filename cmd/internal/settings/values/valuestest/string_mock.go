// Code generated by MockGen. DO NOT EDIT.
// Source: string.go

// Package valuestest is a generated GoMock package.
package valuestest

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	output "github.com/mailchain/mailchain/cmd/internal/settings/output"
)

// MockString is a mock of String interface
type MockString struct {
	ctrl     *gomock.Controller
	recorder *MockStringMockRecorder
}

// MockStringMockRecorder is the mock recorder for MockString
type MockStringMockRecorder struct {
	mock *MockString
}

// NewMockString creates a new mock instance
func NewMockString(ctrl *gomock.Controller) *MockString {
	mock := &MockString{ctrl: ctrl}
	mock.recorder = &MockStringMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockString) EXPECT() *MockStringMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockString) Get() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(string)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockStringMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockString)(nil).Get))
}

// Set mocks base method
func (m *MockString) Set(v string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", v)
}

// Set indicates an expected call of Set
func (mr *MockStringMockRecorder) Set(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockString)(nil).Set), v)
}

// Attribute mocks base method
func (m *MockString) Attribute() output.Attribute {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attribute")
	ret0, _ := ret[0].(output.Attribute)
	return ret0
}

// Attribute indicates an expected call of Attribute
func (mr *MockStringMockRecorder) Attribute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attribute", reflect.TypeOf((*MockString)(nil).Attribute))
}
