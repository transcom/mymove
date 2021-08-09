//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
// nolint:golint
package scenario

import (
	"go.uber.org/zap"

	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// devSeedScenario builds a basic set of data for e2e testing
type devSeedScenario NamedScenario

// DevSeedScenario setup information for the dev seed
var DevSeedScenario = devSeedScenario{
	Name: "dev_seed",
}

// Setup initializes the run setup for the devseed scenario
func (e *devSeedScenario) Setup(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, logger *zap.Logger) {

	moveRouter := moverouter.NewMoveRouter(db, logger)

	// Testdatagen factories will create new random duty stations so let's get the standard ones in the migrations
	var allDutyStations []models.DutyStation
	db.All(&allDutyStations)

	var originDutyStationsInGBLOC []models.DutyStation
	db.Where("transportation_offices.GBLOC = ?", "LKNQ").
		InnerJoin("transportation_offices", "duty_stations.transportation_office_id = transportation_offices.id").
		All(&originDutyStationsInGBLOC)

	/*
		ADD NEW SUB-SCENARIOS HERE
	*/

	// sets the sub-scenarios
	e.SubScenarios = map[string]func(){
		"additional_ppm_users":         subScenarioAdditionalPPMUsers(db, userUploader),
		"diverted_shipments":           subScenarioDivertedShipments(db, userUploader, allDutyStations, originDutyStationsInGBLOC),
		"hhg_onboarding":               subScenarioHHGOnboarding(db, userUploader),
		"hhg_services_counseling":      subScenarioHHGServicesCounseling(db, userUploader, allDutyStations, originDutyStationsInGBLOC),
		"payment_request_calculations": subScenarioPaymentRequestCalculations(db, userUploader, primeUploader, moveRouter, logger),
		"ppm_and_hhg":                  subScenarioPPMAndHHG(db, userUploader, moveRouter),
		"ppm_office_queue":             subScenarioPPMOfficeQueue(db, userUploader, moveRouter),
		"shipment_hhg_cancelled":       subScenarioShipmentHHGCancelled(db, allDutyStations, originDutyStationsInGBLOC),
		"txo_queues":                   subScenarioTXOQueues(db, userUploader, logger),
		"misc":                         subScenarioMisc(db, userUploader, primeUploader, moveRouter),
		"reweighs":                     subScenarioReweighs(db, userUploader, primeUploader, moveRouter),
	}
}

// Run does that data load thing
func (e *devSeedScenario) Run(logger *zap.Logger, namedSubScenario string) {
	// sub-scenario name validation runs before this part is reached
	// run only the specified sub-scenario
	if subScenarioFunc, ok := e.SubScenarios[namedSubScenario]; ok {
		logger.Info("running sub-scenario: " + namedSubScenario)

		subScenarioFunc()

		logger.Info("done running sub-scenario: " + namedSubScenario)
	} else {
		// otherwise, run through all sub-scenarios
		for name, subScenarioFunc := range e.SubScenarios {
			logger.Info("running sub-scenario: " + name)

			subScenarioFunc()

			logger.Info("done running sub-scenario: " + name)
		}
	}
}
