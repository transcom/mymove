package middleware

import (
	"net/http"
	"path/filepath"

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
	f, err := cfs.fs.Open(path)
	logger := cfs.appCtx.Logger()
	logger.Info("Using CustomFileSystem for " + path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := cfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
