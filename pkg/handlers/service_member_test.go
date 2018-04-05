package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth/context"
	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/models"
)

func Test_ServiceMember(t *testing.T) {
	fmt.Println("testing233")
}

func (suite *HandlerSuite) TestShowServiceMemberHandler() {
	t := suite.T()

	// Given: A servicemember and a user
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	newServiceMember := models.ServiceMember{
		UserID: user.ID,
	}
	suite.mustSave(&newServiceMember)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	req = req.WithContext(ctx)

	params := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
	}
	// And: show ServiceMember is queried
	showHandler := ShowServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*servicememberop.ShowServiceMemberOK)
	servicemember := okResponse.Payload

	// And: Returned query to include our added servicemember
	if servicemember.UserID.String() != user.ID.String() {
		t.Errorf("Expected an servicemember to have user ID '%v'. None do.", user.ID)
	}

}

func (suite *HandlerSuite) TestShowServiceMemberHandlerNoUser() {
	t := suite.T()

	// Given: A servicemember with a user that isn't logged in
	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	servicemember := models.ServiceMember{
		UserID: user.ID,
	}
	suite.mustSave(&servicemember)

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	showServiceMemberParams := servicememberop.NewShowServiceMemberParams()
	showServiceMemberParams.HTTPRequest = req

	// And: Show servicemember is queried
	showHandler := ShowServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(showServiceMemberParams)

	// Then: Expect a 401 unauthorized
	_, ok := showResponse.(*servicememberop.ShowServiceMemberUnauthorized)
	if !ok {
		t.Errorf("Expected to get an unauthorized response, but got something else.")
	}
}

func (suite *HandlerSuite) TestShowServiceMemberWrongUser() {
	t := suite.T()

	// Given: A servicemember with a not-logged-in user and a separate logged-in user
	notLoggedInUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&notLoggedInUser)

	loggedInUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&loggedInUser)

	// When: A servicemember is created for not-logged-in-user
	newServiceMember := models.ServiceMember{
		UserID: notLoggedInUser.ID,
	}
	suite.mustSave(&newServiceMember)

	// And: the context contains the auth values for logged-in user
	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, loggedInUser.ID, "fake token")
	req = req.WithContext(ctx)
	showServiceMemberParams := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
	}
	// And: Show servicemember is queried
	showHandler := ShowServiceMemberHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(showServiceMemberParams)

	_, ok := showResponse.(*servicememberop.ShowServiceMemberForbidden)
	if !ok {
		t.Fatalf("Request failed: %#v", showResponse)
	}
}
