package publicapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/mocks"

	accesscodeops "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accesscode"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestValidateAccessCodeHandler_Valid() {
	// create user
	user := testdatagen.MakeDefaultUser(suite.DB())
	selectedMoveType := models.SelectedMoveTypePPM

	// creates access code
	code := "TEST1"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: &selectedMoveType,
	}
	fullCode := fmt.Sprintf("%s-%s", selectedMoveType, code)
	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/accesscode/valid?code=%s", fullCode), nil)
	request = suite.AuthenticateUserRequest(request, user)

	params := accesscodeops.ValidateAccessCodeParams{
		HTTPRequest: request,
		Code:        &fullCode,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	accessCodeValidator := &mocks.AccessCodeValidator{}
	accessCodeValidator.On("ValidateAccessCode",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("models.SelectedMoveType"),
	).Return(&accessCode, true, nil)

	handler := ValidateAccessCodeHandler{context, accessCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validateAccessCodeResponse := response.(*accesscodeops.ValidateAccessCodeOK)
	validateAccessCodePayload := validateAccessCodeResponse.Payload

	suite.NotNil(validateAccessCodePayload.AccessCode)
	suite.True(*validateAccessCodePayload.Valid)
	suite.Assertions.IsType(&accesscodeops.ValidateAccessCodeOK{}, response)
}

func (suite *HandlerSuite) TestValidateAccessCodeHandler_Invalid() {
	// create user
	user := testdatagen.MakeDefaultUser(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHG
	smID, _ := uuid.NewV4()

	// creates access code
	code := "TEST2"
	claimedTime := time.Now()
	invalidAccessCode := models.AccessCode{
		Code:            code,
		MoveType:        &selectedMoveType,
		ServiceMemberID: &smID,
		ClaimedAt:       &claimedTime,
	}
	fullCode := fmt.Sprintf("%s-%s", selectedMoveType, code)

	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/accesscode/valid?code=%s", fullCode), nil)
	request = suite.AuthenticateUserRequest(request, user)

	params := accesscodeops.ValidateAccessCodeParams{
		HTTPRequest: request,
		Code:        &fullCode,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	accessCodeValidator := &mocks.AccessCodeValidator{}
	accessCodeValidator.On("ValidateAccessCode",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("models.SelectedMoveType"),
	).Return(&invalidAccessCode, false, nil)

	handler := ValidateAccessCodeHandler{context, accessCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validateAccessCodeResponse := response.(*accesscodeops.ValidateAccessCodeOK)
	validateAccessCodePayload := validateAccessCodeResponse.Payload

	suite.False(*validateAccessCodePayload.Valid)
	suite.Assertions.IsType(&accesscodeops.ValidateAccessCodeOK{}, response)
}
