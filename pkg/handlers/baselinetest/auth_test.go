package baselinetest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trussworks/httpbaselinetest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type AuthSuite struct {
	*BaselineSuite
}

func TestAuthSuite(t *testing.T) {
	s := &AuthSuite{NewBaselineSuite(t)}
	suite.Run(t, s)
	s.PopTestSuite.TearDown()
}

func (suite *AuthSuite) TestLoginGovRedirect() {
	setupFunc := func(name string, btest *httpbaselinetest.HTTPBaselineTest) error {

		user := testdatagen.MakeDefaultUser(suite.DB())
		// user is in office_users but has never logged into the app
		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				Active: true,
				UserID: &user.ID,
			},
			User: user,
		})

		fakeToken := "some_token"

		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          user.ID,
			IDToken:         fakeToken,
			Hostname:        suite.appNames.OfficeServername,
			Email:           officeUser.Email,
		}

		// login.gov state cookie
		cookieName := authentication.StateCookieName(&session)
		cookie := http.Cookie{
			Name:    cookieName,
			Value:   "some mis-matched hash value",
			Path:    "/",
			Expires: auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes),
		}
		btest.Cookies = []http.Cookie{cookie}

		btest.Handler = suite.RoutingForTest()

		btest.Db = suite.getSqlxDb()

		return nil
	}

	suite.BaselineTestSuite.Run("GET unauthenticated /auth/login-gov", httpbaselinetest.HTTPBaselineTest{
		Setup:  setupFunc,
		Method: http.MethodGet,
		Path:   "/auth/login-gov",
		Host:   suite.appNames.OfficeServername,
		Tables: []string{"office_users", "users"},
	})
}
