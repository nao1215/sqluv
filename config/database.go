package config

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"

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
