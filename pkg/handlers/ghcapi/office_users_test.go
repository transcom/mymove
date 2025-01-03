package ghcapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	officeuserop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) setupOfficeUserCreatorTestScenario() (*mocks.OfficeUserCreator, *mocks.UserRoleAssociator, *mocks.RoleAssociater, *mocks.TransportaionOfficeAssignmentUpdater, *RequestOfficeUserHandler) {
	mockCreator := &mocks.OfficeUserCreator{}
	mockUserRoleAssociator := &mocks.UserRoleAssociator{}
	mockRoleAssociator := &mocks.RoleAssociater{}
	mockTransportaionOfficeAssignmentUpdater := &mocks.TransportaionOfficeAssignmentUpdater{}
	handler := &RequestOfficeUserHandler{
		HandlerConfig:                        suite.HandlerConfig(),
		OfficeUserCreator:                    mockCreator,
		NewQueryFilter:                       query.NewQueryFilter,
		UserRoleAssociator:                   mockUserRoleAssociator,
		RoleAssociater:                       mockRoleAssociator,
		TransportaionOfficeAssignmentUpdater: mockTransportaionOfficeAssignmentUpdater,
	}
	return mockCreator, mockUserRoleAssociator, mockRoleAssociator, mockTransportaionOfficeAssignmentUpdater, handler
}

