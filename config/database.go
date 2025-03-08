package config

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"modernc.org/sqlite"
)

// MemoryDB is *sql.DB for excuting sql.
type MemoryDB *sql.DB

// NewMemoryDB create *sql.DB for SQLite3. SQLite3 store data in memory.
// The return function is the function to close the DB.
func NewMemoryDB() (MemoryDB, func(), error) {
	initSQLite3()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, nil, err
	}
	return MemoryDB(db), func() { db.Close() }, nil
}

// once is an object to ensure that the sqlite3 driver is registered only once.
var once sync.Once

// initSQLite3 registers the sqlite3 driver.
func initSQLite3() {
	once.Do(func() {
		sql.Register("sqlite3", sqliteDriver{Driver: &sqlite.Driver{}})
	})
}

// sqliteDriver is a driver that enables foreign keys.
type sqliteDriver struct {
	*sqlite.Driver
}

// Open opens a database specified by its database driver name and a driver-specific data source name.
func (d sqliteDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c, ok := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if !ok {
		return nil, errors.New("connection does not support Exec method")
	}

	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		if err := conn.Close(); err != nil {
			return nil, fmt.Errorf("failed to close connection: %w", err)
		}
		return nil, fmt.Errorf("failed to enable enable foreign keys: %w", err)
	}
	return conn, nil
}

// DBMS is a common interface for database connections
type DBMS *sql.DB

// MySQLDB is *sql.DB for executing SQL with MySQL.
type MySQLDB = DBMS

// PostgreSQLDB is *sql.DB for executing SQL with PostgreSQL
type PostgreSQLDB = DBMS

// MySQLConfig holds the configuration for MySQL connection.
type MySQLConfig struct {
	host     string
	port     int
	user     string
	password string
	database string
}

// NewMySQLConfig creates MySQLConfig.
func NewMySQLConfig(
	host string,
	port int,
	user string,
	password string,
	database string,
) MySQLConfig {
	return MySQLConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
	}
}

// NewMySQLDB creates *sql.DB for MySQL.
// The return function is the function to close the DB.
func NewMySQLDB(config MySQLConfig) (MySQLDB, func(), error) {
	c := mysql.Config{
		DBName:    config.database,
		User:      config.user,
		Passwd:    config.password,
		Addr:      fmt.Sprintf("%s:%d", config.host, config.port),
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}
	return MySQLDB(db), func() { db.Close() }, nil
}

// PostgreSQLConfig holds the configuration for PostgreSQL connection.
type PostgreSQLConfig struct {
	host     string
	port     int
	user     string
	password string
	database string
	sslMode  string
}

// NewPostgreSQLConfig creates PostgreSQLConfig.
func NewPostgreSQLConfig(
	host string,
	port int,
	user string,
	password string,
	database string,
) PostgreSQLConfig {
	return PostgreSQLConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
		sslMode:  "disable", // Default to disable for development
	}
}

// NewPostgreSQLDB creates *sql.DB for PostgreSQL.
// The return function is the function to close the DB.
func NewPostgreSQLDB(config PostgreSQLConfig) (PostgreSQLDB, func(), error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.host,
		config.port,
		config.user,
		config.password,
		config.database,
		config.sslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	return PostgreSQLDB(db), func() { db.Close() }, nil
}
