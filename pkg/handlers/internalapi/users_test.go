package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
)

func (suite *HandlerSuite) TestUnknownLoggedInUserHandler() {
	unknownUser := factory.BuildUser(nil, nil, nil)

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateUserRequest(req, unknownUser)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}
	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{suite.NewHandlerConfig(), builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.NotNil(okResponse.Payload.Permissions)
	suite.Equal(okResponse.Payload.ID.String(), unknownUser.ID.String())
}

func (suite *HandlerSuite) TestServiceMemberNoTransportationOfficeLoggedInUserHandler() {
	suite.Run("current duty location missing", func() {
		sm := factory.BuildExtendedServiceMember(suite.DB(), nil, nil)

		req := httptest.NewRequest("GET", "/users/logged_in", nil)
		req = suite.AuthenticateRequest(req, sm)

		params := userop.ShowLoggedInUserParams{
			HTTPRequest: req,
		}
		builder := officeuser.NewOfficeUserFetcherPop()
		handler := ShowLoggedInUserHandler{suite.NewHandlerConfig(), builder}

		response := handler.Handle(params)

		okResponse, ok := response.(*userop.ShowLoggedInUserOK)
		suite.True(ok)
		suite.NotNil(okResponse.Payload.Permissions)
		suite.Equal(sm.UserID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("new duty location missing", func() {
		// add orders
		order := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)

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
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handler := ShowLoggedInUserHandler{handlerConfig, builder}

		response := handler.Handle(params)

		okResponse, ok := response.(*userop.ShowLoggedInUserOK)
		suite.True(ok, "Response should be ok")
		suite.NotNil(okResponse, "Response should not be nil")
		suite.NotNil(okResponse.Payload.Permissions)
		suite.Equal(sm.UserID.String(), okResponse.Payload.ID.String())
	})
}

func (suite *HandlerSuite) TestServiceMemberNoMovesLoggedInUserHandler() {
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Show: models.BoolPointer(false),
			},
		},
	}, nil)

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}

	handlerConfig := suite.NewHandlerConfig()

	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{handlerConfig, builder}

	response := handler.Handle(params)

	suite.IsType(&userop.ShowLoggedInUserUnauthorized{}, response)

}

func (suite *HandlerSuite) TestOfficeUserWithCloseoutOfficeHandler() {
	closeoutOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    closeoutOffice,
			LinkOnly: true,
			Type:     &factory.TransportationOffices.CloseoutOffice,
		},
	}, nil)
	orders := move.Orders
	orders.Moves = append(orders.Moves, move)

	suite.Run("will not have deleted roles as inactive", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor, roles.RoleTypeTOO})
		var assignedSCRole roles.Role
		for _, officeUserRole := range officeUser.User.Roles {
			if officeUserRole.RoleType == roles.RoleTypeServicesCounselor {
				assignedSCRole = officeUserRole
			}
		}

		// Mark SC as deleted
		err := suite.DB().RawQuery(
			`UPDATE users_roles SET deleted_at = now() WHERE user_id = ? AND role_id = ?`,
			officeUser.User.ID,
			assignedSCRole.ID,
		).Exec()
		suite.NoError(err)

		req := httptest.NewRequest("GET", "/users/logged_in", nil)
		req = suite.AuthenticateUserRequest(req, officeUser.User)

		params := userop.ShowLoggedInUserParams{
			HTTPRequest: req,
		}
		fakeS3 := storageTest.NewFakeS3Storage(true)
		builder := officeuser.NewOfficeUserFetcherPop()
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)

		handler := ShowLoggedInUserHandler{handlerConfig, builder}

		response := handler.Handle(params)

		okResponse, ok := response.(*userop.ShowLoggedInUserOK)

		suite.True(ok)
		suite.Equal(0, len(okResponse.Payload.InactiveRoles), "With 2 roles and 1 deleted, 1 should be current and 0 should be inactive")
	})
}

func (suite *HandlerSuite) TestServiceMemberWithCloseoutOfficeHandler() {
	closeoutOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    closeoutOffice,
			LinkOnly: true,
			Type:     &factory.TransportationOffices.CloseoutOffice,
		},
	}, nil)
	orders := move.Orders
	orders.Moves = append(orders.Moves, move)

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, orders.ServiceMember)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}
	fakeS3 := storageTest.NewFakeS3Storage(true)
	builder := officeuser.NewOfficeUserFetcherPop()
	handlerConfig := suite.NewHandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)

	handler := ShowLoggedInUserHandler{handlerConfig, builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)

	suite.True(ok)
	suite.NotNil(okResponse.Payload.Permissions)
	suite.Equal(move.CloseoutOffice.ID.String(), okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.ID.String())
	suite.Equal(move.CloseoutOffice.Name, *okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.Name)
	suite.Equal(move.CloseoutOffice.Address.ID.String(), okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.Address.ID.String())
	suite.Equal(move.CloseoutOffice.Gbloc, okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.Gbloc)

}

func (suite *HandlerSuite) TestServiceMemberWithNoCloseoutOfficeHandler() {
	// factory.BuildMove doesn't create a CloseoutOffice unless it's passed in via customizations
	move := factory.BuildMove(suite.DB(), nil, nil)

	orders := move.Orders
	orders.Moves = append(orders.Moves, move)

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateRequest(req, orders.ServiceMember)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}
	fakeS3 := storageTest.NewFakeS3Storage(true)
	builder := officeuser.NewOfficeUserFetcherPop()
	handlerConfig := suite.NewHandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)

	handler := ShowLoggedInUserHandler{handlerConfig, builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)

	suite.True(ok, "Response should be ok")
	suite.NotNil(okResponse, "Response should not be nil")
	suite.NotNil(okResponse.Payload.Permissions)
	suite.Equal(move.ID.String(), okResponse.Payload.ServiceMember.Orders[0].Moves[0].ID.String())
}
