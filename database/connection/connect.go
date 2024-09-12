package connection

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	socketTCP  = "tcp"
	socketUnix = "unix"

	driverNameMySQL   = "mysql"
	dataSourceNameFmt = "%s:%s@%s/%s?parseTime=true&%s"
)

func ConnectTCP(
	user string,
	password string,
	host string,
	port int,
	database string,
	parameters *string,
) (*sql.DB, error) {
	return connect(
		user,
		password,
		fmt.Sprintf("%s(%s:%d)", socketTCP, host, port),
		database,
		parameters,
	)
}

func ConnectUnix(
	user string,
	password string,
	socketDirectory string,
	socketName string,
	database string,
	parameters *string,
) (*sql.DB, error) {
	return connect(
		user,
		password,
		fmt.Sprintf("%s(%s/%s)", socketUnix, socketDirectory, socketName),
		database,
		parameters,
	)
}

func connect(
	user string,
	password string,
	connectionString string,
	database string,
	parameters *string,
) (*sql.DB, error) {
	databaseHandle, err := sql.Open(
		driverNameMySQL,
		fmt.Sprintf(
			dataSourceNameFmt,
			user,
			password,
			connectionString,
			database,
			determineParameters(parameters),
		),
	)
	if err != nil {
		return nil, err
	}

	return configureConnection(databaseHandle), nil
}

func determineParameters(parameters *string) string {
	if parameters == nil {
		return ""
	} else {
		return *parameters
	}
}

func configureConnection(databaseHandle *sql.DB) *sql.DB {
	databaseHandle.SetConnMaxLifetime(time.Minute * 10)
	databaseHandle.SetConnMaxIdleTime(time.Minute * 5)
	databaseHandle.SetMaxOpenConns(25)
	databaseHandle.SetMaxIdleConns(5)
	return databaseHandle
}
