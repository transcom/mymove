package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREContract(dbTx *pop.Connection) error {
	if gre.ContractCode == "" {
		return errors.New("no contract code provided")
	}

	// If no contract name is provided, default to the contract code.
	contractName := gre.ContractName
	if contractName == "" {
		contractName = gre.ContractCode
	}

	// See if contract code already exists.
	exists, err := dbTx.Where("code = ?", gre.ContractCode).Exists(&models.ReContract{})
	if err != nil {
		return errors.Wrapf(err, "could not determine if contract code [%s] existed", gre.ContractCode)
	}
	if exists {
		return fmt.Errorf("the provided contract code [%s] already exists", gre.ContractCode)
	}

	// Contract code is new; insert it.
	contract := models.ReContract{
		Code: gre.ContractCode,
		Name: contractName,
	}
	verrs, err := dbTx.ValidateAndSave(&contract)
	if err != nil {
		return errors.Wrapf(err, "could not save contract: %+v", contract)
	}
	if verrs.HasAny() {
		return errors.Wrapf(verrs, "validation errors when saving contract: %+v", contract)
	}

	gre.contractID = contract.ID

	return nil
}
