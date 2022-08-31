package adminapi

import (
	"errors"

	"github.com/stretchr/testify/mock"

	electronicorderop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_order"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *HandlerSuite) TestGetElectronicOrdersTotalsHandler() {
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	filter := "Issuer.eq:marines"

	suite.Run("successful response", func() {
		params := electronicorderop.GetElectronicOrdersTotalsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/electronic_orders/totals"),
			Filter:      []string{filter},
		}

		electronicOrderCategoryCountFetcher := &mocks.ElectronicOrderCategoryCountFetcher{}
		electronicOrderCategoryCountFetcher.On("FetchElectronicOrderCategoricalCounts",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(map[interface{}]int{models.IssuerArmy: 2}, nil)
		handler := GetElectronicOrdersTotalsHandler{
			HandlerConfig:                       suite.HandlerConfig(),
			ElectronicOrderCategoryCountFetcher: electronicOrderCategoryCountFetcher,
			NewQueryFilter:                      newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&electronicorderop.GetElectronicOrdersTotalsOK{}, response)
	})

	suite.Run("error response", func() {
		params := electronicorderop.GetElectronicOrdersTotalsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/electronic_orders/totals"),
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
			HandlerConfig:                       suite.HandlerConfig(),
			ElectronicOrderCategoryCountFetcher: electronicOrderCategoryCountFetcher,
			NewQueryFilter:                      newQueryFilter,
		}

		handler.Handle(params)
		suite.Error(err, "An error happened")
	})
}
