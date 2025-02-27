// Package config provides functionality to read environment variables,
// runtime arguments, and configuration files.
package config

import (
	"errors"
	"runtime/debug"

	"github.com/spf13/pflag"
)

var (
	// Version is sqluv command version. Version value is assigned by LDFLAGS.
	Version string
)

// Argument represents a runtime argument.
type Argument struct {
	// targetSQL is a target SQL query.
	targetSQL *targetSQL
	// usage represents a usage flag.
	usage *usage
	// version represents a version flag.
	version *version
}

// NewArgument creates a new Argument instance.
func NewArgument(args []string) (*Argument, error) {
	flag := pflag.FlagSet{}
	helpFlag := false
	versionFlag := false

	query := flag.StringP("sql", "s", "", "sql query you want to run")
	flag.BoolVarP(&helpFlag, "help", "h", false, "print help message")
	flag.BoolVarP(&versionFlag, "version", "v", false, "print sqluv version")
	if err := flag.Parse(args[1:]); err != nil {
		return nil, err
	}
	filePaths := flag.Args()

	targetSQL, err := newTargetSQL(*query, filePaths)
	if err != nil {
		return nil, err
	}

	return &Argument{
		targetSQL: targetSQL,
		usage:     newUsage(helpFlag, flag),
		version:   newVersion(versionFlag),
	}, nil
}

// IsTUI returns true if sqluv command runs TUI mode.
func (a *Argument) IsTUI() bool {
	if a.CanUsage() || a.CanVersion() {
		return false
	}
	return !a.targetSQL.hasQuery()
}

// CanUsage returns true if sqluv command can show usage message.
func (a *Argument) CanUsage() bool {
	return a.usage.isOn()
}

// CanVersion returns true if sqluv command can show version message.
func (a *Argument) CanVersion() bool {
	return a.version.isOn()
}

// Version returns sqluv command version.
func (a *Argument) Version() string {
	return a.version.String()
}

// Usage returns sqluv command usage.
func (a *Argument) Usage() string {
	return a.usage.String()
}

// targetSQL represents a target SQL query.
// If user does not specify --sql option, sqluv command run TUI mode.
type targetSQL struct {
	// query is SQL query (for --sql option)
	query string
	// filePaths is a file path list that executes SQL query.
	// The expected file format is csv, tsv, or ltsv.
	filePaths []string
}

// newTargetSQL creates a new targetSQL.
// If user set --sql option, you must set file path.
// If user set file path but not set --sql option, sqluv command run TUI mode.
func newTargetSQL(query string, filePaths []string) (*targetSQL, error) {
	if query != "" && len(filePaths) == 0 {
		return nil, errors.New("if you set --sql option, you must set file path")
	}
	return &targetSQL{
		query:     query,
		filePaths: filePaths,
	}, nil
}

// hasQuery returns true if targetSQL has a query.
func (t *targetSQL) hasQuery() bool {
	return t.query != ""
}

// usage represents a usage flag
type usage struct {
	// on is a flag to show usage message.
	on bool
	// message is a usage message.
	message string
}

// newUsage creates a new usage.
func newUsage(on bool, flag pflag.FlagSet) *usage {
	s := `sqluv - simple terminal UI for multiple DBMS & local CSV/TSV/LTSV.
	
[Usage]
  sqluv [OPTIONS] [FILE_PATHS]

[Example]
  - Execute query for csv file.
    sqluv --sql 'SELECT * FROM sample' ./path/to/sample.csv
  - Run TUI mode.
    sqluv

[OPTIONS]
`
	s += flag.FlagUsages()
	s += `
[LICENSE]
  MIT LICENSE - Copyright (c) 2025 CHIKAMATSU Naohiro
  https://github.com/nao1215/sqluv/blob/main/LICENSE

[CONTACT]
  https://github.com/nao1215/sqluv/issues

[NOTE]
  If you execute SQL queries for CSV/TSV/LTSV files,
  sqluv runs the DB in SQLite3 in-memory mode. So, you can use only SQLite3 syntax.
`
	return &usage{
		on:      on,
		message: s,
	}
}

// usage returns usage message.
func (u *usage) String() string {
	return u.message
}

// isOn returns true if usage flag is on.
func (u *usage) isOn() bool {
	return u.on
}

// version represents a version flag.
type version struct {
	// on is a flag to show version message.
	on bool
	// message is a version message.
	message string
}

// newVersion creates a new version.
func newVersion(on bool) *version {
	return &version{
		on: on,
		message: func() string {
			commandName := "sqluv"
			if Version != "" {
				return commandName + " " + Version
			}
			if buildInfo, ok := debug.ReadBuildInfo(); ok {
				if buildInfo.Main.Version != "" {
					return commandName + " " + buildInfo.Main.Version
				}
			}
			return commandName + "(devel)"
		}(),
	}
}

// version returns version message.
func (v *version) String() string {
	return v.message
}

// isOn returns true if version flag is on.
func (v *version) isOn() bool {
	return v.on
}
