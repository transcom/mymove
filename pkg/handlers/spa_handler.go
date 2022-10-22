package handlers

import (
	"net/http"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
)

// This is straight from github.com/gorilla/mux

// SpaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type SpaHandler struct {
	staticPath string
	indexPath  string
	cfs        CustomFileSystem
}

// NewSpaHandler returns a new handler for a Single Page App
func NewSpaHandler(staticPath string, indexPath string, cfs CustomFileSystem) SpaHandler {
	return SpaHandler{
		staticPath: staticPath,
		indexPath:  indexPath,
		cfs:        cfs,
	}
}

// from https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
type CustomFileSystem struct {
	fs        http.FileSystem
	indexPath string
	logger    *zap.Logger
}

func NewCustomFileSystem(fs http.FileSystem, indexPath string, logger *zap.Logger) CustomFileSystem {
	return CustomFileSystem{
		fs:        fs,
		indexPath: indexPath,
		logger:    logger,
	}
}

func (cfs CustomFileSystem) Open(path string) (http.File, error) {
	f, openErr := cfs.fs.Open(path)
	logger := cfs.logger
	logger.Debug("Using CustomFileSystem for " + path)

	if openErr != nil {
		logger.Error("Error with opening", zap.Error(openErr))
		return nil, openErr
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, cfs.indexPath)
		if _, indexOpenErr := cfs.fs.Open(index); indexOpenErr != nil {
			closeErr := f.Close()
			if closeErr != nil {
				logger.Error("Unable to close ", zap.Error(closeErr))
				return nil, closeErr
			}

			logger.Error("Unable to open index.html in the directory", zap.Error(indexOpenErr))
			return nil, indexOpenErr
		}
	}

	return f, nil
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	logger.Debug("Using SPA Handler for " + r.URL.Path)

	// get the absolute path to prevent directory traversal
	_, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	// use the customFileSystem so that we do not expose directory listings
	http.FileServer(h.cfs).ServeHTTP(w, r)
}

// NewFileHandler serves up a single file
func NewFileHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}
