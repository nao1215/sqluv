package repository

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../../infrastructure/mock/$GOFILE -package mock

type (
	// TableCreator creates a table in the database.
	TableCreator interface {
		CreateTable(ctx context.Context, t *model.Table) error
	}

	// TablesGetter gets tables in the database.
	TablesGetter interface {
		GetTables(ctx context.Context) ([]*model.Table, error)
	}

	// RecordsInserter inserts records in memory.
	RecordsInserter interface {
		InsertRecords(ctx context.Context, t *model.Table) error
	}

	// QueryExecutor executes a query in memory.
	QueryExecutor interface {
		ExecuteQuery(ctx context.Context, sql *model.SQL) (*model.Table, error)
	}

	// StatementExecutor executes a statement in memory.
	StatementExecutor interface {
		ExecuteStatement(ctx context.Context, sql *model.SQL) (int64, error)
	}

	// TableDDLGetter gets a table's DDL information.
	TableDDLGetter interface {
		GetTableDDL(ctx context.Context, tableName string) ([]*model.Table, error)
	}
)
