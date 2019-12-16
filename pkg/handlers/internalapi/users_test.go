package internalapi

import (
	"net/http/httptest"

	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUnknownLoggedInUserHandler() {
	unknownUser := testdatagen.MakeDefaultUser(suite.DB())

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateUserRequest(req, unknownUser)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	handler := ShowLoggedInUserHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), unknownUser.ID.String())
}

func (suite *HandlerSuite) TestServiceMemberLoggedInUserRequiringAccessCodeHandler() {
	firstName := "Joseph"
	smRole := models.Role{
		RoleType: "customer",
	}
	suite.NoError(suite.DB().Save(&smRole))
	sm := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:          &firstName,
			RequiresAccessCode: true,
		},
		User: models.User{
			Roles: []models.Role{smRole},
		},
	})
	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	featureFlag := handlers.FeatureFlag{Name: "requires-access-code", Active: true}
	context.SetFeatureFlag(featureFlag)
	handler := ShowLoggedInUserHandler{context}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), sm.UserID.String())
	suite.Equal("Joseph", *okResponse.Payload.ServiceMember.FirstName)
	suite.Equal("customer", *okResponse.Payload.Roles[0].RoleType)
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

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	featureFlag := handlers.FeatureFlag{Name: "requires-access-code", Active: false}
	context.SetFeatureFlag(featureFlag)
	handler := ShowLoggedInUserHandler{context}

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

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	handler := ShowLoggedInUserHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), sm.UserID.String())
}
