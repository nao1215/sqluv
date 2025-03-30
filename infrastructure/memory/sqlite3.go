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

		// Retrieve column info for the table using PRAGMA table_info
		pragmaQuery := "PRAGMA table_info('" + name + "')"
		colRows, err := tx.QueryContext(ctx, pragmaQuery)
		if err != nil {
			return nil, err
		}

		columns := []string{}
		for colRows.Next() {
			var cid int
			var colName string
			var colType string
			var notnull int
			var dfltValue interface{}
			var pk int
			if err := colRows.Scan(&cid, &colName, &colType, &notnull, &dfltValue, &pk); err != nil {
				colRows.Close()
				return nil, err
			}
			columns = append(columns, colName)
		}
		colRows.Close()

		if err := colRows.Err(); err != nil {
			return nil, err
		}

		header := model.NewHeader(columns)
		// Create table with header (column names) and empty records
		tables = append(tables, model.NewTable(name, header, []model.Record{}))
	}
	if err := rows.Err(); err != nil {
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

	table, err := infrastructure.Query(ctx, tx, sql)
	if err != nil {
		return nil, err
	}
	return model.NewTable(infrastructure.ExtractTableName(sql), table.Header(), table.Records()), nil
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

// _ interface implementation check
var _ repository.TableDDLGetter = (*ddlGetter)(nil)

type ddlGetter struct {
	db *sql.DB
}

// NewTableDDLGetter return ddlGetter
func NewTableDDLGetter(db config.MemoryDB) repository.TableDDLGetter {
	return &ddlGetter{db: db}
}

// GetTableDDL get table DDL in memory
func (d *ddlGetter) GetTableDDL(ctx context.Context, tableName string) ([]*model.Table, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// PRAGMA table_info returns: cid, name, type, notnull, dflt_value, pk
	query := "PRAGMA table_info(" + tableName + ")"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define header for DDL info using the model.Header (CSV/TSV table header)
	header := model.NewHeader([]string{"Column Name", "Type", "Precision", "Nullable", "DefaultValue", "PrimaryKey"})
	var records []model.Record

	for rows.Next() {
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
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	ddlTable := model.NewTable(tableName, header, records)
	return []*model.Table{ddlTable}, nil
}
