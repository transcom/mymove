package ghcapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/jarcoal/httpmock"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	customerops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	customerservice "github.com/transcom/mymove/pkg/services/office_user/customer"
)

const officeProviderName = "officeProvider"

func (suite *HandlerSuite) TestGetCustomerHandlerIntegration() {
	customer := factory.BuildServiceMember(suite.DB(), nil, nil)

	request := httptest.NewRequest("GET", "/customer/{customerID}", nil)
	params := customerops.GetCustomerParams{
		HTTPRequest: request,
		CustomerID:  strfmt.UUID(customer.ID.String()),
	}
	handlerConfig := suite.HandlerConfig()
	handler := GetCustomerHandler{
		handlerConfig,
		customerservice.NewCustomerFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&customerops.GetCustomerOK{}, response)
	getCustomerResponse := response.(*customerops.GetCustomerOK)
	getCustomerPayload := getCustomerResponse.Payload

	// Validate outgoing payload
	suite.NoError(getCustomerPayload.Validate(strfmt.Default))

	suite.Equal(strfmt.UUID(customer.ID.String()), getCustomerPayload.ID)
	suite.Equal(*customer.Edipi, getCustomerPayload.Edipi)
	suite.Equal(strfmt.UUID(customer.UserID.String()), getCustomerPayload.UserID)
	suite.Equal(customer.Affiliation.String(), getCustomerPayload.Agency)
	suite.Equal(customer.PersonalEmail, getCustomerPayload.Email)
	suite.Equal(customer.Telephone, getCustomerPayload.Phone)
	suite.NotZero(getCustomerPayload.CurrentAddress)
}

func (suite *HandlerSuite) TestUpdateCustomerHandler() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeServicesCounselor,
	})

	body := &ghcmessages.UpdateCustomerPayload{
		LastName:  "Newlastname",
		FirstName: "Newfirstname",
		Phone:     handlers.FmtString("223-455-3399"),
		BackupContact: &ghcmessages.BackupContact{
			Name:  handlers.FmtString("New Backup Contact"),
			Phone: handlers.FmtString("445-345-1212"),
			Email: handlers.FmtString("newbackup@mail.com"),
		},
	}
	currentAddress := ghcmessages.Address{
		StreetAddress1: handlers.FmtString("123 New Street"),
		City:           handlers.FmtString("Newcity"),
		State:          handlers.FmtString("MA"),
		PostalCode:     handlers.FmtString("12345"),
	}
	body.CurrentAddress.Address = currentAddress

	customer := factory.BuildExtendedServiceMember(suite.DB(), nil, nil)
	request := httptest.NewRequest("PATCH", "/orders/{customerID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := customerops.UpdateCustomerParams{
		HTTPRequest: request,
		CustomerID:  strfmt.UUID(customer.ID.String()),
		IfMatch:     etag.GenerateEtag(customer.UpdatedAt),
		Body:        body,
	}
	handlerConfig := suite.HandlerConfig()
	handler := UpdateCustomerHandler{
		handlerConfig,
		customerservice.NewCustomerUpdater(),
	}

	// Validate incoming payload
	suite.NoError(params.Body.Validate(strfmt.Default))

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	// TODO: test with actual updated customer?
	// updatedCustomer, _ := models.FetchServiceMember(suite.DB(), customer.ID)
	suite.Assertions.IsType(&customerops.UpdateCustomerOK{}, response)
	updateCustomerResponse := response.(*customerops.UpdateCustomerOK)
	updateCustomerPayload := updateCustomerResponse.Payload

	// Validate outgoing payload
	suite.NoError(updateCustomerPayload.Validate(strfmt.Default))

	suite.Equal(body.FirstName, updateCustomerPayload.FirstName)
	suite.Equal(body.LastName, updateCustomerPayload.LastName)
	suite.Equal(body.Phone, updateCustomerPayload.Phone)
	suite.Equal(body.CurrentAddress.StreetAddress1, updateCustomerPayload.CurrentAddress.StreetAddress1)
	suite.Equal(body.CurrentAddress.City, updateCustomerPayload.CurrentAddress.City)
	suite.Equal(body.CurrentAddress.PostalCode, updateCustomerPayload.CurrentAddress.PostalCode)
	suite.Equal(body.CurrentAddress.State, updateCustomerPayload.CurrentAddress.State)
	suite.Equal(body.BackupContact.Name, updateCustomerPayload.BackupContact.Name)
	suite.Equal(body.BackupContact.Phone, updateCustomerPayload.BackupContact.Phone)
	suite.Equal(body.BackupContact.Email, updateCustomerPayload.BackupContact.Email)
}

