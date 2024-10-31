package movetaskorder_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	m "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderCreatorIntegration() {
	factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeMS)
	factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeCS)
	builder := query.NewQueryBuilder()
	mtoCreator := m.NewMoveTaskOrderCreator(builder)

	order := factory.BuildOrder(suite.DB(), nil, nil)
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
