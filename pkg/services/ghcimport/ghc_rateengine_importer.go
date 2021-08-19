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
	importers := []struct {
		importFunction func(*pop.Connection) error
		action         string
	}{
		// NOTE: Ordering is significant as these functions must run in this order.

		// Reference tables
		{gre.importREContract, "Importing contract"},                          // Also populates gre.ContractID
		{gre.importREContractYears, "Importing contract years"},               // Populates gre.contractYearToIDMap
		{gre.importREDomesticServiceArea, "Importing domestic service areas"}, // Also populates gre.serviceAreaToIDMap
		{gre.importRERateArea, "Importing rate areas"},                        // Also populates gre.domesticRateAreaToIDMap and gre.internationalRateAreaToIDMap
		{gre.mapZipCodesToRERateAreas, "Mapping zip3s and zip5s to rate areas"},
		{gre.loadServiceMap, "Loading service map"}, // Populates gre.serviceToIDMap

		// Non-reference tables
		{gre.importREDomesticLinehaulPrices, "Importing domestic linehaul prices"},
		{gre.importREDomesticServiceAreaPrices, "Importing domestic service area prices"},
		{gre.importREDomesticOtherPrices, "Importing domestic other prices"},
		{gre.importREInternationalPrices, "Importing international prices"},
		{gre.importREInternationalOtherPrices, "Importing international other prices"},
		{gre.importRETaskOrderFees, "Importing task order fees"},
		{gre.importREDomesticAccessorialPrices, "Importing domestic accessorial prices"},
		{gre.importREIntlAccessorialPrices, "Importing international accessorial prices"},
		{gre.importREShipmentTypePrices, "Importing shipment type prices"},
	}

	for _, importer := range importers {
		spinner, err := pterm.DefaultSpinner.Start(importer.action)
		if err != nil {
			return fmt.Errorf("failed to create pterm spinner: %w", err)
		}

		err = importer.importFunction(dbTx)
		if err != nil {
			spinner.Fail()
			return fmt.Errorf("importer failed: %s: %w", importer.action, err)
		}

		spinner.Success()
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
