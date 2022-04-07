package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/cli"
	accesscodeops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/accesscode"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestFetchAccessCodeHandler_Success() {
	// create user
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHG

	// creates access code
	code := "TEST0"
	accessCode := models.AccessCode{
		Code:            code,
		MoveType:        selectedMoveType,
		ServiceMemberID: &serviceMember.ID,
	}

	// makes request
	request := httptest.NewRequest("GET", "/access_codes", nil)
	request = suite.AuthenticateRequest(request, serviceMember)

	params := accesscodeops.FetchAccessCodeParams{
		HTTPRequest: request,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	context.SetFeatureFlag(
		handlers.FeatureFlag{Name: cli.FeatureFlagAccessCode, Active: true},
	)
	accessCodeFetcher := &mocks.AccessCodeFetcher{}
	accessCodeFetcher.On("FetchAccessCode",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&accessCode, nil)

	handler := FetchAccessCodeHandler{context, accessCodeFetcher}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	fetchAccessCodeResponse := response.(*accesscodeops.FetchAccessCodeOK)
	fetchAccessCodePayload := fetchAccessCodeResponse.Payload

	suite.NotNil(fetchAccessCodePayload)
	suite.Assertions.IsType(&accesscodeops.FetchAccessCodeOK{}, response)
	suite.Equal(*fetchAccessCodePayload.Code, code)
}

func (suite *HandlerSuite) TestFetchAccessCodeHandler_Failure() {
	// create user
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	// makes request
	request := httptest.NewRequest("GET", "/access_codes", nil)
	request = suite.AuthenticateRequest(request, serviceMember)

	params := accesscodeops.FetchAccessCodeParams{
		HTTPRequest: request,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	context.SetFeatureFlag(
		handlers.FeatureFlag{Name: cli.FeatureFlagAccessCode, Active: true},
	)
	accessCodeFetcher := &mocks.AccessCodeFetcher{}
	accessCodeFetcher.On("FetchAccessCode",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.AccessCode{}, models.ErrFetchNotFound)

	handler := FetchAccessCodeHandler{context, accessCodeFetcher}
	response := handler.Handle(params)

	fetchAccessCodeResponse := response.(*accesscodeops.FetchAccessCodeNotFound)
	suite.Assertions.IsType(&accesscodeops.FetchAccessCodeNotFound{}, fetchAccessCodeResponse)
}

func (suite *HandlerSuite) TestFetchAccessCodeHandler_FeatureFlagIsOff() {
	// create user
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	// makes request
	request := httptest.NewRequest("GET", "/access_codes", nil)
	request = suite.AuthenticateRequest(request, serviceMember)

	params := accesscodeops.FetchAccessCodeParams{
		HTTPRequest: request,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	context.SetFeatureFlag(
		handlers.FeatureFlag{Name: cli.FeatureFlagAccessCode, Active: false},
	)
	accessCodeFetcher := &mocks.AccessCodeFetcher{}
	handler := FetchAccessCodeHandler{context, accessCodeFetcher}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	suite.Assertions.IsType(&accesscodeops.FetchAccessCodeOK{}, response)
}

func (suite *HandlerSuite) TestClaimAccessCodeHandler_Success() {
	// create user
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHG

	// creates access code
	code := "TEST2"
	claimedAt := time.Now()

	// makes request
	request := httptest.NewRequest("PATCH", "/access_codes/invalid", nil)
	request = suite.AuthenticateRequest(request, serviceMember)

	params := accesscodeops.ClaimAccessCodeParams{
		HTTPRequest: request,
		AccessCode:  accesscodeops.ClaimAccessCodeBody{Code: &code},
	}

	claimedAccessCode := models.AccessCode{
		Code:            code,
		MoveType:        selectedMoveType,
		ClaimedAt:       &claimedAt,
		ServiceMemberID: &serviceMember.ID,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	accessCodeClaimer := &mocks.AccessCodeClaimer{}
	accessCodeClaimer.On("ClaimAccessCode",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&claimedAccessCode, validate.NewErrors(), nil)

	handler := ClaimAccessCodeHandler{context, accessCodeClaimer}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	claimAccessCodeResponse := response.(*accesscodeops.ClaimAccessCodeOK)
	claimAccessCodePayload := claimAccessCodeResponse.Payload

	suite.Assertions.Equal(claimedAccessCode.Code, *claimAccessCodePayload.Code)
	suite.Assertions.Equal(claimedAccessCode.MoveType.String(), *claimAccessCodePayload.MoveType)
	suite.Assertions.Equal(claimAccessCodePayload.ClaimedAt, handlers.FmtDateTime(*claimedAccessCode.ClaimedAt))
	suite.Assertions.Equal(claimAccessCodePayload.ServiceMemberID, *handlers.FmtUUID(*claimedAccessCode.ServiceMemberID))
	suite.Assertions.IsType(&accesscodeops.ClaimAccessCodeOK{}, response)
}
