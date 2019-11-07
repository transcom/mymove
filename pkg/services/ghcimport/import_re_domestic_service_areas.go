package ghcimport

import "github.com/gobuffalo/pop"

type REDomesticServiceAreasImporter struct {
	db     *pop.Connection
	logger Logger
}

func (re REDomesticServiceAreasImporter) Import() error {
	return nil
}

func (re REDomesticServiceAreasImporter) Description() string {
	return "re_domestic_service_area importer"
}
