package interactor

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/sqluv/domain/model"
	infrastructure "github.com/nao1215/sqluv/infrastructure/mock"
	"github.com/nao1215/sqluv/usecase"
	"go.uber.org/mock/gomock"
)

func TestFileReaderRead(t *testing.T) {
	t.Parallel()

	t.Run("success to read csv", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		csvReader := infrastructure.NewMockCSVReader(ctrl)

		// Set up the expected behavior of the mock.
		csvReader.EXPECT().ReadCSV(gomock.Any(), gomock.Any()).Return(
			model.NewTable(
				"test",
				model.Header([]string{"id", "name"}),
				[]model.Record{
					{"1", "foo"},
					{"2", "bar"},
				},
			), nil,
		)

		fileReader := NewFileReader(csvReader, nil, nil)
		file, err := model.NewFile(filepath.Join("testdata", "test.csv"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		got, err := fileReader.Read(t.Context(), file)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := model.NewTable(
			"test",
			model.Header([]string{"id", "name"}),
			[]model.Record{
				{"1", "foo"},
				{"2", "bar"},
			},
		)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("success to read tsv", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		tsvReader := infrastructure.NewMockTSVReader(ctrl)

		// Set up the expected behavior of the mock.
		tsvReader.EXPECT().ReadTSV(gomock.Any(), gomock.Any()).Return(
			model.NewTable(
				"test",
				model.Header([]string{"id", "name"}),
				[]model.Record{
					{"1", "foo"},
					{"2", "bar"},
				},
			), nil,
		)

		fileReader := NewFileReader(nil, tsvReader, nil)
		file, err := model.NewFile(filepath.Join("testdata", "test.tsv"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		got, err := fileReader.Read(t.Context(), file)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := model.NewTable(
			"test",
			model.Header([]string{"id", "name"}),
			[]model.Record{
				{"1", "foo"},
				{"2", "bar"},
			},
		)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("success to read ltsv", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		ltsvReader := infrastructure.NewMockLTSVReader(ctrl)

		// Set up the expected behavior of the mock.
		ltsvReader.EXPECT().ReadLTSV(gomock.Any(), gomock.Any()).Return(
			model.NewTable(
				"test",
				model.Header([]string{"id", "name"}),
				[]model.Record{
					{"1", "foo"},
					{"2", "bar"},
				},
			), nil,
		)

		fileReader := NewFileReader(nil, nil, ltsvReader)
		file, err := model.NewFile(filepath.Join("testdata", "test.ltsv"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		got, err := fileReader.Read(t.Context(), file)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := model.NewTable(
			"test",
			model.Header([]string{"id", "name"}),
			[]model.Record{
				{"1", "foo"},
				{"2", "bar"},
			},
		)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("differs: (-want +got)\n%s", diff)
		}
	})

	t.Run("fail to read unsupported file format", func(t *testing.T) {
		t.Parallel()

		fileReader := NewFileReader(nil, nil, nil)
		file, err := model.NewFile(filepath.Join("testdata", "test.txt"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if _, err := fileReader.Read(t.Context(), file); !errors.Is(err, usecase.ErrNotSupportedFileFormat) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
