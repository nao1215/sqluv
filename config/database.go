package config

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
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
		DBName:               config.database,
		User:                 config.user,
		Passwd:               config.password,
		Addr:                 fmt.Sprintf("%s:%d", config.host, config.port),
		Net:                  "tcp",
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
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

// SQLite3Config holds the configuration for SQLite3 connection.
type SQLite3Config struct {
	filepath string
}

// NewSQLite3Config creates SQLite3Config.
func NewSQLite3Config(filepath string) SQLite3Config {
	return SQLite3Config{
		filepath: filepath,
	}
}

// SQLite3DB is *sql.DB for executing SQL with SQLite3.
type SQLite3DB = DBMS

// NewSQLite3DB creates *sql.DB for SQLite3.
// The return function is the function to close the DB.
func NewSQLite3DB(config SQLite3Config) (SQLite3DB, func(), error) {
	// We're already using the sqlite driver for in-memory DB
	// Make sure it's initialized
	initSQLite3()

	db, err := sql.Open("sqlite3", config.filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to SQLite3: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping SQLite3: %w", err)
	}
	return SQLite3DB(db), func() { db.Close() }, nil
}

// SQLServerConfig holds the configuration for SQL Server connection.
type SQLServerConfig struct {
	host     string
	port     int
	user     string
	password string
	database string
	// Additional SQL Server specific options can be added here
	trustServerCertificate bool
}

// NewSQLServerConfig creates SQLServerConfig.
func NewSQLServerConfig(
	host string,
	port int,
	user string,
	password string,
	database string,
) SQLServerConfig {
	return SQLServerConfig{
		host:                   host,
		port:                   port,
		user:                   user,
		password:               password,
		database:               database,
		trustServerCertificate: true, // For development purposes
	}
}

// SQLServerDB is *sql.DB for executing SQL with SQL Server.
type SQLServerDB = DBMS

// NewSQLServerDB creates *sql.DB for SQL Server.
// The return function is the function to close the DB.
func NewSQLServerDB(config SQLServerConfig) (SQLServerDB, func(), error) {
	// SQL Server connection string format
	connStr := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=true;trustservercertificate=%t",
		config.host,
		config.user,
		config.password,
		config.port,
		config.database,
		config.trustServerCertificate,
	)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to SQL Server: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to ping SQL Server: %w", err)
	}
	return SQLServerDB(db), func() { db.Close() }, nil
}

// HistoryDB is *sql.DB for sqluv shell history.
type HistoryDB *sql.DB

// NewHistoryDB create *sql.DB for history.
// The return function is the function to close the DB.
func NewHistoryDB(cfg *DBConfig) (HistoryDB, func(), error) {
	db, err := sql.Open("sqlite3", cfg.hisotryDBPath)
	if err != nil {
		return nil, nil, err
	}
	return HistoryDB(db), func() { db.Close() }, nil
}
