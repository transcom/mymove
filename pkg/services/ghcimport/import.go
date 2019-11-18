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
	contractID         uuid.UUID
	serviceAreaToIDMap map[string]uuid.UUID
}

func (gre *GHCRateEngineImporter) runImports(dbTx *pop.Connection) error {
	var err error

	// Reference tables
	gre.contractID, err = gre.importREContract(dbTx)
	if err != nil {
		return errors.Wrap(err, "Failed to import re_contract")
	}

	gre.serviceAreaToIDMap, err = gre.importREDomesticServiceArea(dbTx)
	if err != nil {
		return errors.Wrap(err, "Failed to import re_domestic_service_area")
	}

	err = gre.importRERateArea(dbTx)
	if err != nil {
		return errors.Wrap(err, "Failed to import re_rate_area")
	}

	// Non-reference tables
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
