package ghcimport

import (
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

type GHCREImporter interface {
	Import() error
	Description() string
}

type GHCRateEngineImporter struct {
	DB     *pop.Connection
	Logger Logger
}

func (gre *GHCRateEngineImporter) callImporter(importer GHCREImporter) error {
	err := importer.Import()
	if err != nil {
		return errors.Wrapf(err, "GHC Rate Engine Importer failed for <%s>", importer.Description())
	}
	return nil
}

func (gre *GHCRateEngineImporter) runImports() error {

	// re_domestic_service_areas
	err := gre.callImporter(REDomesticServiceAreasImporter{
		db:     gre.DB,
		logger: gre.Logger,
	})
	if err != nil {
		return err
	}

	// re_rate_area
	err = gre.callImporter(RERateAreasImporter{
		db:     gre.DB,
		logger: gre.Logger,
	})
	if err != nil {
		return err
	}

	return nil
}

func (gre *GHCRateEngineImporter) Import() error {

	err := gre.DB.Transaction(func(connection *pop.Connection) error {
		dbError := gre.runImports()
		return dbError
	})
	if err != nil {
		return errors.Wrap(err, "Transaction failed during GHC Rate Engine Import()")
	}
	return nil
}
