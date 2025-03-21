package persistence

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
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
		got, err := c.ReadCSV(t.Context(), file)
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
		got, err := c.ReadTSV(t.Context(), file)
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
		got, err := c.ReadLTSV(t.Context(), file)
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
		if _, err = c.ReadLTSV(t.Context(), file); !errors.Is(err, infrastructure.ErrNoLabel) {
			t.Errorf("error is wrong. got: %v, want: %v", err, infrastructure.ErrNoLabel)
		}
	})
}

func TestIOReaderHTTPS(t *testing.T) {
	t.Parallel()

	t.Run("success to read data from HTTPS", func(t *testing.T) {
		t.Parallel()

		file, err := model.NewFile("https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv")
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}

		reader, cleanup, err := ioReader(t.Context(), file)
		if err != nil {
			t.Fatalf("ioReader returned error: %v", err)
		}
		defer func() {
			if err := cleanup(); err != nil {
				t.Errorf("cleanup returned error: %v", err)
			}
		}()

		data, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read data from reader: %v", err)
		}

		content := string(data)
		if !strings.HasPrefix(content, "actor,total_gross") {
			t.Errorf("unexpected file format, got: %v", content[:20])
		}
	})
}
