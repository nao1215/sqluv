// Code generated by MockGen. DO NOT EDIT.
// Source: persistence.go
//
// Generated by this command:
//
//	mockgen -typed -source=persistence.go -destination=../../infrastructure/mock/persistence.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/nao1215/sqluv/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockQueryToRemoteExecutor is a mock of QueryToRemoteExecutor interface.
type MockQueryToRemoteExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockQueryToRemoteExecutorMockRecorder
	isgomock struct{}
}

// MockQueryToRemoteExecutorMockRecorder is the mock recorder for MockQueryToRemoteExecutor.
type MockQueryToRemoteExecutorMockRecorder struct {
	mock *MockQueryToRemoteExecutor
}

// NewMockQueryToRemoteExecutor creates a new mock instance.
func NewMockQueryToRemoteExecutor(ctrl *gomock.Controller) *MockQueryToRemoteExecutor {
	mock := &MockQueryToRemoteExecutor{ctrl: ctrl}
	mock.recorder = &MockQueryToRemoteExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryToRemoteExecutor) EXPECT() *MockQueryToRemoteExecutorMockRecorder {
	return m.recorder
}

// ExecuteQuery mocks base method.
func (m *MockQueryToRemoteExecutor) ExecuteQuery(ctx context.Context, sql *model.SQL) (*model.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteQuery", ctx, sql)
	ret0, _ := ret[0].(*model.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteQuery indicates an expected call of ExecuteQuery.
func (mr *MockQueryToRemoteExecutorMockRecorder) ExecuteQuery(ctx, sql any) *MockQueryToRemoteExecutorExecuteQueryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteQuery", reflect.TypeOf((*MockQueryToRemoteExecutor)(nil).ExecuteQuery), ctx, sql)
	return &MockQueryToRemoteExecutorExecuteQueryCall{Call: call}
}

// MockQueryToRemoteExecutorExecuteQueryCall wrap *gomock.Call
type MockQueryToRemoteExecutorExecuteQueryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockQueryToRemoteExecutorExecuteQueryCall) Return(arg0 *model.Table, arg1 error) *MockQueryToRemoteExecutorExecuteQueryCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockQueryToRemoteExecutorExecuteQueryCall) Do(f func(context.Context, *model.SQL) (*model.Table, error)) *MockQueryToRemoteExecutorExecuteQueryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockQueryToRemoteExecutorExecuteQueryCall) DoAndReturn(f func(context.Context, *model.SQL) (*model.Table, error)) *MockQueryToRemoteExecutorExecuteQueryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockStatementToRemoteExecutor is a mock of StatementToRemoteExecutor interface.
type MockStatementToRemoteExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockStatementToRemoteExecutorMockRecorder
	isgomock struct{}
}

// MockStatementToRemoteExecutorMockRecorder is the mock recorder for MockStatementToRemoteExecutor.
type MockStatementToRemoteExecutorMockRecorder struct {
	mock *MockStatementToRemoteExecutor
}

// NewMockStatementToRemoteExecutor creates a new mock instance.
func NewMockStatementToRemoteExecutor(ctrl *gomock.Controller) *MockStatementToRemoteExecutor {
	mock := &MockStatementToRemoteExecutor{ctrl: ctrl}
	mock.recorder = &MockStatementToRemoteExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatementToRemoteExecutor) EXPECT() *MockStatementToRemoteExecutorMockRecorder {
	return m.recorder
}

// ExecuteStatement mocks base method.
func (m *MockStatementToRemoteExecutor) ExecuteStatement(ctx context.Context, sql *model.SQL) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteStatement", ctx, sql)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteStatement indicates an expected call of ExecuteStatement.
func (mr *MockStatementToRemoteExecutorMockRecorder) ExecuteStatement(ctx, sql any) *MockStatementToRemoteExecutorExecuteStatementCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteStatement", reflect.TypeOf((*MockStatementToRemoteExecutor)(nil).ExecuteStatement), ctx, sql)
	return &MockStatementToRemoteExecutorExecuteStatementCall{Call: call}
}

// MockStatementToRemoteExecutorExecuteStatementCall wrap *gomock.Call
type MockStatementToRemoteExecutorExecuteStatementCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockStatementToRemoteExecutorExecuteStatementCall) Return(arg0 int64, arg1 error) *MockStatementToRemoteExecutorExecuteStatementCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockStatementToRemoteExecutorExecuteStatementCall) Do(f func(context.Context, *model.SQL) (int64, error)) *MockStatementToRemoteExecutorExecuteStatementCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockStatementToRemoteExecutorExecuteStatementCall) DoAndReturn(f func(context.Context, *model.SQL) (int64, error)) *MockStatementToRemoteExecutorExecuteStatementCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockTablesInRemoteGetter is a mock of TablesInRemoteGetter interface.
type MockTablesInRemoteGetter struct {
	ctrl     *gomock.Controller
	recorder *MockTablesInRemoteGetterMockRecorder
	isgomock struct{}
}

