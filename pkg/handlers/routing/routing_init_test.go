package routing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
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

	suite.EqualDefaultIndex(rr)
}

func (suite *RoutingSuite) TestServeGHC() {

	serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
	routingConfig := suite.RoutingConfig()
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	// make the request with auth
	req := suite.NewAuthenticatedMilRequest("GET", fmt.Sprintf("/ghc/v1/customer/%s", serviceMember.ID.String()), nil, serviceMember)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	// the GHC API is not available to the Mil app, so the default
	// route is served for GET
	suite.EqualDefaultIndex(rr)

	// make the request without auth
	req = suite.NewMilRequest("GET", fmt.Sprintf("/ghc/v1/customer/%s", serviceMember.ID.String()), nil)
	rr = httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	// the GHC API is not available to the Mil app, so the default
	// route is served for GET
	suite.EqualDefaultIndex(rr)

	// make the request with GHC routing turned off
	routingConfig.ServeGHC = false
	noghcHandler := suite.SetupCustomSiteHandler(routingConfig)
	req = suite.NewMilRequest("GET", fmt.Sprintf("/ghc/v1/customer/%s", serviceMember.ID.String()), nil)
	rr = httptest.NewRecorder()
	noghcHandler.ServeHTTP(rr, req)
	// if the API is not enabled, the routing will be served by the
	// SPA handler, sending back the index page, which will have the
	// javascript SPA routing
	suite.EqualDefaultIndex(rr)
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

func (suite *RoutingSuite) TestBasicStorageRouting() {
	routingConfig := suite.RoutingConfig()
	routingConfig.LocalStorageWebRoot = "storage"

	// If LocalStorageRoot is "/path", and LocalStorageWebRoot is
	// "storage", the file handler expects to serve /storage/foo
	// from "/path/storage/foo"
	routingConfig.LocalStorageRoot = "path"
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	storateFileContents := "sample\nstorage\n"
	storageFileName := "test_storage.txt"
	storagePath := filepath.Join(routingConfig.LocalStorageRoot,
		routingConfig.LocalStorageWebRoot, storageFileName)
	suite.CreateFileWithContent(storagePath, storateFileContents)

	rpath := fmt.Sprintf("/%s/%s", routingConfig.LocalStorageWebRoot, storageFileName)
	req := suite.NewMilRequest("GET", rpath, nil)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(storateFileContents, rr.Body.String())
}

func (suite *RoutingSuite) TestBasicHealthRouting() {
	siteHandler := suite.SetupSiteHandler()
	req := suite.NewMilRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualDataString := rr.Body.String()
	suite.Contains(actualDataString, `"database"`)
	suite.Contains(actualDataString, `"gitBranch"`)
	suite.Contains(actualDataString, `"gitCommit"`)

	// test health check with IP host as that is what requests from
	// the AWS ELB look like
	req = httptest.NewRequest("GET", "http://1.2.3.4:8443/health", nil)
	rr = httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualDataString = rr.Body.String()
	suite.Contains(actualDataString, `"database"`)
	suite.Contains(actualDataString, `"gitBranch"`)
	suite.Contains(actualDataString, `"gitCommit"`)
}

func (suite *RoutingSuite) TestBasicStaticRouting() {
	routingConfig := suite.RoutingConfig()
	staticCSSFilename := "static.css"
	cssPath := filepath.Join(routingConfig.BuildRoot, "static", staticCSSFilename)
	cssContent := "body { font: inherit; }"
	suite.CreateFileWithContent(cssPath, cssContent)

	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	req := suite.NewMilRequest("GET", fmt.Sprintf("/static/%s", staticCSSFilename), nil)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(cssContent, rr.Body.String())
}

func (suite *RoutingSuite) TestBasicDownloadsRouting() {
	routingConfig := suite.RoutingConfig()
	downloadFilename := "download.txt"
	downloadPath := filepath.Join(routingConfig.BuildRoot, "downloads", downloadFilename)
	downloadContent := "some\ndownload\ncontent"
	suite.CreateFileWithContent(downloadPath, downloadContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	req := suite.NewMilRequest("GET", fmt.Sprintf("/downloads/%s", downloadFilename), nil)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(downloadContent, rr.Body.String())
}

func (suite *RoutingSuite) TestBasicAuthLoginRouting() {
	siteHandler := suite.SetupSiteHandler()

	serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
	// an authenticted request will redirect and we just want to check
	// that the route is set up correctly
	req := suite.NewAuthenticatedMilRequest("GET", "/auth/login-gov", nil, serviceMember)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusTemporaryRedirect, rr.Code)
	locationURL := rr.Result().Header.Get("Location")
	u, err := url.Parse(locationURL)
	suite.NoError(err)
	// The redirect contains both the host and the port and we're just
	// checking that this is the redirect we expected from the routing setup
	suite.Contains(u.Host, suite.RoutingConfig().HandlerConfig.AppNames().MilServername)
	suite.Equal("/", u.Path)
}

func (suite *RoutingSuite) TestBasicAuthLogoutRouting() {
	siteHandler := suite.SetupSiteHandler()

	serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
	// an authenticted request will redirect and we just want to check
	// that the route is set up correctly
	req := suite.NewAuthenticatedMilRequest("POST", "/auth/logout", nil, serviceMember)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	u, err := url.Parse(rr.Body.String())
	suite.NoError(err)
	suite.Contains(u.Host, suite.RoutingConfig().HandlerConfig.AppNames().MilServername)
}
