package ghcimport

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREContract(dbTx *pop.Connection) (uuid.UUID, error) {
	if gre.ContractCode == "" {
		return uuid.Nil, errors.New("No contract code provided")
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
		return uuid.Nil, errors.Wrapf(err, "Could not save contract: %+v", contract)
	}
	if verrs.HasAny() {
		return uuid.Nil, errors.Wrapf(verrs, "Validation errors when saving contract: %+v", contract)
	}

	return contract.ID, nil
}
