package config

import (
	"github.com/nao1215/sqluv/domain/model"
	"github.com/spf13/pflag"
)

var (
	// Version is sqluv command version. Version value is assigned by LDFLAGS.
	Version string
)

// Argument represents a runtime argument.
type Argument struct {
	// files is CSV/TSV/LTSV file path list that import to SQLite3 in-memory mode.
	files []*model.File
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

	flag.BoolVarP(&helpFlag, "help", "h", false, "print help message")
	flag.BoolVarP(&versionFlag, "version", "v", false, "print sqluv version")
	if err := flag.Parse(args[1:]); err != nil {
		return nil, err
	}

	files := make([]*model.File, 0, len(flag.Args()))
	for _, filePath := range flag.Args() {
		f, err := model.NewFile(filePath)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return &Argument{
		files:   files,
		usage:   newUsage(helpFlag, flag),
		version: newVersion(versionFlag),
	}, nil
}

// Files returns CSV/TSV/LTSV file path list.
func (a *Argument) Files() []*model.File {
	return a.files
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
			return commandName + " (devel)"
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
