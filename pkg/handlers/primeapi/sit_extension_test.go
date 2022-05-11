package primeapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	sitextensionservice "github.com/transcom/mymove/pkg/services/sit_extension"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) CreateSITExtensionHandler() {
	// Make an available move
	move := testdatagen.MakeAvailableMove(suite.DB())

	// Make a shipment on the available MTO
	mtoShipment1 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

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

	// Create handler
	handler := CreateSITExtensionHandler{
		handlers.NewHandlerContext(suite.DB(), suite.Logger()),
		sitextensionservice.NewSitExtensionCreator(moveRouter),
	}

	suite.T().Run("Success 201 - Creat SIT extension", func(t *testing.T) {
		// Testcase:   sitExtension is created
		// Expected:   Success response 201

		// Create request params
		req := httptest.NewRequest("POST", fmt.Sprintf("/mto-shipments/%s/sit-extensions", mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.CreateSITExtensionParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(mtoShipment1.ID),
			Body:          sitExtension,
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)

		// Check response type
		suite.IsType(&mtoshipmentops.CreateSITExtensionCreated{}, response)

		// Check values
		sitExtensionResponse := response.(*mtoshipmentops.CreateSITExtensionCreated).Payload

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

	suite.T().Run("Failure 422 - Shipment not found, invalid parameter", func(t *testing.T) {
		// Testcase:   Shipment ID is not found
		// Expected:   Success response 422

		// Update with verification reason\
		badID, _ := uuid.NewV4()
		req := httptest.NewRequest("POST", fmt.Sprintf("/mto-shipments/%s/sit-extensions", mtoShipment1.ID.String()), nil)
		params := mtoshipmentops.CreateSITExtensionParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(badID),
			Body:          sitExtension,
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)

		// Check response type
		suite.IsType(&mtoshipmentops.CreateSITExtensionUnprocessableEntity{}, response)
	})

}
