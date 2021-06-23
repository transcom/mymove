package ghcimport

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pterm/pterm"
)

// GHCRateEngineImporter is the rate engine importer for GHC
type GHCRateEngineImporter struct {
	Logger            Logger
	ContractCode      string
	ContractName      string
	ContractStartDate time.Time
	// TODO: add reference maps here as needed for dependencies between tables
	ContractID                   uuid.UUID
	serviceAreaToIDMap           map[string]uuid.UUID
	domesticRateAreaToIDMap      map[string]uuid.UUID
	internationalRateAreaToIDMap map[string]uuid.UUID
	serviceToIDMap               map[string]uuid.UUID
	contractYearToIDMap          map[string]uuid.UUID
}

func (gre *GHCRateEngineImporter) runImports(dbTx *pop.Connection) error {
	// TODO: Check error from spinner calls
	// TODO: Can we wrap these functions somehow to abstract the spinner part?

	// Reference tables
	spinner, _ := pterm.DefaultSpinner.Start("Importing contract")
	err := gre.importREContract(dbTx) // Also populates gre.ContractID
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_contract: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing contract years")
	err = gre.importREContractYears(dbTx) // Populates gre.contractYearToIDMap
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_contract_years: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing domestic service areas")
	err = gre.importREDomesticServiceArea(dbTx) // Also populates gre.serviceAreaToIDMap
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_domestic_service_area: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing rate areas")
	err = gre.importRERateArea(dbTx) // Also populates gre.domesticRateAreaToIDMap and gre.internationalRateAreaToIDMap
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_rate_area: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Mapping zip3s and zip5s to rate areas")
	err = gre.mapZipCodesToRERateAreas(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to map zip3s and zip5s to re_rate_areas: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Loading service map")
	err = gre.loadServiceMap(dbTx) // Populates gre.serviceToIDMap
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to load service map: %w", err)
	}
	spinner.Success()

	// Non-reference tables
	spinner, _ = pterm.DefaultSpinner.Start("Importing domestic linehaul prices")
	err = gre.importREDomesticLinehaulPrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_domestic_linehaul_prices: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing domestic service area prices")
	err = gre.importREDomesticServiceAreaPrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_domestic_service_area_prices: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing domestic other prices")
	err = gre.importREDomesticOtherPrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_domestic_other_prices: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing international prices")
	err = gre.importREInternationalPrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_intl_prices: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing international other prices")
	err = gre.importREInternationalOtherPrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_intl_other_prices: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing task order fees")
	err = gre.importRETaskOrderFees(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_task_order_fees: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing domestic accessorial prices")
	err = gre.importREDomesticAccessorialPrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_domestic_accessorial_prices: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing international accessorial prices")
	err = gre.importREIntlAccessorialPrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_intl_accessorial_prices: %w", err)
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Importing shipment type prices")
	err = gre.importREShipmentTypePrices(dbTx)
	if err != nil {
		spinner.Fail()
		return fmt.Errorf("failed to import re_shipment_type_prices: %w", err)
	}
	spinner.Success()

	return nil
}

// Import runs the import
func (gre *GHCRateEngineImporter) Import(db *pop.Connection) error {
	err := db.Transaction(func(connection *pop.Connection) error {
		dbTxError := gre.runImports(connection)
		return dbTxError
	})
	if err != nil {
		return fmt.Errorf("transaction failed during GHC Rate Engine Import(): %w", err)
	}
	return nil
}
