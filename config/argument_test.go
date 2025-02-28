// Package config provides functionality to read environment variables,
// runtime arguments, and configuration files.
package config

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/sqluv/domain/model"
)

func TestNewArgument(t *testing.T) {
	t.Parallel()

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Argument
		wantErr bool
	}{
		{
			name: "If user set --sql option but not set file path, return error",
			args: args{
				args: []string{"sqluv", "--sql", "select * from users"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewArgument(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewArgument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if cmp.Diff(got, tt.want) != "" {
				t.Errorf("NewArgument() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgumentVersion(t *testing.T) {
	t.Run("If user set --version option and developer set sqluv version, return version string", func(t *testing.T) {
		Version = "0.1.0"

		a, err := NewArgument([]string{"sqluv", "--version"})
		if err != nil {
			t.Fatalf("NewArgument() = %v, want nil", err)
		}

		if !a.CanVersion() {
			t.Errorf("CanVersion() = %v, want %v", a.CanUsage(), true)
		}

		want := "sqluv 0.1.0"
		if got := a.Version(); got != want {
			t.Errorf("Version() = %v, want %v", got, want)
		}
	})

	t.Run("If user set --version option and developer does not set sqluv version, return empty string", func(t *testing.T) {
		Version = ""

		a, err := NewArgument([]string{"sqluv", "--version"})
		if err != nil {
			t.Fatalf("NewArgument() = %v, want nil", err)
		}

		if !a.CanVersion() {
			t.Errorf("CanVersion() = %v, want %v", a.CanUsage(), true)
		}

		want := "sqluv (devel)"
		if got := a.Version(); got != want {
			t.Errorf("Version() = %v, want %v", got, want)
		}
	})
}

func TestArgumentUsage(t *testing.T) {
	t.Run("If user set --help option, return usage string", func(t *testing.T) {
		a, err := NewArgument([]string{"sqluv", "--help"})
		if err != nil {
			t.Fatalf("NewArgument() = %v, want nil", err)
		}

		if !a.CanUsage() {
			t.Errorf("CanUsage() = %v, want %v", a.CanUsage(), true)
		}

		want := `sqluv - simple terminal UI for multiple DBMS & local CSV/TSV/LTSV.
	
[Usage]
  sqluv [OPTIONS] [FILE_PATHS]

[OPTIONS]
  -h, --help      print help message
  -v, --version   print sqluv version

[LICENSE]
  MIT LICENSE - Copyright (c) 2025 CHIKAMATSU Naohiro
  https://github.com/nao1215/sqluv/blob/main/LICENSE

[CONTACT]
  https://github.com/nao1215/sqluv/issues

[NOTE]
  If you execute SQL queries for CSV/TSV/LTSV files,
  sqluv runs the DB in SQLite3 in-memory mode. So, you can use only SQLite3 syntax.
`
		if diff := cmp.Diff(a.Usage(), want); diff != "" {
			t.Errorf("Usage() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestArgumentFiles(t *testing.T) {
	t.Parallel()

	type fields struct {
		files   []*model.File
		usage   *usage
		version *version
	}
	tests := []struct {
		name   string
		fields fields
		want   []*model.File
	}{
		{
			name: "If user set file paths, return file paths",
			fields: fields{
				files: []*model.File{
					model.NewFile("users.csv"),
					model.NewFile("users.tsv"),
				},
			},
			want: []*model.File{
				model.NewFile("users.csv"),
				model.NewFile("users.tsv"),
			},
		},
		{
			name: "If user does not set file paths, return empty slice",
			fields: fields{
				files: []*model.File{},
			},
			want: []*model.File{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := &Argument{
				files:   tt.fields.files,
				usage:   tt.fields.usage,
				version: tt.fields.version,
			}
			if got := a.Files(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Argument.FilePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
