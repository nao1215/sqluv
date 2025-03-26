package interactor

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/usecase"
)

// _ interface implementation check
var _ usecase.HistoryTableCreator = (*historyTableCreator)(nil)

type historyTableCreator struct {
	r repository.HistoryTableCreator
}

// NewHistoryTableCreator return HistoryTableCreator
func NewHistoryTableCreator(r repository.HistoryTableCreator) usecase.HistoryTableCreator {
	return &historyTableCreator{r: r}
}

// CreateTable create table for sqluv history.
func (hi *historyTableCreator) CreateTable(ctx context.Context) error {
	return hi.r.CreateTable(ctx)
}

// _ interface implementation check
var _ usecase.HistoryCreator = (*historyCreator)(nil)

type historyCreator struct {
	r repository.HistoryCreator
}

// NewHistoryCreator return HistoryCreator
func NewHistoryCreator(r repository.HistoryCreator) usecase.HistoryCreator {
	return &historyCreator{r: r}
}

// Create create history record.
func (hi *historyCreator) Create(ctx context.Context, history model.History) error {
	h := model.Histories{history}
	return hi.r.Create(ctx, h.ToTable())
}

// _ interface implementation check
var _ usecase.HistoryLister = (*historyLister)(nil)

type historyLister struct {
	r repository.HistoryLister
}

// NewHistoryLister return HistoryLister
func NewHistoryLister(r repository.HistoryLister) usecase.HistoryLister {
	return &historyLister{r: r}
}

// List get all sqluv history.
func (hi *historyLister) List(ctx context.Context) (model.Histories, error) {
	return hi.r.List(ctx)
}
