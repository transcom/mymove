package handlers

import (
	"net/http/httptest"

	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth/context"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUnknownLoggedInUserHandler() {
	t := suite.T()

	unknownUser, err := testdatagen.MakeUser(suite.db)
	if err != nil {
		t.Fatal("couldn't create a user")
	}

	params := userop.NewShowLoggedInUserParams()
	req := httptest.NewRequest("GET", "/users/logged_in", nil)

	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, unknownUser.ID, "fake token")

	params.HTTPRequest = req.WithContext(ctx)

	handler := ShowLoggedInUserHandler(NewHandlerContext(suite.db, suite.logger))

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	if *okResponse.Payload.ID != *fmtUUID(unknownUser.ID) || okResponse.Payload.Type != "UNKNOWN" {
		t.Fatalf("Didn't get back what we wanted. %#v", okResponse)
	}

}

func (suite *HandlerSuite) TestServiceMemberLoggedInUserHandler() {
	t := suite.T()

	smUser, err := testdatagen.MakeUser(suite.db)
	if err != nil {
		t.Fatal("couldn't create a user")
	}

	firstName := "Joseph"
	serviceMember := models.ServiceMember{
		UserID:    smUser.ID,
		User:      smUser,
		FirstName: &firstName,
	}

	verrs, err := models.CreateServiceMemberWithAddresses(suite.db, &smUser, &serviceMember)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
		t.Fatal("Couldnt create theSM")
	}

	params := userop.NewShowLoggedInUserParams()
	req := httptest.NewRequest("GET", "/users/logged_in", nil)

	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, smUser.ID, "fake token")

	params.HTTPRequest = req.WithContext(ctx)

	handler := ShowLoggedInUserHandler(NewHandlerContext(suite.db, suite.logger))

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	if *okResponse.Payload.ID != *fmtUUID(smUser.ID) || okResponse.Payload.Type != "SERVICE_MEMBER" {
		t.Fatalf("Didn't get back what we wanted. %#v", okResponse.Payload)
	}

	if *okResponse.Payload.ServiceMember.FirstName != "Joseph" {
		t.Fatalf("Didn't get the SM right. %#v", okResponse.Payload.ServiceMember)
	}

}

func (suite *HandlerSuite) TestNoLoggedInUserHandler() {
	t := suite.T()

	_, err := testdatagen.MakeUser(suite.db)
	if err != nil {
		t.Fatal("couldn't create a user")
	}

	badUUID := uuid.Must(uuid.NewV4())

	params := userop.NewShowLoggedInUserParams()
	req := httptest.NewRequest("GET", "/users/logged_in", nil)

	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, badUUID, "fake token")

	params.HTTPRequest = req.WithContext(ctx)

	handler := ShowLoggedInUserHandler(NewHandlerContext(suite.db, suite.logger))

	response := handler.Handle(params)

	_, ok := response.(*userop.ShowLoggedInUserUnauthorized)
	if !ok {
		t.Fatalf("Request should not have succeeded: %#v", response)
	}

}
