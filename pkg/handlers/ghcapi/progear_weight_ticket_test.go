package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	progearops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	progear "github.com/transcom/mymove/pkg/services/progear_weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateProGearWeightTicketHandler() {
	// Reusable objects
	progearUpdater := progear.NewCustomerProgearWeightTicketUpdater()

	type progearUpdateSubtestData struct {
		ppmShipment models.PPMShipment
		progear     models.ProgearWeightTicket
		params      progearops.UpdateProGearWeightTicketParams
		handler     UpdateProgearWeightTicketHandler
	}
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, _ bool) (subtestData progearUpdateSubtestData) {
		db := appCtx.DB()

		// Use fake data:
		subtestData.progear = factory.BuildProgearWeightTicket(db, nil, nil)
		subtestData.ppmShipment = subtestData.progear.PPMShipment

		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets/%s", subtestData.ppmShipment.ID.String(), subtestData.progear.ID.String())
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		req := httptest.NewRequest("PATCH", endpoint, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		eTag := etag.GenerateEtag(subtestData.progear.UpdatedAt)

		subtestData.params = progearops.UpdateProGearWeightTicketParams{
			HTTPRequest:           req,
			PpmShipmentID:         *handlers.FmtUUID(subtestData.ppmShipment.ID),
			ProGearWeightTicketID: *handlers.FmtUUID(subtestData.progear.ID),
			IfMatch:               eTag,
		}

		subtestData.handler = UpdateProgearWeightTicketHandler{
			suite.createS3HandlerConfig(),
			progearUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

		params := subtestData.params

		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.progear.Document,
				LinkOnly: true,
			},
		}, nil)

		hasWeightTickets := true
		belongsToSelf := true
		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{
			HasWeightTickets: hasWeightTickets,
			Weight:           handlers.FmtInt64(4000),
			BelongsToSelf:    belongsToSelf,
		}

		// Validate incoming payload: no body to validate
		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketOK{}, response)

		updatedProgear := response.(*progearops.UpdateProGearWeightTicketOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedProgear.Validate(strfmt.Default))

		suite.Equal(subtestData.progear.ID.String(), updatedProgear.ID.String())
		suite.Equal(params.UpdateProGearWeightTicket.Weight, updatedProgear.Weight)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{}
		wrongUUIDString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("3ce0e367-a337-46e3-b4cf-f79aebc4f6c8"))
		params.ProGearWeightTicketID = *wrongUUIDString

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{}
		params.IfMatch = "wrong-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.ProgearWeightTicketUpdater{}
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		ownsTrailer := true
		hasWeightTickets := true

		params.UpdateProGearWeightTicket = &ghcmessages.UpdateProGearWeightTicket{
			Weight:           handlers.FmtInt64(1000),
			BelongsToSelf:    ownsTrailer,
			HasWeightTickets: hasWeightTickets,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateProgearWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.ProgearWeightTicket"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		handler := UpdateProgearWeightTicketHandler{
			suite.HandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&progearops.UpdateProGearWeightTicketInternalServerError{}, response)
	})
}
