package ghcapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *GhcAPISuite) TestSwaggerYaml() {
	routingConfig := suite.RoutingConfig()
	routingConfig.GHCSwaggerPath = "foo/bar/baz"
	swaggerContent := "some\nswagger\ncontent\n"
	suite.CreateFileWithContent(routingConfig.GHCSwaggerPath, swaggerContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
	req := suite.NewAuthenticatedOfficeRequest("GET", "/ghc/v1/swagger.yaml", nil, officeUser)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerContent, string(actualData))
}

func (suite *GhcAPISuite) TestSwaggerUI() {
	routingConfig := suite.RoutingConfig()
	swaggerUIContent := "some\nswaggerUI\ncontent\n"
	swaggerUIPath := filepath.Join(routingConfig.BuildRoot, "swagger-ui", "ghc.html")
	suite.CreateFileWithContent(swaggerUIPath, swaggerUIContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
	req := suite.NewAuthenticatedOfficeRequest("GET", "/ghc/v1/docs", nil, officeUser)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerUIContent, string(actualData))
}
