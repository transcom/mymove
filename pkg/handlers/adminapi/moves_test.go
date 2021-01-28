package adminapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/services/pagination"

	"github.com/transcom/mymove/pkg/services"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	moveop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexMovesHandler() {
	// test that everything is wired up correctly
	m := testdatagen.MakeDefaultMove(suite.DB())
	req := httptest.NewRequest("GET", "/moves", nil)

	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := moveop.IndexMovesParams{
			HTTPRequest: req,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexMovesHandler{
			HandlerContext:  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:  query.NewQueryFilter,
			MoveListFetcher: move.NewMoveListFetcher(queryBuilder),
			NewPagination:   pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&moveop.IndexMovesOK{}, response)
		okResponse := response.(*moveop.IndexMovesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(m.ID.String(), okResponse.Payload[0].ID.String())
	})
	suite.T().Run("test failed response", func(t *testing.T) {
		params := moveop.IndexMovesParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		moveListFetcher := &mocks.MoveListFetcher{}
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
		moveListFetcher.On("FetchMoveList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexMovesHandler{
			HandlerContext:  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:  newQueryFilter,
			MoveListFetcher: moveListFetcher,
			NewPagination:   pagination.NewPagination,
		}
		response := handler.Handle(params)

		suite.CheckErrorResponse(response, http.StatusNotFound, expectedError.Error())
	})
}

func (suite *HandlerSuite) TestIndexMovesHandlerHelpers() {
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := IndexMovesHandler{
		HandlerContext:  handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		NewQueryFilter:  query.NewQueryFilter,
		MoveListFetcher: move.NewMoveListFetcher(queryBuilder),
		NewPagination:   pagination.NewPagination,
	}

	suite.T().Run("test filters present", func(t *testing.T) {

		s := `{"locator":"TEST123"}`
		qfs := handler.generateQueryFilters(&s, suite.TestLogger())
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("locator", "=", "TEST123"),
		}
		suite.Equal(expectedFilters, qfs)
	})
}

func (suite *HandlerSuite) TestUpdateMoveHandler() {
	defaultMove := testdatagen.MakeDefaultMove(suite.DB())

	// Create handler and request:
	builder := query.NewQueryBuilder(suite.DB())
	handler := UpdateMoveHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), builder, mtoserviceitem.NewMTOServiceItemCreator(builder)),
	}
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/moves/%s", defaultMove.ID), nil)

	// Case: Move is successfully updated
	suite.T().Run("200 - OK response", func(t *testing.T) {
		params := moveop.UpdateMoveParams{
			HTTPRequest: req,
			MoveID:      *handlers.FmtUUID(defaultMove.ID),
			Move: moveop.UpdateMoveBody{
				Show: true,
			},
		}
		// Run swagger validations
		suite.NoError(params.Move.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&moveop.UpdateMoveOK{}, response)

		// Check values
		moveOK := response.(*moveop.UpdateMoveOK)
		suite.Equal(moveOK.Payload.ID.String(), defaultMove.ID.String())
		suite.Equal(*moveOK.Payload.Show, params.Move.Show)
	})

	// Case: Move is not found
	suite.T().Run("404 - Move not found", func(t *testing.T) {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		params := moveop.UpdateMoveParams{
			HTTPRequest: req,
			MoveID:      *handlers.FmtUUID(badUUID),
			Move: moveop.UpdateMoveBody{
				Show: true,
			},
		}
		// Run swagger validations
		suite.NoError(params.Move.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&moveop.UpdateMoveNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestGetMoveHandler() {
	// test that everything is wired up correctly
	defaultMove := testdatagen.MakeDefaultMove(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s", defaultMove.ID), nil)

	suite.T().Run("200 - OK response", func(t *testing.T) {
		params := moveop.GetMoveParams{
			HTTPRequest: req,
			MoveID:      *handlers.FmtUUID(defaultMove.ID),
		}
		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		}

		response := handler.Handle(params)

		suite.IsType(&moveop.GetMoveOK{}, response)
		okResponse := response.(*moveop.GetMoveOK)
		suite.Equal(defaultMove.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("404 - Move not found", func(t *testing.T) {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		badReq := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s", badUUID), nil)
		params := moveop.GetMoveParams{
			HTTPRequest: badReq,
			MoveID:      *handlers.FmtUUID(badUUID),
		}

		handler := GetMoveHandler{
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		}

		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveInternalServerError{}, response)
	})
}
