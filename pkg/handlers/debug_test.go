package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type HandlersSuite struct {
	BaseHandlerTestSuite
}

func TestHandlersSuite(t *testing.T) {
	hs := &HandlersSuite{
		BaseHandlerTestSuite: NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *HandlersSuite) TestDebugSessionsHandler() {
	user := testdatagen.MakeDefaultUser(suite.DB())

	handlerConfig := suite.HandlerConfig()
	sessionManagers := auth.SetupSessionManagers(nil, false,
		time.Duration(180), time.Duration(180))
	handlerConfig.SetSessionManagers(sessionManagers)

	milHost := handlerConfig.appNames.MilServername

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        milHost,
		UserID:          user.ID,
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/debug/sessions", milHost), nil)
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	h := NewDebugSessionsHandler(handlerConfig)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	body := rr.Body.Bytes()
	// as this is for debugging, just make sure it returns something
	suite.True(len(body) > 0, body)
}
