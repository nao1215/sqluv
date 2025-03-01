package model

import (
	"os"
	"strings"
)

// File represents file.
type File struct {
	// path is file path.
	path string
}

// NewFile create new File.
func NewFile(
	path string,
) *File {
	return &File{
		path: path,
	}
}

// IsCSV return true if file is csv.
func (f *File) IsCSV() bool {
	return strings.HasSuffix(f.path, ".csv")
}

// IsTSV return true if file is tsv.
func (f *File) IsTSV() bool {
	return strings.HasSuffix(f.path, ".tsv")
}

// IsLTSV return true if file is ltsv.
func (f *File) IsLTSV() bool {
	return strings.HasSuffix(f.path, ".ltsv")
}

// Open open file.
func (f *File) Open() (*os.File, error) {
	return os.Open(f.path)
}
