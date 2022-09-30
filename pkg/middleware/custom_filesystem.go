package middleware

import (
	"net/http"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

type CustomFileSystem struct {
	fs     http.FileSystem
	appCtx appcontext.AppContext
}

func NewCustomFileSystem(fs http.FileSystem, appCtx appcontext.AppContext) CustomFileSystem {
	return CustomFileSystem{
		fs:     fs,
		appCtx: appCtx,
	}
}

func (cfs CustomFileSystem) Open(path string) (http.File, error) {
	f, openErr := cfs.fs.Open(path)
	logger := cfs.appCtx.Logger()
	logger.Debug("Using CustomFileSystem for " + path)

	if openErr != nil {
		logger.Error("Error with opening", zap.Error(openErr))
		return nil, openErr
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
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
