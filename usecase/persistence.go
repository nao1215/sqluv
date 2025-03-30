package usecase

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../interactor/mock/$GOFILE -package mock

type (
	// ExecuteQueryOutput is the output of ExecuteSQL.
	ExecuteQueryOutput struct {
		table        *model.Table
		rowsAffected int64
	}

	// QueryExecutor executes a query in database.
	QueryExecutor interface {
		ExecuteQuery(ctx context.Context, sql *model.SQL) (*ExecuteQueryOutput, error)
	}

	// TablesGetter gets tables in database.
	TablesGetter interface {
		GetTables(ctx context.Context) ([]*model.Table, error)
	}

	// TableDDLGetter is an interface for retrieving a table's DDL information,
	// including columns, data types, precision, nullability, default values, primary key status, etc.
	TableDDLGetter interface {
		GetTableDDL(ctx context.Context, tableName string) ([]*model.Table, error)
	}
)

// NewExecuteQueryOutput creates a new ExecuteSQLOutput.
// table may be nil.
func NewExecuteQueryOutput(table *model.Table, rowsAffected int64) *ExecuteQueryOutput {
	return &ExecuteQueryOutput{
		table:        table,
		rowsAffected: rowsAffected,
	}
}

// RowsAffected returns the number of rows affected.
func (e ExecuteQueryOutput) RowsAffected() int64 {
	return e.rowsAffected
}

// Table returns the table.
func (e ExecuteQueryOutput) Table() *model.Table {
	return e.table
}

// HasTable returns true if the table is not nil.
func (e ExecuteQueryOutput) HasTable() bool {
	return e.table != nil
}
