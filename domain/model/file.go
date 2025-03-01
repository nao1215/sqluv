package model

import (
	"errors"
	"os"
	"path/filepath"
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
) (*File, error) {
	if path == "" {
		return nil, errors.New("file path is empty")
	}
	return &File{
		path: path,
	}, nil
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

// NameWithoutExt return file name without extension.
func (f *File) NameWithoutExt() string {
	filename := filepath.Base(f.path)
	ext := filepath.Ext(filename)

	// Handle hidden files (starting with a dot)
	if strings.HasPrefix(filename, ".") {
		// If the filename is just a dot (like ".gitignore"), keep it as is
		if filename == ext {
			return filename
		}
	}
	return strings.TrimSuffix(filename, ext)
}
