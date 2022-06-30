package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// TODO not really a real test, just wanted to have an insert/query
func MakeEvaluationReport(db *pop.Connection, assertions Assertions) models.EvaluationReport {
	evaluationReport := models.EvaluationReport{
		ID: uuid.Must(uuid.NewV4()),
	}
	if !isZeroUUID(assertions.MTOShipment.ID) {
		evaluationReport.ShipmentID = &assertions.MTOShipment.ID
		evaluationReport.Shipment = &assertions.MTOShipment
	}
	mergeModels(&evaluationReport, assertions.EvaluationReport)
	mustCreate(db, &evaluationReport, assertions.Stub)

	return evaluationReport
}
