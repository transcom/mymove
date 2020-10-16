package adminapi

import (
	"errors"

	"github.com/stretchr/testify/mock"

	electronicorderop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_order"

	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetElectronicOrdersTotalsHandler() {
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/electronic_orders/totals", nil)
	req = suite.AuthenticateAdminRequest(req, requestUser)

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	filter := "Issuer.eq:marines"

	suite.T().Run("successful response", func(t *testing.T) {
		params := electronicorderop.GetElectronicOrdersTotalsParams{
			HTTPRequest: req,
			Filter:      []string{filter},
		}

		electronicOrderCategoryCountFetcher := &mocks.ElectronicOrderCategoryCountFetcher{}
		electronicOrderCategoryCountFetcher.On("FetchElectronicOrderCategoricalCounts",
			mock.Anything,
			mock.Anything,
		).Return(map[interface{}]int{models.IssuerArmy: 2}, nil)
		handler := GetElectronicOrdersTotalsHandler{
			HandlerContext:                      handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			ElectronicOrderCategoryCountFetcher: electronicOrderCategoryCountFetcher,
			NewQueryFilter:                      newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&electronicorderop.GetElectronicOrdersTotalsOK{}, response)
	})

	suite.T().Run("error response", func(t *testing.T) {
		params := electronicorderop.GetElectronicOrdersTotalsParams{
			HTTPRequest: req,
			Filter:      []string{filter},
		}

		err := errors.New("An error happened")

		electronicOrderCategoryCountFetcher := &mocks.ElectronicOrderCategoryCountFetcher{}
		electronicOrderCategoryCountFetcher.On("FetchElectronicOrderCategoricalCounts",
			mock.Anything,
			mock.Anything,
		).Return(nil, err)
		handler := GetElectronicOrdersTotalsHandler{
			HandlerContext:                      handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			ElectronicOrderCategoryCountFetcher: electronicOrderCategoryCountFetcher,
			NewQueryFilter:                      newQueryFilter,
		}

		handler.Handle(params)
		suite.Error(err, "An error happened")
	})

}
