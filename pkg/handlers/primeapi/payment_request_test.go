package primeapi

import (
	"github.com/stretchr/testify/mock"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/mocks"
	"testing"
)

func (suite *ghcapi.HandlerSuite) TestCreatePaymentRequestHandler() {

	suite.T().Run("Not implemented", func(t *testing.T) {
		paymentRequestCreator := &mocks.PaymentRequestCreator{}

		paymentRequestCreator.On("CreatePaymentRequest",
			&paymentRequest,
			mock.Anything).Return(&paymentRequest, nil, nil).Once()

		handler := CreatePaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&paymentrequestop.CreateServiceItemCreated{}, response)
	})
}