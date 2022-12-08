package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	movingexpenseservice "github.com/transcom/mymove/pkg/services/moving_expense"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestUpdateMovingExpenseHandler() {
	// Create Reusable objects
	movingExpenseUpdater := movingexpenseservice.NewMovingExpenseUpdater()

	type movingExpenseUpdateSubtestData struct {
		ppmShipment   models.PPMShipment
		movingExpense models.MovingExpense
		params        movingexpenseops.UpdateMovingExpenseParams
		handler       UpdateMovingExpenseHandler
	}
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData movingExpenseUpdateSubtestData) {
		db := appCtx.DB()

		// Fake data:
		subtestData.movingExpense = testdatagen.MakeMovingExpense(db, testdatagen.Assertions{})
		subtestData.ppmShipment = subtestData.movingExpense.PPMShipment

		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense/%s", subtestData.ppmShipment.ID.String(), subtestData.movingExpense.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		eTag := etag.GenerateEtag(subtestData.movingExpense.UpdatedAt)

		subtestData.params = movingexpenseops.UpdateMovingExpenseParams{
			HTTPRequest:     req,
			PpmShipmentID:   *handlers.FmtUUID(subtestData.ppmShipment.ID),
			MovingExpenseID: *handlers.FmtUUID(subtestData.movingExpense.ID),
			IfMatch:         eTag,
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = UpdateMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			movingExpenseUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Moving Expense - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params
		amount := unit.Cents(5000)

		params.UpdateMovingExpense = &ghcmessages.UpdateMovingExpense{
			Amount: *handlers.FmtCost(&amount),
		}

		// Validate incoming payload: no body to validate

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response)
		updatedMovingExpense := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedMovingExpense.Validate(strfmt.Default))

		suite.Equal(subtestData.movingExpense.ID.String(), updatedMovingExpense.ID.String())
		suite.Equal(params.UpdateMovingExpense.Amount, *updatedMovingExpense.Amount)
	})
}
