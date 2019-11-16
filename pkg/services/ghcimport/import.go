package ghcimport

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type GHCRateEngineImporter struct {
	Logger       Logger
	ContractCode string
	ContractName string
	// TODO: add reference maps here as needed for dependencies between tables
	// like UUID maps for domestic service areas
	// domesticServiceAreaUUIDs map[string]uuid.UUID
	contractID         uuid.UUID
	serviceAreaToIDMap map[string]uuid.UUID
}

func (gre *GHCRateEngineImporter) runImports(dbTx *pop.Connection) error {
	err := gre.importREContract(dbTx)
	if err != nil {
		return errors.Wrap(err, "Failed to import re_contract")
	}

	err = gre.importRERateArea(dbTx)
	if err != nil {
		return errors.Wrap(err, "Failed to import re_rate_area")
	}

	err = gre.importREDomesticServiceArea(dbTx)
	if err != nil {
		return errors.Wrap(err, "Failed to import re_domestic_service_area")
	}

	err = gre.importREDomesticLinehaulPrices(dbTx)
	if err != nil {
		return errors.Wrap(err, "Failed to import re_domestic_linehaul_prices")
	}

	return nil
}

func (gre *GHCRateEngineImporter) Import(db *pop.Connection) error {

	err := db.Transaction(func(connection *pop.Connection) error {
		dbTxError := gre.runImports(connection)
		return dbTxError
	})
	if err != nil {
		return errors.Wrap(err, "Transaction failed during GHC Rate Engine Import()")
	}
	return nil
}
