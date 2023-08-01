package internalapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *InternalAPISuite) TestSwaggerYaml() {
	routingConfig := suite.RoutingConfig()
	routingConfig.APIInternalSwaggerPath = "foo/bar/baz"
	swaggerContent := "some\nswagger\ncontent\n"
	suite.CreateFileWithContent(routingConfig.APIInternalSwaggerPath, swaggerContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
	req := suite.NewAuthenticatedMilRequest("GET", "/internal/swagger.yaml", nil, serviceMember)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerContent, string(actualData))
}

func (suite *InternalAPISuite) TestSwaggerUI() {
	routingConfig := suite.RoutingConfig()
	swaggerUIContent := "some\nswaggerUI\ncontent\n"
	swaggerUIPath := filepath.Join(routingConfig.BuildRoot, "swagger-ui", "internal.html")
	suite.CreateFileWithContent(swaggerUIPath, swaggerUIContent)
	siteHandler := suite.SetupCustomSiteHandler(routingConfig)

	serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
	req := suite.NewAuthenticatedMilRequest("GET", "/internal/docs", nil, serviceMember)
	rr := httptest.NewRecorder()
	siteHandler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	actualData, err := io.ReadAll(rr.Body)
	suite.NoError(err)
	suite.Equal(swaggerUIContent, string(actualData))
}
