package ghcapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models/roles"
	customerservice "github.com/transcom/mymove/pkg/services/office_user/customer"

	customerops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestGetCustomerHandlerIntegration() {
	customer := testdatagen.MakeDefaultServiceMember(suite.DB())

	request := httptest.NewRequest("GET", "/customer/{customerID}", nil)
	params := customerops.GetCustomerParams{
		HTTPRequest: request,
		CustomerID:  strfmt.UUID(customer.ID.String()),
	}
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := GetCustomerHandler{
		handlerConfig,
		customerservice.NewCustomerFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	getCustomerResponse := response.(*customerops.GetCustomerOK)
	getCustomerPayload := getCustomerResponse.Payload
	suite.Assertions.IsType(&customerops.GetCustomerOK{}, response)
	suite.Equal(strfmt.UUID(customer.ID.String()), getCustomerPayload.ID)
	suite.Equal(*customer.Edipi, getCustomerPayload.DodID)
	suite.Equal(strfmt.UUID(customer.UserID.String()), getCustomerPayload.UserID)
	suite.Equal(customer.Affiliation.String(), getCustomerPayload.Agency)
	suite.Equal(customer.PersonalEmail, getCustomerPayload.Email)
	suite.Equal(customer.Telephone, getCustomerPayload.Phone)
	suite.NotZero(getCustomerPayload.CurrentAddress)
}

func (suite *HandlerSuite) TestUpdateCustomerHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeServicesCounselor,
	})
	body := &ghcmessages.UpdateCustomerPayload{
		LastName:  "Newlastname",
		FirstName: "Newfirstname",
		Phone:     handlers.FmtString("123-455-3399"),
		CurrentAddress: &ghcmessages.Address{
			StreetAddress1: handlers.FmtString("123 New Street"),
			City:           handlers.FmtString("Newcity"),
			State:          handlers.FmtString("MA"),
			PostalCode:     handlers.FmtString("12345"),
		},
		BackupContact: &ghcmessages.BackupContact{
			Name:  handlers.FmtString("New Backup Contact"),
			Phone: handlers.FmtString("445-345-1212"),
			Email: handlers.FmtString("newbackup@mail.com"),
		},
	}
	customer := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{})
	request := httptest.NewRequest("PATCH", "/orders/{customerID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := customerops.UpdateCustomerParams{
		HTTPRequest: request,
		CustomerID:  strfmt.UUID(customer.ID.String()),
		IfMatch:     etag.GenerateEtag(customer.UpdatedAt),
		Body:        body,
	}
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := UpdateCustomerHandler{
		handlerConfig,
		customerservice.NewCustomerUpdater(),
	}
	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	// TODO: test with actual updated customer?
	// updatedCustomer, _ := models.FetchServiceMember(suite.DB(), customer.ID)
	updateCustomerResponse := response.(*customerops.UpdateCustomerOK)
	updateCustomerPayload := updateCustomerResponse.Payload
	suite.Assertions.IsType(&customerops.UpdateCustomerOK{}, response)
	suite.Equal(body.FirstName, updateCustomerPayload.FirstName)
	suite.Equal(body.LastName, updateCustomerPayload.LastName)
	suite.Equal(body.Phone, updateCustomerPayload.Phone)
	suite.Equal(body.CurrentAddress.StreetAddress1, updateCustomerPayload.CurrentAddress.StreetAddress1)
	suite.Equal(body.CurrentAddress.City, updateCustomerPayload.CurrentAddress.City)
	suite.Equal(body.CurrentAddress.PostalCode, updateCustomerPayload.CurrentAddress.PostalCode)
	suite.Equal(body.CurrentAddress.State, updateCustomerPayload.CurrentAddress.State)
	suite.Equal(body.BackupContact.Name, updateCustomerPayload.BackupContact.Name)
	suite.Equal(body.BackupContact.Phone, updateCustomerPayload.BackupContact.Phone)
	suite.Equal(body.BackupContact.Email, updateCustomerPayload.BackupContact.Email)
}
