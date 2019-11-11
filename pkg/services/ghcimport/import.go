package ghcimport

import (
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

type GHCRateEngineImporter struct {
	DB     *pop.Connection
	Logger Logger
}

func (gre *GHCRateEngineImporter) runImports() error {

	err := gre.importRERateArea()
	if err != nil {
		return errors.Wrap(err, "Failed to import re_rate_area")
	}

	err = gre.importREDomesticServiceArea()
	if err != nil {
		return errors.Wrap(err, "Failed to import re_domestic_service_area")
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
