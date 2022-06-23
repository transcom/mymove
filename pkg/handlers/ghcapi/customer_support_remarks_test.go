package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	customersupportremarksop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
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
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

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

	suite.Run("Successful PATCH", func() {
		locator := models.GenerateLocator()
		remarkID := uuid.Must(uuid.NewV4())

		remark := models.CustomerSupportRemark{
			ID: remarkID,
		}

		updatedRemarkText := "This is an updated customer support remark."
		id := strfmt.UUID(remarkID.String())
		body := &ghcmessages.UpdateCustomerSupportRemarkPayload{
			Content: &updatedRemarkText,
			ID:      &id,
		}

		updater := &mocks.CustomerSupportRemarkUpdater{}
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		handler := UpdateCustomerSupportRemarkHandler{handlerConfig, updater}

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/moves/%s/customer-support-remarks/", locator), nil)

		params := customersupportremarksop.UpdateCustomerSupportRemarkForMoveParams{
			HTTPRequest: request,
			Locator:     locator,
			Body:        body,
		}

		updater.On("UpdateCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(&remark, nil).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.UpdateCustomerSupportRemarkForMoveOK{}, response)
	})

	suite.Run("unsuccessful PATCH", func() {

		locator := models.GenerateLocator()
		remarkID := uuid.Must(uuid.NewV4())

		updatedRemarkText := "This is an updated customer support remark."
		id := strfmt.UUID(remarkID.String())
		body := &ghcmessages.UpdateCustomerSupportRemarkPayload{
			Content: &updatedRemarkText,
			ID:      &id,
		}

		updater := &mocks.CustomerSupportRemarkUpdater{}
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		handler := UpdateCustomerSupportRemarkHandler{handlerConfig, updater}

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/moves/%s/customer-support-remarks/", locator), nil)

		params := customersupportremarksop.UpdateCustomerSupportRemarkForMoveParams{
			HTTPRequest: request,
			Locator:     locator,
			Body:        body,
		}

		updater.On("UpdateCustomerSupportRemark",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, fmt.Errorf("error")).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&customersupportremarksop.UpdateCustomerSupportRemarkForMoveInternalServerError{}, response)
	})
}
