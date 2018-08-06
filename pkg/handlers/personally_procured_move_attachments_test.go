package handlers

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestCreatePPMHandler() {
	user1 := testdatagen.MakeDefaultServiceMember(suite.db)
	orders := testdatagen.MakeDefaultOrder(suite.db)
	var selectedType = internalmessages.SelectedMoveTypeCOMBO

	move, verrs, locErr := orders.CreateNewMove(suite.db, &selectedType)
	suite.False(verrs.HasAny(), "failed to create new move")
	suite.Nil(locErr)

	request := httptest.NewRequest("POST", "/fake/path", nil)
	request = suite.authenticateRequest(request, orders.ServiceMember)

	newPPMPayload := internalmessages.CreatePersonallyProcuredMovePayload{
		WeightEstimate:   swag.Int64(12),
		PickupPostalCode: swag.String("00112"),
		DaysInStorage:    swag.Int64(3),
	}

	newPPMParams := ppmop.CreatePersonallyProcuredMoveParams{
		MoveID: strfmt.UUID(move.ID.String()),
		CreatePersonallyProcuredMovePayload: &newPPMPayload,
		HTTPRequest:                         request,
	}

	handler := CreatePersonallyProcuredMoveHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(newPPMParams)
	// assert we got back the 201 response
	createdResponse := response.(*ppmop.CreatePersonallyProcuredMoveCreated)
	createdIssuePayload := createdResponse.Payload
	suite.NotNil(createdIssuePayload.ID)

	// Next try the wrong user
	request = suite.authenticateRequest(request, user1)
	newPPMParams.HTTPRequest = request

	badUserResponse := handler.Handle(newPPMParams)
	suite.checkResponseForbidden(badUserResponse)

	// Now try a bad move
	newPPMParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newPPMParams)
	suite.checkResponseNotFound(badMoveResponse)

}
