package persistence

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/infrastructure"
)

// _ interface implementation check
var _ repository.CSVReader = (*csvReader)(nil)

type csvReader struct{}

// NewCSVReader return new CSVReader.
func NewCSVReader() repository.CSVReader {
	return &csvReader{}
}

// ReadCSV read records from CSV files and return them as model.CSV.
func (c *csvReader) ReadCSV(f *os.File) (*model.Table, error) {
	r := csv.NewReader(f)

	header := model.Header{}
	records := []model.Record{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if len(header) == 0 {
			header = row
			continue
		}
		records = append(records, model.NewRecord(row))
	}
	return model.NewTable(filepath.Base(f.Name()), model.NewHeader(header), records), nil
}

// _ interface implementation check
var _ repository.TSVReader = (*tsvReader)(nil)

type tsvReader struct{}

// NewTSVReader return new TSVReader.
func NewTSVReader() repository.TSVReader {
	return &tsvReader{}
}

// ReadTSV read records from TSV files and return them as model.TSV.
func (t *tsvReader) ReadTSV(f *os.File) (*model.Table, error) {
	r := csv.NewReader(f)
	r.Comma = '\t'

	header := model.Header{}
	records := []model.Record{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if len(header) == 0 {
			header = row
			continue
		}
		records = append(records, model.NewRecord(row))
	}
	return model.NewTable(filepath.Base(f.Name()), model.NewHeader(header), records), nil
}

// _ interface implementation check
var _ repository.LTSVReader = (*ltsvReader)(nil)

type ltsvReader struct{}

// NewLTSVReader return new LTSVReader.
func NewLTSVReader() repository.LTSVReader {
	return &ltsvReader{}
}

// ReadLTSV read records from LTSV files and return them as model.LTSV.
func (l *ltsvReader) ReadLTSV(f *os.File) (*model.Table, error) {
	r := csv.NewReader(f)
	r.Comma = '\t'

	label := model.Label{}
	records := []model.Record{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if len(label) == 0 {
			for _, v := range row {
				l, _, err := l.labelAndData(v)
				if err != nil {
					return nil, err
				}
				label = append(label, l)
			}
		}

		r := model.Record{}
		for _, v := range row {
			_, data, _ := l.labelAndData(v) //nolint:errcheck // error is already checked.
			r = append(r, data)
		}
		records = append(records, r)
	}
	return model.NewTable(filepath.Base(f.Name()), model.NewHeader(label), records), nil
}

// labelAndData split label and data.
func (l *ltsvReader) labelAndData(s string) (string, string, error) {
	idx := strings.Index(s, ":")
	if idx == -1 || idx == 0 {
		return "", "", infrastructure.ErrNoLabel
	}
	return s[:idx], s[idx+1:], nil
}
