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
		f, err := NewFile(path)
		if err != nil {
			t.Errorf("error should be nil. got: %v", err)
		}

		if f.path != path {
			t.Errorf("path is wrong. got: %s, want: %s", f.path, path)
		}
		if f.protocol != "file://" {
			t.Errorf("protocol is wrong. got: %s, want: file://", f.protocol)
		}
	})

	t.Run("success to create File with protocol", func(t *testing.T) {
		t.Parallel()

		path := "file:///path/to/file"
		f, err := NewFile(path)
		if err != nil {
			t.Errorf("error should be nil. got: %v", err)
		}

		if f.path != "/path/to/file" {
			t.Errorf("path is wrong. got: %s, want: %s", f.path, path)
		}
		if f.protocol != "file://" {
			t.Errorf("protocol is wrong. got: %s, want: ", f.protocol)
		}
	})

	t.Run("fail to create File", func(t *testing.T) {
		t.Parallel()

		path := ""
		_, err := NewFile(path)
		if err == nil {
			t.Error("error should not be nil")
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

		f, err := NewFile(filepath.Join("testdata", "sample.txt"))
		if err != nil {
			t.Errorf("error should be nil. got: %v", err)
		}
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

		f, err := NewFile("notfound.txt")
		if err != nil {
			t.Errorf("error should be nil. got: %v", err)
		}
		_, err = f.Open()
		if err == nil {
			t.Error("error should not be nil")
		}
	})
}

func TestFileNameWithoutExt(t *testing.T) {
	t.Parallel()

	type fields struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "file with extension",
			fields: fields{
				path: "test.csv",
			},
			want: "test",
		},
		{
			name: "file without extension",
			fields: fields{
				path: "test",
			},
			want: "test",
		},
		{
			name: "file with multiple dots",
			fields: fields{
				path: "test.data.csv",
			},
			want: "test.data",
		},
		{
			name: "hidden file",
			fields: fields{
				path: ".test",
			},
			want: ".test",
		},
		{
			name: "file with path",
			fields: fields{
				path: "/path/to/test.csv",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &File{
				path: tt.fields.path,
			}
			if got := f.NameWithoutExt(); got != tt.want {
				t.Errorf("File.NameWithoutExt() = %v, want %v: %s", got, tt.want, tt.name)
			}
		})
	}
}

func TestFile_IsFileProtocol(t *testing.T) {
	t.Parallel()

	type fields struct {
		path     string
		protocol string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "file protocol",
			fields: fields{
				path:     "path/to/file",
				protocol: "file://",
			},
			want: true,
		},
		{
			name: "not file protocol",
			fields: fields{
				path:     "path/to/file",
				protocol: "http://",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &File{
				path:     tt.fields.path,
				protocol: tt.fields.protocol,
			}
			if got := f.IsFileProtocol(); got != tt.want {
				t.Errorf("File.IsFileProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_IsHTTPProtocol(t *testing.T) {
	t.Parallel()

	type fields struct {
		path     string
		protocol string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "http protocol",
			fields: fields{
				path:     "path/to/file",
				protocol: "http://",
			},
			want: true,
		},
		{
			name: "https protocol",
			fields: fields{
				path:     "path/to/file",
				protocol: "https://",
			},
			want: true,
		},
		{
			name: "not http protocol",
			fields: fields{
				path:     "path/to/file",
				protocol: "file://",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &File{
				path:     tt.fields.path,
				protocol: tt.fields.protocol,
			}
			if got := f.IsHTTPProtocol(); got != tt.want {
				t.Errorf("File.IsHTTPProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_FullURL(t *testing.T) {
	t.Parallel()

	t.Run("get full URL", func(t *testing.T) {
		t.Parallel()

		f := &File{
			path:     "/path/to/file",
			protocol: "http://",
		}
		want := "http:///path/to/file"
		if got := f.FullURL(); got != want {
			t.Errorf("File.FullURL() = %v, want %v", got, want)
		}
	})
}
