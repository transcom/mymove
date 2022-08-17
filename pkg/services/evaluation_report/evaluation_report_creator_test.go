package evaluationreport

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EvaluationReportSuite) TestEvaluationReportCreator() {
	creator := NewEvaluationReportCreator()

	suite.Run("Can create customer support report successfully", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		report := &models.EvaluationReport{ShipmentID: &shipment.ID, Type: models.EvaluationReportTypeShipment, OfficeUserID: officeUser.ID}
		createdEvaluationReport, err := creator.CreateEvaluationReport(suite.AppContextForTest(), report, move.Locator)

		suite.Nil(err)
		suite.NotNil(createdEvaluationReport)
		suite.NotNil(createdEvaluationReport.MoveID)
		suite.Equal(createdEvaluationReport.OfficeUserID, officeUser.ID)
		suite.Equal(createdEvaluationReport.Type, models.EvaluationReportTypeShipment)
		suite.Equal(*createdEvaluationReport.ShipmentID, shipment.ID)
		suite.NotNil(createdEvaluationReport.OfficeUserID)
		suite.NotNil(createdEvaluationReport.ShipmentID)
		suite.NotNil(createdEvaluationReport.CreatedAt)
		suite.NotNil(createdEvaluationReport.UpdatedAt)
		suite.Nil(createdEvaluationReport.DeletedAt)
		suite.Nil(createdEvaluationReport.SubmittedAt)
	})

	suite.Run("Shipment evaluation report requires valid shipmnet", func() {

		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		move := testdatagen.MakeDefaultMove(suite.DB())
		badID := uuid.Must(uuid.NewV4())
		report := &models.EvaluationReport{ShipmentID: &badID, Type: models.EvaluationReportTypeShipment, OfficeUserID: officeUser.ID}
		createdEvaluationReport, err := creator.CreateEvaluationReport(suite.AppContextForTest(), report, move.Locator)

		suite.Error(err)
		suite.Nil(createdEvaluationReport)
		suite.Equal("FETCH_NOT_FOUND", err.Error())
	})
}
