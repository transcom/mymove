package ghcapi

import (
	"errors"
	"fmt"
	"net/http"
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
	serviceItem := models.MTOServiceItem{ID: serviceItemID, MoveTaskOrderID: moveTaskOrderID, ReServiceID: reServiceID}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	req := httptest.NewRequest("POST", fmt.Sprintf("/move_task_orders/%s/mto_service_items", moveTaskOrderID.String()), nil)
	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoserviceitemop.CreateMTOServiceItemParams{
		HTTPRequest:     req,
		MoveTaskOrderID: serviceItem.MoveTaskOrderID.String(),
		ReServiceID:     serviceItem.ReServiceID.String(),
	}

	serviceItemCreator := &mocks.MTOServiceItemCreator{}

	suite.T().Run("Successful create", func(t *testing.T) {
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(&serviceItem, nil, nil).Once()

		handler := CreateServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemCreated{}, response)
	})

	suite.T().Run("Failed create", func(t *testing.T) {
		err := errors.New("cannot create service item")
		serviceItemCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, validate.NewErrors(), err).Once()

		handler := CreateServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			serviceItemCreator,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.CreateMTOServiceItemInternalServerError{}, response)
	})
}
