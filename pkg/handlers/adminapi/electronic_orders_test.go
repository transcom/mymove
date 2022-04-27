package adminapi

import (
	"errors"
	"net/http"

	"github.com/stretchr/testify/mock"

	electronicorderop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_order"

	"net/http/httptest"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetElectronicOrdersTotalsHandler() {
	setupRequest := func() *http.Request {
		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		req := httptest.NewRequest("GET", "/electronic_orders/totals", nil)
		return suite.AuthenticateAdminRequest(req, requestUser)
	}

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	filter := "Issuer.eq:marines"

	suite.Run("successful response", func() {
		params := electronicorderop.GetElectronicOrdersTotalsParams{
			HTTPRequest: setupRequest(),
			Filter:      []string{filter},
		}

		electronicOrderCategoryCountFetcher := &mocks.ElectronicOrderCategoryCountFetcher{}
		electronicOrderCategoryCountFetcher.On("FetchElectronicOrderCategoricalCounts",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(map[interface{}]int{models.IssuerArmy: 2}, nil)
		handler := GetElectronicOrdersTotalsHandler{
			HandlerContext:                      handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ElectronicOrderCategoryCountFetcher: electronicOrderCategoryCountFetcher,
			NewQueryFilter:                      newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&electronicorderop.GetElectronicOrdersTotalsOK{}, response)
	})

	suite.Run("error response", func() {
		params := electronicorderop.GetElectronicOrdersTotalsParams{
			HTTPRequest: setupRequest(),
			Filter:      []string{filter},
		}

		err := errors.New("An error happened")

		electronicOrderCategoryCountFetcher := &mocks.ElectronicOrderCategoryCountFetcher{}
		electronicOrderCategoryCountFetcher.On("FetchElectronicOrderCategoricalCounts",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, err)
		handler := GetElectronicOrdersTotalsHandler{
			HandlerContext:                      handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ElectronicOrderCategoryCountFetcher: electronicOrderCategoryCountFetcher,
			NewQueryFilter:                      newQueryFilter,
		}

		handler.Handle(params)
		suite.Error(err, "An error happened")
	})
}
