package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildProofOfServiceDoc() {
	suite.Run("Successful creation of default ProofOfServiceDoc", func() {
		// Under test:      BuildProofOfServiceDoc
		// Set up:          Create a ProofOfServiceDoc with no customizations or traits
		// Expected outcome: ProofOfServiceDoc should be created with default values

		// SETUP
		proofOfServiceDoc := BuildProofOfServiceDoc(suite.DB(), nil, nil)

		suite.False(proofOfServiceDoc.ID.IsNil())
		suite.False(proofOfServiceDoc.PaymentRequest.MoveTaskOrderID.IsNil())
		suite.NotNil(proofOfServiceDoc.PaymentRequest.MoveTaskOrder)
		suite.False(proofOfServiceDoc.PaymentRequest.MoveTaskOrder.ID.IsNil())
	})

	suite.Run("Successful creation of custom ProofOfServiceDoc", func() {
		// Under test:      BuildProofOfServiceDoc
		// Set up:          Create a ProofOfServiceDoc and pass custom fields
		// Expected outcome: ProofOfServiceDoc should be created with custom values

		// SETUP
		customMove := models.Move{
			Locator: "AAAA",
		}
		customPaymentrequest := models.PaymentRequest{
			Status: models.PaymentRequestStatusPaid,
		}

		// CALL FUNCTION UNDER TEST
		proofOfServiceDoc := BuildProofOfServiceDoc(suite.DB(), []Customization{
			{Model: customMove},
			{Model: customPaymentrequest},
		}, nil)

		suite.False(proofOfServiceDoc.ID.IsNil())
		suite.False(proofOfServiceDoc.PaymentRequest.MoveTaskOrderID.IsNil())
		suite.NotNil(proofOfServiceDoc.PaymentRequest.MoveTaskOrder)
		suite.False(proofOfServiceDoc.PaymentRequest.MoveTaskOrder.ID.IsNil())
		suite.Equal(customMove.Locator, proofOfServiceDoc.PaymentRequest.MoveTaskOrder.Locator)
		suite.False(proofOfServiceDoc.PaymentRequestID.IsNil())
		suite.NotNil(proofOfServiceDoc.PaymentRequest)
		suite.False(proofOfServiceDoc.PaymentRequest.ID.IsNil())
		suite.Equal(customPaymentrequest.Status, proofOfServiceDoc.PaymentRequest.Status)
	})

	suite.Run("Successful creation of stubbed ProofOfServiceDoc", func() {
		// Under test:      BuildProofOfServiceDoc
		// Set up:          Create a stubbed ProofOfServiceDoc
		// Expected outcome:No new ProofOfServiceDoc should be created
		precount, err := suite.DB().Count(&models.ProofOfServiceDoc{})
		suite.NoError(err)

		proofOfServiceDoc := BuildProofOfServiceDoc(nil, nil, nil)

		// VALIDATE RESULTS
		suite.True(proofOfServiceDoc.ID.IsNil())
		suite.True(proofOfServiceDoc.PaymentRequestID.IsNil())
		suite.NotNil(proofOfServiceDoc.PaymentRequest)
		suite.True(proofOfServiceDoc.PaymentRequest.ID.IsNil())

		// Count how many notification are in the DB, no new
		// notifications should have been created
		count, err := suite.DB().Count(&models.ProofOfServiceDoc{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful return of linkOnly ProofOfServiceDoc", func() {
		// Under test:       BuildProofOfServiceDoc
		// Set up:           Pass in a linkOnly ProofOfServiceDoc
		// Expected outcome: No new ProofOfServiceDoc should be created

		// Check num ProofOfServiceDoc records
		precount, err := suite.DB().Count(&models.ProofOfServiceDoc{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		proofOfServiceDoc := BuildProofOfServiceDoc(suite.DB(), []Customization{
			{
				Model: models.ProofOfServiceDoc{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.ProofOfServiceDoc{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, proofOfServiceDoc.ID)
	})
}
