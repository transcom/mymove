package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ConfigSuite struct {
	*testingsuite.PopTestSuite
}

func TestConfigSuite(t *testing.T) {
	ts := &ConfigSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ConfigSuite) TestConfigHandler() {
	suite.Run("AuditableAppContextFromRequestWithErrors executes the function argument it is passed", func() {

		appCtx := suite.AppContextForTest()
		sessionManagers := auth.SetupSessionManagers(nil, false, time.Duration(180*time.Second), time.Duration(180*time.Second))
		handler := NewHandlerConfig(appCtx.DB(), nil, "", nil, nil, nil, nil, nil, false, nil, nil, false, ApplicationTestServername(), sessionManagers, nil)
		req, err := http.NewRequest("GET", "/", nil)
		suite.NoError(err)
		myMethodCalled := false
		myFunction := func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			myMethodCalled = true
			return nil, nil
		}
		handler.AuditableAppContextFromRequestWithErrors(req, myFunction)
		suite.Equal(myMethodCalled, true)
	})
}
