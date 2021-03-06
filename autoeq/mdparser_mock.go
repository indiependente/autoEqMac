// Code generated by MockGen. DO NOT EDIT.
// Source: mdparser.go

// Package autoeq is a generated GoMock package.
package autoeq

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMarkDownParser is a mock of MarkDownParser interface
type MockMarkDownParser struct {
	ctrl     *gomock.Controller
	recorder *MockMarkDownParserMockRecorder
}

// MockMarkDownParserMockRecorder is the mock recorder for MockMarkDownParser
type MockMarkDownParserMockRecorder struct {
	mock *MockMarkDownParser
}

// NewMockMarkDownParser creates a new mock instance
func NewMockMarkDownParser(ctrl *gomock.Controller) *MockMarkDownParser {
	mock := &MockMarkDownParser{ctrl: ctrl}
	mock.recorder = &MockMarkDownParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMarkDownParser) EXPECT() *MockMarkDownParserMockRecorder {
	return m.recorder
}

// ParseMetadata mocks base method
func (m *MockMarkDownParser) ParseMetadata(arg0 []byte) ([]EQMetadata, error) {
	ret := m.ctrl.Call(m, "ParseMetadata", arg0)
	ret0, _ := ret[0].([]EQMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseMetadata indicates an expected call of ParseMetadata
func (mr *MockMarkDownParserMockRecorder) ParseMetadata(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseMetadata", reflect.TypeOf((*MockMarkDownParser)(nil).ParseMetadata), arg0)
}
