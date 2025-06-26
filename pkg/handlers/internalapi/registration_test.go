package internalapi

import (
	"fmt"
	"net/http/httptest"
	"os"

	"github.com/jarcoal/httpmock"
	"github.com/markbates/goth"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	registrationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/registration"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
)

const milProviderName = "milProvider"

func (suite *HandlerSuite) TestCustomerRegistrationHandler() {
	affiliation := internalmessages.AffiliationNAVY
	suite.Run("Successful registration with existing exact matched okta user", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		// mocking the okta customer group id env variable
		originalGroupID := os.Getenv("OKTA_CUSTOMER_GROUP_ID")
		os.Setenv("OKTA_CUSTOMER_GROUP_ID", "notrealcustomergroupId")
		defer os.Setenv("OKTA_CUSTOMER_GROUP_ID", originalGroupID)

		// these mocked endpoints fetch an exact user
		mockAndActivateOktaGETEndpointExistingUserNoError(provider)
		mockAndActivateOktaGroupGETEndpointNoError(provider)
		mockAndActivateOktaGroupAddEndpointNoError(provider)

		affiliation := internalmessages.AffiliationARMY
		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "First",
			MiddleInitial:    models.StringPointer("M"),
			LastName:         "Last",
			Telephone:        "918-867-5309",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{
			handlerConfig,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationCreated{}, response)
	})

	suite.Run("Successful okta creation but existing user found when dodid_unique flag is on", func() {
		// creating a SM with an existing EDIPI
		factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Edipi: models.StringPointer("1234567890"),
				},
			},
		}, nil)

		// setting the dodid_unique flag to true
		os.Setenv("FEATURE_FLAG_DODID_UNIQUE", "true")
		originalGroupID := os.Getenv("OKTA_CUSTOMER_GROUP_ID")
		os.Setenv("OKTA_CUSTOMER_GROUP_ID", "notrealcustomergroupId")
		defer os.Setenv("OKTA_CUSTOMER_GROUP_ID", originalGroupID)

		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		// these mocked endpoints fetch an exact user
		mockAndActivateOktaGETEndpointExistingUserNoError(provider)
		mockAndActivateOktaGroupGETEndpointNoError(provider)
		mockAndActivateOktaGroupAddEndpointNoError(provider)

		affiliation := internalmessages.AffiliationARMY
		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "First",
			MiddleInitial:    models.StringPointer("M"),
			LastName:         "Last",
			Telephone:        "918-867-5309",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{
			handlerConfig,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)
		errResponse, ok := response.(*registrationop.CustomerRegistrationUnprocessableEntity)
		suite.True(ok)
		suite.Contains(*errResponse.Payload.Detail, "there is already an existing MilMove user with this DoD ID - an Okta account has also been found or created, please try signing into MilMove instead")
	})

	suite.Run("Successful registration when Okta user is created", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		mockAndActivateOktaGETEndpointNoUserNoError(provider)
		mockAndActivateOktaPOSTEndpointsNoError(provider)

		affiliation := internalmessages.AffiliationAIRFORCE
		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "First",
			MiddleInitial:    models.StringPointer("M"),
			LastName:         "Last",
			Telephone:        "918-867-5309",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationCreated{}, response)
	})

	suite.Run("Fail when email matches but EDIPI does not - single okta user", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		mockAndActivateOktaGETEndpointMatchingEmailOnlyNoError(provider)

		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "New",
			LastName:         "User",
			Telephone:        "555-555-1234",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)

		errResponse, ok := response.(*registrationop.CustomerRegistrationUnprocessableEntity)
		suite.True(ok)
		suite.Contains(*errResponse.Payload.Detail, "there is an existing Okta account with that email - please update the DoD ID (EDIPI) in your Okta profile to match your registration DoD ID and try registering again")
	})

	suite.Run("Fail when email matches but EDIPI does not - two okta users", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		mockAndActivateOktaGETEndpointMatchingEmailOnlyNoErrorTwoUsers(provider)

		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "New",
			LastName:         "User",
			Telephone:        "555-555-1234",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)

		errResponse, ok := response.(*registrationop.CustomerRegistrationUnprocessableEntity)
		suite.True(ok)
		suite.Contains(*errResponse.Payload.Detail, "email and DoD IDs match different users - please open up a help desk ticket")
	})

	suite.Run("Fail when EDIPI matches but email does not - single okta user", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		mockAndActivateOktaGETEndpointMatchingEdipiOnlyNoError(provider)

		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "New",
			LastName:         "User",
			Telephone:        "555-555-1234",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)

		errResponse, ok := response.(*registrationop.CustomerRegistrationUnprocessableEntity)
		suite.True(ok)
		suite.Contains(*errResponse.Payload.Detail, "there is an existing Okta account with that DoD ID (EDIPI) - please update the email in your Okta profile to match your registration email and try registering again")
	})

	suite.Run("Fail when email matches but EDIPI does not - two okta users", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		mockAndActivateOktaGETEndpointMatchingEdipiOnlyNoErrorTwoUsers(provider)

		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "New",
			LastName:         "User",
			Telephone:        "555-555-1234",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)

		errResponse, ok := response.(*registrationop.CustomerRegistrationUnprocessableEntity)
		suite.True(ok)
		suite.Contains(*errResponse.Payload.Detail, "email and DoD IDs match different users - please open up a help desk ticket")
	})

	suite.Run("Fail when email matches AND EDIPI matches but on different users - two okta users", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		mockAndActivateOktaGETEndpointMismatchedNoErrorTwoUsers(provider)

		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "New",
			LastName:         "User",
			Telephone:        "555-555-1234",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("1234567890"),
			Email:            "email@email.com",
			PhoneIsPreferred: true,
			EmailIsPreferred: false,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)

		errResponse, ok := response.(*registrationop.CustomerRegistrationUnprocessableEntity)
		suite.True(ok)
		suite.Contains(*errResponse.Payload.Detail, "there are multiple Okta accounts with that email and DoD ID - please open up a help desk ticket")
	})

	suite.Run("Error when Okta user creation fails", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		// mocking an error response when attempting to create an Okta user
		httpmock.RegisterResponder("POST", provider.GetCreateUserURL("true"),
			httpmock.NewStringResponder(400, `{ "error": "Invalid request" }`))
		httpmock.Activate()

		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "ErrorTest",
			LastName:         "User",
			Telephone:        "555-555-5555",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("0987654321"),
			Email:            *handlers.FmtString("error@example.com"),
			PhoneIsPreferred: false,
			EmailIsPreferred: true,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)
	})

	suite.Run("Fail when request does not come from MilMove app", func() {
		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.OfficeApp, // we only allow requests from the cust app
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest: request.WithContext(ctx),
		}

		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)
	})

	suite.Run("Fail when required fields are missing", func() {
		provider, err := factory.BuildOktaProvider(milProviderName)
		suite.NoError(err)

		body := &internalmessages.CreateOktaAndMilMoveUser{
			FirstName:        "", // missing first name
			LastName:         "", // missing last name
			Telephone:        "555-555-5555",
			Affiliation:      &affiliation,
			Edipi:            models.StringPointer("0987654321"),
			Email:            *handlers.FmtString("error@example.com"),
			PhoneIsPreferred: false,
			EmailIsPreferred: true,
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		request := httptest.NewRequest("POST", "/open/register", nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
		}
		ctx := auth.SetSessionInRequestContext(request, session)
		params := registrationop.CustomerRegistrationParams{
			HTTPRequest:  request.WithContext(ctx),
			Registration: body,
		}
		handlerConfig := suite.NewHandlerConfig()
		handler := CustomerRegistrationHandler{handlerConfig}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.Assertions.IsType(&registrationop.CustomerRegistrationUnprocessableEntity{}, response)
	})
}

