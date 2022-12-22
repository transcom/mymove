package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/factory"
	userop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/models"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUnknownLoggedInUserHandler() {
	unknownUser := factory.BuildUser(nil, nil, nil)

	req := httptest.NewRequest("GET", "/users/logged_in", nil)
	req = suite.AuthenticateUserRequest(req, unknownUser)

	params := userop.ShowLoggedInUserParams{
		HTTPRequest: req,
	}
	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{suite.HandlerConfig(), builder}

	response := handler.Handle(params)

	okResponse, ok := response.(*userop.ShowLoggedInUserOK)
	suite.True(ok)
	suite.Equal(okResponse.Payload.ID.String(), unknownUser.ID.String())
}

func (suite *HandlerSuite) TestServiceMemberNoTransportationOfficeLoggedInUserHandler() {
	suite.Run("current duty location missing", func() {
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
		handler := ShowLoggedInUserHandler{suite.HandlerConfig(), builder}

		response := handler.Handle(params)

		okResponse, ok := response.(*userop.ShowLoggedInUserOK)
		suite.True(ok)
		suite.Equal(sm.UserID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("new duty location missing", func() {
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
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handler := ShowLoggedInUserHandler{handlerConfig, builder}

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

	handlerConfig := suite.HandlerConfig()

	builder := officeuser.NewOfficeUserFetcherPop()

	handler := ShowLoggedInUserHandler{handlerConfig, builder}

	response := handler.Handle(params)

	suite.IsType(&userop.ShowLoggedInUserUnauthorized{}, response)

}

func (suite *HandlerSuite) TestServiceMemberWithCloseoutOfficeHandler() {
	suite.Run("service member with ppm closeout office", func() {
		closeoutOffice := testdatagen.MakeTransportationOffice(suite.DB(), testdatagen.Assertions{})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				CloseoutOffice: &closeoutOffice,
			},
		})
		orders := move.Orders
		orders.Moves = append(orders.Moves, move)

		req := httptest.NewRequest("GET", "/users/logged_in", nil)
		req = suite.AuthenticateRequest(req, orders.ServiceMember)

		params := userop.ShowLoggedInUserParams{
			HTTPRequest: req,
		}
		fakeS3 := storageTest.NewFakeS3Storage(true)
		builder := officeuser.NewOfficeUserFetcherPop()
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)

		handler := ShowLoggedInUserHandler{handlerConfig, builder}

		response := handler.Handle(params)

		okResponse, ok := response.(*userop.ShowLoggedInUserOK)

		suite.True(ok)
		suite.Equal(move.CloseoutOffice.ID.String(), okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.ID.String())
		suite.Equal(move.CloseoutOffice.Name, *okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.Name)
		suite.Equal(move.CloseoutOffice.Address.ID.String(), okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.Address.ID.String())
		suite.Equal(move.CloseoutOffice.Gbloc, okResponse.Payload.ServiceMember.Orders[0].Moves[0].CloseoutOffice.Gbloc)
	})
}
