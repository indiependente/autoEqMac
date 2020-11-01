// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package server is a generated GoMock package.
package server

import (
	gomock "github.com/golang/mock/gomock"
	autoeq "github.com/indiependente/autoEqMac/autoeq"
	eqmac "github.com/indiependente/autoEqMac/eqmac"
	io "io"
	reflect "reflect"
)

// MockServer is a mock of Server interface
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// ListEQsMetadata mocks base method
func (m *MockServer) ListEQsMetadata() ([]autoeq.EQMetadata, error) {
	ret := m.ctrl.Call(m, "ListEQsMetadata")
	ret0, _ := ret[0].([]autoeq.EQMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEQsMetadata indicates an expected call of ListEQsMetadata
func (mr *MockServerMockRecorder) ListEQsMetadata() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEQsMetadata", reflect.TypeOf((*MockServer)(nil).ListEQsMetadata))
}

// GetFixedBandEQPreset mocks base method
func (m *MockServer) GetFixedBandEQPreset(id string) (eqmac.EQPreset, error) {
	ret := m.ctrl.Call(m, "GetFixedBandEQPreset", id)
	ret0, _ := ret[0].(eqmac.EQPreset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFixedBandEQPreset indicates an expected call of GetFixedBandEQPreset
func (mr *MockServerMockRecorder) GetFixedBandEQPreset(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFixedBandEQPreset", reflect.TypeOf((*MockServer)(nil).GetFixedBandEQPreset), id)
}

// GetEQMetadataByName mocks base method
func (m *MockServer) GetEQMetadataByName(name string) (autoeq.EQMetadata, error) {
	ret := m.ctrl.Call(m, "GetEQMetadataByName", name)
	ret0, _ := ret[0].(autoeq.EQMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEQMetadataByName indicates an expected call of GetEQMetadataByName
func (mr *MockServerMockRecorder) GetEQMetadataByName(name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEQMetadataByName", reflect.TypeOf((*MockServer)(nil).GetEQMetadataByName), name)
}

// WritePreset mocks base method
func (m *MockServer) WritePreset(w io.Writer, p eqmac.EQPreset) error {
	ret := m.ctrl.Call(m, "WritePreset", w, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// WritePreset indicates an expected call of WritePreset
func (mr *MockServerMockRecorder) WritePreset(w, p interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WritePreset", reflect.TypeOf((*MockServer)(nil).WritePreset), w, p)
}