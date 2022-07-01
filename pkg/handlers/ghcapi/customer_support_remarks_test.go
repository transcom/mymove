package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	customersupportremarksop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	remarksservice "github.com/transcom/mymove/pkg/services/customer_support_remarks"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListCustomerRemarksForMoveHandler() {

	setupTestData := func() (services.CustomerSupportRemarksFetcher, models.CustomerSupportRemark) {

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
		expectedCustomerSupportRemark.Move = move

		return fetcher, expectedCustomerSupportRemark
	}

	suite.Run("Successful list fetch", func() {
		fetcher, remark := setupTestData()
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/customer-support-remarks/", remark.Move.Locator), nil)
		params := customersupportremarksop.GetCustomerSupportRemarksForMoveParams{
			HTTPRequest: request,
			Locator:     remark.Move.Locator,
		}
		handlerConfig := suite.HandlerConfig()
		handler := ListCustomerSupportRemarksHandler{
			HandlerConfig:                 handlerConfig,
			CustomerSupportRemarksFetcher: fetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&customersupportremarksop.GetCustomerSupportRemarksForMoveOK{}, response)
		responsePayload := response.(*customersupportremarksop.GetCustomerSupportRemarksForMoveOK)
		suite.Equal(remark.ID.String(), responsePayload.Payload[0].ID.String())
		suite.Equal(remark.OfficeUserID.String(), responsePayload.Payload[0].OfficeUserID.String())
	})

	suite.Run("404 fetch response", func() {
		fetcher, remark := setupTestData()
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/customer-support-remarks/", remark.Move.Locator), nil)
		params := customersupportremarksop.GetCustomerSupportRemarksForMoveParams{
			HTTPRequest: request,
			Locator:     "ZZZZZZZZ",
		}
		handlerConfig := suite.HandlerConfig()
		handler := ListCustomerSupportRemarksHandler{
			HandlerConfig:                 handlerConfig,
			CustomerSupportRemarksFetcher: fetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&customersupportremarksop.GetCustomerSupportRemarksForMoveNotFound{}, response)
	})

}

