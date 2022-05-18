package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	customersupportremarksop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	remarksservice "github.com/transcom/mymove/pkg/services/customer_support_remarks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListCustomerRemarksForMoveHandler() {
	fetcher := remarksservice.NewCustomerSupportRemarks()
	move := testdatagen.MakeDefaultMove(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	expectedCustomerSupportRemark := testdatagen.MakeCustomerSupportRemark(suite.DB(), testdatagen.Assertions{
		CustomerSupportRemark: models.CustomerSupportRemark{
			Content:      "This is a customer support remark.",
			OfficeUserID: officeUser.ID,
			MoveID:       move.ID,
		},
	})

	suite.T().Run("Successful list fetch", func(t *testing.T) {
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/customer-support-remarks/", move.Locator), nil)
		params := customersupportremarksop.GetCustomerSupportRemarksForMoveParams{
			HTTPRequest: request,
			Locator:     move.Locator,
		}
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		handler := ListCustomerSupportRemarksHandler{
			HandlerConfig:                 handlerConfig,
			CustomerSupportRemarksFetcher: fetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&customersupportremarksop.GetCustomerSupportRemarksForMoveOK{}, response)
		responsePayload := response.(*customersupportremarksop.GetCustomerSupportRemarksForMoveOK)
		suite.Equal(expectedCustomerSupportRemark.ID.String(), responsePayload.Payload[0].ID.String())
		suite.Equal(officeUser.ID.String(), responsePayload.Payload[0].OfficeUserID.String())
	})

	suite.T().Run("404 fetch response", func(t *testing.T) {
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/customer-support-remarks/", move.Locator), nil)
		params := customersupportremarksop.GetCustomerSupportRemarksForMoveParams{
			HTTPRequest: request,
			Locator:     "ZZZZZZZZ",
		}
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		handler := ListCustomerSupportRemarksHandler{
			HandlerConfig:                 handlerConfig,
			CustomerSupportRemarksFetcher: fetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&customersupportremarksop.GetCustomerSupportRemarksForMoveNotFound{}, response)
	})

}
