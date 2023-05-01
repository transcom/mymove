package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildProofOfServiceDoc creates ProofOfServiceDoc.
// Also creates, if not provided
// - PaymentRequest and associated set relationships
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildProofOfServiceDoc(db *pop.Connection, customs []Customization, traits []Trait) models.ProofOfServiceDoc {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cProofOfServiceDoc models.ProofOfServiceDoc
	if result := findValidCustomization(customs, ProofOfServiceDoc); result != nil {
		cProofOfServiceDoc = result.Model.(models.ProofOfServiceDoc)

		if result.LinkOnly {
			return cProofOfServiceDoc
		}
	}

	paymentRequest := BuildPaymentRequest(db, customs, traits)

	ProofOfServiceDoc := models.ProofOfServiceDoc{
		PaymentRequest:   paymentRequest,
		PaymentRequestID: paymentRequest.ID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&ProofOfServiceDoc, cProofOfServiceDoc)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &ProofOfServiceDoc)
	}

	return ProofOfServiceDoc
}
