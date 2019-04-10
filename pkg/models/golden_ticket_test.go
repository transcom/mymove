package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type GoldenTicketSuite struct {
	testingsuite.PopTestSuite
}

func (suite *GoldenTicketSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestGoldenTicketSuite(t *testing.T) {
	hs := &GoldenTicketSuite{PopTestSuite: testingsuite.NewPopTestSuite()}
	suite.Run(t, hs)
}

func (suite *GoldenTicketSuite) TestMakeGoldenTicket() {
	_, verrs, err := models.MakeGoldenTicket(suite.DB(), models.SelectedMoveTypeHHG)

	suite.Nil(err)
	suite.False(verrs.HasAny())
	gt := models.GoldenTicket{}
	err = suite.DB().First(&gt)
	suite.Nil(err)
}

func (suite *GoldenTicketSuite) TestValidGoldenTicket() {
	moveType := models.SelectedMoveTypeHHG
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{SelectedMoveType: &moveType},
	})
	gt := &models.GoldenTicket{}
	gt, verrs, err := models.MakeGoldenTicket(suite.DB(), moveType)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(gt)

	_, isValid := models.ValidateGoldenTicket(suite.DB(), gt.Code, move)
	suite.True(isValid)
}

func (suite *GoldenTicketSuite) TestGoldenTicketInvalidMoveType() {
	moveType := models.SelectedMoveTypeHHG
	gt := &models.GoldenTicket{}
	gt, verrs, err := models.MakeGoldenTicket(suite.DB(), moveType)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(gt)

	invalidMoveType := models.SelectedMoveTypePPM
	moveWithInvalidMoveType := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{SelectedMoveType: &invalidMoveType},
	})

	_, isValid := models.ValidateGoldenTicket(suite.DB(), gt.Code, moveWithInvalidMoveType)
	suite.False(isValid)
}

func (suite *GoldenTicketSuite) TestGoldenTicketInvalidCode() {
	moveType := models.SelectedMoveTypeHHG
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{SelectedMoveType: &moveType},
	})
	gt := &models.GoldenTicket{}
	gt, verrs, err := models.MakeGoldenTicket(suite.DB(), moveType)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(gt)

	_, isValid := models.ValidateGoldenTicket(suite.DB(), "INVALID CODE", move)
	suite.False(isValid)
}

func (suite *GoldenTicketSuite) TestUseGoldenTicket() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	gt, verrs, err := models.MakeGoldenTicket(suite.DB(), *move.SelectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	suite.NotNil(gt)

	gt, verrs, err = models.UseGoldenTicket(suite.DB(), gt.Code, move)
	suite.Nil(err)
	suite.False(verrs.HasAny())
	updatedGT := &models.GoldenTicket{}
	err = suite.DB().Find(updatedGT, gt.ID)
	suite.Nil(err)
	suite.Equal(move.ID, *updatedGT.MoveID)
}

func (suite *GoldenTicketSuite) TestMultipleTickets() {
	gtc := models.GoldenTicketCounts{models.SelectedMoveTypePPM: 66, models.SelectedMoveTypeHHG: 33}
	_, verrs, err := models.MakeGoldenTickets(suite.DB(), gtc)
	suite.Nil(err)
	suite.False(verrs.HasAny())

	hhgs, err := suite.DB().Where("move_type = ?", "HHG").Count(&models.GoldenTicket{})
	ppms, err := suite.DB().Where("move_type = ?", "PPM").Count(&models.GoldenTicket{})

	suite.Equal(33, hhgs)
	suite.Equal(66, ppms)
}
