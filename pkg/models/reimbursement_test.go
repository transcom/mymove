package models_test

import (
	"time"

	"github.com/pkg/errors"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReimbursementStateMachine() {
	reimbursement := BuildDraftReimbursement(1200, MethodOfReceiptOTHERDD)

	err := reimbursement.Request()
	suite.Nil(err)
	suite.Equal(reimbursement.Status, ReimbursementStatusREQUESTED, "expected Requested")

	err = reimbursement.Approve()
	suite.Nil(err)
	suite.Equal(reimbursement.Status, ReimbursementStatusAPPROVED, "expected Approved")

	err = reimbursement.Pay()
	suite.Nil(err)
	suite.Equal(reimbursement.Status, ReimbursementStatusPAID, "expected Paid")

	err = reimbursement.Reject()
	suite.Equal(ErrInvalidTransition, errors.Cause(err))

	reimbursement.Status = ReimbursementStatusDRAFT // NEVER do this outside of a test.

	err = reimbursement.Pay()
	suite.Equal(ErrInvalidTransition, errors.Cause(err))

	err = reimbursement.Approve()
	suite.Equal(ErrInvalidTransition, errors.Cause(err))

	err = reimbursement.Reject()
	suite.Equal(ErrInvalidTransition, errors.Cause(err))

}

func (suite *ModelSuite) TestBasicReimbursement() {
	reimbursement := BuildDraftReimbursement(1200, MethodOfReceiptOTHERDD)

	reimbursement.Request()

	verrs, err := suite.db.ValidateAndCreate(&reimbursement)
	suite.Nil(err)
	suite.False(verrs.HasAny())

	suite.NotNil(reimbursement.ID)

	since := time.Now().Sub(*reimbursement.RequestedDate)
	if since > 1*time.Second {
		suite.T().Fail()
	}

}
