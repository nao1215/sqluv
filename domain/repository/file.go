// Package repository manage sqluv repository logic interface.
package repository

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../../infrastructure/mock/$GOFILE -package mock

type (
	// CSVReader is an interface for reading records from CSV files and returning them as model.Table.
	CSVReader interface {
		ReadCSV(ctx context.Context, file *model.File) (*model.Table, error)
	}

	// CSVWriter is an interface for writing records to CSV files.
	CSVWriter interface {
		WriteCSV(ctx context.Context, file *model.File, table *model.Table) error
	}

	// TSVReader is an interface for reading records from TSV files and returning them as model.Table.
	TSVReader interface {
		ReadTSV(ctx context.Context, file *model.File) (*model.Table, error)
	}

	// TSVWriter is an interface for writing records to TSV files.
	TSVWriter interface {
		WriteTSV(ctx context.Context, file *model.File, table *model.Table) error
	}

	// LTSVReader is an interface for reading records from LTSV files and returning them as model.Table.
	LTSVReader interface {
		ReadLTSV(ctx context.Context, file *model.File) (*model.Table, error)
	}

	// LTSVWriter is an interface for writing records to LTSV files.
	LTSVWriter interface {
		WriteLTSV(ctx context.Context, file *model.File, table *model.Table) error
	}
)
