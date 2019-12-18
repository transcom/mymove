package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMTOServiceItemHandler() {
	moveTaskOrderID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	reServiceID, _ := uuid.NewV4()
	mtoShipmentID, _ := uuid.NewV4()
	metaID, _ := uuid.NewV4()
	serviceItem := models.MTOServiceItem{
		ID: serviceItemID, MoveTaskOrderID: moveTaskOrderID, ReServiceID: reServiceID, MTOShipmentID: mtoShipmentID, MetaID: metaID, MetaType: "unknown",
	}

	req := httptest.NewRequest("POST", fmt.Sprintf("/move_task_orders/%s/mto_service_items", moveTaskOrderID.String()), nil)
	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoserviceitemop.CreateMTOServiceItemParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(serviceItem.MoveTaskOrderID),
		CreateMTOServiceItemBody: mtoserviceitemop.CreateMTOServiceItemBody{
			ReServiceID:   handlers.FmtUUID(serviceItem.ReServiceID),
			MtoShipmentID: handlers.FmtUUID(serviceItem.MTOShipmentID),
			MetaID:        handlers.FmtUUID(serviceItem.MetaID),
			MetaType:      handlers.FmtString(serviceItem.MetaType),
		},
	}

	serviceItemCreator := &mocks.MTOServiceItemCreator{}

	suite.T().Run("Successful create", func(t *testing.T) {
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(&serviceItem, nil, nil).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemCreated{}, response)
	})

	suite.T().Run("Failed create: InternalServiceError", func(t *testing.T) {
		err := errors.New("cannot create service item")
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, validate.NewErrors(), err).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemInternalServerError{}, response)
	})

	suite.T().Run("Failed create: UnprocessableEntity", func(t *testing.T) {
		verrs := validate.NewErrors()
		verrs.Add("error", "error test")
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, verrs, nil).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.T().Run("Failed create: UnprocessableEntity - UUID parsing error", func(t *testing.T) {
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		newParams := mtoserviceitemop.CreateMTOServiceItemParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItem.MoveTaskOrderID),
			CreateMTOServiceItemBody: mtoserviceitemop.CreateMTOServiceItemBody{
				ReServiceID:   handlers.FmtUUID(serviceItem.ReServiceID),
				MtoShipmentID: handlers.FmtUUID(serviceItem.MTOShipmentID),
				MetaID:        handlers.FmtUUID(serviceItem.MetaID),
				MetaType:      handlers.FmtString(serviceItem.MetaType),
			},
		}
		newParams.MoveTaskOrderID = "blah"

		response := handler.Handle(newParams)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.T().Run("Failed create: UnprocessableEntity - Violates foreign key constraints", func(t *testing.T) {
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, validate.NewErrors(), errors.New(models.ViolatesForeignKeyConstraint)).Once()

		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemNotFound{}, response)
	})
}
