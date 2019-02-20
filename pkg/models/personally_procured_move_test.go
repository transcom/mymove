package models_test

import (
	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestPPMValidation() {
	ppm := &PersonallyProcuredMove{}

	expErrors := map[string][]string{
		"status": {"Status can not be blank."},
	}

	suite.verifyValidationErrors(ppm, expErrors)
}

func (suite *ModelSuite) TestPPMAdvance() {

	move := testdatagen.MakeDefaultMove(suite.DB())
	serviceMember := move.Orders.ServiceMember

	advance := BuildDraftReimbursement(1000, MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.Nil(err)
	suite.False(verrs.HasAny())

	advance.Request()
	SavePersonallyProcuredMove(suite.DB(), ppm)
	session := auth.Session{
		UserID:          serviceMember.User.ID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: serviceMember.ID,
	}
	fetchedPPM, err := FetchPersonallyProcuredMove(suite.DB(), &session, ppm.ID)
	suite.Nil(err)
	suite.Equal(fetchedPPM.Advance.Status, ReimbursementStatusREQUESTED, "expected Requested")
}

func (suite *ModelSuite) TestPPMAdvanceNoGTCC() {
	move := testdatagen.MakeDefaultMove(suite.DB())

	advance := BuildDraftReimbursement(1000, MethodOfReceiptGTCC)

	_, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.Nil(err)
	suite.True(verrs.HasAny())
}

func (suite *ModelSuite) TestPPMStateMachine() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())
	orders.Status = OrderStatusSUBMITTED // NEVER do this outside of a test.
	suite.MustSave(&orders)

	selectedMoveType := SelectedMoveTypeHHGPPM

	move, verrs, err := orders.CreateNewMove(suite.DB(), &selectedMoveType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders

	advance := BuildDraftReimbursement(1000, MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, true, &advance)
	suite.Nil(err)
	suite.False(verrs.HasAny())

	ppm.Status = PPMStatusSUBMITTED // NEVER do this outside of a test.

	// Can cancel ppm
	err = ppm.Cancel()
	suite.Nil(err)
	suite.Equal(PPMStatusCANCELED, ppm.Status, "expected Canceled")
}
