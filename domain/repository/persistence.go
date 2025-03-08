package repository

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../../infrastructure/mock/$GOFILE -package mock

type (
	// QueryToRemoteExecutor executes a query in database.
	QueryToRemoteExecutor interface {
		ExecuteQuery(ctx context.Context, sql *model.SQL) (*model.Table, error)
	}

	// StatementToRemoteExecutor executes a statement in memory.
	StatementToRemoteExecutor interface {
		ExecuteStatement(ctx context.Context, sql *model.SQL) (int64, error)
	}

	// TablesInRemoteGetter gets tables in database.
	TablesInRemoteGetter interface {
		GetTables(ctx context.Context) ([]*model.Table, error)
	}
)
