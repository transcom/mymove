package models_test

import (
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestCreateNewMoveShow() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())

	selectedMoveType := SelectedMoveTypeHHG

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	_, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Len(moves, 1)
}

func (suite *ModelSuite) TestCreateNewMoveShowFalse() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())

	selectedMoveType := SelectedMoveTypeHHG

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(false),
	}
	_, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Empty(moves)
}

func (suite *ModelSuite) TestShowPPMQueue() {
	// PPMs should only show statuses in the queue:
	// approved, payment requested and completed

	// Make PPMs with different statuses
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusAPPROVED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusPAYMENTREQUESTED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusCOMPLETED,
		},
	})

	// Expected 3 moves for PPM queue returned
	moves, err := GetMoveQueueItems(suite.DB(), "ppm")
	suite.NoError(err)
	suite.Len(moves, 3)
}

func (suite *ModelSuite) TestShowPPMQueueStatusDraftSubmittedCanceled() {
	// PPMs should only show statuses in the queue:
	// approved, payment requested and completed

	// PPMs not in approved, payment requested or completed states are not returned
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusDRAFT,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusSUBMITTED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Status: models.PPMStatusCANCELED,
		},
	})

	// Expected 0 moves for PPM queue returned
	moves, err := GetMoveQueueItems(suite.DB(), "ppm")
	suite.NoError(err)
	suite.Len(moves, 0)
}

func (suite *ModelSuite) TestQueueNotFound() {
	moves, moveErrs := GetMoveQueueItems(suite.DB(), "queue_not_found")
	suite.Equal(ErrFetchNotFound, moveErrs, "Expected not to find move queue items")
	suite.Empty(moves)
}

func (suite *ModelSuite) TestActivePPMQueue() {
	testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			Status: models.ShipmentStatusINTRANSIT,
		},
	})

	testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAPPROVED,
		},
	})

	now := time.Now()
	testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			Status:                models.ShipmentStatusACCEPTED,
			PmSurveyConductedDate: &now,
		},
	})

	testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			Status: models.ShipmentStatusACCEPTED,
		},
	})

	moves, err := GetMoveQueueItems(suite.DB(), "hhg_active")
	suite.NoError(err)
	suite.Len(moves, 3)
}
