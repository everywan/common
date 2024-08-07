// Code generated by MockGen. DO NOT EDIT.
// Source: source.go

// Package mock_configs is a generated GoMock package.
package mock_configs

import (
	reflect "reflect"

	configs "github.com/everywan/common/configs"
	gomock "github.com/golang/mock/gomock"
)

// MockISource is a mock of ISource interface.
type MockISource struct {
	ctrl     *gomock.Controller
	recorder *MockISourceMockRecorder
}

// MockISourceMockRecorder is the mock recorder for MockISource.
type MockISourceMockRecorder struct {
	mock *MockISource
}

// NewMockISource creates a new mock instance.
func NewMockISource(ctrl *gomock.Controller) *MockISource {
	mock := &MockISource{ctrl: ctrl}
	mock.recorder = &MockISourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISource) EXPECT() *MockISourceMockRecorder {
	return m.recorder
}

// BatchGet mocks base method.
func (m *MockISource) BatchGet(group string, keys ...string) (map[string]string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{group}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BatchGet", varargs...)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchGet indicates an expected call of BatchGet.
func (mr *MockISourceMockRecorder) BatchGet(group interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{group}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchGet", reflect.TypeOf((*MockISource)(nil).BatchGet), varargs...)
}

// DefaultGroup mocks base method.
func (m *MockISource) DefaultGroup() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultGroup")
	ret0, _ := ret[0].(string)
	return ret0
}

// DefaultGroup indicates an expected call of DefaultGroup.
func (mr *MockISourceMockRecorder) DefaultGroup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultGroup", reflect.TypeOf((*MockISource)(nil).DefaultGroup))
}

// Get mocks base method.
func (m *MockISource) Get(group, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", group, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockISourceMockRecorder) Get(group, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockISource)(nil).Get), group, key)
}

// MonitorChange mocks base method.
func (m *MockISource) MonitorChange(group, key string, fn configs.ConfigChangeCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MonitorChange", group, key, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// MonitorChange indicates an expected call of MonitorChange.
func (mr *MockISourceMockRecorder) MonitorChange(group, key, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MonitorChange", reflect.TypeOf((*MockISource)(nil).MonitorChange), group, key, fn)
}

// Name mocks base method.
func (m *MockISource) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockISourceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockISource)(nil).Name))
}
