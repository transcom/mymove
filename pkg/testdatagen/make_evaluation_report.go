package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func MakeEvaluationReport(db *pop.Connection, assertions Assertions) models.EvaluationReport {
	officeUser := assertions.OfficeUser
	if isZeroUUID(assertions.OfficeUser.ID) {
		officeUser = MakeOfficeUser(db, assertions)
	}

	move := assertions.Move
	if isZeroUUID(assertions.Move.ID) {
		move = MakeMove(db, assertions)
	}

	reportType := assertions.EvaluationReport.Type
	if reportType == "" {
		// If no report type is specified, default to Counseling
		reportType = models.EvaluationReportTypeCounseling
	}

	evaluationReport := models.EvaluationReport{
		ID:           uuid.Must(uuid.NewV4()),
		Type:         reportType,
		OfficeUserID: officeUser.ID,
		OfficeUser:   officeUser,
		MoveID:       move.ID,
		Move:         move,
	}

	if !isZeroUUID(assertions.MTOShipment.ID) {
		evaluationReport.ShipmentID = &assertions.MTOShipment.ID
		evaluationReport.Shipment = &assertions.MTOShipment
		evaluationReport.Type = models.EvaluationReportTypeShipment
	}
	mergeModels(&evaluationReport, assertions.EvaluationReport)
	mustCreate(db, &evaluationReport, assertions.Stub)

	return evaluationReport
}
