package services

import "github.com/transcom/mymove/pkg/appcontext"

// TableFromSliceCreator creates and populates a table based on a model slice
//go:generate mockery --name TableFromSliceCreator --disable-version-string
type TableFromSliceCreator interface {
	CreateTableFromSlice(appCtx appcontext.AppContext, slice interface{}) error
}
