// Package memory handle sqlite3 in memory mode
package memory

import (
	"context"
	"database/sql"

	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/infrastructure"
)

// _ interface implementation check
var _ repository.TableCreator = (*tableCreator)(nil)

type tableCreator struct {
	db *sql.DB
}

// NewTableCreator return tableCreator
func NewTableCreator(db config.MemoryDB) repository.TableCreator {
	return &tableCreator{db: db}
}

// CreateTable create table in memory
func (c *tableCreator) CreateTable(ctx context.Context, t *model.Table) error {
	if err := t.Valid(); err != nil {
		return err
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, infrastructure.GenerateCreateTableStatement(t))
	if err != nil {
		return err
	}
	return tx.Commit()
}

// _ interface implementation check
var _ repository.TablesGetter = (*tableGetter)(nil)

type tableGetter struct {
	db *sql.DB
}

// NewTableGetter return tableGetter
func NewTableGetter(db config.MemoryDB) repository.TablesGetter {
	return &tableGetter{db: db}
}

// GetTables get tables in memory
func (g *tableGetter) GetTables(ctx context.Context) ([]*model.Table, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx,
		"SELECT name FROM sqlite_master WHERE type = 'table'")
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
var _ repository.RecordsInserter = (*recordInserter)(nil)

type recordInserter struct {
	db *sql.DB
}

// NewRecordInserter return recordInserter
func NewRecordInserter(db config.MemoryDB) repository.RecordsInserter {
	return &recordInserter{db: db}
}

// InsertRecords insert records in memory
func (r *recordInserter) InsertRecords(ctx context.Context, t *model.Table) error {
	if err := t.Valid(); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, v := range t.Records() {
		if _, err := tx.ExecContext(ctx, infrastructure.GenerateInsertStatement(t.Name(), v)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// _ interface implementation check
var _ repository.QueryExecutor = (*queryExecutor)(nil)

type queryExecutor struct {
	db *sql.DB
}

// NewQueryExecutor return queryExecutor
func NewQueryExecutor(db config.MemoryDB) repository.QueryExecutor {
	return &queryExecutor{db: db}
}

// ExecuteQueryInMemory execute query in memory
func (e *queryExecutor) ExecuteQuery(ctx context.Context, sql *model.SQL) (*model.Table, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	table, err := query(ctx, tx, sql.String())
	if err != nil {
		return nil, err
	}
	return model.NewTable(extractTableName(sql.String()), table.Header(), table.Records()), nil
}

// _ interface implementation check
var _ repository.StatementExecutor = (*statementExecutor)(nil)

type statementExecutor struct {
	db *sql.DB
}

// NewStatementExecutor return statementExecutor
func NewStatementExecutor(db config.MemoryDB) repository.StatementExecutor {
	return &statementExecutor{db: db}
}

// ExecuteStatementInMemory execute statement in memory
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