func (suite *HandlerSuite) TestCreateCustomerWithOktaOptionHandler() {
	suite.Run("Successful customer creation", func() {
		// in order to call the endpoint, we need to be an authenticated office user that's a SC
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeServicesCounselor,
		})

		// Build provider
		provider, err := factory.BuildOktaProvider(officeProviderName)
		suite.NoError(err)

		mockAndActivateOktaEndpoints(provider)

		residentialAddress := ghcmessages.Address{
			StreetAddress1: handlers.FmtString("123 New Street"),
			City:           handlers.FmtString("Newcity"),
			State:          handlers.FmtString("MA"),
			PostalCode:     handlers.FmtString("02110"),
		}

		backupAddress := ghcmessages.Address{
			StreetAddress1: handlers.FmtString("123 Backup Street"),
			City:           handlers.FmtString("Backupcity"),
			State:          handlers.FmtString("MA"),
			PostalCode:     handlers.FmtString("02115"),
		}

		affiliation := ghcmessages.AffiliationARMY

		body := &ghcmessages.CreateCustomerPayload{
			LastName:      "Last",
			FirstName:     "First",
			Telephone:     handlers.FmtString("223-455-3399"),
			Affiliation:   &affiliation,
			Edipi:         "",
			Emplid:        handlers.FmtString(""),
			PersonalEmail: *handlers.FmtString("email@email.com"),
			BackupContact: &ghcmessages.BackupContact{
				Name:  handlers.FmtString("New Backup Contact"),
				Phone: handlers.FmtString("445-345-1212"),
				Email: handlers.FmtString("newbackup@mail.com"),
			},
			ResidentialAddress: struct {
				ghcmessages.Address
			}{
				Address: residentialAddress,
			},
			BackupMailingAddress: struct {
				ghcmessages.Address
			}{
				Address: backupAddress,
			},
			CreateOktaAccount: true,
			// when CacUser is false, this indicates a non-CAC user so CacValidated is set to true
			CacUser: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/customer", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := customerops.CreateCustomerWithOktaOptionParams{
			HTTPRequest: request,
			Body:        body,
		}
		handlerConfig := suite.HandlerConfig()
		handler := CreateCustomerWithOktaOptionHandler{
			handlerConfig,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&customerops.CreateCustomerWithOktaOptionOK{}, response)
		createdCustomerResponse := response.(*customerops.CreateCustomerWithOktaOptionOK)
		createdCustomerPayload := createdCustomerResponse.Payload

		suite.NoError(createdCustomerPayload.Validate(strfmt.Default))

		suite.Equal(body.FirstName, createdCustomerPayload.FirstName)
		suite.Equal(body.LastName, createdCustomerPayload.LastName)
		suite.Equal(body.Telephone, createdCustomerPayload.Telephone)
		suite.Equal(body.BackupContact.Name, createdCustomerPayload.BackupContact.Name)
		suite.Equal(body.BackupContact.Phone, createdCustomerPayload.BackupContact.Phone)
		suite.Equal(body.BackupContact.Email, createdCustomerPayload.BackupContact.Email)
		// when CacUser is false, this indicates a non-CAC user so CacValidated is set to true
		suite.Equal(true, createdCustomerPayload.CacValidated)
	})

	suite.Run("Unable to create customer when using an existing DODID", func() {
		// in order to call the endpoint, we need to be an authenticated office user that's a SC
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeServicesCounselor,
		})

		// creating a customer so we can try and use their existing DODID
		customer := factory.BuildServiceMember(suite.DB(), nil, nil)

		// setting the flag to true
		os.Setenv("FEATURE_FLAG_DODID_UNIQUE", "true")

		// Build provider
		provider, err := factory.BuildOktaProvider(officeProviderName)
		suite.NoError(err)

		mockAndActivateOktaEndpoints(provider)

		residentialAddress := ghcmessages.Address{
			StreetAddress1: handlers.FmtString("123 New Street"),
			City:           handlers.FmtString("Newcity"),
			State:          handlers.FmtString("MA"),
			PostalCode:     handlers.FmtString("02110"),
		}

		backupAddress := ghcmessages.Address{
			StreetAddress1: handlers.FmtString("123 Backup Street"),
			City:           handlers.FmtString("Backupcity"),
			State:          handlers.FmtString("MA"),
			PostalCode:     handlers.FmtString("02115"),
		}

		affiliation := ghcmessages.AffiliationARMY

		body := &ghcmessages.CreateCustomerPayload{
			LastName:      "Last",
			FirstName:     "First",
			Telephone:     handlers.FmtString("223-455-3399"),
			Affiliation:   &affiliation,
			Edipi:         *customer.Edipi,
			PersonalEmail: *handlers.FmtString("email@email.com"),
			BackupContact: &ghcmessages.BackupContact{
				Name:  handlers.FmtString("New Backup Contact"),
				Phone: handlers.FmtString("445-345-1212"),
				Email: handlers.FmtString("newbackup@mail.com"),
			},
			ResidentialAddress: struct {
				ghcmessages.Address
			}{
				Address: residentialAddress,
			},
			BackupMailingAddress: struct {
				ghcmessages.Address
			}{
				Address: backupAddress,
			},
			CreateOktaAccount: true,
			// when CacUser is false, this indicates a non-CAC user so CacValidated is set to true
			CacUser: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/customer", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := customerops.CreateCustomerWithOktaOptionParams{
			HTTPRequest: request,
			Body:        body,
		}
		handlerConfig := suite.HandlerConfig()
		handler := CreateCustomerWithOktaOptionHandler{
			handlerConfig,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&customerops.CreateCustomerWithOktaOptionUnprocessableEntity{}, response)
	})

	suite.Run("Unable to create customer of affiliation Coast Guard with no EMPLID", func() {
		// in order to call the endpoint, we need to be an authenticated office user that's a SC
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeServicesCounselor,
		})

		// Build provider
		provider, err := factory.BuildOktaProvider(officeProviderName)
		suite.NoError(err)

		mockAndActivateOktaEndpoints(provider)

		residentialAddress := ghcmessages.Address{
			StreetAddress1: handlers.FmtString("123 New Street"),
			City:           handlers.FmtString("Newcity"),
			State:          handlers.FmtString("MA"),
			PostalCode:     handlers.FmtString("02110"),
		}

		backupAddress := ghcmessages.Address{
			StreetAddress1: handlers.FmtString("123 Backup Street"),
			City:           handlers.FmtString("Backupcity"),
			State:          handlers.FmtString("MA"),
			PostalCode:     handlers.FmtString("02115"),
		}

		affiliation := ghcmessages.AffiliationCOASTGUARD

		body := &ghcmessages.CreateCustomerPayload{
			LastName:      "Last",
			FirstName:     "First",
			Telephone:     handlers.FmtString("223-455-3399"),
			Affiliation:   &affiliation,
			Edipi:         "1234567890",
			PersonalEmail: *handlers.FmtString("email@email.com"),
			BackupContact: &ghcmessages.BackupContact{
				Name:  handlers.FmtString("New Backup Contact"),
				Phone: handlers.FmtString("445-345-1212"),
				Email: handlers.FmtString("newbackup@mail.com"),
			},
			ResidentialAddress: struct {
				ghcmessages.Address
			}{
				Address: residentialAddress,
			},
			BackupMailingAddress: struct {
				ghcmessages.Address
			}{
				Address: backupAddress,
			},
			CreateOktaAccount: true,
			// when CacUser is false, this indicates a non-CAC user so CacValidated is set to true
			CacUser: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/customer", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := customerops.CreateCustomerWithOktaOptionParams{
			HTTPRequest: request,
			Body:        body,
		}
		handlerConfig := suite.HandlerConfig()
		handler := CreateCustomerWithOktaOptionHandler{
			handlerConfig,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&customerops.CreateCustomerWithOktaOptionUnprocessableEntity{}, response)
		failedToCreateCustomerPayload := response.(*customerops.CreateCustomerWithOktaOptionUnprocessableEntity).Payload.ClientError.Detail
		suite.Equal("ID: 00000000-0000-0000-0000-000000000000 is in a conflicting state Service members from the Coast Guard require an EMPLID for creation.", *failedToCreateCustomerPayload)
	})
}

