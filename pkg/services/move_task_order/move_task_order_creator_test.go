package movetaskorder_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderCreatorIntegration() {
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"),
			Code: "MS",
		},
	})

	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			Code: "CS",
		},
	})

	builder := query.NewQueryBuilder(suite.DB())
	mtoCreator := NewMoveTaskOrderCreator(builder, suite.DB())

	order := testdatagen.MakeDefaultOrder(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())
	contractorID := contractor.ID
	newMto := models.Move{
		OrdersID:     order.ID,
		ContractorID: &contractorID,
		Status:       models.MoveStatusDRAFT,
		Locator:      models.GenerateLocator(),
	}
	actualMTO, verrs, err := mtoCreator.CreateMoveTaskOrder(&newMto)
	suite.NoError(err)
	suite.Empty(verrs)

	suite.NotZero(actualMTO.ID, "move task order is created")
}
