package primeapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *PrimeAPISuite) TestSwaggerYaml() {
	routingConfig := suite.RoutingConfig()
	routingConfig.PrimeSwaggerPath = "foo/bar/baz"
	swaggerContent := "some\nswagger\ncontent\n"
	suite.CreateFileWithContent(routingConfig.PrimeSwaggerPath, swaggerContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	cert := factory.BuildPrimeClientCert(suite.DB())
	req := suite.NewAuthenticatedPrimeRequest("GET", "/prime/v1/swagger.yaml", nil, cert)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerContent, string(actualData))
}

func (suite *PrimeAPISuite) TestSwaggerUI() {
	routingConfig := suite.RoutingConfig()
	swaggerUIContent := "some\nswaggerUI\ncontent\n"
	swaggerUIPath := filepath.Join(routingConfig.BuildRoot, "swagger-ui", "prime.html")
	suite.CreateFileWithContent(swaggerUIPath, swaggerUIContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	cert := factory.BuildPrimeClientCert(suite.DB())
	req := suite.NewAuthenticatedPrimeRequest("GET", "/prime/v1/docs", nil, cert)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerUIContent, string(actualData))
}
