package ghcimport

import (
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREContract(dbTx *pop.Connection) error {
	if gre.ContractCode == "" {
		return errors.New("No contract code provided")
	}

	// If no contract name is provided, default to the contract code.
	contractName := gre.ContractName
	if contractName == "" {
		contractName = gre.ContractCode
	}

	contract := models.ReContract{
		Code: gre.ContractCode,
		Name: contractName,
	}
	verrs, err := dbTx.ValidateAndSave(&contract)
	if err != nil {
		return errors.Wrapf(err, "Could not save contract: %+v", contract)
	}
	if verrs.HasAny() {
		return errors.Wrapf(verrs, "Validation errors when saving contract: %+v", contract)
	}

	gre.contractID = contract.ID

	return nil
}
