package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestListReports() {
	moveId, _ := uuid.NewV4()
	counselingCompleted := time.Now()

	testMove := models.Move{
		ID:                           moveId,
		ServiceCounselingCompletedAt: &counselingCompleted,
	}

	suite.Run("Success - Returns a basic report", func() {
		report := ListReport(suite.AppContextForTest(), &testMove)

		suite.IsType(&pptasmessages.ListReport{}, report)
		suite.Equal(strfmt.UUID(testMove.ID.String()), report.ShipmentID)
	})
}
