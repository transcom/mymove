package pptasapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *PPTASAPISuite) TestSwaggerYaml() {
	routingConfig := suite.RoutingConfig()
	routingConfig.PPTASSwaggerPath = "foo/bar/baz"
	swaggerContent := "some\nswagger\ncontent\n"
	suite.CreateFileWithContent(routingConfig.PPTASSwaggerPath, swaggerContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	cert := factory.BuildPrimeClientCert(suite.DB())
	req := suite.NewAuthenticatedPrimeRequest("GET", "/pptas/v1/swagger.yaml", nil, cert)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerContent, string(actualData))
}
