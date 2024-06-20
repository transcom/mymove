package models_test

import (
	"time"

	"github.com/pkg/errors"

	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReimbursementStateMachine() {
	reimbursement := m.BuildDraftReimbursement(1200, m.MethodOfReceiptOTHERDD)

	err := reimbursement.Request()
	suite.NoError(err)
	suite.Equal(m.ReimbursementStatusREQUESTED, reimbursement.Status, "expected Requested")

	err = reimbursement.Approve()
	suite.NoError(err)
	suite.Equal(m.ReimbursementStatusAPPROVED, reimbursement.Status, "expected Approved")

	err = reimbursement.Pay()
	suite.NoError(err)
	suite.Equal(m.ReimbursementStatusPAID, reimbursement.Status, "expected Paid")

	err = reimbursement.Reject()
	suite.Equal(m.ErrInvalidTransition, errors.Cause(err))

	reimbursement.Status = m.ReimbursementStatusDRAFT // NEVER do this outside of a test.

	err = reimbursement.Pay()
	suite.Equal(m.ErrInvalidTransition, errors.Cause(err))

	err = reimbursement.Approve()
	suite.Equal(m.ErrInvalidTransition, errors.Cause(err))

	err = reimbursement.Reject()
	suite.Equal(m.ErrInvalidTransition, errors.Cause(err))

}

func (suite *ModelSuite) TestBasicReimbursement() {
	reimbursement := m.BuildDraftReimbursement(1200, m.MethodOfReceiptOTHERDD)

	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
	//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
	//RA: in which this would be considered a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	reimbursement.Request()

	verrs, err := reimbursement.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.NotNil(reimbursement.ID)

	since := time.Since(*reimbursement.RequestedDate)
	if since > 1*time.Second {
		suite.T().Fail()
	}

}
