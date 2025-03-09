package repository

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../../infrastructure/mock/$GOFILE -package mock

type (
	// HistoryTableCreator is a repository that creates a table for sqly shell history.
	HistoryTableCreator interface {
		// CreateTable create a DB table for sqly shell history
		CreateTable(ctx context.Context) error
	}

	// HistoryCreator is a repository that creates a history record in DB.
	HistoryCreator interface {
		// Create set history record in DB
		Create(ctx context.Context, t *model.Table) error
	}

	// HistoryLister is a repository that gets sql shell all history.
	HistoryLister interface {
		// List get sql shell all history.
		List(ctx context.Context) (model.Histories, error)
	}
)
