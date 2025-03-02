// Package repository manage sqluv repository logic interface.
package repository

import (
	"github.com/nao1215/sqluv/domain/model"
)

//go:generate mockgen -typed -source=$GOFILE -destination=../../infrastructure/mock/$GOFILE -package mock

type (
	// CSVReader is an interface for reading records from CSV files and returning them as model.Table.
	CSVReader interface {
		ReadCSV(file *model.File) (*model.Table, error)
	}

	// TSVReader is an interface for reading records from TSV files and returning them as model.Table.
	TSVReader interface {
		ReadTSV(file *model.File) (*model.Table, error)
	}

	// LTSVReader is an interface for reading records from LTSV files and returning them as model.Table.
	LTSVReader interface {
		ReadLTSV(file *model.File) (*model.Table, error)
	}
)
