package persistence

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
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
func (c *csvReader) ReadCSV(ctx context.Context, file *model.File) (*model.Table, error) {
	ioReader, closer, err := ioReader(ctx, file)
	if err != nil {
		return nil, err
	}
	defer closer()

	r := csv.NewReader(ioReader)
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
	return model.NewTable(filepath.Base(file.NameWithoutExt()), model.NewHeader(header), records), nil
}

// _ interface implementation check
var _ repository.TSVReader = (*tsvReader)(nil)

type tsvReader struct{}

// NewTSVReader return new TSVReader.
func NewTSVReader() repository.TSVReader {
	return &tsvReader{}
}

// ReadTSV read records from TSV files and return them as model.TSV.
func (t *tsvReader) ReadTSV(ctx context.Context, file *model.File) (*model.Table, error) {
	ioReader, closer, err := ioReader(ctx, file)
	if err != nil {
		return nil, err
	}
	defer closer()

	r := csv.NewReader(ioReader)
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
	return model.NewTable(filepath.Base(file.NameWithoutExt()), model.NewHeader(header), records), nil
}

// _ interface implementation check
var _ repository.LTSVReader = (*ltsvReader)(nil)

type ltsvReader struct{}

// NewLTSVReader return new LTSVReader.
func NewLTSVReader() repository.LTSVReader {
	return &ltsvReader{}
}

// ReadLTSV read records from LTSV files and return them as model.LTSV.
func (l *ltsvReader) ReadLTSV(ctx context.Context, file *model.File) (*model.Table, error) {
	ioReader, closer, err := ioReader(ctx, file)
	if err != nil {
		return nil, err
	}
	defer closer()

	r := csv.NewReader(ioReader)
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
	return model.NewTable(filepath.Base(file.NameWithoutExt()), model.NewHeader(label), records), nil
}

// labelAndData split label and data.
func (l *ltsvReader) labelAndData(s string) (string, string, error) {
	idx := strings.Index(s, ":")
	if idx == -1 || idx == 0 {
		return "", "", infrastructure.ErrNoLabel
	}
	return s[:idx], s[idx+1:], nil
}

// ioReader returns io.Reader, closer and error.
// If file is HTTP protocol, it returns io.Reader from HTTP response body.
// If file is not HTTP protocol, it returns io.Reader from file.
func ioReader(ctx context.Context, file *model.File) (io.Reader, func() error, error) {
	if file.IsHTTPProtocol() {
		// Assume FullURL returns the full URL (e.g. "https://example.com/path/to/file").
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, file.FullURL(), nil)
		if err != nil {
			return nil, nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, nil, fmt.Errorf("remote file request failed with status: %s", resp.Status)
		}
		return resp.Body, resp.Body.Close, nil
	}
	f, err := file.Open()
	if err != nil {
		return nil, nil, err
	}
	return f, f.Close, nil
}
