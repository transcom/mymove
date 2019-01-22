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
	suite.Equal(ReimbursementStatusREQUESTED, reimbursement.Status, "expected Requested")

	err = reimbursement.Approve()
	suite.Nil(err)
	suite.Equal(ReimbursementStatusAPPROVED, reimbursement.Status, "expected Approved")

	err = reimbursement.Pay()
	suite.Nil(err)
	suite.Equal(ReimbursementStatusPAID, reimbursement.Status, "expected Paid")

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

	verrs, err := suite.DB().ValidateAndCreate(&reimbursement)
	suite.Nil(err)
	suite.False(verrs.HasAny())

	suite.NotNil(reimbursement.ID)

	since := time.Now().Sub(*reimbursement.RequestedDate)
	if since > 1*time.Second {
		suite.T().Fail()
	}

}
