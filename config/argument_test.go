// Package config provides functionality to read environment variables,
// runtime arguments, and configuration files.
package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestArgumentIsTUI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  []string
		want bool
	}{
		{
			name: "If user set no argument, sqluv command is tui mode",
			arg:  []string{"sqluv"},
			want: true,
		},
		{
			name: "If user set --sql option, sqluv command is not tui mode",
			arg:  []string{"sqluv", "--sql", "select * from users", "path/to/file.csv"},
			want: false,
		},
		{
			name: "If user set -s option, sqluv command is not tui mode",
			arg:  []string{"sqluv", "-s", "select * from users", "path/to/file.csv"},
			want: false,
		},
		{
			name: "If user set --help option, sqluv command is not tui mode",
			arg:  []string{"sqluv", "--help"},
			want: false,
		},
		{
			name: "If user set -h option, sqluv command is not tui mode",
			arg:  []string{"sqluv", "-h"},
			want: false,
		},
		{
			name: "If user set --version option, sqluv command is not tui mode",
			arg:  []string{"sqluv", "--version"},
			want: false,
		},
		{
			name: "If user set -v option, sqluv command is not tui mode",
			arg:  []string{"sqluv", "-v"},
			want: false,
		},
		{
			name: "If user set --help and --version option, sqluv command is not tui mode",
			arg:  []string{"sqluv", "--help", "--version"},
			want: false,
		},
		{
			name: "If user set only file path, sqluv command is tui mode",
			arg:  []string{"sqluv", "path/to/file.csv"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a, err := NewArgument(tt.arg)
			if err != nil {
				t.Fatalf("NewArgument() = %v, want nil", err)
			}
			if tt.want != a.IsTUI() {
				t.Errorf("IsTUI() = %v, want %v", a.IsTUI(), tt.want)
			}
		})
	}
}

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

[Example]
  - Execute query for csv file.
    sqluv --sql 'SELECT * FROM sample' ./path/to/sample.csv
  - Run TUI mode.
    sqluv

[OPTIONS]
  -s, --sql string   sql query you want to run
  -h, --help         print help message
  -v, --version      print sqluv version

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