func (suite *HandlerSuite) TestSearchCustomersHandler() {
	var requestUser models.User
	setupTestData := func() *http.Request {
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/customer/#{customer.id}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		return req
	}

	suite.Run("Successful customer search by DOD ID", func() {
		req := setupTestData()
		customer := factory.BuildServiceMember(suite.DB(), nil, nil)
		customers := models.ServiceMemberSearchResults{customer.ToSearchResult()}

		mockSearcher := mocks.CustomerSearcher{}

		handler := SearchCustomersHandler{
			HandlerConfig:    suite.HandlerConfig(),
			CustomerSearcher: &mockSearcher,
		}
		mockSearcher.On("SearchCustomers",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchCustomersParams) bool {
				return *params.Edipi == *customer.Edipi &&
					params.CustomerName == nil
			}),
		).Return(customers, 1, nil)

		params := customerops.SearchCustomersParams{
			HTTPRequest: req,
			Body: customerops.SearchCustomersBody{
				Edipi: customer.Edipi,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&customerops.SearchCustomersOK{}, response)
		payload := response.(*customerops.SearchCustomersOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(customer.ID.String(), (*payload).SearchCustomers[0].ID.String())
	})

	suite.Run("Successful customer search by name", func() {
		req := setupTestData()
		customer := factory.BuildServiceMember(suite.DB(), nil, nil)
		customers := models.ServiceMemberSearchResults{customer.ToSearchResult()}

		mockSearcher := mocks.CustomerSearcher{}

		handler := SearchCustomersHandler{
			HandlerConfig:    suite.HandlerConfig(),
			CustomerSearcher: &mockSearcher,
		}
		mockSearcher.On("SearchCustomers",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchCustomersParams) bool {
				return *params.CustomerName == *customer.FirstName &&
					params.Edipi == nil
			}),
		).Return(customers, 1, nil)

		params := customerops.SearchCustomersParams{
			HTTPRequest: req,
			Body: customerops.SearchCustomersBody{
				CustomerName: customer.FirstName,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&customerops.SearchCustomersOK{}, response)
		payload := response.(*customerops.SearchCustomersOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(customer.FirstName, (*payload).SearchCustomers[0].FirstName)
	})
}

// Generate and activate Okta endpoints that will be using during the auth handlers.
func mockAndActivateOktaEndpoints(provider *okta.Provider) {

	activate := "true"
	createUserEndpoint := provider.GetCreateUserURL(activate)
	oktaID := "fakeSub"

	httpmock.RegisterResponder("POST", createUserEndpoint,
		httpmock.NewStringResponder(200, fmt.Sprintf(`{
		"id": "%s",
		"profile": {
			"firstName": "First",
			"lastName": "Last",
			"email": "email@email.com",
			"login": "email@email.com"
		}
	}`, oktaID)))

	httpmock.Activate()
}
