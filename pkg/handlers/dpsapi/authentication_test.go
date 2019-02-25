package dpsapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/dpsapi/dpsoperations/dps"
	"github.com/transcom/mymove/pkg/gen/dpsmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetUserHandler() {
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetIWSPersonLookup(iws.TestingPersonLookup{})
	dpsParams := dpsauth.Params{
		CookieSecret:  []byte("cookie secret"),
		CookieExpires: 60,
	}
	context.SetDPSAuthParams(dpsParams)
	handler := GetUserHandler{context}

	affiliation := models.AffiliationARMY
	middleName := "Test"
	suffix := "II"
	smData := models.ServiceMember{
		Affiliation: &affiliation,
		Edipi:       models.StringPointer("1234567890"),
		MiddleName:  &middleName,
		Suffix:      &suffix,
		Telephone:   models.StringPointer("555-555-5555"),
	}
	assertions := testdatagen.Assertions{ServiceMember: smData}
	serviceMember := testdatagen.MakeServiceMember(suite.DB(), assertions)
	loginGovID := serviceMember.User.LoginGovUUID.String()
	cookie, err := dpsauth.LoginGovIDToCookie(loginGovID, dpsParams.CookieSecret, dpsParams.CookieExpires)
	suite.Nil(err)

	request := httptest.NewRequest("GET", "/dps/v0/authentication/user", nil)
	params := dps.GetUserParams{Token: cookie.Value}
	params.HTTPRequest = request

	// Missing client certificate
	response := handler.Handle(params)
	_, ok := response.(*dps.GetUserUnauthorized)
	suite.True(ok)

	// With client certificate
	digest := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	subject := "/C=US/ST=DC/L=Washington/O=Test/OU=Test Cert/CN=localhost"
	clientCert := models.ClientCert{
		Sha256Digest:    digest,
		Subject:         subject,
		AllowDpsAuthAPI: true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	requestContext := authentication.SetClientCertInRequestContext(request, &clientCert)
	params.HTTPRequest = request.WithContext(requestContext)
	response = handler.Handle(params)
	okResponse := response.(*dps.GetUserOK)
	suite.Equal(*okResponse.Payload.Affiliation, dpsmessages.AffiliationArmy)
	suite.Equal(okResponse.Payload.Email, serviceMember.User.LoginGovEmail)
	suite.Equal(okResponse.Payload.FirstName, *serviceMember.FirstName)
	suite.Equal(okResponse.Payload.LastName, *serviceMember.LastName)
	suite.Equal(okResponse.Payload.LoginGovID, strfmt.UUID(loginGovID))
	suite.Equal(*okResponse.Payload.MiddleName, *serviceMember.MiddleName)
	suite.Equal(okResponse.Payload.SocialSecurityNumber, strfmt.SSN(iws.SSN))
	suite.Equal(*okResponse.Payload.Suffix, *serviceMember.Suffix)
	suite.Equal(*okResponse.Payload.Telephone, *serviceMember.Telephone)
}
