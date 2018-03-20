package handlers

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/satori/go.uuid"

	authctx "github.com/transcom/mymove/pkg/auth/context"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestCreatePPMHandler() {
	t := suite.T()

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}

	user2 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "someoneelse@example.com",
	}

	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}
	verrs, err = suite.db.ValidateAndCreate(&user2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: internalmessages.SelectedMoveTypeCOMBO,
	}
	verrs, err = suite.db.ValidateAndCreate(&move)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	fmt.Println(user1.ID, user2.ID, move.UserID)

	request := httptest.NewRequest("POST", "/fake/path", nil)
	ctx := request.Context()
	ctx = authctx.PopulateAuthContext(ctx, user1.ID, "faketoken")
	request = request.WithContext(ctx)

	newPPMPayload := internalmessages.CreatePersonallyProcuredMovePayload{WeightEstimate: swag.Int64(12)}

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

	fmt.Println(createdIssuePayload)

	// Next try the wrong user
	ctx = authctx.PopulateAuthContext(ctx, user2.ID, "faketoken")
	request = request.WithContext(ctx)
	newPPMParams.HTTPRequest = request

	badUserResponse := handler.Handle(newPPMParams)
	// assert we got unauthorized
	notAuthorizedResponse := badUserResponse.(*ppmop.CreatePersonallyProcuredMoveForbidden)

	fmt.Println(notAuthorizedResponse)

	// Now try a bad move
	newPPMParams.MoveID = strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	badMoveResponse := handler.Handle(newPPMParams)
	// assert we get the 404
	notFoundResponse := badMoveResponse.(*ppmop.CreatePersonallyProcuredMoveNotFound)

	fmt.Println(notFoundResponse)

}

func (suite *HandlerSuite) TestIndexPPMHandler() {
	t := suite.T()

	user1 := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}

	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	move := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: internalmessages.SelectedMoveTypeCOMBO,
	}
	verrs, err = suite.db.ValidateAndCreate(&move)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	move2 := models.Move{
		UserID:           user1.ID,
		SelectedMoveType: internalmessages.SelectedMoveTypeCOMBO,
	}
	verrs, err = suite.db.ValidateAndCreate(&move2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	ppm1 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		WeightEstimate: swag.Int64(1),
	}
	ppm2 := models.PersonallyProcuredMove{
		MoveID:         move.ID,
		WeightEstimate: swag.Int64(2),
	}
	otherPPM := models.PersonallyProcuredMove{
		MoveID:         move2.ID,
		WeightEstimate: swag.Int64(4),
	}

	verrs, err = suite.db.ValidateAndCreate(&ppm1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	verrs, err = suite.db.ValidateAndCreate(&ppm2)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	verrs, err = suite.db.ValidateAndCreate(&otherPPM)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	request := httptest.NewRequest("GET", "/fake/path", nil)
	ctx := request.Context()
	ctx = authctx.PopulateAuthContext(ctx, user1.ID, "faketoken")
	request = request.WithContext(ctx)

	indexPPMParams := ppmop.IndexPersonallyProcuredMovesParams{
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: request,
	}

	handler := IndexPersonallyProcuredMovesHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(indexPPMParams)

	// assert we got back the 201 response
	okResponse := response.(*ppmop.IndexPersonallyProcuredMovesOK)
	indexPPMPayload := okResponse.Payload

	fmt.Println(indexPPMPayload)
	for _, ppm := range indexPPMPayload {
		if *ppm.ID == *fmtUUID(otherPPM.ID) {
			t.Error("We should only have got back ppms associated with this move")
		}
	}
	if len(indexPPMPayload) != 2 {
		t.Error("We should have gotten back two good ones. ")
	}

}
