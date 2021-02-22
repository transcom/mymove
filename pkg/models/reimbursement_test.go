package models_test

import (
	"time"

	"github.com/pkg/errors"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReimbursementStateMachine() {
	reimbursement := BuildDraftReimbursement(1200, MethodOfReceiptOTHERDD)

	err := reimbursement.Request()
	suite.NoError(err)
	suite.Equal(ReimbursementStatusREQUESTED, reimbursement.Status, "expected Requested")

	err = reimbursement.Approve()
	suite.NoError(err)
	suite.Equal(ReimbursementStatusAPPROVED, reimbursement.Status, "expected Approved")

	err = reimbursement.Pay()
	suite.NoError(err)
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

	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
	//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
	//RA: in which this would be considered a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	reimbursement.Request() // nolint:errcheck

	verrs, err := suite.DB().ValidateAndCreate(&reimbursement)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	suite.NotNil(reimbursement.ID)

	since := time.Since(*reimbursement.RequestedDate)
	if since > 1*time.Second {
		suite.T().Fail()
	}

}
