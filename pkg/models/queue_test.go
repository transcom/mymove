package models_test

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestCreateMoveWithPPMShow() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: orders,
		Move: Move{
			ID:   uuid.FromStringOrNil("7024c8c5-52ca-4639-bf69-dd8238308c98"),
			Show: swag.Bool(true),
		},
	})

	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		ServiceMember: move.Orders.ServiceMember,
		PersonallyProcuredMove: PersonallyProcuredMove{
			Move:   move,
			MoveID: move.ID,
		},
	})

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Len(moves, 1)
}

func (suite *ModelSuite) TestCreateMoveWithPPMNoShow() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: orders,
		Move: Move{
			ID:   uuid.FromStringOrNil("7024c8c5-52ca-4639-bf69-dd8238308c98"),
			Show: swag.Bool(false),
		},
	})

	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		ServiceMember: move.Orders.ServiceMember,
		PersonallyProcuredMove: PersonallyProcuredMove{
			Move:   move,
			MoveID: move.ID,
		},
	})

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Empty(moves)

}

func (suite *ModelSuite) TestCreateNewMoveWithNoPPMShow() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())

	moveOptions := MoveOptions{
		Show: swag.Bool(true),
	}
	_, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Empty(moves)
}

func (suite *ModelSuite) TestShowPPMQueue() {
	all := map[string]bool{
		string(PPMStatusAPPROVED):         true,
		string(PPMStatusPAYMENTREQUESTED): true,
		string(PPMStatusCOMPLETED):        true,
		string(PPMStatusSUBMITTED):        true,
		string(PPMStatusDRAFT):            true,
	}

	new := map[string]bool{
		string(PPMStatusSUBMITTED): true,
		string(PPMStatusDRAFT):     true,
	}

	tests := []struct {
		input      string
		movesCount int
		want       map[string]bool
	}{
		{input: "new", movesCount: 2, want: new},
		{input: "ppm_payment_requested", movesCount: 1, want: map[string]bool{string(PPMStatusPAYMENTREQUESTED): true}},
		{input: "ppm_completed", movesCount: 1, want: map[string]bool{string(PPMStatusCOMPLETED): true}},
		{input: "ppm_approved", movesCount: 1, want: map[string]bool{string(PPMStatusAPPROVED): true}},
		{input: "all", movesCount: 5, want: all},
	}

	// Make PPMs with different statuses
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: PersonallyProcuredMove{
			Status: PPMStatusAPPROVED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: PersonallyProcuredMove{
			Status: PPMStatusPAYMENTREQUESTED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: PersonallyProcuredMove{
			Status: PPMStatusCOMPLETED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		Move: Move{
			Status: MoveStatusSUBMITTED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		Move: Move{
			Status: MoveStatusAPPROVED,
		},
		PersonallyProcuredMove: PersonallyProcuredMove{
			Status: PPMStatusSUBMITTED,
		},
	})

	for _, tc := range tests {
		moves, err := GetMoveQueueItems(suite.DB(), tc.input)

		suite.NoError(err)
		suite.Len(moves, tc.movesCount)
		for _, move := range moves {
			suite.True(tc.want[*move.PpmStatus])
		}
	}

}

func (suite *ModelSuite) TestQueueNotFound() {
	moves, moveErrs := GetMoveQueueItems(suite.DB(), "queue_not_found")
	suite.Equal(ErrFetchNotFound, moveErrs, "Expected not to find move queue items")
	suite.Empty(moves)
}
