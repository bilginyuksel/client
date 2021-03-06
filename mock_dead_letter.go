// Code generated by MockGen. DO NOT EDIT.
// Source: letter.go

// Package client is a generated GoMock package.
package client

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDeadLetter is a mock of DeadLetter interface.
type MockDeadLetter struct {
	ctrl     *gomock.Controller
	recorder *MockDeadLetterMockRecorder
}

// MockDeadLetterMockRecorder is the mock recorder for MockDeadLetter.
type MockDeadLetterMockRecorder struct {
	mock *MockDeadLetter
}

// NewMockDeadLetter creates a new mock instance.
func NewMockDeadLetter(ctrl *gomock.Controller) *MockDeadLetter {
	mock := &MockDeadLetter{ctrl: ctrl}
	mock.recorder = &MockDeadLetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeadLetter) EXPECT() *MockDeadLetterMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockDeadLetter) Save(letter *Letter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", letter)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockDeadLetterMockRecorder) Save(letter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockDeadLetter)(nil).Save), letter)
}
