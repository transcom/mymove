package ghcimport

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREContract(appCtx appcontext.AppContext) error {
	if gre.ContractCode == "" {
		return errors.New("no contract code provided")
	}

	// If no contract name is provided, default to the contract code.
	contractName := gre.ContractName
	if contractName == "" {
		contractName = gre.ContractCode
	}

	// See if contract code already exists.
	exists, err := appCtx.DB().Where("code = ?", gre.ContractCode).Exists(&models.ReContract{})
	if err != nil {
		return fmt.Errorf("could not determine if contract code [%s] existed: %w", gre.ContractCode, err)
	}
	if exists {
		return fmt.Errorf("the provided contract code [%s] already exists", gre.ContractCode)
	}

	// Contract code is new; insert it.
	contract := models.ReContract{
		Code: gre.ContractCode,
		Name: contractName,
	}
	verrs, err := appCtx.DB().ValidateAndSave(&contract)
	if verrs.HasAny() {
		return fmt.Errorf("validation errors when saving contract [%+v]: %w", contract, verrs)
	}
	if err != nil {
		return fmt.Errorf("could not save contract [%+v]: %w", contract, err)
	}

	gre.ContractID = contract.ID

	return nil
}
