package usecase

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../interactor/mock/$GOFILE -package mock

type (
	// TableCreator creates a table in the database.
	TableCreator interface {
		CreateTable(ctx context.Context, t *model.Table) error
	}

	// ExecuteSQLOutput is the output of ExecuteSQL.
	ExecuteSQLOutput struct {
		table        *model.Table
		rowsAffected int64
	}

	// RecordsInserter inserts records in memory.
	RecordsInserter interface {
		InsertRecords(ctx context.Context, t *model.Table) error
	}

	// SQLExecutor executes a SQL statement.
	SQLExecutor interface {
		ExecuteSQL(ctx context.Context, sql *model.SQL) (*ExecuteSQLOutput, error)
	}
)

// NewExecuteSQLOutput creates a new ExecuteSQLOutput.
// table may be nil.
func NewExecuteSQLOutput(table *model.Table, rowsAffected int64) *ExecuteSQLOutput {
	return &ExecuteSQLOutput{
		table:        table,
		rowsAffected: rowsAffected,
	}
}

// RowsAffected returns the number of rows affected.
func (e ExecuteSQLOutput) RowsAffected() int64 {
	return e.rowsAffected
}

// Table returns the table.
func (e ExecuteSQLOutput) Table() *model.Table {
	return e.table
}

// HasTable returns true if the table is not nil.
func (e ExecuteSQLOutput) HasTable() bool {
	return e.table != nil
}
