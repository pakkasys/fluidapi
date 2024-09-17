package connection

import (
	"database/sql"
	"fmt"
	"time"
)

// ConnectionType holds the type of the database connection.
type ConnectionType string

const (
	ConnectionTCP  ConnectionType = "tcp"  // TCP connection type
	ConnectionUnix ConnectionType = "unix" // Unix socket connection type
)

// DatabaseConfigurator defines methods needed to configure a database
// connection.
type DatabaseConfigurator interface {
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
}

// SQLDB wraps *sql.DB to implement the DatabaseConfigurator interface.
type SQLDB struct {
	*sql.DB
}

// SetConnMaxLifetime sets the maximum time a connection may be reused.
func (db *SQLDB) SetConnMaxLifetime(d time.Duration) {
	db.DB.SetConnMaxLifetime(d)
}

// SetConnMaxIdleTime sets the maximum time a connection may remain idle.
func (db *SQLDB) SetConnMaxIdleTime(d time.Duration) {
	db.DB.SetConnMaxIdleTime(d)
}

// SetMaxOpenConns sets the maximum number of open connections.
func (db *SQLDB) SetMaxOpenConns(n int) {
	db.DB.SetMaxOpenConns(n)
}

// SetMaxIdleConns sets the maximum number of idle connections.
func (db *SQLDB) SetMaxIdleConns(n int) {
	db.DB.SetMaxIdleConns(n)
}

// SQLConnectorFunc is a function type that opens a new database connection.
type SQLConnectorFunc func(driverName, dataSourceName string) (*sql.DB, error)

// DefaultSQLConnector is the default function used to connect to the database.
var DefaultSQLConnector SQLConnectorFunc = sql.Open

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
//   - user: database user
//   - password: database password
//   - database: database name
//   - driverName: database driver name
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
//   - user: database user
//   - password: database password
//   - database: database name
//   - socketDirectory: Unix socket directory
//   - socketName: Unix socket name
//   - driverName: database driver name
func NewDefaultUnixConfig(
	user string,
	password string,
	database string,
	socketDirectory string,
	socketName string,
	driverName string,
) *Config {
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
//   - cfg: database connection configuration
func Connect(cfg *Config) (*sql.DB, error) {
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

	db, err := sql.Open(cfg.DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	configureConnection(db, cfg)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func configureConnection(db DatabaseConfigurator, cfg *Config) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}
