package dpsapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/dpsapi/dpsoperations/dps"
	"github.com/transcom/mymove/pkg/gen/dpsmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetUserHandler() {
	context := handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())
	context.SetIWSPersonLookup(iws.TestingPersonLookup{})
	handler := GetUserHandler{context}

	affiliation := models.AffiliationARMY
	middleName := "Test"
	suffix := "II"
	smData := models.ServiceMember{
		Affiliation: &affiliation,
		Edipi:       models.StringPointer("1234567890"),
		MiddleName:  &middleName,
		Suffix:      &suffix,
		Telephone:   models.StringPointer("555-555-5555"),
	}
	assertions := testdatagen.Assertions{ServiceMember: smData}
	serviceMember := testdatagen.MakeServiceMember(suite.TestDB(), assertions)
	loginGovID := serviceMember.User.LoginGovUUID.String()
	cookie, err := dpsauth.LoginGovIDToCookie(loginGovID)
	suite.Nil(err)

	request := httptest.NewRequest("GET", "/dps/v0/authentication/user", nil)
	params := dps.GetUserParams{Token: cookie.Value}
	params.HTTPRequest = request

	response := handler.Handle(params)
	okResponse := response.(*dps.GetUserOK)
	suite.Equal(*okResponse.Payload.Affiliation, dpsmessages.AffiliationArmy)
	suite.Equal(okResponse.Payload.Email, serviceMember.User.LoginGovEmail)
	suite.Equal(okResponse.Payload.FirstName, *serviceMember.FirstName)
	suite.Equal(okResponse.Payload.LastName, *serviceMember.LastName)
	suite.Equal(okResponse.Payload.LoginGovID, strfmt.UUID(loginGovID))
	suite.Equal(*okResponse.Payload.MiddleName, *serviceMember.MiddleName)
	suite.Equal(okResponse.Payload.SocialSecurityNumber, strfmt.SSN(iws.SSN))
	suite.Equal(*okResponse.Payload.Suffix, *serviceMember.Suffix)
	suite.Equal(*okResponse.Payload.Telephone, *serviceMember.Telephone)
}
