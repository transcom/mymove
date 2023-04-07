package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestCreateMoveWithPPMShow() {
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: Move{
				Show: BoolPointer(true),
			},
		},
	}, nil)
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
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: Move{
				Show: BoolPointer(false),
			},
		},
	}, nil)

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
	orders := factory.BuildOrder(suite.DB(), nil, nil)
	factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	moveOptions := MoveOptions{
		Show: BoolPointer(true),
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
