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
	// protocol is file protocol.
	protocol string
}

// NewFile create new File.
// If path is empty, return error.
// If path does not contain protocol, add file:// protocol.
func NewFile(
	path string,
) (*File, error) {
	if path == "" {
		return nil, errors.New("file path is empty")
	}

	protocol := ""
	if !strings.Contains(path, "://") {
		protocol = "file://"
	} else {
		protocol = strings.Split(path, "://")[0] + "://"
		path = strings.Split(path, "://")[1]
	}

	return &File{
		path:     path,
		protocol: protocol,
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

// IsFileProtocol return true if file protocol is file://.
func (f *File) IsFileProtocol() bool {
	return f.protocol == "file://"
}

// IsHTTPProtocol return true if file protocol is http:// or https://.
func (f *File) IsHTTPProtocol() bool {
	return f.protocol == "http://" || f.protocol == "https://"
}

// FullURL return full URL.
func (f *File) FullURL() string {
	return f.protocol + f.path
}

// IsS3Protocol return true if file protocol is s3://.
func (f *File) IsS3Protocol() bool {
	return f.protocol == "s3://"
}

// BucketAndKey return bucket and key.
func (f *File) BucketAndKey() (string, string) {
	if !f.IsS3Protocol() {
		return "", ""
	}
	if strings.Contains(f.path, "/") {
		return strings.Split(f.path, "/")[0], strings.Split(f.path, "/")[1]
	}
	return strings.Split(f.path, "/")[0], ""
}
