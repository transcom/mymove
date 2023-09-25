package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	oktaop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/okta_profile"
)

func (suite *HandlerSuite) TestGetOktaProfileHandler() {

	// Given: A logged-in user
	user := factory.BuildServiceMember(suite.DB(), nil, nil)
	suite.MustSave(&user)

	req := httptest.NewRequest("GET", "/okta_profile", nil)
	req = suite.AuthenticateRequest(req, user)

	params := oktaop.ShowOktaInfoParams{
		HTTPRequest: req,
	}

	handler := GetOktaProfileHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	// TODO figure out how to write this test correctly
	suite.Assertions.IsType(nil, response)
}
