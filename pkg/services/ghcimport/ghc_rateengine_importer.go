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
	Logger                       Logger
	ContractCode                 string
	ContractName                 string
	ContractStartDate            time.Time
	ContractID                   uuid.UUID
	serviceAreaToIDMap           map[string]uuid.UUID
	domesticRateAreaToIDMap      map[string]uuid.UUID
	internationalRateAreaToIDMap map[string]uuid.UUID
	serviceToIDMap               map[string]uuid.UUID
	contractYearToIDMap          map[string]uuid.UUID
}

func (gre *GHCRateEngineImporter) runImports(dbTx *pop.Connection) error {
	// Reference tables
	err := withSpinner(gre.importREContract, dbTx, "Importing contract") // Also populates gre.ContractID
	if err != nil {
		return fmt.Errorf("failed to import re_contract: %w", err)
	}

	err = withSpinner(gre.importREContractYears, dbTx, "Importing contract years") // Populates gre.contractYearToIDMap
	if err != nil {
		return fmt.Errorf("failed to import re_contract_years: %w", err)
	}

	err = withSpinner(gre.importREDomesticServiceArea, dbTx, "Importing domestic service areas") // Also populates gre.serviceAreaToIDMap
	if err != nil {
		return fmt.Errorf("failed to import re_domestic_service_area: %w", err)
	}

	err = withSpinner(gre.importRERateArea, dbTx, "Importing rate areas") // Also populates gre.domesticRateAreaToIDMap and gre.internationalRateAreaToIDMap
	if err != nil {
		return fmt.Errorf("failed to import re_rate_area: %w", err)
	}

	err = withSpinner(gre.mapZipCodesToRERateAreas, dbTx, "Mapping zip3s and zip5s to rate areas")
	if err != nil {
		return fmt.Errorf("failed to map zip3s and zip5s to re_rate_areas: %w", err)
	}

	err = withSpinner(gre.loadServiceMap, dbTx, "Loading service map") // Populates gre.serviceToIDMap
	if err != nil {
		return fmt.Errorf("failed to load service map: %w", err)
	}

	// Non-reference tables
	err = withSpinner(gre.importREDomesticLinehaulPrices, dbTx, "Importing domestic linehaul prices")
	if err != nil {
		return fmt.Errorf("failed to import re_domestic_linehaul_prices: %w", err)
	}

	err = withSpinner(gre.importREDomesticServiceAreaPrices, dbTx, "Importing domestic service area prices")
	if err != nil {
		return fmt.Errorf("failed to import re_domestic_service_area_prices: %w", err)
	}

	err = withSpinner(gre.importREDomesticOtherPrices, dbTx, "Importing domestic other prices")
	if err != nil {
		return fmt.Errorf("failed to import re_domestic_other_prices: %w", err)
	}

	err = withSpinner(gre.importREInternationalPrices, dbTx, "Importing international prices")
	if err != nil {
		return fmt.Errorf("failed to import re_intl_prices: %w", err)
	}

	err = withSpinner(gre.importREInternationalOtherPrices, dbTx, "Importing international other prices")
	if err != nil {
		return fmt.Errorf("failed to import re_intl_other_prices: %w", err)
	}

	err = withSpinner(gre.importRETaskOrderFees, dbTx, "Importing task order fees")
	if err != nil {
		return fmt.Errorf("failed to import re_task_order_fees: %w", err)
	}

	err = withSpinner(gre.importREDomesticAccessorialPrices, dbTx, "Importing domestic accessorial prices")
	if err != nil {
		return fmt.Errorf("failed to import re_domestic_accessorial_prices: %w", err)
	}

	err = withSpinner(gre.importREIntlAccessorialPrices, dbTx, "Importing international accessorial prices")
	if err != nil {
		return fmt.Errorf("failed to import re_intl_accessorial_prices: %w", err)
	}

	err = withSpinner(gre.importREShipmentTypePrices, dbTx, "Importing shipment type prices")
	if err != nil {
		return fmt.Errorf("failed to import re_shipment_type_prices: %w", err)
	}

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

func withSpinner(f func(*pop.Connection) error, dbTx *pop.Connection, description string) error {
	spinner, err := pterm.DefaultSpinner.Start(description)
	if err != nil {
		return err
	}

	err = f(dbTx)
	if err != nil {
		spinner.Fail()
		return err
	}

	spinner.Success()
	return nil
}
