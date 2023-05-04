package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models/roles"
)

type RoutingSuite struct {
	BaseRoutingSuite
}

func TestRoutingSuite(t *testing.T) {
	hs := &RoutingSuite{
		NewBaseRoutingSuite(),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *RoutingSuite) TestBasicRoutingInit() {

	req := suite.NewMilRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	suite.SetupSiteHandler().ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(suite.indexContent, rr.Body.String())
}

func (suite *RoutingSuite) TestServeGHC() {

	user := factory.BuildUser(suite.DB(), nil, nil)
	routingConfig := suite.RoutingConfig()
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	// make the request without auth
	req := suite.NewMilRequest("GET", fmt.Sprintf("/ghc/v1/customer/%s", user.ID.String()), nil)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusUnauthorized, rr.Code)

	// make the request with GHC routing turned off
	routingConfig.ServeGHC = false
	noghcHandler := suite.SetupCustomSiteHandler(routingConfig)
	req = suite.NewMilRequest("GET", fmt.Sprintf("/ghc/v1/customer/%s", user.ID.String()), nil)
	rr = httptest.NewRecorder()
	noghcHandler.ServeHTTP(rr, req)
	// if the API is not enabled, the routing will be served by the
	// SPA handler, sending back the index page, which will have the
	// javascript SPA routing
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(suite.indexContent, rr.Body.String())
}

func (suite *RoutingSuite) TestOfficeLoggedInEndpoint() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
	req := suite.NewAuthenticatedOfficeRequest("GET", "/internal/users/logged_in", nil, officeUser)

	rr := httptest.NewRecorder()
	suite.SetupSiteHandler().ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	var userPayload internalmessages.LoggedInUserPayload
	suite.NoError(json.Unmarshal(rr.Body.Bytes(), &userPayload))
	suite.Equal(officeUser.UserID.String(), userPayload.ID.String())
	suite.NotNil(userPayload.OfficeUser)
	suite.Equal(officeUser.ID.String(), userPayload.OfficeUser.ID.String())
	suite.NotEmpty(userPayload.Permissions)
}
