package ghcimport

import "github.com/gobuffalo/pop"

type RERateAreasImporter struct {
	db     *pop.Connection
	logger Logger
}

func (re RERateAreasImporter) Import() error {
	return nil
}

func (re RERateAreasImporter) Description() string {
	return "re_rate_area importer"
}
