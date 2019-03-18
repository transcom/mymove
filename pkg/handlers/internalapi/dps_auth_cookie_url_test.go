package internalapi

import (
	"fmt"
	"net/http/httptest"
	"net/url"

	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/dps_auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestDPSAuthCookieURLHandler() {
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	dpsAuthParams := dpsauth.Params{
		SDDCProtocol:   "http",
		SDDCHostname:   "testhost",
		SDDCPort:       "100",
		SecretKey:      "secretkey",
		DPSRedirectURL: "http://example.com",
		CookieName:     "test",
	}
	context.SetDPSAuthParams(dpsAuthParams)
	handler := DPSAuthGetCookieURLHandler{context}

	// Normal service member (not a DPS user) happy path
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	request := httptest.NewRequest("POST", "/dps_auth/cookie_url", nil)
	request = suite.AuthenticateRequest(request, serviceMember)

	params := dps_auth.NewGetCookieURLParams()
	params.HTTPRequest = request

	response := handler.Handle(params)
	// TODO: Once the TODO is resolved in DPSAuthGetCookieURLHandler::Handle,
	// this should be changed to GetCookieURLOK
	_, ok := response.(*dps_auth.GetCookieURLForbidden)
	suite.True(ok)

	// Normal service member (not a DPS user) permission error
	dpsRedirectURL := "http://example.com"
	params.DpsRedirectURL = &dpsRedirectURL
	response = handler.Handle(params)
	_, ok = response.(*dps_auth.GetCookieURLForbidden)
	suite.True(ok)

	// Make the service member a DPS user, should no longer get a permission error when setting params
	dpsUser := testdatagen.MakeDpsUser(suite.DB(), testdatagen.Assertions{User: serviceMember.User})
	request = suite.AuthenticateDpsRequest(request, serviceMember, dpsUser)
	params.HTTPRequest = request
	response = handler.Handle(params)
	okResponse := response.(*dps_auth.GetCookieURLOK)
	url, _ := url.Parse(okResponse.Payload.CookieURL.String())
	suite.Equal(url.Scheme, dpsAuthParams.SDDCProtocol)
	suite.Equal(url.Host, fmt.Sprintf("%s:%s", dpsAuthParams.SDDCHostname, dpsAuthParams.SDDCPort))
	suite.Contains(url.Query(), "token")
	suite.True(ok)
}
