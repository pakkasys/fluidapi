package mock

import (
	"github.com/pakkasys/fluidapi/database/util"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the DB interface.
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Prepare(query string) (util.Stmt, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(util.Stmt), args.Error(1)
	}
}

// MockTx is a mock implementation of the Tx interface.
type MockTx struct {
	MockDB
}

func (m *MockTx) Commit() error {
	return m.Called().Error(0)
}

func (m *MockTx) Rollback() error {
	return m.Called().Error(0)
}

// MockStmt is a mock implementation of the Stmt interface.
type MockStmt struct {
	mock.Mock
}

func (m *MockStmt) Close() error {
	return m.Called().Error(0)
}

func (m *MockStmt) QueryRow(args ...any) util.Row {
	argsCalled := m.Called(args)
	return argsCalled.Get(0).(util.Row)
}

func (m *MockStmt) Exec(args ...any) (util.Result, error) {
	argsCalled := m.Called(args)
	if argsCalled.Get(0) == nil {
		return nil, argsCalled.Error(1)
	}
	return argsCalled.Get(0).(util.Result), argsCalled.Error(1)
}

func (m *MockStmt) Query(args ...any) (util.Rows, error) {
	argsCalled := m.Called(args)
	if argsCalled.Get(0) == nil {
		return nil, argsCalled.Error(1)
	}
	return argsCalled.Get(0).(util.Rows), argsCalled.Error(1)
}

// MockRows is a mock implementation of the Rows interface.
type MockRows struct {
	mock.Mock
}

func (m *MockRows) Scan(dest ...any) error {
	return m.Called(dest).Error(0)
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Close() error {
	return m.Called().Error(0)
}

func (m *MockRows) Err() error {
	return m.Called().Error(0)
}

// MockRow is a mock implementation of the Row interface.
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	return m.Called(dest).Error(0)
}

func (m *MockRow) Err() error {
	return m.Called().Error(0)
}

// MockResult is a mock implementation of the Result interface.
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
