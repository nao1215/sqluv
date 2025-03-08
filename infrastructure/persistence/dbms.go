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

// GetTables gets tables in the database
func (g *tablesGetter) GetTables(ctx context.Context) ([]*model.Table, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var rows *sql.Rows
	// Use the appropriate query based on database type
	switch g.dbmsType {
	case config.MySQL:
		// Get all tables in the current database for MySQL
		query := "SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = ?"
		rows, err = tx.QueryContext(ctx, query, g.database)
	case config.PostgreSQL:
		// Get all tables in the current database for PostgreSQL
		query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
		rows, err = tx.QueryContext(ctx, query)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", g.dbmsType)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := []*model.Table{}
	var name string
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, model.NewTable(name, model.Header{}, []model.Record{}))
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return tables, nil
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
