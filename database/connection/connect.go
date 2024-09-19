package connection

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pakkasys/fluidapi/database/util"
)

// ConnectionType holds the type of the database connection.
type ConnectionType string

const (
	ConnectionTCP  ConnectionType = "tcp"  // TCP connection type
	ConnectionUnix ConnectionType = "unix" // Unix socket connection type
)

// DBInterface abstracts database operations.
type DBInterface interface {
	Ping() error
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	Prepare(query string) (util.Stmt, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (util.Tx, error)
	Exec(query string, args ...any) (util.Result, error)
	Query(query string, args ...any) (util.Rows, error)
	Close() error
}

// SQLDB wraps *sql.DB to implement DBInterface.
type SQLDB struct {
	*sql.DB
}

// NewDB creates a new database connection.
func (db *SQLDB) Ping() error {
	return db.DB.Ping()
}

// SetConnMaxLifetime sets the maximum time a connection may be reused.
func (db *SQLDB) SetConnMaxLifetime(d time.Duration) {
	db.DB.SetConnMaxLifetime(d)
}

// SetConnMaxIdleTime sets the maximum time an idle connection may be reused.
func (db *SQLDB) SetConnMaxIdleTime(d time.Duration) {
	db.DB.SetConnMaxIdleTime(d)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (db *SQLDB) SetMaxOpenConns(n int) {
	db.DB.SetMaxOpenConns(n)
}

// SetMaxIdleConns sets the maximum number of idle connections to the database.
func (db *SQLDB) SetMaxIdleConns(n int) {
	db.DB.SetMaxIdleConns(n)
}

// Prepare creates a prepared statement for later queries or executions.
func (db *SQLDB) Prepare(query string) (util.Stmt, error) {
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &util.RealStmt{Stmt: stmt}, nil
}

// BeginTx creates a transaction and returns it.
func (db *SQLDB) BeginTx(
	ctx context.Context,
	opts *sql.TxOptions,
) (util.Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &util.RealTx{Tx: tx}, nil
}

// Exec executes a query without returning rows.
func (db *SQLDB) Exec(query string, args ...any) (util.Result, error) {
	res, err := db.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &util.RealResult{Result: res}, nil
}

// Query executes a query that returns rows.
func (db *SQLDB) Query(query string, args ...any) (util.Rows, error) {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &util.RealRows{Rows: rows}, nil
}

// ConnectOptions holds the options for the database connection.
type ConnectOptions struct {
	DBFactory func(driver string, dsn string) (DBInterface, error)
}

// Config holds the configuration for the database connection.
type Config struct {
	User            string         // Database user
	Password        string         // Database password
	Host            string         // Database host
	Port            int            // Database port
	Database        string         // Database name
	SocketDirectory string         // Unix socket directory
	SocketName      string         // Unix socket name
	Parameters      string         // Connection parameters
	ConnectionType  ConnectionType // Connection type
	ConnMaxLifetime time.Duration  // Connection max lifetime
	ConnMaxIdleTime time.Duration  // Connection max idle time
	MaxOpenConns    int            // Max open connections
	MaxIdleConns    int            // Max idle connections
	DriverName      string         // Driver name
	DSNFormat       string         // Custom DSN format
}

// NewDefaultTCPConfig returns a Config with default settings for TCP
// connections.
//
//   - user: Database user
//   - password: Database password
//   - database: Database name
//   - driverName: Driver name
func NewDefaultTCPConfig(
	user string,
	password string,
	database string,
	driverName string,
) *Config {
	return &Config{
		User:            user,
		Password:        password,
		Database:        database,
		ConnectionType:  ConnectionTCP,
		Host:            "localhost",
		Port:            3306,
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		DriverName:      driverName,
		DSNFormat:       "%s:%s@tcp(%s:%d)/%s?parseTime=true&%s",
	}
}

// NewDefaultUnixConfig returns a Config with default settings for Unix socket
// connections.
//
//   - user: Database user
//   - password: Database password
//   - database: Database name
//   - socketDirectory: Unix socket directory
//   - socketName: Unix socket name
//   - driverName: Driver name
func NewDefaultUnixConfig(
	user string,
	password string,
	database string,
	socketDirectory string,
	socketName string,
	driverName string) *Config {
	return &Config{
		User:            user,
		Password:        password,
		Database:        database,
		ConnectionType:  ConnectionUnix,
		SocketDirectory: socketDirectory,
		SocketName:      socketName,
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		DriverName:      driverName,
		DSNFormat:       "%s:%s@unix(%s/%s)/%s?parseTime=true&%s",
	}
}

// Connect establishes a connection to the database using the provided
// configuration.
//
//   - cfg: Database configuration
//   - connectOpts: Optional ConnectOptions
func Connect(
	cfg *Config,
	connectOpts ...ConnectOptions,
) (DBInterface, error) {
	useConnectOpts := determineConnectOpts(connectOpts)

	var dsn string
	switch cfg.ConnectionType {
	case ConnectionTCP:
		dsn = fmt.Sprintf(
			cfg.DSNFormat,
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Parameters,
		)
	case ConnectionUnix:
		dsn = fmt.Sprintf(
			cfg.DSNFormat,
			cfg.User,
			cfg.Password,
			cfg.SocketDirectory,
			cfg.SocketName,
			cfg.Database,
			cfg.Parameters,
		)
	default:
		return nil, fmt.Errorf(
			"unsupported connection type: %s",
			cfg.ConnectionType,
		)
	}

	db, err := useConnectOpts.DBFactory(cfg.DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	configureConnection(db, cfg)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func newSQLDB(driver string, dsn string) (DBInterface, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &SQLDB{DB: db}, nil
}

func determineConnectOpts(connectOpts []ConnectOptions) ConnectOptions {
	if len(connectOpts) == 0 {
		return ConnectOptions{DBFactory: newSQLDB}
	} else {
		return connectOpts[0]
	}
}

func configureConnection(db DBInterface, cfg *Config) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}
