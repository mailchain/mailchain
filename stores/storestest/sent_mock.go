// Code generated by MockGen. DO NOT EDIT.
// Source: sent.go

// Package storestest is a generated GoMock package.
package storestest

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mail "github.com/mailchain/mailchain/internal/mail"
)

// MockSent is a mock of Sent interface.
type MockSent struct {
	ctrl     *gomock.Controller
	recorder *MockSentMockRecorder
}

// MockSentMockRecorder is the mock recorder for MockSent.
type MockSentMockRecorder struct {
	mock *MockSent
}

// NewMockSent creates a new mock instance.
func NewMockSent(ctrl *gomock.Controller) *MockSent {
	mock := &MockSent{ctrl: ctrl}
	mock.recorder = &MockSentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSent) EXPECT() *MockSentMockRecorder {
	return m.recorder
}

// Key mocks base method.
func (m *MockSent) Key(messageID mail.ID, contentsHash, msg []byte) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Key", messageID, contentsHash, msg)
	ret0, _ := ret[0].(string)
	return ret0
}

// Key indicates an expected call of Key.
func (mr *MockSentMockRecorder) Key(messageID, contentsHash, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Key", reflect.TypeOf((*MockSent)(nil).Key), messageID, contentsHash, msg)
}

// PutMessage mocks base method.
func (m *MockSent) PutMessage(messageID mail.ID, contentsHash, msg []byte, headers map[string]string) (string, string, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMessage", messageID, contentsHash, msg, headers)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// PutMessage indicates an expected call of PutMessage.
func (mr *MockSentMockRecorder) PutMessage(messageID, contentsHash, msg, headers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMessage", reflect.TypeOf((*MockSent)(nil).PutMessage), messageID, contentsHash, msg, headers)
}
