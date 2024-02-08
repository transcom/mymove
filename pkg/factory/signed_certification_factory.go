package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildSignedCertification creates a single SignedCertification.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildSignedCertification(db *pop.Connection, customs []Customization, traits []Trait) models.SignedCertification {
	customs = setupCustomizations(customs, traits)

	var cSignedCertification models.SignedCertification
	if result := findValidCustomization(customs, SignedCertification); result != nil {
		cSignedCertification = result.Model.(models.SignedCertification)
		if result.LinkOnly {
			return cSignedCertification
		}
	}

	move := BuildMove(db, customs, traits)

	certificationType := models.SignedCertificationTypeSHIPMENT
	signedCertification := models.SignedCertification{
		MoveID:            move.ID,
		SubmittingUserID:  move.Orders.ServiceMember.UserID,
		CertificationType: &certificationType,
		CertificationText: "LEGAL TEXT",
		Signature:         "SIGNATURE",
		Date:              testdatagen.NextValidMoveDate,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&signedCertification, cSignedCertification)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &signedCertification)
	}

	return signedCertification
}
