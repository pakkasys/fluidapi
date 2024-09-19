package connection

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/pakkasys/fluidapi/database/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the DBInterface.
type MockDB struct{ mock.Mock }

func (m *MockDB) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) SetConnMaxLifetime(d time.Duration) { m.Called(d) }
func (m *MockDB) SetConnMaxIdleTime(d time.Duration) { m.Called(d) }
func (m *MockDB) SetMaxOpenConns(n int)              { m.Called(n) }
func (m *MockDB) SetMaxIdleConns(n int)              { m.Called(n) }

func (m *MockDB) Prepare(query string) (util.Stmt, error) {
	args := m.Called(query)
	return args.Get(0).(util.Stmt), args.Error(1)
}

func (m *MockDB) BeginTx(
	ctx context.Context,
	opts *sql.TxOptions,
) (util.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(util.Tx), args.Error(1)
}

func (m *MockDB) Exec(query string, args ...any) (util.Result, error) {
	calledArgs := m.Called(query, args)
	return nil, calledArgs.Get(1).(error)
}

func (m *MockDB) Query(query string, args ...any) (util.Rows, error) {
	calledArgs := m.Called(query, args)
	return nil, calledArgs.Get(1).(error)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockDriver is a stub driver that satisfies the driver.Driver interface.
type MockDriver struct{}

func (m *MockDriver) Open(name string) (driver.Conn, error) {
	return &MockConn{}, nil
}

// MockConn is a stub connection that satisfies the driver.Conn interface.
type MockConn struct{}

func (c *MockConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (c *MockConn) Close() error {
	return nil
}

func (c *MockConn) Begin() (driver.Tx, error) {
	return nil, nil
}

// MockDBInterface is a mock implementation of the DBInterface.
type MockDBInterface struct {
	mock.Mock
}

func (m *MockDBInterface) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDBInterface) SetConnMaxLifetime(d time.Duration) {
	m.Called(d)
}

func (m *MockDBInterface) SetConnMaxIdleTime(d time.Duration) {
	m.Called(d)
}

func (m *MockDBInterface) SetMaxOpenConns(n int) {
	m.Called(n)
}

func (m *MockDBInterface) SetMaxIdleConns(n int) {
	m.Called(n)
}

func (m *MockDBInterface) Prepare(query string) (*sql.Stmt, error) {
	args := m.Called(query)
	return nil, args.Error(1)
}

func (m *MockDBInterface) BeginTx(
	ctx context.Context,
	opts *sql.TxOptions,
) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

// TestNewDefaultTCPConfig tests default TCP configuration creation.
func TestNewDefaultTCPConfig(t *testing.T) {
	user := "testuser"
	password := "testpass"
	database := "testdb"
	driverName := "mysql"

	cfg := NewDefaultTCPConfig(user, password, database, driverName)
	assert.NotNil(t, cfg)
	assert.Equal(t, user, cfg.User)
	assert.Equal(t, password, cfg.Password)
	assert.Equal(t, database, cfg.Database)
	assert.Equal(t, driverName, cfg.DriverName)
	assert.Equal(t, ConnectionTCP, cfg.ConnectionType)
}

// TestNewDefaultUnixConfig tests default Unix socket configuration creation.
func TestNewDefaultUnixConfig(t *testing.T) {
	user := "testuser"
	password := "testpass"
	database := "testdb"
	socketDirectory := "/var/run/mysqld"
	socketName := "mysqld.sock"
	driverName := "mysql"

	cfg := NewDefaultUnixConfig(
		user,
		password,
		database,
		socketDirectory,
		socketName,
		driverName,
	)
	assert.NotNil(t, cfg)
	assert.Equal(t, user, cfg.User)
	assert.Equal(t, password, cfg.Password)
	assert.Equal(t, database, cfg.Database)
	assert.Equal(t, socketDirectory, cfg.SocketDirectory)
	assert.Equal(t, socketName, cfg.SocketName)
	assert.Equal(t, driverName, cfg.DriverName)
	assert.Equal(t, ConnectionUnix, cfg.ConnectionType)
}

// TestConnect_Success tests the successful connection scenario.
func TestConnect_Success(t *testing.T) {
	cfg := NewDefaultTCPConfig("user", "password", "database", "mysql")

	mockDB := new(MockDB)
	mockDB.On("Ping").Return(nil)
	mockDB.On("SetConnMaxLifetime", cfg.ConnMaxLifetime).Return()
	mockDB.On("SetConnMaxIdleTime", cfg.ConnMaxIdleTime).Return()
	mockDB.On("SetMaxOpenConns", cfg.MaxOpenConns).Return()
	mockDB.On("SetMaxIdleConns", cfg.MaxIdleConns).Return()

	// Mock factory function to return the mockDB
	dbFactory := func(driver string, dsn string) (DBInterface, error) {
		return mockDB, nil
	}

	db, err := Connect(cfg, ConnectOptions{DBFactory: dbFactory})

	assert.NoError(t, err)
	assert.NotNil(t, db)
	mockDB.AssertExpectations(t)
}

// TestConnect_FailedOpen tests the scenario where opening the connection fails.
func TestConnect_FailedOpen(t *testing.T) {
	cfg := NewDefaultTCPConfig("user", "password", "database", "mysql")

	// Mock factory function to return an error
	dbFactory := func(driver string, dsn string) (DBInterface, error) {
		return nil, errors.New("failed to open database")
	}

	db, err := Connect(cfg, ConnectOptions{DBFactory: dbFactory})

	assert.Error(t, err)
	assert.Nil(t, db)
}

// TestConnect_FailedPing tests the scenario where pinging the database fails.
func TestConnect_FailedPing(t *testing.T) {
	cfg := NewDefaultTCPConfig("user", "password", "database", "mysql")

	mockDB := new(MockDB)
	mockDB.On("Ping").Return(errors.New("failed to ping database"))
	mockDB.On("SetConnMaxLifetime", cfg.ConnMaxLifetime).Return()
	mockDB.On("SetConnMaxIdleTime", cfg.ConnMaxIdleTime).Return()
	mockDB.On("SetMaxOpenConns", cfg.MaxOpenConns).Return()
	mockDB.On("SetMaxIdleConns", cfg.MaxIdleConns).Return()

	dbFactory := func(driver string, dsn string) (DBInterface, error) {
		return mockDB, nil
	}

	db, err := Connect(cfg, ConnectOptions{DBFactory: dbFactory})

	assert.Error(t, err)
	assert.Nil(t, db)
}

// TestConnect_Unix_Success tests the Unix socket connection case.
func TestConnect_Unix_Success(t *testing.T) {
	cfg := NewDefaultUnixConfig(
		"user",
		"password",
		"database",
		"/var/run/mysqld",
		"mysqld.sock",
		"mysql",
	)

	mockDB := new(MockDB)
	mockDB.On("Ping").Return(nil)
	mockDB.On("SetConnMaxLifetime", cfg.ConnMaxLifetime).Return()
	mockDB.On("SetConnMaxIdleTime", cfg.ConnMaxIdleTime).Return()
	mockDB.On("SetMaxOpenConns", cfg.MaxOpenConns).Return()
	mockDB.On("SetMaxIdleConns", cfg.MaxIdleConns).Return()

	dbFactory := func(driver string, dsn string) (DBInterface, error) {
		return mockDB, nil
	}

	db, err := Connect(cfg, ConnectOptions{DBFactory: dbFactory})

	assert.NoError(t, err)
	assert.NotNil(t, db)
	mockDB.AssertExpectations(t)
}

// TestConnect_UnsupportedConnectionType unsupported connection type case.
func TestConnect_UnsupportedConnectionType(t *testing.T) {
	cfg := &Config{
		User:           "user",
		Password:       "password",
		Database:       "database",
		ConnectionType: "unsupported",
		DriverName:     "mysql",
	}

	dbFactory := func(driver string, dsn string) (DBInterface, error) {
		assert.Fail(t, "should not be called")
		return nil, nil
	}

	db, err := Connect(cfg, ConnectOptions{DBFactory: dbFactory})

	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Equal(t, "unsupported connection type: unsupported", err.Error())
}

// TestDetermineConnectOpts_Empty tests the case without connection options.
func TestDetermineConnectOpts_Empty(t *testing.T) {
	opts := determineConnectOpts(nil)

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.DBFactory)

	_, err := opts.DBFactory("nonexistent", "user:password@/dbname")
	assert.EqualError(
		t,
		err,
		"sql: unknown driver \"nonexistent\" (forgotten import?)",
	)
}

