package adminapi

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	moveop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/moves"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestIndexMovesHandler() {
	suite.Run("integration test ok response", func() {
		m := factory.BuildMove(suite.DB(), nil, nil)
		params := moveop.IndexMovesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/moves"),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := IndexMovesHandler{
			HandlerConfig:   suite.HandlerConfig(),
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
	suite.Run("test failed response", func() {
		params := moveop.IndexMovesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/moves"),
		}
		expectedError := models.ErrFetchNotFound
		moveListFetcher := &mocks.MoveListFetcher{}
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
		moveListFetcher.On("FetchMoveList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		handler := IndexMovesHandler{
			HandlerConfig:   suite.HandlerConfig(),
			NewQueryFilter:  newQueryFilter,
			MoveListFetcher: moveListFetcher,
			NewPagination:   pagination.NewPagination,
		}
		response := handler.Handle(params)

		suite.CheckErrorResponse(response, http.StatusNotFound, expectedError.Error())
	})
}

func (suite *HandlerSuite) TestIndexMovesHandlerHelpers() {
	suite.Run("test filters present", func() {
		s := `{"locator":"TEST123"}`
		qfs := generateQueryFilters(suite.Logger(), &s, locatorFilterConverters)
		expectedFilters := []services.QueryFilter{
			query.NewQueryFilter("locator", "=", "TEST123"),
		}
		suite.Equal(expectedFilters, qfs)
	})
}

func (suite *HandlerSuite) TestUpdateMoveHandler() {
	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}
	setupHandler := func() UpdateMoveHandler {
		builder := query.NewQueryBuilder()
		moveRouter := move.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		return UpdateMoveHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderUpdater(
				builder,
				mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
				moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil),
			),
		}
	}
	show := true

	// Case: Move is successfully updated
	suite.Run("200 - OK response", func() {
		defaultMove := factory.BuildMove(suite.DB(), nil, nil)
		params := moveop.UpdateMoveParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/moves/%s", defaultMove.ID)),
			MoveID:      *handlers.FmtUUID(defaultMove.ID),
			Move: &adminmessages.MoveUpdate{
				Show: &show,
			},
		}
		// Run swagger validations
		suite.NoError(params.Move.Validate(strfmt.Default))

		// Run handler and check response
		response := setupHandler().Handle(params)
		suite.IsType(&moveop.UpdateMoveOK{}, response)

		// Check values
		moveOK := response.(*moveop.UpdateMoveOK)
		suite.Equal(defaultMove.ID.String(), moveOK.Payload.ID.String())
		suite.Equal(*params.Move.Show, *moveOK.Payload.Show)
	})

	// Case: Move is not found
	suite.Run("404 - Move not found", func() {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		params := moveop.UpdateMoveParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/moves/%s", badUUID)),
			MoveID:      *handlers.FmtUUID(badUUID),
			Move: &adminmessages.MoveUpdate{
				Show: &show,
			},
		}
		// Run swagger validations
		suite.NoError(params.Move.Validate(strfmt.Default))

		// Run handler and check response
		response := setupHandler().Handle(params)
		suite.IsType(&moveop.UpdateMoveNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestGetMoveHandler() {
	suite.Run("200 - OK response", func() {
		defaultMove := factory.BuildMove(suite.DB(), nil, nil)
		params := moveop.GetMoveParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/moves/%s", defaultMove.ID)),
			MoveID:      *handlers.FmtUUID(defaultMove.ID),
		}
		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
		}

		response := handler.Handle(params)

		suite.IsType(&moveop.GetMoveOK{}, response)
		okResponse := response.(*moveop.GetMoveOK)
		suite.Equal(defaultMove.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("500 - Internal Server Error for No SQL Rows Returned", func() {
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		params := moveop.GetMoveParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/moves/%s", badUUID)),
			MoveID:      *handlers.FmtUUID(badUUID),
		}

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
		}

		response := handler.Handle(params)
		suite.IsType(&moveop.GetMoveInternalServerError{}, response)
	})
}
