package storage

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Filesystem is a storage backend that uses the local filesystem. It is intended only
// for use in development to avoid dependency on an external service.
type Filesystem struct {
	root    string
	webRoot string
	logger  *zap.Logger
}

// NewFilesystem creates a new S3 using the provided AWS session.
func NewFilesystem(root string, webRoot string, logger *zap.Logger) *Filesystem {
	return &Filesystem{root, webRoot, logger}
}

// Store stores the content from an io.ReadSeeker at the specified key.
func (fs *Filesystem) Store(key string, data io.ReadSeeker, checksum string) (*StoreResult, error) {
	joined := filepath.Join(fs.root, key)
	dir := filepath.Dir(joined)

	/*
		#nosec - filesystem storage is only used for local development.
	*/
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "could not create parent directory")
	}

	file, err := os.Create(joined)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return nil, errors.Wrap(err, "write to disk failed")
	}
	return &StoreResult{}, nil
}

// Delete deletes the file at the specified key
func (fs *Filesystem) Delete(key string) error {
	joined := filepath.Join(fs.root, key)

	return os.Remove(joined)
}

// Key returns a joined key
func (fs *Filesystem) Key(args ...string) string {
	return path.Join(args...)
}

// PresignedURL returns a URL that provides access to a file for 15 mintes.
func (fs *Filesystem) PresignedURL(key, contentType string) (string, error) {
	values := url.Values{}
	values.Add("contentType", contentType)
	url := fs.webRoot + "/" + key + "?" + values.Encode()
	return url, nil
}

// Fetch retrieves a copy of a file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (fs *Filesystem) Fetch(key string) (string, error) {
	outputFile, err := ioutil.TempFile(os.TempDir(), "filesystem")
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer outputFile.Close()

	sourcePath := filepath.Join(fs.root, key)
	// #nosec - this code is only used in development as a fallback for S3
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer sourceFile.Close()

	_, err = io.Copy(outputFile, sourceFile)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return outputFile.Name(), nil
}

// NewFilesystemHandler returns an Handler that adds a Content-Type header so that
// files are handled properly by the browser.
func NewFilesystemHandler(root string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.URL.Query().Get("contentType")
		if contentType != "" {
			w.Header().Add("Content-Type", contentType)
		}

		input := filepath.Join(root, filepath.FromSlash(path.Clean("/"+r.URL.Path)))
		http.ServeFile(w, r, input)
	})
}
