package interactor

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/usecase"
)

// _ interface implementation check
var _ usecase.QueryExecutor = (*queryExecutor)(nil)

type queryExecutor struct {
	repository.QueryToRemoteExecutor
	repository.StatementToRemoteExecutor
}

// NewQueryExecutor creates a new QueryExecutor.
func NewQueryExecutor(
	qe repository.QueryExecutor,
	se repository.StatementToRemoteExecutor,
) usecase.QueryExecutor {
	return &queryExecutor{
		QueryToRemoteExecutor:     qe,
		StatementToRemoteExecutor: se,
	}
}

// ExecuteQuery executes a query in MySQL database.
func (m *queryExecutor) ExecuteQuery(ctx context.Context, sql *model.SQL) (*usecase.ExecuteQueryOutput, error) {
	if sql.IsSelect() || sql.IsExplain() || sql.IsWith() {
		table, err := m.QueryToRemoteExecutor.ExecuteQuery(ctx, sql)
		if err != nil {
			return nil, err
		}
		return usecase.NewExecuteQueryOutput(table, 0), nil
	}
	rowsAffected, err := m.StatementToRemoteExecutor.ExecuteStatement(ctx, sql)
	if err != nil {
		return nil, err
	}
	return usecase.NewExecuteQueryOutput(nil, rowsAffected), nil
}

// _ interface implementation check
var _ usecase.TablesGetter = (*tablesGetter)(nil)

type tablesGetter struct {
	repository.TablesInRemoteGetter
}

// NewTablesGetter creates a new MySQLTablesGetter.
func NewTablesGetter(
	tg repository.TablesInRemoteGetter,
) usecase.TablesGetter {
	return &tablesGetter{
		TablesInRemoteGetter: tg,
	}
}

// GetTables gets tables in MySQL database.
func (m *tablesGetter) GetTables(ctx context.Context) ([]*model.Table, error) {
	return m.TablesInRemoteGetter.GetTables(ctx)
}
