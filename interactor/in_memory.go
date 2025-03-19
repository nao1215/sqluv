package interactor

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/usecase"
)

// _ interface implementation check
var _ usecase.TableCreator = (*tableCreator)(nil)

type tableCreator struct {
	repository.TableCreator
}

// NewTableCreator create new TableCreator.
func NewTableCreator(
	c repository.TableCreator,
) usecase.TableCreator {
	return &tableCreator{
		TableCreator: c,
	}
}

// CreateTable creates a table in the database.
func (r *tableCreator) CreateTable(ctx context.Context, t *model.Table) error {
	return r.TableCreator.CreateTable(ctx, t)
}

var _ usecase.TablesGetter = (*localTablesGetter)(nil)

type localTablesGetter struct {
	repository.TablesGetter
}

// NewLocalTablesGetter create new TablesGetter.
func NewLocalTablesGetter(
	g repository.TablesGetter,
) usecase.TablesGetter {
	return &localTablesGetter{
		TablesGetter: g,
	}
}

// GetTables gets tables in the database.
func (r *localTablesGetter) GetTables(ctx context.Context) ([]*model.Table, error) {
	return r.TablesGetter.GetTables(ctx)
}

// _ interface implementation check
var _ usecase.RecordsInserter = (*recordsInserter)(nil)

type recordsInserter struct {
	repository.RecordsInserter
}

// NewRecordsInserter create new RecordsInserter.
func NewRecordsInserter(
	i repository.RecordsInserter,
) usecase.RecordsInserter {
	return &recordsInserter{
		RecordsInserter: i,
	}
}

// InsertRecords inserts records in memory.
func (r *recordsInserter) InsertRecords(ctx context.Context, t *model.Table) error {
	return r.RecordsInserter.InsertRecords(ctx, t)
}

// _ interface implementation check
var _ usecase.SQLExecutor = (*sqlExecutor)(nil)

// sqlExecutor executes a SQL statement.
type sqlExecutor struct {
	repository.QueryExecutor
	repository.StatementExecutor
}

// NewSQLExecutor create new SQLExecutor.
func NewSQLExecutor(
	q repository.QueryExecutor,
	s repository.StatementExecutor,
) usecase.SQLExecutor {
	return &sqlExecutor{
		QueryExecutor:     q,
		StatementExecutor: s,
	}
}

// ExecuteSQL executes a SQL statement.
func (r *sqlExecutor) ExecuteSQL(ctx context.Context, sql *model.SQL) (*usecase.ExecuteSQLOutput, error) {
	if sql.IsSelect() || sql.IsExplain() {
		table, err := r.QueryExecutor.ExecuteQuery(ctx, sql)
		if err != nil {
			return nil, err
		}
		return usecase.NewExecuteSQLOutput(table, 0), nil
	}
	rowsAffected, err := r.StatementExecutor.ExecuteStatement(ctx, sql)
	if err != nil {
		return nil, err
	}
	return usecase.NewExecuteSQLOutput(nil, rowsAffected), nil
}
