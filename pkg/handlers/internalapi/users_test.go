package internalapi

import (
	"github.com/go-openapi/swag"

	officeuser "github.com/transcom/mymove/pkg/services/office_user"

	"github.com/transcom/mymove/pkg/models/roles"

	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUnknownLoggedInUserHandler() {
	unknownUser := testdatagen.MakeStubbedUser(suite.DB())

	req := suite.NewRequestWithContext("GET", "/users/logged_in", nil)
	req = suite.AuthenticateUserRequest(req, unknownUser)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}
	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{handlers.NewHandlerConfig(suite.DB(), suite.TestLogger()), builder}

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
	req := suite.NewRequestWithContext("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	hConfig := handlers.NewHandlerConfig(suite.DB(), suite.TestLogger())
	featureFlag := handlers.FeatureFlag{Name: "requires-access-code", Active: true}
	hConfig.SetFeatureFlag(featureFlag)
	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{hConfig, builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), sm.UserID.String())
	suite.Equal("Joseph", *okResponse.Payload.ServiceMember.FirstName)
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

	req := suite.NewRequestWithContext("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	hConfig := handlers.NewHandlerConfig(suite.DB(), suite.TestLogger())
	featureFlag := handlers.FeatureFlag{Name: "requires-access-code", Active: false}
	hConfig.SetFeatureFlag(featureFlag)
	builder := officeuser.NewOfficeUserFetcherPop()
	handler := ShowLoggedInUserHandler{hConfig, builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), sm.UserID.String())
	suite.Equal("Jane", *okResponse.Payload.ServiceMember.FirstName)
	suite.False(okResponse.Payload.ServiceMember.RequiresAccessCode)
}

func (suite *HandlerSuite) TestServiceMemberNoTransportationOfficeLoggedInUserHandler() {
	firstName := "Joseph"
	sm := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: &firstName,
		},
	})

	// Remove transportation office info from current station
	station := sm.DutyStation
	station.TransportationOfficeID = nil
	suite.MustSave(&station)

	req := suite.NewRequestWithContext("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}
	builder := officeuser.NewOfficeUserFetcherPop()
	handler := ShowLoggedInUserHandler{handlers.NewHandlerConfig(suite.DB(), suite.TestLogger()), builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), sm.UserID.String())
}

func (suite *HandlerSuite) TestServiceMemberNoMovesLoggedInUserHandler() {

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Show: swag.Bool(false),
		},
	})

	req := suite.NewRequestWithContext("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	hConfig := handlers.NewHandlerConfig(suite.DB(), suite.TestLogger())

	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{hConfig, builder}

	response := handler.Handle(params)

	suite.IsType(&userop.ShowLoggedInUserUnauthorized{}, response)

}
