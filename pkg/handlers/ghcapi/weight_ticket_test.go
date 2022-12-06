package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	weightticketops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateWeightTicketHandler() {
	// Reusable objects
	ppmShipmentUpdater := mocks.PPMShipmentUpdater{}
	weightTicketUpdater := weightticket.NewCustomerWeightTicketUpdater(&ppmShipmentUpdater)

	type weightTicketUpdateSubtestData struct {
		ppmShipment  models.PPMShipment
		weightTicket models.WeightTicket
		params       weightticketops.UpdateWeightTicketParams
		handler      UpdateWeightTicketHandler
	}
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData weightTicketUpdateSubtestData) {
		db := appCtx.DB()

		// Use fake data:
		subtestData.weightTicket = testdatagen.MakeWeightTicket(db, testdatagen.Assertions{})
		subtestData.ppmShipment = subtestData.weightTicket.PPMShipment
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		eTag := etag.GenerateEtag(subtestData.weightTicket.UpdatedAt)

		subtestData.params = weightticketops.UpdateWeightTicketParams{
			HTTPRequest:    req,
			PpmShipmentID:  *handlers.FmtUUID(subtestData.ppmShipment.ID),
			WeightTicketID: *handlers.FmtUUID(subtestData.weightTicket.ID),
			IfMatch:        eTag,
		}

		subtestData.handler = UpdateWeightTicketHandler{
			suite.createS3HandlerConfig(),
			weightTicketUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		// An upload must exist if trailer is owned and qualifies to be claimed
		testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &subtestData.weightTicket.ProofOfTrailerOwnershipDocumentID,
				Document:   subtestData.weightTicket.ProofOfTrailerOwnershipDocument,
			},
		})

		// Add vehicleDescription
		params.UpdateWeightTicketPayload = &ghcmessages.UpdateWeightTicket{
			EmptyWeight: handlers.FmtInt64(1),
			FullWeight:  handlers.FmtInt64(4000),
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketOK{}, response)

		updatedWeightTicket := response.(*weightticketops.UpdateWeightTicketOK).Payload
		suite.NoError(updatedWeightTicket.Validate(strfmt.Default))
		suite.Equal(subtestData.weightTicket.ID.String(), updatedWeightTicket.ID.String())
		suite.Equal(params.UpdateWeightTicketPayload.FullWeight, updatedWeightTicket.FullWeight)
		suite.Equal(params.UpdateWeightTicketPayload.EmptyWeight, updatedWeightTicket.EmptyWeight)
	})
}
