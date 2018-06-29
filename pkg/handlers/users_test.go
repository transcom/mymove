package handlers

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/auth"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUnknownLoggedInUserHandler() {
	t := suite.T()

	unknownUser := testdatagen.MakeDefaultUser(suite.db)

	params := userop.NewShowLoggedInUserParams()
	req := httptest.NewRequest("GET", "/users/logged_in", nil)

	session := &auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          unknownUser.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, session)

	params.HTTPRequest = req.WithContext(ctx)

	handler := ShowLoggedInUserHandler(NewHandlerContext(suite.db, suite.logger))

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	if *okResponse.Payload.ID != *fmtUUID(unknownUser.ID) {
		t.Fatalf("Didn't get back what we wanted. %#v", okResponse)
	}

}

func (suite *HandlerSuite) TestServiceMemberLoggedInUserHandler() {
	t := suite.T()

	smUser := testdatagen.MakeDefaultUser(suite.db)

	firstName := "Joseph"
	serviceMember := models.ServiceMember{
		UserID:    smUser.ID,
		User:      smUser,
		FirstName: &firstName,
	}

	verrs, err := models.SaveServiceMember(suite.db, &serviceMember)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
		t.Fatal("Couldnt create theSM")
	}

	params := userop.NewShowLoggedInUserParams()
	req := httptest.NewRequest("GET", "/users/logged_in", nil)

	session := &auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          smUser.ID,
		ServiceMemberID: serviceMember.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, session)

	params.HTTPRequest = req.WithContext(ctx)

	handler := ShowLoggedInUserHandler(NewHandlerContext(suite.db, suite.logger))

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	if *okResponse.Payload.ID != *fmtUUID(smUser.ID) {
		t.Fatalf("Didn't get back what we wanted. %#v", okResponse.Payload)
	}

	if *okResponse.Payload.ServiceMember.FirstName != "Joseph" {
		t.Fatalf("Didn't get the SM right. %#v", okResponse.Payload.ServiceMember)
	}

}
