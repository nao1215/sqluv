package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

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

	table, err := infrastructure.Query(ctx, tx, sql)
	if err != nil {
		return nil, err
	}
	return model.NewTable(infrastructure.ExtractTableName(sql), table.Header(), table.Records()), nil
}

// _ interface implementation check
var _ repository.TablesInRemoteGetter = (*tablesGetter)(nil)

type tablesGetter struct {
	db       *sql.DB
	database string
	user     string
	dbmsType config.DBMSType
}

// NewTablesGetter returns tablesGetter
func NewTablesGetter(db config.DBMS, conf *config.DBConnection) repository.TablesInRemoteGetter {
	return &tablesGetter{
		db:       db,
		database: conf.Database,
		user:     conf.User,
		dbmsType: conf.Type,
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
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema = $1 OR table_schema = 'public'"
		rows, err = g.db.QueryContext(ctx, query, g.user)
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
		columnQuery = "SELECT column_name FROM information_schema.columns WHERE table_schema = $1 OR table_schema = 'public' AND table_name = $2"
		colRows, err = g.db.QueryContext(ctx, columnQuery, g.user, tableName)
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

type ddlGetter struct {
	db       *sql.DB
	dbmsType config.DBMSType
	database string
	user     string
}

// NewTableDDLGetter returns a new TableDDLGetter.
func NewTableDDLGetter(
	db *sql.DB,
	conf *config.DBConnection,
) repository.TableDDLInRemoteGetter {
	return &ddlGetter{
		db:       db,
		database: conf.Database,
		user:     conf.User,
		dbmsType: conf.Type,
	}
}

// GetTableDDL retrieves the DDL info of the given table as a Table struct.
func (d *ddlGetter) GetTableDDL(ctx context.Context, tableName string) ([]*model.Table, error) {
	var query string
	var rows *sql.Rows
	var err error

	switch d.dbmsType {
	case config.MySQL:
		query = `
            SELECT COLUMN_NAME, DATA_TYPE, IFNULL(CHARACTER_MAXIMUM_LENGTH, 0),
                   IS_NULLABLE, IFNULL(COLUMN_DEFAULT, ''), COLUMN_KEY
            FROM information_schema.columns
            WHERE table_schema=? AND table_name=?`
		rows, err = d.db.QueryContext(ctx, query, d.database, tableName)
	case config.PostgreSQL:
		query = `
            SELECT column_name, data_type, COALESCE(character_maximum_length, 0),
                   is_nullable, COALESCE(column_default, ''), ''
            FROM information_schema.columns
            WHERE table_schema='public' OR table_schema=$1 AND table_name=$2::text`
		rows, err = d.db.QueryContext(ctx, query, d.user, tableName)
	case config.SQLite3:
		// PRAGMA table_info returns: cid, name, type, notnull, dflt_value, pk
		query = "PRAGMA table_info(" + tableName + ")"
		rows, err = d.db.QueryContext(ctx, query)
	case config.SQLServer:
		query = `
            SELECT COLUMN_NAME, DATA_TYPE, ISNULL(CHARACTER_MAXIMUM_LENGTH, 0),
                   IS_NULLABLE, ISNULL(COLUMN_DEFAULT, ''), '' as column_key
            FROM INFORMATION_SCHEMA.COLUMNS 
            WHERE TABLE_CATALOG=? AND TABLE_NAME=?`
		rows, err = d.db.QueryContext(ctx, query, d.database, tableName)
	default:
		return nil, fmt.Errorf("unsupported DBMS type: %v", d.dbmsType)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define header for DDL info using the model.Header (CSV/TSV table header)
	header := model.NewHeader([]string{"Column Name", "Type", "Precision", "Nullable", "DefaultValue", "PrimaryKey"})
	var records []model.Record

	for rows.Next() {
		switch d.dbmsType {
		case config.SQLite3:
			var cid int
			var colName, colType string
			var notnull, pk int
			var dflt sql.NullString
			if err := rows.Scan(&cid, &colName, &colType, &notnull, &dflt, &pk); err != nil {
				return nil, err
			}
			dfltValue := ""
			if dflt.Valid {
				dfltValue = dflt.String
			}
			nullable := "YES"
			if notnull != 0 {
				nullable = "NO"
			}
			key := ""
			if pk > 0 {
				key = "PRI"
			}
			records = append(records, model.NewRecord([]string{
				colName, colType, "0", nullable, dfltValue, key,
			}))
		default:
			var colName, dataType, nullable, dfltValue, columnKey string
			var precision int
			if err := rows.Scan(&colName, &dataType, &precision, &nullable, &dfltValue, &columnKey); err != nil {
				return nil, err
			}
			records = append(records, model.NewRecord([]string{
				colName, dataType, strconv.Itoa(precision), nullable, dfltValue, columnKey,
			}))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	ddlTable := model.NewTable(tableName, header, records)
	return []*model.Table{ddlTable}, nil
}
