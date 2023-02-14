package movetaskorder_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderCreatorIntegration() {
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"),
			Code: models.ReServiceCodeMS,
		},
	})

	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			Code: models.ReServiceCodeCS,
		},
	})

	builder := query.NewQueryBuilder()
	mtoCreator := NewMoveTaskOrderCreator(builder)

	order := testdatagen.MakeDefaultOrder(suite.DB())
	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
	contractorID := contractor.ID
	newMto := models.Move{
		OrdersID:     order.ID,
		ContractorID: &contractorID,
		Status:       models.MoveStatusDRAFT,
		Locator:      models.GenerateLocator(),
	}
	actualMTO, verrs, err := mtoCreator.CreateMoveTaskOrder(suite.AppContextForTest(), &newMto)
	suite.NoError(err)
	suite.Empty(verrs)

	suite.NotZero(actualMTO.ID, "move task order is created")
}
