package services

// TableFromSliceCreator creates and populates a table based on a model slice
//go:generate mockery --name TableFromSliceCreator --disable-version-string
type TableFromSliceCreator interface {
	CreateTableFromSlice(slice interface{}) error
}
