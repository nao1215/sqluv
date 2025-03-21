package interactor

import (
	"context"

	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/usecase"
)

// _ interface implementation check
var _ usecase.FileReader = (*fileReader)(nil)

type fileReader struct {
	repository.CSVReader
	repository.TSVReader
	repository.LTSVReader
}

// NewFileReader create new FileReader.
func NewFileReader(
	csvReader repository.CSVReader,
	tsvReader repository.TSVReader,
	ltsvReader repository.LTSVReader,
) usecase.FileReader {
	return &fileReader{
		CSVReader:  csvReader,
		TSVReader:  tsvReader,
		LTSVReader: ltsvReader,
	}
}

// Read read records from CSV/TSV/LTSV files and return them as model.Table.
func (r *fileReader) Read(ctx context.Context, file *model.File) (*model.Table, error) {
	switch {
	case file.IsCSV():
		return r.CSVReader.ReadCSV(ctx, file)
	case file.IsTSV():
		return r.TSVReader.ReadTSV(ctx, file)
	case file.IsLTSV():
		return r.LTSVReader.ReadLTSV(ctx, file)
	default:
		return nil, usecase.ErrNotSupportedFileFormat
	}
}
