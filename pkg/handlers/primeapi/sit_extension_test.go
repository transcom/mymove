package primeapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	sitextensionservice "github.com/transcom/mymove/pkg/services/sit_extension"
)

func (suite *HandlerSuite) CreateSITExtensionHandler() {

	// Make sit extension
	daysRequested := int64(30)
	remarks := "We need an extension"
	reason := "AWAITING_COMPLETION_OF_RESIDENCE"

	sitExtension := &primemessages.CreateSITExtension{
		RequestedDays:     &daysRequested,
		ContractorRemarks: &remarks,
		RequestReason:     &reason,
	}

	// Create move router for SitExtension Createor
	moveRouter := moverouter.NewMoveRouter()
	setupTestData := func() (CreateSITExtensionHandler, models.MTOShipment) {

		// Make an available move
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// Make a shipment on the available MTO
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Create handler
		handler := CreateSITExtensionHandler{
			suite.HandlerConfig(),
			sitextensionservice.NewSitExtensionCreator(moveRouter),
		}
		return handler, shipment
	}

	suite.Run("Success 201 - Creat SIT extension", func() {
		// Testcase:   sitExtension is created
		// Expected:   Success response 201
		handler, shipment := setupTestData()
		// Create request params
		req := httptest.NewRequest("POST", fmt.Sprintf("/mto-shipments/%s/sit-extensions", shipment.ID.String()), nil)
		params := mtoshipmentops.CreateSITExtensionParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          sitExtension,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)

		// Check response type
		suite.IsType(&mtoshipmentops.CreateSITExtensionCreated{}, response)
		sitExtensionResponse := response.(*mtoshipmentops.CreateSITExtensionCreated).Payload

		// Validate outgoing payload
		suite.NoError(sitExtensionResponse.Validate(strfmt.Default))

		// Check values
		suite.Equal(daysRequested, sitExtensionResponse.RequestedDays)
		suite.Equal(models.SITExtensionStatusPending, sitExtensionResponse.Status)
		suite.Equal(daysRequested, sitExtensionResponse.RequestedDays)
		suite.Equal(models.SITExtensionRequestReasonAwaitingCompletionOfResidence, sitExtensionResponse.RequestReason)
		suite.Equal(remarks, sitExtensionResponse.ContractorRemarks)
		suite.NotNil(sitExtensionResponse.ID)
		suite.NotNil(sitExtensionResponse.CreatedAt)
		suite.NotNil(sitExtensionResponse.UpdatedAt)
		suite.NotNil(sitExtensionResponse.ETag)
	})

	suite.Run("Failure 422 - Shipment not found, invalid parameter", func() {
		// Testcase:   Shipment ID is not found
		// Expected:   Success response 422
		handler, shipment := setupTestData()

		// Update with verification reason\
		badID, _ := uuid.NewV4()
		req := httptest.NewRequest("POST", fmt.Sprintf("/mto-shipments/%s/sit-extensions", shipment.ID.String()), nil)
		params := mtoshipmentops.CreateSITExtensionParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(badID),
			Body:          sitExtension,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)

		// Check response type
		suite.IsType(&mtoshipmentops.CreateSITExtensionUnprocessableEntity{}, response)
		responsePayload := response.(*mtoshipmentops.CreateSITExtensionUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}
