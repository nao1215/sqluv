package model

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewFile(t *testing.T) {
	t.Parallel()

	t.Run("success to create File", func(t *testing.T) {
		t.Parallel()

		path := "path/to/file"
		f := NewFile(path)

		if f.path != path {
			t.Errorf("path is wrong. got: %s, want: %s", f.path, path)
		}
	})
}

func TestFileIsCSV(t *testing.T) {
	type fields struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "file is csv",
			fields: fields{
				path: "test.csv",
			},
			want: true,
		},
		{
			name: "file is not csv",
			fields: fields{
				path: "test.tsv",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				path: tt.fields.path,
			}
			if got := f.IsCSV(); got != tt.want {
				t.Errorf("File.IsCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileIsTSV(t *testing.T) {
	type fields struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "file is tsv",
			fields: fields{
				path: "test.tsv",
			},
			want: true,
		},
		{
			name: "file is not tsv",
			fields: fields{
				path: "test.ltsv",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				path: tt.fields.path,
			}
			if got := f.IsTSV(); got != tt.want {
				t.Errorf("File.IsTSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileIsLTSV(t *testing.T) {
	type fields struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "file is ltsv",
			fields: fields{
				path: "test.ltsv",
			},
			want: true,
		},
		{
			name: "file is not ltsv",
			fields: fields{
				path: "test.csv",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				path: tt.fields.path,
			}
			if got := f.IsLTSV(); got != tt.want {
				t.Errorf("File.IsLTSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileOpen(t *testing.T) {
	t.Parallel()

	t.Run("success to open file", func(t *testing.T) {
		t.Parallel()

		f := NewFile(filepath.Join("testdata", "sample.txt"))
		file, err := f.Open()
		if err != nil {
			t.Errorf("error should be nil. got: %v", err)
		}
		t.Cleanup(func() {
			file.Close()
		})

		data := make([]byte, 6)
		_, err = file.Read(data)
		if err != nil {
			t.Errorf("error should be nil. got: %v", err)
		}

		want := "sample"
		got := string(data)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("value is mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("fail to open file", func(t *testing.T) {
		t.Parallel()

		f := NewFile("notfound.txt")
		_, err := f.Open()
		if err == nil {
			t.Error("error should not be nil")
		}
	})
}
