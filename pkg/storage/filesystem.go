package storage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

// Filesystem is a storage backend that uses the local filesystem. It is intended only
// for use in development to avoid dependency on an external service.
type Filesystem struct {
	root    string
	webRoot string
	logger  *zap.Logger
	fs      *afero.Afero
	tempFs  *afero.Afero
}

// FilesystemParams contains parameter for instantiating a Filesystem storage backend
type FilesystemParams struct {
	root    string
	webRoot string
	logger  *zap.Logger
}

// NewFilesystemParams returns FilesystemParams after checking path
func NewFilesystemParams(localStorageRoot string, localStorageWebRoot string, logger *zap.Logger) FilesystemParams {
	absTmpPath, err := filepath.Abs(localStorageRoot)
	if err != nil {
		log.Fatalln(fmt.Errorf("could not get absolute path for %s", localStorageRoot))
	}
	storagePath := path.Join(absTmpPath, localStorageWebRoot)
	webRoot := "/" + localStorageWebRoot

	return FilesystemParams{
		root:    storagePath,
		webRoot: webRoot,
		logger:  logger,
	}
}

// NewFilesystem creates a new Filesystem struct using the provided FilesystemParams
func NewFilesystem(params FilesystemParams) *Filesystem {
	var fs = afero.NewOsFs()
	var tempFs = afero.NewMemMapFs()

	return &Filesystem{
		root:    params.root,
		webRoot: params.webRoot,
		logger:  params.logger,
		fs:      &afero.Afero{Fs: fs},
		tempFs:  &afero.Afero{Fs: tempFs},
	}
}

// Store stores the content from an io.ReadSeeker at the specified key.
func (fs *Filesystem) Store(key string, data io.ReadSeeker, checksum string) (*StoreResult, error) {
	if key == "" {
		return nil, errors.New("A valid StorageKey must be set before data can be uploaded")
	}

	joined := filepath.Join(fs.root, key)
	dir := filepath.Dir(joined)

	err := fs.fs.MkdirAll(dir, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "could not create parent directory")
	}

	file, err := fs.fs.Create(joined)
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
	return errors.Wrap(fs.fs.Remove(joined), "could not remove file")
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
func (fs *Filesystem) Fetch(key string) (io.ReadCloser, error) {
	sourcePath := filepath.Join(fs.root, key)
	f, err := fs.fs.Open(sourcePath)
	return f, errors.Wrap(err, "could not open file")
}

// FileSystem returns the underlying afero filesystem
func (fs *Filesystem) FileSystem() *afero.Afero {
	return fs.fs
}

// TempFileSystem returns the temporary afero filesystem
func (fs *Filesystem) TempFileSystem() *afero.Afero {
	return fs.tempFs
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