// MockTablesInRemoteGetterMockRecorder is the mock recorder for MockTablesInRemoteGetter.
type MockTablesInRemoteGetterMockRecorder struct {
	mock *MockTablesInRemoteGetter
}

// NewMockTablesInRemoteGetter creates a new mock instance.
func NewMockTablesInRemoteGetter(ctrl *gomock.Controller) *MockTablesInRemoteGetter {
	mock := &MockTablesInRemoteGetter{ctrl: ctrl}
	mock.recorder = &MockTablesInRemoteGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTablesInRemoteGetter) EXPECT() *MockTablesInRemoteGetterMockRecorder {
	return m.recorder
}

// GetTables mocks base method.
func (m *MockTablesInRemoteGetter) GetTables(ctx context.Context) ([]*model.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTables", ctx)
	ret0, _ := ret[0].([]*model.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTables indicates an expected call of GetTables.
func (mr *MockTablesInRemoteGetterMockRecorder) GetTables(ctx any) *MockTablesInRemoteGetterGetTablesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTables", reflect.TypeOf((*MockTablesInRemoteGetter)(nil).GetTables), ctx)
	return &MockTablesInRemoteGetterGetTablesCall{Call: call}
}

// MockTablesInRemoteGetterGetTablesCall wrap *gomock.Call
type MockTablesInRemoteGetterGetTablesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTablesInRemoteGetterGetTablesCall) Return(arg0 []*model.Table, arg1 error) *MockTablesInRemoteGetterGetTablesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTablesInRemoteGetterGetTablesCall) Do(f func(context.Context) ([]*model.Table, error)) *MockTablesInRemoteGetterGetTablesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTablesInRemoteGetterGetTablesCall) DoAndReturn(f func(context.Context) ([]*model.Table, error)) *MockTablesInRemoteGetterGetTablesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockTableDDLInRemoteGetter is a mock of TableDDLInRemoteGetter interface.
type MockTableDDLInRemoteGetter struct {
	ctrl     *gomock.Controller
	recorder *MockTableDDLInRemoteGetterMockRecorder
	isgomock struct{}
}

// MockTableDDLInRemoteGetterMockRecorder is the mock recorder for MockTableDDLInRemoteGetter.
type MockTableDDLInRemoteGetterMockRecorder struct {
	mock *MockTableDDLInRemoteGetter
}

// NewMockTableDDLInRemoteGetter creates a new mock instance.
func NewMockTableDDLInRemoteGetter(ctrl *gomock.Controller) *MockTableDDLInRemoteGetter {
	mock := &MockTableDDLInRemoteGetter{ctrl: ctrl}
	mock.recorder = &MockTableDDLInRemoteGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTableDDLInRemoteGetter) EXPECT() *MockTableDDLInRemoteGetterMockRecorder {
	return m.recorder
}

// GetTableDDL mocks base method.
func (m *MockTableDDLInRemoteGetter) GetTableDDL(ctx context.Context, tableName string) ([]*model.Table, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTableDDL", ctx, tableName)
	ret0, _ := ret[0].([]*model.Table)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTableDDL indicates an expected call of GetTableDDL.
func (mr *MockTableDDLInRemoteGetterMockRecorder) GetTableDDL(ctx, tableName any) *MockTableDDLInRemoteGetterGetTableDDLCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTableDDL", reflect.TypeOf((*MockTableDDLInRemoteGetter)(nil).GetTableDDL), ctx, tableName)
	return &MockTableDDLInRemoteGetterGetTableDDLCall{Call: call}
}

// MockTableDDLInRemoteGetterGetTableDDLCall wrap *gomock.Call
type MockTableDDLInRemoteGetterGetTableDDLCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableDDLInRemoteGetterGetTableDDLCall) Return(arg0 []*model.Table, arg1 error) *MockTableDDLInRemoteGetterGetTableDDLCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableDDLInRemoteGetterGetTableDDLCall) Do(f func(context.Context, string) ([]*model.Table, error)) *MockTableDDLInRemoteGetterGetTableDDLCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableDDLInRemoteGetterGetTableDDLCall) DoAndReturn(f func(context.Context, string) ([]*model.Table, error)) *MockTableDDLInRemoteGetterGetTableDDLCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
