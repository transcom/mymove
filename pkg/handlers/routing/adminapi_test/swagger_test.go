package adminapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *AdminAPISuite) TestSwaggerYaml() {
	routingConfig := suite.RoutingConfig()
	routingConfig.AdminSwaggerPath = "foo/bar/baz"
	swaggerContent := "some\nswagger\ncontent\n"
	suite.CreateFileWithContent(routingConfig.AdminSwaggerPath, swaggerContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	adminUser := factory.BuildAdminUser(suite.DB(), factory.GetTraitActiveAdminUser(), nil)
	req := suite.NewAuthenticatedAdminRequest("GET", "/admin/v1/swagger.yaml", nil, adminUser)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerContent, string(actualData))
}

func (suite *AdminAPISuite) TestSwaggerUI() {
	routingConfig := suite.RoutingConfig()
	swaggerUIContent := "some\nswaggerUI\ncontent\n"
	swaggerUIPath := filepath.Join(routingConfig.BuildRoot, "swagger-ui", "admin.html")
	suite.CreateFileWithContent(swaggerUIPath, swaggerUIContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	adminUser := factory.BuildAdminUser(suite.DB(), factory.GetTraitActiveAdminUser(), nil)
	req := suite.NewAuthenticatedAdminRequest("GET", "/admin/v1/docs", nil, adminUser)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerUIContent, string(actualData))
}
