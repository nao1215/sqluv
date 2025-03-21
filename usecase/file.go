// Package usecase provides usecase layer.
package usecase

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../interactor/mock/$GOFILE -package mock

type (
	// FileReader is an interface for reading records from CSV/TSV/LTSV files and returning them as model.Table.
	FileReader interface {
		Read(ctx context.Context, file *model.File) (*model.Table, error)
	}
)
