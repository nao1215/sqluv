package persistence

import (
	"compress/bzip2"
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/infrastructure"
	"github.com/ulikunitz/xz"
)

// _ interface implementation check
var _ repository.CSVReader = (*csvReader)(nil)

type csvReader struct {
	awsClient S3Client
}

// NewCSVReader return new CSVReader.
func NewCSVReader(awsClient S3Client) repository.CSVReader {
	return &csvReader{awsClient: awsClient}
}

// ReadCSV read records from CSV files and return them as model.CSV.
func (c *csvReader) ReadCSV(ctx context.Context, file *model.File) (*model.Table, error) {
	ioReader, closer, err := ioReader(ctx, file, c.awsClient)
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

type tsvReader struct {
	awsClient S3Client
}

// NewTSVReader return new TSVReader.
func NewTSVReader(awsClient S3Client) repository.TSVReader {
	return &tsvReader{awsClient: awsClient}
}

// ReadTSV read records from TSV files and return them as model.TSV.
func (t *tsvReader) ReadTSV(ctx context.Context, file *model.File) (*model.Table, error) {
	ioReader, closer, err := ioReader(ctx, file, t.awsClient)
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

type ltsvReader struct {
	awsClient S3Client
}

// NewLTSVReader return new LTSVReader.
func NewLTSVReader(awsClient S3Client) repository.LTSVReader {
	return &ltsvReader{awsClient: awsClient}
}

// ReadLTSV read records from LTSV files and return them as model.LTSV.
func (l *ltsvReader) ReadLTSV(ctx context.Context, file *model.File) (*model.Table, error) {
	ioReader, closer, err := ioReader(ctx, file, l.awsClient)
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
// If the file is gzip compressed (as indicated by File.IsGZ),
// it wraps the underlying reader with gzip.NewReader.
func ioReader(ctx context.Context, file *model.File, s3Client S3Client) (io.Reader, func() error, error) {
	var reader io.Reader
	var closer func() error
	var err error

	if file.IsS3Protocol() {
		reader, closer, err = ioReaderFromS3(ctx, file, s3Client)
	} else if file.IsHTTPProtocol() {
		reader, closer, err = ioReaderFromHTTP(ctx, file)
	} else {
		reader, closer, err = ioReaderFromFile(file)
	}
	if err != nil {
		return nil, nil, err
	}

	return wrapCompressedReader(file, reader, closer)
}

func wrapCompressedReader(file *model.File, reader io.Reader, closer func() error) (io.Reader, func() error, error) {
	if file.IsGZ() {
		return wrapGZReader(reader, closer)
	}
	if file.IsBZ2() {
		return wrapBZ2Reader(reader, closer)
	}
	if file.IsXZ() {
		return wrapXZReader(reader, closer)
	}
	if file.IsZSTD() {
		return wrapZstdReader(reader, closer)
	}
	return reader, closer, nil
}

func wrapGZReader(reader io.Reader, closer func() error) (io.Reader, func() error, error) {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		closer()
		return nil, nil, err
	}
	newCloser := func() error {
		if err := gzReader.Close(); err != nil {
			return err
		}
		return closer()
	}
	return gzReader, newCloser, nil
}

func wrapBZ2Reader(reader io.Reader, closer func() error) (io.Reader, func() error, error) {
	bz2Reader := bzip2.NewReader(reader)
	return bz2Reader, closer, nil
}

func wrapXZReader(reader io.Reader, closer func() error) (io.Reader, func() error, error) {
	xzReader, err := xz.NewReader(reader)
	if err != nil {
		closer()
		return nil, nil, err
	}
	// xzReader doesn't require a close; use original closer.
	return xzReader, closer, nil
}

func wrapZstdReader(reader io.Reader, closer func() error) (io.Reader, func() error, error) {
	decoder, err := zstd.NewReader(reader)
	if err != nil {
		closer()
		return nil, nil, err
	}
	newCloser := func() error {
		decoder.Close()
		return closer()
	}
	return decoder, newCloser, nil
}

func ioReaderFromS3(ctx context.Context, file *model.File, s3Client S3Client) (io.Reader, func() error, error) {
	// Assume file.BucketAndKey splits the path correctly.
	bucket, key := file.BucketAndKey()
	rc, err := s3Client.GetObject(ctx, bucket, key)
	if err != nil {
		return nil, nil, err
	}
	return rc, func() error { return rc.Close() }, nil
}

func ioReaderFromHTTP(ctx context.Context, file *model.File) (io.Reader, func() error, error) {
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

func ioReaderFromFile(file *model.File) (io.Reader, func() error, error) {
	f, err := file.Open()
	if err != nil {
		return nil, nil, err
	}
	return f, f.Close, nil
}

// _ interface implementation check
var _ repository.CSVWriter = (*csvWriter)(nil)

type csvWriter struct{}

// NewCSVWriter return new CSVWriter.
func NewCSVWriter() repository.CSVWriter {
	return &csvWriter{}
}

// WriteCSV write records to CSV files.
func (c *csvWriter) WriteCSV(_ context.Context, file *model.File, table *model.Table) error {
	f, err := file.Create()
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	records := [][]string{
		table.Header(),
	}
	for _, v := range table.Records() {
		records = append(records, v)
	}
	return w.WriteAll(records)
}

// _ interface implementation check
var _ repository.TSVWriter = (*tsvWriter)(nil)

type tsvWriter struct{}

// NewTSVWriter return new TSVWriter.
func NewTSVWriter() repository.TSVWriter {
	return &tsvWriter{}
}

// WriteTSV write records to TSV files.
func (t *tsvWriter) WriteTSV(_ context.Context, file *model.File, table *model.Table) error {
	f, err := file.Create()
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = '\t'
	records := [][]string{
		table.Header(),
	}
	for _, v := range table.Records() {
		records = append(records, v)
	}
	return w.WriteAll(records)
}

// _ interface implementation check
var _ repository.LTSVWriter = (*ltsvWriter)(nil)

type ltsvWriter struct{}

// NewLTSVWriter return new LTSVWriter.
func NewLTSVWriter() repository.LTSVWriter {
	return &ltsvWriter{}
}

// WriteLTSV write records to LTSV files.
func (l *ltsvWriter) WriteLTSV(_ context.Context, file *model.File, table *model.Table) error {
	f, err := file.Create()
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = '\t'
	records := [][]string{}
	for _, v := range table.Records() {
		r := model.Record{}
		for i, data := range v {
			r = append(r, table.Header()[i]+":"+data)
		}
		records = append(records, r)
	}
	return w.WriteAll(records)
}
