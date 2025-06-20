package adminapi

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	rejectedofficeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/rejected_office_users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	rejectedofficeusers "github.com/transcom/mymove/pkg/services/rejected_office_users"
)

func (suite *HandlerSuite) TestIndexRejectedOfficeUsersHandler() {
	// test that everything is wired up
	suite.Run("rejected users result in ok response", func() {
		// building two office user with rejected status
		rejectedOfficeUsers := models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})}

		params := rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRejectedOfficeUsersHandler{
			HandlerConfig:                 suite.NewHandlerConfig(),
			NewQueryFilter:                query.NewQueryFilter,
			RejectedOfficeUserListFetcher: rejectedofficeusers.NewRejectedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                 pagination.NewPagination,
		}

		response := handler.Handle(params)

		// should get an ok response
		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse := response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Equal(len(okResponse.Payload), len(rejectedOfficeUsers))

		actualID := []string{okResponse.Payload[0].ID.String(), okResponse.Payload[1].ID.String()}
		expected := []string{rejectedOfficeUsers[0].ID.String(), rejectedOfficeUsers[1].ID.String()}
		for _, expectedID := range expected {
			suite.True(slices.Contains(actualID, expectedID))
		}
	})

	suite.Run("able to search by name and filter", func() {
		status := models.OfficeUserStatusREJECTED
		rejectionReason := "Test rejection Reason"
		rejectionReason2 := "Test rejection2 Reason"
		rejectedOn := time.Date(2025, 03, 05, 1, 1, 1, 1, time.Local)
		rejectedOn2 := time.Date(2024, 03, 07, 1, 1, 1, 1, time.Local)
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "JPPO Test Office",
				},
			},
		}, nil)
		transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "PPO Rome Test Office",
				},
			},
		}, nil)
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Angelina",
					LastName:               "Jolie",
					Email:                  "laraCroft@mail.mil",
					Status:                 &status,
					TransportationOfficeID: transportationOffice2.ID,
					RejectionReason:        &rejectionReason2,
					RejectedOn:             &rejectedOn2,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Billy",
					LastName:               "Bob",
					Email:                  "bigBob@mail.mil",
					Status:                 &status,
					TransportationOfficeID: transportationOffice2.ID,
					RejectionReason:        &rejectionReason,
					RejectedOn:             &rejectedOn2,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTIO})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Nick",
					LastName:               "Cage",
					Email:                  "conAirKilluh@mail.mil",
					Status:                 &status,
					TransportationOfficeID: transportationOffice2.ID,
					RejectionReason:        &rejectionReason,
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Nick",
					LastName:               "Cage",
					Email:                  "conAirKilluh2@mail.mil",
					Status:                 &status,
					TransportationOfficeID: transportationOffice.ID,
					RejectionReason:        &rejectionReason,
					RejectedOn:             &rejectedOn,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		// partial search
		nameSearch := "ic"
		filterJSON := fmt.Sprintf("{\"search\":\"%s\"}", nameSearch)
		params := rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRejectedOfficeUsersHandler{
			HandlerConfig:                 suite.NewHandlerConfig(),
			NewQueryFilter:                query.NewQueryFilter,
			RejectedOfficeUserListFetcher: rejectedofficeusers.NewRejectedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                 pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse := response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Contains(*okResponse.Payload[0].FirstName, nameSearch)
		suite.Contains(*okResponse.Payload[1].FirstName, nameSearch)

		// email search
		emailSearch := "AirKilluh2"
		filterJSON = fmt.Sprintf("{\"emails\":\"%s\"}", emailSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].Email, emailSearch)

		// firstName search
		firstSearch := "Angel"
		filterJSON = fmt.Sprintf("{\"firstName\":\"%s\"}", firstSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].FirstName, firstSearch)

		// lastName search
		lastSearch := "Jo"
		filterJSON = fmt.Sprintf("{\"lastName\":\"%s\"}", lastSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].LastName, lastSearch)

		// transportation office search
		filterJSON = "{\"offices\":\"JP\"}"
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(strfmt.UUID(transportationOffice.ID.String()), *okResponse.Payload[0].TransportationOfficeID)

		// rejection reason search
		reasonSearch := "n2"
		filterJSON = fmt.Sprintf("{\"rejectionReason\":\"%s\"}", reasonSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].RejectionReason, reasonSearch)

		// rejectedOn search
		rejectedOnSearch := "5"
		filterJSON = fmt.Sprintf("{\"rejectedOn\":\"%s\"}", rejectedOnSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(okResponse.Payload[0].RejectedOn.String(), rejectedOnSearch)

		// roles search
		roleSearch := "Services Counselor"
		filterJSON = fmt.Sprintf("{\"roles\":\"%s\"}", roleSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Contains(*okResponse.Payload[0].Roles[0].RoleName, roleSearch)
		suite.Contains(*okResponse.Payload[1].Roles[0].RoleName, roleSearch)

	})
}

func (suite *HandlerSuite) TestGetRejectedOfficeUserHandler() {
	suite.Run("integration test ok response", func() {
		rejectedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := rejectedofficeuserop.GetRejectedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/rejected_office_users/%s", rejectedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(rejectedOfficeUser.ID.String()),
		}

		mockRoleFetcher := &mocks.RoleFetcher{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleFetcher.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		queryBuilder := query.NewQueryBuilder()
		handler := GetRejectedOfficeUserHandler{
			suite.NewHandlerConfig(),
			rejectedofficeusers.NewRejectedOfficeUserFetcher(queryBuilder),
			mockRoleFetcher,
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.GetRejectedOfficeUserOK{}, response)
		okResponse := response.(*rejectedofficeuserop.GetRejectedOfficeUserOK)
		suite.Equal(rejectedOfficeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("successful response", func() {
		rejectedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := rejectedofficeuserop.GetRejectedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/rejected_office_users/%s", rejectedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(rejectedOfficeUser.ID.String()),
		}

		rejectedOfficeUserFetcher := &mocks.RejectedOfficeUserFetcher{}
		rejectedOfficeUserFetcher.On("FetchRejectedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(rejectedOfficeUser, nil).Once()

		mockRoleFetcher := &mocks.RoleFetcher{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleFetcher.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		handler := GetRejectedOfficeUserHandler{
			suite.NewHandlerConfig(),
			rejectedOfficeUserFetcher,
			mockRoleFetcher,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.GetRejectedOfficeUserOK{}, response)
		okResponse := response.(*rejectedofficeuserop.GetRejectedOfficeUserOK)
		suite.Equal(rejectedOfficeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		rejectedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := rejectedofficeuserop.GetRejectedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/rejected_office_users/%s", rejectedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(rejectedOfficeUser.ID.String()),
		}

		expectedError := models.ErrFetchNotFound
		rejectedOfficeUserFetcher := &mocks.RejectedOfficeUserFetcher{}
		rejectedOfficeUserFetcher.On("FetchRejectedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(models.OfficeUser{}, expectedError).Once()

		mockRoleFetcher := &mocks.RoleFetcher{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleFetcher.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		handler := GetRejectedOfficeUserHandler{
			suite.NewHandlerConfig(),
			rejectedOfficeUserFetcher,
			mockRoleFetcher,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
