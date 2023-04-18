package evaluationreport

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *EvaluationReportSuite) TestEvaluationReportCreator() {
	creator := NewEvaluationReportCreator()

	suite.Run("Can create customer support report successfully", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
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
		suite.Nil(createdEvaluationReport.SubmittedAt)
	})

	suite.Run("Shipment evaluation report requires valid shipmnet", func() {

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		move := factory.BuildMove(suite.DB(), nil, nil)
		badID := uuid.Must(uuid.NewV4())
		report := &models.EvaluationReport{ShipmentID: &badID, Type: models.EvaluationReportTypeShipment, OfficeUserID: officeUser.ID}
		createdEvaluationReport, err := creator.CreateEvaluationReport(suite.AppContextForTest(), report, move.Locator)

		suite.Error(err)
		suite.Nil(createdEvaluationReport)
		suite.Equal("FETCH_NOT_FOUND", err.Error())
	})
}
