package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeSignedCertification creates a single signed certification
func MakeSignedCertification(db *pop.Connection, assertions Assertions) models.SignedCertification {
	moveID := assertions.SignedCertification.MoveID
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.SignedCertification.MoveID) {
		moveID = MakeMove(db, assertions).ID
	}
	userID := assertions.SignedCertification.SubmittingUserID
	if isZeroUUID(assertions.SignedCertification.SubmittingUserID) {
		sm := MakeServiceMember(db, assertions)
		userID = sm.UserID
	}

	certificationType := models.SignedCertificationTypePPM
	signedCertification := models.SignedCertification{
		MoveID:                   moveID,
		SubmittingUserID:         userID,
		PersonallyProcuredMoveID: nil,
		ShipmentID:               nil,
		CertificationType:        &certificationType,
		CertificationText:        "LEGAL TEXT",
		Signature:                "SIGNATURE",
		Date:                     NextValidMoveDate,
	}

	// Overwrite values with those from assertions
	mergeModels(&signedCertification, assertions.SignedCertification)

	mustCreate(db, &signedCertification)

	return signedCertification
}

// MakeDefaultSignedCertification returns a MoveDocument with default values
func MakeDefaultSignedCertification(db *pop.Connection) models.SignedCertification {
	return MakeSignedCertification(db, Assertions{})
}