// TestDetermineConnectOpts_WithOptions tests the case with connection options.
func TestDetermineConnectOpts_WithOptions(t *testing.T) {
	mockFactoryCalled := false
	mockFactory := func(driver string, dsn string) (DBInterface, error) {
		mockFactoryCalled = true
		return nil, nil
	}
	inputOpts := ConnectOptions{DBFactory: mockFactory}

	opts := determineConnectOpts([]ConnectOptions{inputOpts})

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.DBFactory)

	_, err := opts.DBFactory("mysql", "user:password@/dbname")
	assert.NoError(t, err)
	assert.True(t, mockFactoryCalled, "Expected mock factory to be called")
}

// TestNewSQLDB_Success tests the scenario where newSQLDB succeeds.
func TestNewSQLDB_Success(t *testing.T) {
	sql.Register("mockDriver", &MockDriver{})
	db, err := newSQLDB("mockDriver", "user:password@/dbname")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

// TestNewSQLDB_FailedOpen tests the scenario where newSQLDB fails.
func TestNewSQLDB_FailedOpen(t *testing.T) {
	// Use an invalid driver name to force an error
	db, err := newSQLDB("invalid_driver", "user:password@/dbname")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, db)
}

// // TestConfigureConnection tests the configureConnection function.
// func TestConfigureConnection(t *testing.T) {
// 	mockDB := new(MockDBInterface)
// 	cfg := &Config{
// 		ConnMaxLifetime: 15 * time.Minute,
// 		ConnMaxIdleTime: 10 * time.Minute,
// 		MaxOpenConns:    30,
// 		MaxIdleConns:    10,
// 	}

// 	mockDB.On("SetConnMaxLifetime", cfg.ConnMaxLifetime).Return()
// 	mockDB.On("SetConnMaxIdleTime", cfg.ConnMaxIdleTime).Return()
// 	mockDB.On("SetMaxOpenConns", cfg.MaxOpenConns).Return()
// 	mockDB.On("SetMaxIdleConns", cfg.MaxIdleConns).Return()

// 	configureConnection(mockDB, cfg)
// 	mockDB.AssertExpectations(t)
// }
