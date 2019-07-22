package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/gen/adminmessages"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/user"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func newMockQueryFilterBuilder(filter *mocks.QueryFilter) services.NewQueryFilter {
	return func(column string, comparator string, value interface{}) services.QueryFilter {
		return filter
	}
}

func (suite *HandlerSuite) TestIndexOfficeUsersHandler() {
	// replace this with generated UUID when filter param is built out
	uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID: id,
		},
	}
	testdatagen.MakeOfficeUser(suite.DB(), assertions)
	testdatagen.MakeDefaultOfficeUser(suite.DB())

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", "/office_users", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: req,
		}

		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexOfficeUsersHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:        query.NewQueryFilter,
			OfficeUserListFetcher: user.NewOfficeUserListFetcher(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		officeUser := models.OfficeUser{ID: id}
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: req,
		}
		officeUserListFetcher := &mocks.OfficeUserListFetcher{}
		officeUserListFetcher.On("FetchOfficeUserList",
			mock.Anything,
		).Return(models.OfficeUsers{officeUser}, nil).Once()
		handler := IndexOfficeUsersHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:        newQueryFilter,
			OfficeUserListFetcher: officeUserListFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		officeUserListFetcher := &mocks.OfficeUserListFetcher{}
		officeUserListFetcher.On("FetchOfficeUserList",
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexOfficeUsersHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:        newQueryFilter,
			OfficeUserListFetcher: officeUserListFetcher,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestCreateOfficeUserHandler() {

	transportationOfficeID, _ := uuid.NewV4()
	officeUserID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	officeUser := models.OfficeUser{ID: officeUserID, TransportationOfficeID: transportationOfficeID, UserID: nil}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	req := httptest.NewRequest("POST", "/office_users", nil)
	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := officeuserop.CreateOfficeUserParams{
		HTTPRequest: req,
		OfficeUser: &adminmessages.OfficeUserCreatePayload{
			FirstName:              officeUser.FirstName,
			LastName:               officeUser.LastName,
			Telephone:              officeUser.Telephone,
			TransportationOfficeID: strfmt.UUID(officeUser.TransportationOfficeID.String()),
		},
	}

	suite.T().Run("Successful create", func(t *testing.T) {
		officeUserCreator := &mocks.OfficeUserCreator{}

		officeUserCreator.On("CreateOfficeUser",
			&officeUser,
			mock.Anything).Return(&officeUser, nil, nil).Once()

		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})

	suite.T().Run("Failed create", func(t *testing.T) {
		officeUserCreator := &mocks.OfficeUserCreator{}

		officeUserCreator.On("CreateOfficeUser",
			&officeUser,
			mock.Anything).Return(&officeUser, nil, nil).Once()

		handler := CreateOfficeUserHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			officeUserCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})

	officeUserCreator := &mocks.OfficeUserCreator{}
	err := validate.NewErrors()

	officeUserCreator.On("CreateOfficeUser",
		&officeUser,
		mock.Anything).Return(nil, err, nil).Once()

	handler := CreateOfficeUserHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		officeUserCreator,
		newQueryFilter,
	}

	handler.Handle(params)
	suite.Error(err, "Error saving user")

}
