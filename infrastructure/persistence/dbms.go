package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/infrastructure"
)

// _ interface implementation check
var _ repository.QueryToRemoteExecutor = (*queryExecutor)(nil)

type queryExecutor struct {
	db *sql.DB
}

// NewQueryExecutor returns queryExecutor
func NewQueryExecutor(db config.DBMS) repository.QueryToRemoteExecutor {
	return &queryExecutor{db: db}
}

// ExecuteQuery executes query in a database
func (e *queryExecutor) ExecuteQuery(ctx context.Context, sql *model.SQL) (*model.Table, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	table, err := infrastructure.Query(ctx, tx, sql.String())
	if err != nil {
		return nil, err
	}
	return model.NewTable(infrastructure.ExtractTableName(sql.String()), table.Header(), table.Records()), nil
}

// _ interface implementation check
var _ repository.TablesInRemoteGetter = (*tablesGetter)(nil)

type tablesGetter struct {
	db       *sql.DB
	database string
	dbmsType config.DBMSType
}

// NewTablesGetter returns tablesGetter
func NewTablesGetter(db config.DBMS, database string, dbmsType config.DBMSType) repository.TablesInRemoteGetter {
	return &tablesGetter{
		db:       db,
		database: database,
		dbmsType: dbmsType,
	}
}

// GetTables gets tables in the database and their columns.
func (g *tablesGetter) GetTables(ctx context.Context) ([]*model.Table, error) {
	var err error
	var query string
	var rows *sql.Rows
	switch g.dbmsType {
	case config.MySQL:
		query = "SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = ?"
		rows, err = g.db.QueryContext(ctx, query, g.database)
	case config.PostgreSQL:
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
		rows, err = g.db.QueryContext(ctx, query)
	case config.SQLite3:
		query = "SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'"
		rows, err = g.db.QueryContext(ctx, query)
	case config.SQLServer:
		query = "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_CATALOG = ?"
		rows, err = g.db.QueryContext(ctx, query, g.database)
	default:
		return nil, fmt.Errorf("unsupported dbms type: %v", g.dbmsType)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := []*model.Table{}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		columns, err := g.getColumns(ctx, tableName)
		if err != nil {
			return nil, err
		}
		header := model.NewHeader(columns)
		tbl := model.NewTable(tableName, header, []model.Record{})
		tables = append(tables, tbl)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

// getColumns retrieves the column names for a given table.
func (g *tablesGetter) getColumns(ctx context.Context, tableName string) ([]string, error) {
	var columnQuery string
	var colRows *sql.Rows
	var err error

	switch g.dbmsType {
	case config.MySQL:
		columnQuery = "SELECT column_name FROM information_schema.columns WHERE table_schema = ? AND table_name = ?"
		colRows, err = g.db.QueryContext(ctx, columnQuery, g.database, tableName)
	case config.PostgreSQL:
		columnQuery = "SELECT column_name FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2"
		colRows, err = g.db.QueryContext(ctx, columnQuery, "public", tableName)
	case config.SQLite3:
		columnQuery = fmt.Sprintf("PRAGMA table_info(%s)", tableName)
		colRows, err = g.db.QueryContext(ctx, columnQuery)
	case config.SQLServer:
		columnQuery = "SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_CATALOG = ? AND TABLE_NAME = ?"
		colRows, err = g.db.QueryContext(ctx, columnQuery, g.database, tableName)
	default:
		return nil, fmt.Errorf("unsupported dbms type: %v", g.dbmsType)
	}
	if err != nil {
		return nil, err
	}
	defer colRows.Close()

	columns := []string{}
	switch g.dbmsType {
	case config.SQLite3:
		var cid int
		var name, colType string
		var notnull, pk int
		var dfltValue sql.NullString
		for colRows.Next() {
			if err := colRows.Scan(&cid, &name, &colType, &notnull, &dfltValue, &pk); err != nil {
				return nil, err
			}
			columns = append(columns, name)
		}
	default:
		var colName string
		for colRows.Next() {
			if err := colRows.Scan(&colName); err != nil {
				return nil, err
			}
			columns = append(columns, colName)
		}
	}
	if err = colRows.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}

// _ interface implementation check
var _ repository.StatementToRemoteExecutor = (*statementExecutor)(nil)

type statementExecutor struct {
	db *sql.DB
}

// NewStatementExecutor return statementExecutor
func NewStatementExecutor(db config.DBMS) repository.StatementToRemoteExecutor {
	return &statementExecutor{db: db}
}

// ExecuteStatement execute statement
func (e *statementExecutor) ExecuteStatement(ctx context.Context, sql *model.SQL) (int64, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, sql.String())
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
