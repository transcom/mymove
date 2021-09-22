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

	"github.com/transcom/mymove/pkg/appcontext"
	moverouter "github.com/transcom/mymove/pkg/services/move"

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
func (e *devSeedScenario) Setup(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	db := appCtx.DB()
	logger := appCtx.Logger()
	moveRouter := moverouter.NewMoveRouter()

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
		"additional_ppm_users":         subScenarioAdditionalPPMUsers(appCtx, userUploader),
		"diverted_shipments":           subScenarioDivertedShipments(appCtx, userUploader, allDutyStations, originDutyStationsInGBLOC),
		"hhg_onboarding":               subScenarioHHGOnboarding(appCtx, userUploader),
		"hhg_services_counseling":      subScenarioHHGServicesCounseling(appCtx, userUploader, allDutyStations, originDutyStationsInGBLOC),
		"payment_request_calculations": subScenarioPaymentRequestCalculations(appCtx, userUploader, primeUploader, moveRouter),
		"ppm_and_hhg":                  subScenarioPPMAndHHG(appCtx, userUploader, moveRouter),
		"ppm_office_queue":             subScenarioPPMOfficeQueue(appCtx, userUploader, moveRouter),
		"shipment_hhg_cancelled":       subScenarioShipmentHHGCancelled(appCtx, allDutyStations, originDutyStationsInGBLOC),
		"txo_queues":                   subScenarioTXOQueues(appCtx, userUploader, logger),
		"misc":                         subScenarioMisc(appCtx, userUploader, primeUploader, moveRouter),
		"reweighs":                     subScenarioReweighs(appCtx, userUploader, primeUploader, moveRouter),
		"sit_extensions":               subScenarioSITExtensions(appCtx, userUploader, primeUploader),
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
