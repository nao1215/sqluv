package persistence

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/infrastructure"
)

func TestCSVReaderReadCSV(t *testing.T) {
	t.Parallel()

	t.Run("success to read CSV", func(t *testing.T) {
		t.Parallel()

		file, err := model.NewFile(filepath.Join("testdata", "sample.csv"))
		if err != nil {
			t.Fatal(err)
		}

		c := NewCSVReader()
		got, err := c.ReadCSV(file)
		if err != nil {
			t.Fatal(err)
		}

		want := model.NewTable(
			"sample",
			model.NewHeader([]string{"id", "first_name", "last_name"}),
			[]model.Record{
				model.NewRecord([]string{"1", "John", "Doe"}),
				model.NewRecord([]string{"2", "Jane", "Doe"}),
				model.NewRecord([]string{"3", "John", "Smith"}),
			},
		)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("value is mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestTSVReaderReadTSV(t *testing.T) {
	t.Parallel()

	t.Run("success to read TSV", func(t *testing.T) {
		t.Parallel()

		file, err := model.NewFile(filepath.Join("testdata", "sample.tsv"))
		if err != nil {
			t.Fatal(err)
		}

		c := NewTSVReader()
		got, err := c.ReadTSV(file)
		if err != nil {
			t.Fatal(err)
		}

		want := model.NewTable(
			"sample",
			model.NewHeader([]string{"id", "first_name", "last_name"}),
			[]model.Record{
				model.NewRecord([]string{"1", "John", "Doe"}),
				model.NewRecord([]string{"2", "Jane", "Doe"}),
				model.NewRecord([]string{"3", "John", "Smith"}),
			},
		)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("value is mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestLTSVReaderReadLTSV(t *testing.T) {
	t.Parallel()

	t.Run("success to read LTSV", func(t *testing.T) {
		t.Parallel()

		file, err := model.NewFile(filepath.Join("testdata", "sample.ltsv"))
		if err != nil {
			t.Fatal(err)
		}

		c := NewLTSVReader()
		got, err := c.ReadLTSV(file)
		if err != nil {
			t.Fatal(err)
		}

		want := model.NewTable(
			"sample",
			model.NewHeader([]string{"id", "first_name", "last_name"}),
			[]model.Record{
				model.NewRecord([]string{"1", "John", "Doe"}),
				model.NewRecord([]string{"2", "Jane", "Doe"}),
				model.NewRecord([]string{"3", "John", "Smith"}),
			},
		)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("value is mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("fail to read LTSV: no label", func(t *testing.T) {
		t.Parallel()

		file, err := model.NewFile(filepath.Join("testdata", "no_label.ltsv"))
		if err != nil {
			t.Fatal(err)
		}

		c := NewLTSVReader()
		if _, err = c.ReadLTSV(file); !errors.Is(err, infrastructure.ErrNoLabel) {
			t.Errorf("error is wrong. got: %v, want: %v", err, infrastructure.ErrNoLabel)
		}
	})
}
