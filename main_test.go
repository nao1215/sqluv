// Package main is the entry point of the sqluv command.
package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		want       int
		wantStdout string
		wantStderr string
	}{
		{
			name: "If user set --help option, return usage string",
			args: []string{"sqluv", "--help"},
			want: 0,
			wantStdout: `sqluv - simple terminal UI for multiple DBMS & local CSV/TSV/LTSV.
	
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
`,
			wantStderr: "",
		},
		{
			name:       "If user set --version option, return version string",
			args:       []string{"sqluv", "--version"},
			want:       0,
			wantStdout: "sqluv (devel)",
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			if got := run(stdout, stderr, tt.args); got != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
			if diff := cmp.Diff(stdout.String(), tt.wantStdout); diff != "" {
				t.Errorf("stdout mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(stderr.String(), tt.wantStderr); diff != "" {
				t.Errorf("stderr mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
