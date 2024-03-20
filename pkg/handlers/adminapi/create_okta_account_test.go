package adminapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/jarcoal/httpmock"
	"github.com/markbates/goth"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/okta"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/okta"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	oktaAuth "github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
)

const DummyRSAModulus = "0OtoQx0UQHbkrlEA8YsZ-tW20S4_YgQZkRtN61tzzZ5Es63KH_crZymNi19gwD2kq_9RJu376oqL81YONxJXxRyQawrJCali6YYn7-qqBl9acLDwP0W_jAan7TFNWau1AvRIrP0o3tkBse5NNiaEMvkfxD_5EKtQdKeP6grUe90"
const adminProviderName = "adminProvider"

func (suite *HandlerSuite) TestCreateOktaAccountHandler2() {
	adminUser := factory.BuildAdminUser(suite.DB(), []factory.Customization{
		{
			Model: models.AdminUser{
				Active: true,
			},
		},
	}, []factory.Trait{
		factory.GetTraitActiveUser,
		factory.GetTraitAdminUserEmail,
	})
	user := adminUser.User

	// Build provider
	provider, err := factory.BuildOktaProvider(adminProviderName)
	suite.NoError(err)

	mockAndActivateOktaEndpoints(provider)

	firstName := "Micheal"
	lastName := "Jackson"
	email := "MJ2000@example.com"
	login := "MJ2000@example.com"
	cacEdipi := "1234567890"
	mobilePhone := "462-940-8555"
	gsaID := "string"

	body := &adminmessages.CreateOktaAccount{
		FirstName:   &firstName,
		LastName:    &lastName,
		Email:       &email,
		Login:       &login,
		CacEdipi:    &cacEdipi,
		MobilePhone: &mobilePhone,
		GsaID:       &gsaID,
		GroupID:     []string{},
	}

	defer goth.ClearProviders()
	goth.UseProviders(provider)

	request := httptest.NewRequest("POST", "/create-okta-account", nil)
	request = suite.AuthenticateAdminRequest(request, user)

	params := userop.CreateOktaAccountParams{
		HTTPRequest:              request,
		CreateOktaAccountPayload: body,
	}
	handlerConfig := suite.HandlerConfig()
	handler := CreateOktaAccount{
		handlerConfig,
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.Assertions.IsType(&okta.CreateOktaAccountOK{}, response)

	suite.Assertions.IsType(&okta.CreateOktaAccountOK{}, response)
	createAccountResponse := response.(*userop.CreateOktaAccountOK)
	createAccountPayload := createAccountResponse.Payload

	// Validate outgoing payload
	suite.NoError(createAccountPayload.Validate(strfmt.Default))

	suite.Equal(body.FirstName, createAccountPayload.FirstName)
	suite.Equal(body.LastName, createAccountPayload.LastName)
	suite.Equal(body.MobilePhone, createAccountPayload.MobilePhone)
	suite.Equal(body.Email, createAccountPayload.Email)
}

// Generate and activate Okta endpoints that will be using during the handler
func mockAndActivateOktaEndpoints(provider *oktaAuth.Provider) {
	activate := "true"

	createAccountEndpoint := provider.GetCreateAccountURL(activate)
	oktaID := "fakeSub"

	httpmock.RegisterResponder("POST", createAccountEndpoint,
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
