package connection

import (
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

type DBFactory func(driver string, dsn string) (util.DB, error)

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

// NewDefaultTCPConfig returns a Config with default settings for TCP connections.
func NewDefaultTCPConfig(user, password, database, driverName string) *Config {
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

// NewDefaultUnixConfig returns a Config with default settings for Unix socket connections.
func NewDefaultUnixConfig(user, password, database, socketDirectory, socketName, driverName string) *Config {
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

// Connect establishes a connection to the database using the provided configuration.
func Connect(cfg *Config, dbFactory DBFactory) (util.DB, error) {
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
		return nil, fmt.Errorf("unsupported connection type: %s", cfg.ConnectionType)
	}

	db, err := dbFactory(cfg.DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	configureConnection(db, cfg)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// configureConnection configures the connection pool parameters.
func configureConnection(db util.DB, cfg *Config) {
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
}
