// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk (interfaces: ClientWithResponsesInterface)

// Package prmocks is a generated GoMock package.
package prmocks

import (
	context "context"
	reflect "reflect"

	providerregistrysdk "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	gomock "github.com/golang/mock/gomock"
)

// MockClientWithResponsesInterface is a mock of ClientWithResponsesInterface interface.
type MockClientWithResponsesInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientWithResponsesInterfaceMockRecorder
}

// MockClientWithResponsesInterfaceMockRecorder is the mock recorder for MockClientWithResponsesInterface.
type MockClientWithResponsesInterfaceMockRecorder struct {
	mock *MockClientWithResponsesInterface
}

// NewMockClientWithResponsesInterface creates a new mock instance.
func NewMockClientWithResponsesInterface(ctrl *gomock.Controller) *MockClientWithResponsesInterface {
	mock := &MockClientWithResponsesInterface{ctrl: ctrl}
	mock.recorder = &MockClientWithResponsesInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientWithResponsesInterface) EXPECT() *MockClientWithResponsesInterfaceMockRecorder {
	return m.recorder
}

// GetProviderSetupDocsWithResponse mocks base method.
func (m *MockClientWithResponsesInterface) GetProviderSetupDocsWithResponse(arg0 context.Context, arg1, arg2, arg3 string, arg4 ...providerregistrysdk.RequestEditorFn) (*providerregistrysdk.GetProviderSetupDocsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProviderSetupDocsWithResponse", varargs...)
	ret0, _ := ret[0].(*providerregistrysdk.GetProviderSetupDocsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderSetupDocsWithResponse indicates an expected call of GetProviderSetupDocsWithResponse.
func (mr *MockClientWithResponsesInterfaceMockRecorder) GetProviderSetupDocsWithResponse(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderSetupDocsWithResponse", reflect.TypeOf((*MockClientWithResponsesInterface)(nil).GetProviderSetupDocsWithResponse), varargs...)
}

// GetProviderUsageDocWithResponse mocks base method.
func (m *MockClientWithResponsesInterface) GetProviderUsageDocWithResponse(arg0 context.Context, arg1, arg2, arg3 string, arg4 ...providerregistrysdk.RequestEditorFn) (*providerregistrysdk.GetProviderUsageDocResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProviderUsageDocWithResponse", varargs...)
	ret0, _ := ret[0].(*providerregistrysdk.GetProviderUsageDocResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderUsageDocWithResponse indicates an expected call of GetProviderUsageDocWithResponse.
func (mr *MockClientWithResponsesInterfaceMockRecorder) GetProviderUsageDocWithResponse(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderUsageDocWithResponse", reflect.TypeOf((*MockClientWithResponsesInterface)(nil).GetProviderUsageDocWithResponse), varargs...)
}

// GetProviderWithResponse mocks base method.
func (m *MockClientWithResponsesInterface) GetProviderWithResponse(arg0 context.Context, arg1, arg2, arg3 string, arg4 ...providerregistrysdk.RequestEditorFn) (*providerregistrysdk.GetProviderResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProviderWithResponse", varargs...)
	ret0, _ := ret[0].(*providerregistrysdk.GetProviderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderWithResponse indicates an expected call of GetProviderWithResponse.
func (mr *MockClientWithResponsesInterfaceMockRecorder) GetProviderWithResponse(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderWithResponse", reflect.TypeOf((*MockClientWithResponsesInterface)(nil).GetProviderWithResponse), varargs...)
}

// ListAllProvidersWithResponse mocks base method.
func (m *MockClientWithResponsesInterface) ListAllProvidersWithResponse(arg0 context.Context, arg1 ...providerregistrysdk.RequestEditorFn) (*providerregistrysdk.ListAllProvidersResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListAllProvidersWithResponse", varargs...)
	ret0, _ := ret[0].(*providerregistrysdk.ListAllProvidersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllProvidersWithResponse indicates an expected call of ListAllProvidersWithResponse.
func (mr *MockClientWithResponsesInterfaceMockRecorder) ListAllProvidersWithResponse(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllProvidersWithResponse", reflect.TypeOf((*MockClientWithResponsesInterface)(nil).ListAllProvidersWithResponse), varargs...)
}
