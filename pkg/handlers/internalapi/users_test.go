package internalapi

import (
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/swag"

	officeuser "github.com/transcom/mymove/pkg/services/office_user"

	"github.com/transcom/mymove/pkg/models/roles"

	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUnknownLoggedInUserHandler() {
	unknownUser := testdatagen.MakeStubbedUser(suite.DB())

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateUserRequest(req, unknownUser)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}
	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), unknownUser.ID.String())
}

func (suite *HandlerSuite) TestServiceMemberLoggedInUserRequiringAccessCodeHandler() {
	firstName := "Joseph"
	smRole := roles.Role{
		RoleType: roles.RoleTypeCustomer,
	}

	user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
		User: models.User{
			Roles: []roles.Role{smRole},
		},
	})

	suite.NoError(suite.DB().Save(&smRole))
	sm := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:          &firstName,
			RequiresAccessCode: true,
			UserID:             user.ID,
			User:               user,
		},
	})
	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	featureFlag := handlers.FeatureFlag{Name: "requires-access-code", Active: true}
	context.SetFeatureFlag(featureFlag)
	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{context, builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), sm.UserID.String())
	suite.Equal(firstName, *okResponse.Payload.ServiceMember.FirstName)
	suite.Equal(string(roles.RoleTypeCustomer), *okResponse.Payload.Roles[0].RoleType)
	suite.Equal(1, len(okResponse.Payload.Roles))
	suite.True(okResponse.Payload.ServiceMember.RequiresAccessCode)
}

func (suite *HandlerSuite) TestServiceMemberLoggedInUserNotRequiringAccessCodeHandler() {
	firstName := "Jane"
	sm := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: &firstName,
		},
	})

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	featureFlag := handlers.FeatureFlag{Name: "requires-access-code", Active: false}
	context.SetFeatureFlag(featureFlag)
	builder := officeuser.NewOfficeUserFetcherPop()
	handler := ShowLoggedInUserHandler{context, builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(sm.UserID.String(), okResponse.Payload.ID.String())
	suite.Equal(firstName, *okResponse.Payload.ServiceMember.FirstName)
	suite.False(okResponse.Payload.ServiceMember.RequiresAccessCode)
}

func (suite *HandlerSuite) TestServiceMemberNoTransportationOfficeLoggedInUserHandler() {
	suite.T().Run("current duty location missing", func(t *testing.T) {
		sm := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{})

		// Remove transportation office info from current duty location
		dutyLocation := sm.DutyLocation
		dutyLocation.TransportationOfficeID = nil
		suite.MustSave(&dutyLocation)

		req := httptest.NewRequest("GET", "/users/logged_in", nil)
		req = suite.AuthenticateRequest(req, sm)

		params := userop.ShowLoggedInUserParams{
			HTTPRequest: req,
		}
		builder := officeuser.NewOfficeUserFetcherPop()
		handler := ShowLoggedInUserHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), builder}

		response := handler.Handle(params)

		okResponse, ok := response.(*userop.ShowLoggedInUserOK)
		suite.True(ok)
		suite.Equal(sm.UserID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("new duty location missing", func(t *testing.T) {
		// add orders
		order := testdatagen.MakeOrderWithoutDefaults(suite.DB(), testdatagen.Assertions{})

		sm := order.ServiceMember

		// Remove transportation office info from new duty location
		// happens when a customer is not done
		dutyLocation := order.NewDutyLocation
		dutyLocation.TransportationOfficeID = nil
		suite.MustSave(&dutyLocation)

		req := httptest.NewRequest("GET", "/users/logged_in", nil)
		req = suite.AuthenticateRequest(req, sm)

		params := userop.ShowLoggedInUserParams{
			HTTPRequest: req,
		}
		fakeS3 := storageTest.NewFakeS3Storage(true)
		builder := officeuser.NewOfficeUserFetcherPop()
		context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		context.SetFileStorer(fakeS3)
		handler := ShowLoggedInUserHandler{context, builder}

		response := handler.Handle(params)

		okResponse, ok := response.(*userop.ShowLoggedInUserOK)
		suite.True(ok, "Response should be ok")
		suite.NotNil(okResponse, "Response should not be nil")
		suite.Equal(sm.UserID.String(), okResponse.Payload.ID.String())
	})
}

func (suite *HandlerSuite) TestServiceMemberNoMovesLoggedInUserHandler() {

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Show: swag.Bool(false),
		},
	})

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())

	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{context, builder}

	response := handler.Handle(params)

	suite.IsType(&userop.ShowLoggedInUserUnauthorized{}, response)

}