func (suite *HandlerSuite) TestCreateCustomerSupportRemarksHandler() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	suite.Run("Successful POST", func() {
		handlerConfig := suite.HandlerConfig()

		creator := &mocks.CustomerSupportRemarksCreator{}
		handler := CreateCustomerSupportRemarksHandler{handlerConfig, creator}

		request := httptest.NewRequest("POST", fmt.Sprintf("/moves/%s/customer-support-remarks/", move.Locator), nil)

		params := customersupportremarksop.CreateCustomerSupportRemarkForMoveParams{
			HTTPRequest: request,
			Locator:     move.Locator,
		}

		remarkID := uuid.Must(uuid.NewV4())
		returnRemark := models.CustomerSupportRemark{
			ID:           remarkID,
			MoveID:       move.ID,
			Move:         move,
			Content:      "This is a customer support remark.",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			OfficeUser:   officeUser,
			OfficeUserID: officeUser.ID,
		}

		creator.On("CreateCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.CustomerSupportRemark"),
			mock.AnythingOfType("string"),
		).Return(&returnRemark, nil).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.CreateCustomerSupportRemarkForMoveOK{}, response)
	})

	suite.Run("unsuccessful POST", func() {
		handlerConfig := suite.HandlerConfig()

		creator := &mocks.CustomerSupportRemarksCreator{}
		handler := CreateCustomerSupportRemarksHandler{handlerConfig, creator}

		request := httptest.NewRequest("POST", fmt.Sprintf("/moves/%s/customer-support-remarks/", move.Locator), nil)

		params := customersupportremarksop.CreateCustomerSupportRemarkForMoveParams{
			HTTPRequest: request,
			Locator:     move.Locator,
		}

		creator.On("CreateCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.CustomerSupportRemark"),
			mock.AnythingOfType("string"),
		).Return(nil, fmt.Errorf("error")).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.CreateCustomerSupportRemarkForMoveInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateCustomerSupportRemarksHandler() {

	setupTestData := func() (*mocks.CustomerSupportRemarkUpdater, models.CustomerSupportRemark, models.CustomerSupportRemark) {

		updater := mocks.CustomerSupportRemarkUpdater{}
		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		originalRemark := testdatagen.MakeCustomerSupportRemark(suite.DB(), testdatagen.Assertions{
			CustomerSupportRemark: models.CustomerSupportRemark{
				Content:      "This is a customer support remark.",
				OfficeUserID: officeUser.ID,
				MoveID:       move.ID,
			},
			Move: move,
		})

		updatedRemark := originalRemark
		updatedRemark.Content = "Changed my mind"

		return &updater, originalRemark, updatedRemark

	}

	suite.Run("Successful PATCH", func() {
		updater, ogRemark, updatedRemark := setupTestData()

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/customer-support-remarks/%s", &ogRemark.ID), nil)
		payload := ghcmessages.UpdateCustomerSupportRemarkPayload{
			Content: &updatedRemark.Content,
		}

		params := customersupportremarksop.UpdateCustomerSupportRemarkForMoveParams{
			HTTPRequest:             request,
			Body:                    &payload,
			CustomerSupportRemarkID: strfmt.UUID(ogRemark.ID.String()),
		}

		handlerConfig := suite.HandlerConfig()
		handler := UpdateCustomerSupportRemarkHandler{handlerConfig, updater}

		updater.On("UpdateCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("customer_support_remarks.UpdateCustomerSupportRemarkForMoveParams"),
		).Return(&updatedRemark, nil)

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.UpdateCustomerSupportRemarkForMoveOK{}, response)
	})

	suite.Run("unsuccessful PATCH", func() {
		updater, _, updatedRemark := setupTestData()
		badRemarkID := uuid.Must(uuid.NewV4())

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/customer-support-remarks/%s", badRemarkID), nil)
		payload := ghcmessages.UpdateCustomerSupportRemarkPayload{
			Content: &updatedRemark.Content,
		}

		handlerConfig := suite.HandlerConfig()
		handler := UpdateCustomerSupportRemarkHandler{handlerConfig, updater}

		params := customersupportremarksop.UpdateCustomerSupportRemarkForMoveParams{
			HTTPRequest: request,
			Body:        &payload,
		}

		updater.On("UpdateCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("customer_support_remarks.UpdateCustomerSupportRemarkForMoveParams"),
		).Return(nil, fmt.Errorf("error"))

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.UpdateCustomerSupportRemarkForMoveInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestDeleteCustomerSupportRemarksHandler() {
	suite.Run("Successful DELETE", func() {
		remarkID := uuid.Must(uuid.NewV4())

		deleter := &mocks.CustomerSupportRemarkDeleter{}
		handlerConfig := suite.HandlerConfig()
		handler := DeleteCustomerSupportRemarkHandler{handlerConfig, deleter}

		request := httptest.NewRequest("DELETE", fmt.Sprintf("/customer-support-remarks/%s/", remarkID.String()), nil)

		params := customersupportremarksop.DeleteCustomerSupportRemarkParams{
			HTTPRequest:             request,
			CustomerSupportRemarkID: *handlers.FmtUUID(remarkID),
		}

		deleter.On("DeleteCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.DeleteCustomerSupportRemarkNoContent{}, response)
	})

	suite.Run("unsuccessful DELETE", func() {
		remarkID := uuid.Must(uuid.NewV4())

		deleter := &mocks.CustomerSupportRemarkDeleter{}
		handlerConfig := suite.HandlerConfig()
		handler := DeleteCustomerSupportRemarkHandler{handlerConfig, deleter}

		request := httptest.NewRequest("DELETE", fmt.Sprintf("/customer-support-remarks/%s/", remarkID.String()), nil)

		params := customersupportremarksop.DeleteCustomerSupportRemarkParams{
			HTTPRequest:             request,
			CustomerSupportRemarkID: *handlers.FmtUUID(remarkID),
		}

		deleter.On("DeleteCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(fmt.Errorf("error")).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.DeleteCustomerSupportRemarkInternalServerError{}, response)
	})
}
