package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildEvaluationReport creates an unsubmitted EvaluationReport.
// Also creates, if not provided
// - OfficeUser
// - Move
//
// Notes:
// if evaluation report Type == models.EvaluationReportTypeShipment, it also creates
//   - MTOShipment
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildEvaluationReport(db *pop.Connection, customs []Customization, traits []Trait) models.EvaluationReport {
	customs = setupCustomizations(customs, traits)

	// Find EvaluationReport customization and extract the custom model
	var cEvaluationReport models.EvaluationReport
	if result := findValidCustomization(customs, EvaluationReport); result != nil {
		cEvaluationReport = result.Model.(models.EvaluationReport)
		if result.LinkOnly {
			return cEvaluationReport
		}
	}

	officeUser := BuildOfficeUser(db, customs, traits)
	move := BuildMove(db, customs, traits)

	// create default EvaluationReport
	evaluationReport := models.EvaluationReport{
		Type:         models.EvaluationReportTypeCounseling,
		OfficeUserID: officeUser.ID,
		OfficeUser:   officeUser,
		MoveID:       move.ID,
		Move:         move,
	}

	// If Type is EvaluationReportTypeShipment, Shipment/ShipmentID can't be null
	if cEvaluationReport.Type == models.EvaluationReportTypeShipment {
		// If a shipment needs to be created, use the move created above as a LinkOnly customization
		if db != nil {
			// can only do LinkOnly if we have an ID, which we won't have
			// for a stubbed evaluation report
			customs = replaceCustomization(customs, Customization{
				Model:    move,
				LinkOnly: true,
			})
		}

		shipment := BuildMTOShipment(db, customs, nil)

		evaluationReport.Shipment = &shipment
		evaluationReport.ShipmentID = &shipment.ID
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&evaluationReport, cEvaluationReport)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &evaluationReport)
	}

	return evaluationReport
}
