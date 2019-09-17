package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/transcom/mymove/pkg/services"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/gofrs/uuid"

	accesscodeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/accesscode"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexAccessCodesHandler() {
	// test that everything is wired up correctly
	m := testdatagen.MakeDefaultMove(suite.DB())

	id, _ := uuid.NewV4()
	sm := m.Orders.ServiceMember
	ac := models.AccessCode{
		ID:              id,
		ServiceMemberID: &sm.ID,
		ServiceMember:   sm,
		Code:            "ABCXYZ",
		MoveType:        m.SelectedMoveType,
	}
	assertions := testdatagen.Assertions{
		AccessCode: ac,
	}
	testdatagen.MakeAccessCode(suite.DB(), assertions)
	req := httptest.NewRequest("GET", "/access_codes", nil)
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := accesscodeop.IndexAccessCodesParams{
			HTTPRequest: req,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexAccessCodesHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:        query.NewQueryFilter,
			AccessCodeListFetcher: accesscode.NewAccessCodeListFetcher(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&accesscodeop.IndexAccessCodesOK{}, response)
		okResponse := response.(*accesscodeop.IndexAccessCodesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(ac.ID.String(), okResponse.Payload[0].ID.String())
		suite.Equal(ac.Code, okResponse.Payload[0].Code)
		suite.Equal(string(*ac.MoveType), okResponse.Payload[0].MoveType)
		suite.Equal(m.Locator, okResponse.Payload[0].Locator)
	})
	suite.T().Run("test failed response", func(t *testing.T) {
		params := accesscodeop.IndexAccessCodesParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		accessCodeListFetcher := &mocks.AccessCodeListFetcher{}
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
		accessCodeListFetcher.On("FetchAccessCodeList",
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexAccessCodesHandler{
			HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:        newQueryFilter,
			AccessCodeListFetcher: accessCodeListFetcher,
		}
		response := handler.Handle(params)

		suite.CheckErrorResponse(response, http.StatusNotFound, expectedError.Error())
	})
}

func (suite *HandlerSuite) TestIndexAccessCodesHandlerHelpers() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := IndexAccessCodesHandler{
		HandlerContext:        handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		NewQueryFilter:        query.NewQueryFilter,
		AccessCodeListFetcher: accesscode.NewAccessCodeListFetcher(queryBuilder),
	}

	suite.T().Run("test both filters present", func(t *testing.T) {
		qfs := handler.generateQueryFilters([]string{`{"move_type": "PPM"}`, `{"code": "XYZBCS"}`}, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("move_type", "=", "PPM"),
			query.NewQueryFilter("code", "=", "XYZBCS"),
		}
		suite.Equal(expectedFilters, qfs)
	})
	suite.T().Run("test only move_type present", func(t *testing.T) {
		qfs := handler.generateQueryFilters([]string{`{"move_type": "PPM"}`}, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("move_type", "=", "PPM"),
		}
		suite.Equal(expectedFilters, qfs)
	})
	suite.T().Run("test only code present", func(t *testing.T) {
		qfs := handler.generateQueryFilters([]string{`{"code": "XYZBCS"}`}, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("code", "=", "XYZBCS"),
		}
		suite.Equal(expectedFilters, qfs)
	})
}
