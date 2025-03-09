package usecase

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../interactor/mock/$GOFILE -package mock

type (
	// HistoryTableCreator create table for sqluv history.
	HistoryTableCreator interface {
		// CreateTable create table for sqluv history.
		CreateTable(ctx context.Context) error
	}

	// HistoryCreator create history record.
	HistoryCreator interface {
		// Create create history record.
		Create(ctx context.Context, history model.History) error
	}

	// HistoryLister get all sqluv history.
	HistoryLister interface {
		// List get all sqluv history.
		List(ctx context.Context) (model.Histories, error)
	}
)
