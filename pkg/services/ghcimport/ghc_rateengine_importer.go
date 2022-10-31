package ghcimport

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/theckman/yacspin"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// GHCRateEngineImporter is the rate engine importer for GHC
type GHCRateEngineImporter struct {
	ContractCode                 string
	ContractName                 string
	ContractStartDate            time.Time
	ContractID                   uuid.UUID
	serviceAreaToIDMap           map[string]uuid.UUID
	domesticRateAreaToIDMap      map[string]uuid.UUID
	internationalRateAreaToIDMap map[string]uuid.UUID
	serviceToIDMap               map[models.ReServiceCode]uuid.UUID
	contractYearToIDMap          map[string]uuid.UUID
}

func (gre *GHCRateEngineImporter) runImports(appCtx appcontext.AppContext) error {
	importers := []struct {
		importFunction func(appcontext.AppContext) error
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
		cfg := yacspin.Config{
			Frequency:         100 * time.Millisecond,
			CharSet:           yacspin.CharSets[11],
			Suffix:            " ", // puts a least one space between the animating spinner and the Message
			Message:           importer.action,
			SuffixAutoColon:   true,
			ColorAll:          true,
			Colors:            []string{"fgYellow"},
			StopCharacter:     "✓",
			StopColors:        []string{"fgGreen"},
			StopMessage:       "SUCCESS",
			StopFailCharacter: "✗",
			StopFailColors:    []string{"fgRed"},
			StopFailMessage:   "FAILED",
		}

		spinner, err := yacspin.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to make spinner from struct: %w", err)
		}

		if err = spinner.Start(); err != nil {
			return fmt.Errorf("failed to start spinner: %w", err)
		}

		/*
			spinner, err := pterm.DefaultSpinner.Start(importer.action)
			if err != nil {
				return fmt.Errorf("failed to create pterm spinner: %w", err)
			}*/

		err = importer.importFunction(appCtx)
		if err != nil {
			stopErr := spinner.StopFail()
			if stopErr != nil {
				return fmt.Errorf("failed to stop spinner: %w", stopErr)
			}

			return fmt.Errorf("importer failed: %s: %w", importer.action, err)
		}

		err = spinner.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop spinner: %w", err)
		}
	}

	return nil
}

// Import runs the import
func (gre *GHCRateEngineImporter) Import(appCtx appcontext.AppContext) error {
	err := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		dbTxError := gre.runImports(txnAppCtx)
		return dbTxError
	})
	if err != nil {
		return fmt.Errorf("transaction failed during GHC Rate Engine Import(): %w", err)
	}
	return nil
}
