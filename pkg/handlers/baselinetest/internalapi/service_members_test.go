package internalapi

import (
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/suite"
	"github.com/trussworks/httpbaselinetest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/baselinetest"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type HandlerSuite struct {
	*baselinetest.BaselineSuite
}

func TestHandlerSuite(t *testing.T) {
	s := &HandlerSuite{baselinetest.NewBaselineSuite(t)}
	suite.Run(t, s)
	s.PopTestSuite.TearDown()
}

func fakeAddressPayload() *internalmessages.Address {
	return &internalmessages.Address{
		StreetAddress1: swag.String("An address"),
		StreetAddress2: swag.String("Apt. 2"),
		StreetAddress3: swag.String("address line 3"),
		City:           swag.String("Happytown"),
		State:          swag.String("AL"),
		PostalCode:     swag.String("01234"),
	}
}

func (suite *HandlerSuite) TestSubmitServiceMemberHandlerAllValues() {
	setupFunc := func(name string, btest *httpbaselinetest.HTTPBaselineTest) error {

		btest.Handler = suite.RoutingForTest()

		btest.Db = suite.GetSqlxDb()

		// Given: A logged-in user
		user := testdatagen.MakeDefaultUser(suite.DB())
		btest.Path = "/internal/service_members"
		btest.Cookies = suite.CookiesForUser(user)
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		for i := range btest.Cookies {
			if btest.Cookies[i].Name == auth.MaskedGorillaCSRFToken {
				headers["X-CSRF-Token"] = btest.Cookies[i].Value
			}
		}
		btest.Headers = headers

		// When: a new ServiceMember is posted
		newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
			UserID:               strfmt.UUID(user.ID.String()),
			Edipi:                swag.String("9999999999"),
			FirstName:            swag.String("random string bla"),
			MiddleName:           swag.String("random string bla"),
			LastName:             swag.String("random string bla"),
			Suffix:               swag.String("random string bla"),
			Telephone:            swag.String("555-555-5555"),
			SecondaryTelephone:   swag.String("555-555-1234"),
			PersonalEmail:        swag.String("wml@example.com"),
			PhoneIsPreferred:     swag.Bool(false),
			EmailIsPreferred:     swag.Bool(true),
			ResidentialAddress:   fakeAddressPayload(),
			BackupMailingAddress: fakeAddressPayload(),
		}
		body, err := newServiceMemberPayload.MarshalBinary()
		if err != nil {
			return err
		}
		btest.Body = string(body)
		return nil
	}

	suite.BaselineTestSuite.Run("POST /service_members", httpbaselinetest.HTTPBaselineTest{
		Setup:  setupFunc,
		Method: http.MethodPost,
		Host:   suite.AppNames.MilServername,
		Tables: []string{"addresses", "service_members"},
	})
}
