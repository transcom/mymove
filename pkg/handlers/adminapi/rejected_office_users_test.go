package adminapi

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
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
			HandlerConfig:                 suite.HandlerConfig(),
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

	suite.Run("able to search by name & email", func() {
		status := models.OfficeUserStatusREJECTED
		rejectionReason := "Test rejection Reason"
		rejectedOn := time.Date(2025, 03, 05, 1, 1, 1, 1, time.Local)
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "JPPO Test Office",
				},
			},
		}, nil)
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Angelina",
					LastName:  "Jolie",
					Email:     "laraCroft@mail.mil",
					Status:    &status,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Billy",
					LastName:  "Bob",
					Email:     "bigBob@mail.mil",
					Status:    &status,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTIO})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Nick",
					LastName:  "Cage",
					Email:     "conAirKilluh@mail.mil",
					Status:    &status,
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

		// partial first name search
		nameSearch := "Nick"
		filterJSON := fmt.Sprintf("{\"search\":\"%s\"}", nameSearch)
		params := rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRejectedOfficeUsersHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			NewQueryFilter:                query.NewQueryFilter,
			RejectedOfficeUserListFetcher: rejectedofficeusers.NewRejectedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                 pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse := response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(nameSearch, *okResponse.Payload[0].FirstName)
		suite.Equal(nameSearch, *okResponse.Payload[1].FirstName)

		// email search
		emailSearch := "conAirKilluh2"
		filterJSON = fmt.Sprintf("{\"emails\":\"%s\"}", emailSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)

		respEmail := *okResponse.Payload[0].Email
		suite.Equal(emailSearch, respEmail[0:len(emailSearch)])
		suite.Equal(emailSearch, respEmail[0:len(emailSearch)])

		// firstName search
		firstSearch := "Angelina"
		filterJSON = fmt.Sprintf("{\"firstName\":\"%s\"}", firstSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(firstSearch, *okResponse.Payload[0].FirstName)

		// lastName search
		lastSearch := "Cage"
		filterJSON = fmt.Sprintf("{\"lastName\":\"%s\"}", lastSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(lastSearch, *okResponse.Payload[0].LastName)
		suite.Equal(lastSearch, *okResponse.Payload[1].LastName)

		// transportation office search
		filterJSON = "{\"offices\":\"JPPO\"}"
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
		reasonSearch := "Test rejection"
		filterJSON = fmt.Sprintf("{\"rejectionReason\":\"%s\"}", reasonSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		respRejection := *okResponse.Payload[0].RejectionReason
		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(reasonSearch, respRejection[0:len(reasonSearch)])

		// rejectedOn search
		rejectedOnSearch := "03"
		filterJSON = fmt.Sprintf("{\"rejectedOn\":\"%s\"}", rejectedOnSearch)
		params = rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse = response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)

		rejectedOnResp := okResponse.Payload[0].RejectedOn.String()
		actualYear := strings.Split(rejectedOnResp, "-")[0]
		actualMonth, err := strconv.Atoi(strings.Split(rejectedOnResp, "-")[1])

		expectedYear := strconv.Itoa(rejectedOn.Year())
		expectedMonth := int(rejectedOn.Month())

		suite.NoError(err)
		suite.Equal(expectedYear, actualYear)
		suite.Equal(expectedMonth, actualMonth)

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
		suite.Equal(roleSearch, *okResponse.Payload[0].Roles[0].RoleName)
		suite.Equal(roleSearch, *okResponse.Payload[1].Roles[0].RoleName)

	})
}

func (suite *HandlerSuite) TestGetRejectedOfficeUserHandler() {
	suite.Run("integration test ok response", func() {
		rejectedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := rejectedofficeuserop.GetRejectedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/rejected_office_users/%s", rejectedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(rejectedOfficeUser.ID.String()),
		}

		mockRoleAssociator := &mocks.RoleAssociater{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		queryBuilder := query.NewQueryBuilder()
		handler := GetRejectedOfficeUserHandler{
			suite.HandlerConfig(),
			rejectedofficeusers.NewRejectedOfficeUserFetcher(queryBuilder),
			mockRoleAssociator,
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

		mockRoleAssociator := &mocks.RoleAssociater{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		handler := GetRejectedOfficeUserHandler{
			suite.HandlerConfig(),
			rejectedOfficeUserFetcher,
			mockRoleAssociator,
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

		mockRoleAssociator := &mocks.RoleAssociater{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		handler := GetRejectedOfficeUserHandler{
			suite.HandlerConfig(),
			rejectedOfficeUserFetcher,
			mockRoleAssociator,
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
