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

// IsCSV returns true if the file is a CSV.
// It now also returns true for files ending with .csv.gz or .csv.xz.
func (f *File) IsCSV() bool {
	if f.IsGZ() {
		base := strings.TrimSuffix(f.path, ".gz")
		return strings.HasSuffix(base, ".csv")
	}
	if f.IsBZ2() {
		base := strings.TrimSuffix(f.path, ".bz2")
		return strings.HasSuffix(base, ".csv")
	}
	if f.IsXZ() {
		base := strings.TrimSuffix(f.path, ".xz")
		return strings.HasSuffix(base, ".csv")
	}
	if f.IsZSTD() {
		base := strings.TrimSuffix(f.path, ".zst")
		return strings.HasSuffix(base, ".csv")
	}
	return strings.HasSuffix(f.path, ".csv")
}

// IsTSV returns true if the file is a TSV.
// It now also returns true for files ending with .tsv.gz or .tsv.xz.
func (f *File) IsTSV() bool {
	if f.IsGZ() {
		base := strings.TrimSuffix(f.path, ".gz")
		return strings.HasSuffix(base, ".tsv")
	}
	if f.IsBZ2() {
		base := strings.TrimSuffix(f.path, ".bz2")
		return strings.HasSuffix(base, ".tsv")
	}
	if f.IsXZ() {
		base := strings.TrimSuffix(f.path, ".xz")
		return strings.HasSuffix(base, ".tsv")
	}
	if f.IsZSTD() {
		base := strings.TrimSuffix(f.path, ".zst")
		return strings.HasSuffix(base, ".tsv")
	}
	return strings.HasSuffix(f.path, ".tsv")
}

// IsLTSV returns true if the file is a LTSV.
// It now also returns true for files ending with .ltsv.gz or .ltsv.xz.
func (f *File) IsLTSV() bool {
	if f.IsGZ() {
		base := strings.TrimSuffix(f.path, ".gz")
		return strings.HasSuffix(base, ".ltsv")
	}
	if f.IsBZ2() {
		base := strings.TrimSuffix(f.path, ".bz2")
		return strings.HasSuffix(base, ".ltsv")
	}
	if f.IsXZ() {
		base := strings.TrimSuffix(f.path, ".xz")
		return strings.HasSuffix(base, ".ltsv")
	}
	if f.IsZSTD() {
		base := strings.TrimSuffix(f.path, ".zst")
		return strings.HasSuffix(base, ".ltsv")
	}
	return strings.HasSuffix(f.path, ".ltsv")
}

// Open open file.
func (f *File) Open() (*os.File, error) {
	return os.Open(f.path)
}

// Create create file.
func (f *File) Create() (*os.File, error) {
	return os.Create(f.path)
}

// NameWithoutExt return file name without extension.
// e.g. "/home/nao/test.csv" -> "test"、"test.csv.gz" -> "test", ".gitignore" -> ".gitignore"
func (f *File) NameWithoutExt() string {
	base := filepath.Base(f.path)
	if base[0] == '.' {
		return base
	}
	if idx := strings.Index(base, "."); idx != -1 {
		return base[:idx]
	}
	return base
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

// IsGZ returns true if the file has a .gz extension.
func (f *File) IsGZ() bool {
	return strings.HasSuffix(f.path, ".gz")
}

// IsBZ2 returns true if the file has a .bz2 extension.
func (f *File) IsBZ2() bool {
	return strings.HasSuffix(f.path, ".bz2")
}

// IsXZ returns true if the file is compressed with xz (.xz).
func (f *File) IsXZ() bool {
	return strings.HasSuffix(f.path, ".xz")
}

// IsZSTD returns true if the file has a .zstd extension.
func (f *File) IsZSTD() bool {
	return strings.HasSuffix(f.path, ".zst")
}