// Generate and activate Okta endpoints that will be using during the auth handlers.
func mockAndActivateOktaGETEndpointNoUserNoError(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	response := "[]"

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaGETEndpointExistingUserNoError(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	oktaID := "fakeSub"
	response := fmt.Sprintf(`[
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567890"
			}
		}
	]`, oktaID)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaGETEndpointMatchingEmailOnlyNoError(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	oktaID := "fakeSub"
	response := fmt.Sprintf(`[
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "0987654321"
			}
		}
	]`, oktaID)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaGETEndpointMatchingEmailOnlyNoErrorTwoUsers(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	oktaID := "fakeSub"
	oktaID2 := "fakeSub2"
	response := fmt.Sprintf(`[
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": ""
			}
		},
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567891"
			}
		}
	]`, oktaID, oktaID2)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaGETEndpointMatchingEdipiOnlyNoError(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	oktaID := "fakeSub"
	response := fmt.Sprintf(`[
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "notYourEmail@email.com",
				"email": "notYourEmail@email.com",
				"cac_edipi": "1234567890"
			}
		}
	]`, oktaID)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaGETEndpointMatchingEdipiOnlyNoErrorTwoUsers(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	oktaID := "fakeSub"
	oktaID2 := "fakeSub2"
	response := fmt.Sprintf(`[
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "anotherEmail@email.com",
				"email": "anotherEmail@email.com",
				"cac_edipi": "1234567890"
			}
		},
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email2@email.com",
				"email": "email2@email.com",
				"cac_edipi": "1234567891"
			}
		}
	]`, oktaID, oktaID2)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaGETEndpointMismatchedNoErrorTwoUsers(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	oktaID := "fakeSub"
	oktaID2 := "fakeSub2"
	response := fmt.Sprintf(`[
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "anotherEmail@email.com",
				"email": "anotherEmail@email.com",
				"cac_edipi": "1234567890"
			}
		},
		{
			"id": "%s",
			"status": "PROVISIONED",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567891"
			}
		}
	]`, oktaID, oktaID2)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaPOSTEndpointsNoError(provider *okta.Provider) {

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

func mockAndActivateOktaGroupGETEndpointNoError(provider *okta.Provider) {

	oktaID := "fakeSub"
	getGroupsEndpoint := provider.GetUserGroupsURL(oktaID)

	httpmock.RegisterResponder("GET", getGroupsEndpoint,
		httpmock.NewStringResponder(200, `[]`))

	httpmock.Activate()
}

func mockAndActivateOktaGroupAddEndpointNoError(provider *okta.Provider) {

	oktaID := "fakeSub"
	groupID := "notrealcustomergroupId"
	addGroupEndpoint := provider.AddUserToGroupURL(groupID, oktaID)

	httpmock.RegisterResponder("PUT", addGroupEndpoint,
		httpmock.NewStringResponder(204, ""))

	httpmock.Activate()
}