// Services Counselor. Task Ordering Officer (TOO), Task Invoicing Officer (TIO),
// Quality Assurance Evaluator (QAE), and Customer Service Representative (CSR)
// Are all roles allowed to request office user (They authenticate with AuthenticateOfficeRequest)
func (suite *HandlerSuite) TestRequestOfficeUserHandler() {
	suite.Run("Successfully requests the creation of an office user", func() {
		mockCreator, mockRoleAssociator, mockRoleFetcher, mockTransportaionOfficeAssignmentUpdater, handler := suite.setupOfficeUserCreatorTestScenario()

		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeServicesCounselor, roles.RoleTypeTIO, roles.RoleTypeQae})
		request := httptest.NewRequest("POST", "/requested-office-users", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := officeuserop.CreateRequestedOfficeUserParams{
			HTTPRequest: request,
			OfficeUser: &ghcmessages.OfficeUserCreate{
				FirstName:              "John",
				LastName:               "Doe",
				Telephone:              "555-555-5555",
				Email:                  "johndoe@example.com",
				Edipi:                  models.StringPointer("1234567890"),
				TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
				Roles:                  []*ghcmessages.OfficeUserRole{{RoleType: handlers.FmtString(string(roles.RoleTypeTOO))}},
			},
		}

		status := models.OfficeUserStatusREQUESTED
		// Mock successful creation in the database
		mockCreator.On(
			"CreateOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(&models.OfficeUser{
			ID:                     uuid.Must(uuid.NewV4()),
			UserID:                 models.UUIDPointer(uuid.Must(uuid.NewV4())),
			FirstName:              "John",
			LastName:               "Doe",
			Telephone:              "555-555-5555",
			Email:                  "johndoe@example.com",
			EDIPI:                  models.StringPointer("1234567890"),
			TransportationOfficeID: transportationOffice.ID,
			Status:                 &status,
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		}, nil, nil).Once()

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		// Mock successful role association
		mockRoleAssociator.On(
			"UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()
		// Mock successful role return
		mockRoleFetcher.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		mockTransportationAssignments := models.TransportationOfficeAssignments{
			models.TransportationOfficeAssignment{
				ID:                     officeUser.ID,
				TransportationOfficeID: officeUser.TransportationOfficeID,
				PrimaryOffice:          models.BoolPointer(true),
				CreatedAt:              time.Now(),
				UpdatedAt:              time.Now(),
			},
		}
		mockTransportaionOfficeAssignmentUpdater.On(
			"UpdateTransportaionOfficeAssignments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(mockTransportationAssignments, nil)

		// Handle params with mocked services
		response := handler.Handle(params)

		suite.IsType(&officeuserop.CreateRequestedOfficeUserCreated{}, response)
		createdResponse := response.(*officeuserop.CreateRequestedOfficeUserCreated)
		suite.Equal("John", *createdResponse.Payload.FirstName)
		suite.Equal("REQUESTED", *createdResponse.Payload.Status)
		suite.Equal(1, len(createdResponse.Payload.TransportationOfficeAssignments))

		// Ensure that the mock assertions are met
		mockCreator.AssertExpectations(suite.T())
		mockRoleAssociator.AssertExpectations(suite.T())
	})

	suite.Run("Responds proper validation errors", func() {
		mockCreator, _, _, _, handler := suite.setupOfficeUserCreatorTestScenario()

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeServicesCounselor})
		transportationOfficeID, _ := uuid.NewV4()
		request := httptest.NewRequest("POST", "/requested-office-users", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := officeuserop.CreateRequestedOfficeUserParams{
			HTTPRequest: request,
			OfficeUser: &ghcmessages.OfficeUserCreate{
				FirstName:              "John",
				LastName:               "Doe",
				Telephone:              "555-555-5555",
				Email:                  "johndoeexample.com",
				Edipi:                  models.StringPointer("1234567890"),
				TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
				Roles:                  []*ghcmessages.OfficeUserRole{{RoleType: handlers.FmtString(string(roles.RoleTypeTOO))}},
			},
		}

		// Mock validation error of faulty email format
		verrs := validate.NewErrors()
		verrs.Add("email", "Invalid email format")
		// Mock the "On CreateOfficeUser" -> return verrs as the email format was wrong
		mockCreator.On("CreateOfficeUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, verrs, nil)
		// User role update mock not required as this function will error out before it is called (Expected behavior)

		// Trigger the mocks
		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateRequestedOfficeUserUnprocessableEntity{}, response)
		verrResponse, ok := response.(*officeuserop.CreateRequestedOfficeUserUnprocessableEntity)
		suite.True(ok)
		suite.NotEmpty(verrResponse.Payload.InvalidFields, "expected validation errors")
		// Since we mocked an email verr, make sure it's here
		suite.Contains(verrResponse.Payload.InvalidFields, "email", "expected error on 'email' field")

		// Ensure that the mock assertion is met
		mockCreator.AssertExpectations(suite.T())
	})

	suite.Run("Bad transportation office ID", func() {
		_, _, _, _, handler := suite.setupOfficeUserCreatorTestScenario()
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		transportationOfficeID := "Not a UUID"
		request := httptest.NewRequest("POST", "/requested-office-users", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := officeuserop.CreateRequestedOfficeUserParams{
			HTTPRequest: request,
			OfficeUser: &ghcmessages.OfficeUserCreate{
				FirstName:              "John",
				LastName:               "Doe",
				Telephone:              "555-555-5555",
				Email:                  "johndoe@example.com",
				Edipi:                  models.StringPointer("1234567890"),
				TransportationOfficeID: strfmt.UUID(transportationOfficeID),
				Roles:                  []*ghcmessages.OfficeUserRole{{RoleType: handlers.FmtString(string(roles.RoleTypeTOO))}},
			},
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateRequestedOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("No payload roles", func() {
		_, _, _, _, handler := suite.setupOfficeUserCreatorTestScenario()
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		transportationOfficeID, _ := uuid.NewV4()
		request := httptest.NewRequest("POST", "/requested-office-users", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := officeuserop.CreateRequestedOfficeUserParams{
			HTTPRequest: request,
			OfficeUser: &ghcmessages.OfficeUserCreate{
				FirstName:              "John",
				LastName:               "Doe",
				Telephone:              "555-555-5555",
				Email:                  "johndoe@example.com",
				Edipi:                  models.StringPointer("1234567890"),
				TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
			},
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateRequestedOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Bad payload roles", func() {
		_, _, _, _, handler := suite.setupOfficeUserCreatorTestScenario()
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{})
		transportationOfficeID := "Not a UUID"
		request := httptest.NewRequest("POST", "/requested-office-users", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := officeuserop.CreateRequestedOfficeUserParams{
			HTTPRequest: request,
			OfficeUser: &ghcmessages.OfficeUserCreate{
				FirstName:              "John",
				LastName:               "Doe",
				Telephone:              "555-555-5555",
				Email:                  "johndoe@example.com",
				Edipi:                  models.StringPointer("1234567890"),
				TransportationOfficeID: strfmt.UUID(transportationOfficeID),
				Roles:                  []*ghcmessages.OfficeUserRole{{RoleType: handlers.FmtString(string(roles.RoleTypeTOO))}},
			},
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateRequestedOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Enforces identification rule", func() {
		_, _, _, _, handler := suite.setupOfficeUserCreatorTestScenario()

		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeServicesCounselor, roles.RoleTypeTIO, roles.RoleTypeQae})
		request := httptest.NewRequest("POST", "/requested-office-users", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		// EDIPI and other unique ID missing
		params := officeuserop.CreateRequestedOfficeUserParams{
			HTTPRequest: request,
			OfficeUser: &ghcmessages.OfficeUserCreate{
				FirstName:              "John",
				LastName:               "Doe",
				Telephone:              "555-555-5555",
				Email:                  "johndoe@example.com",
				TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
				Roles:                  []*ghcmessages.OfficeUserRole{{RoleType: handlers.FmtString(string(roles.RoleTypeTOO))}},
			},
		}

		//Our handler will fail before any mock services are needed
		response := handler.Handle(params)

		suite.IsType(&officeuserop.CreateRequestedOfficeUserUnprocessableEntity{}, response)

		verrResponse, ok := response.(*officeuserop.CreateRequestedOfficeUserUnprocessableEntity)
		suite.True(ok)
		suite.NotEmpty(verrResponse.Payload.ClientError, "expected validation errors from missing identification param")
		verrDetail := "Data received from requester is bad: BAD_DATA: Either an EDIPI or Other Unique ID must be provided"
		suite.Contains(*verrResponse.Payload.ClientError.Detail, verrDetail)
	})
}
