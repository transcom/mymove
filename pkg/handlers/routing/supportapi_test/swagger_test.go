package supportapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

func (suite *SupportAPISuite) TestSwaggerYaml() {
	routingConfig := suite.RoutingConfig()
	routingConfig.SupportSwaggerPath = "foo/bar/baz"
	swaggerContent := "some\nswagger\ncontent\n"
	suite.CreateFileWithContent(routingConfig.SupportSwaggerPath, swaggerContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	req := suite.NewTLSAuthenticatedPrimeRequest("GET", "/support/v1/swagger.yaml", nil)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerContent, string(actualData))
}

func (suite *SupportAPISuite) TestSwaggerUI() {
	routingConfig := suite.RoutingConfig()
	swaggerUIContent := "some\nswaggerUI\ncontent\n"
	swaggerUIPath := filepath.Join(routingConfig.BuildRoot, "swagger-ui", "support.html")
	suite.CreateFileWithContent(swaggerUIPath, swaggerUIContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	req := suite.NewTLSAuthenticatedPrimeRequest("GET", "/support/v1/docs", nil)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerUIContent, string(actualData))
}
