package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

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

	certificationType := models.SignedCertificationTypeSHIPMENT
	signedCertification := models.SignedCertification{
		MoveID:                   moveID,
		SubmittingUserID:         userID,
		PersonallyProcuredMoveID: nil,
		CertificationType:        &certificationType,
		CertificationText:        "LEGAL TEXT",
		Signature:                "SIGNATURE",
		Date:                     NextValidMoveDate,
	}

	// Overwrite values with those from assertions
	mergeModels(&signedCertification, assertions.SignedCertification)

	mustCreate(db, &signedCertification, assertions.Stub)

	return signedCertification
}

// MakeSignedCertificationForPPM creates a single signed certification
func MakeSignedCertificationForPPM(db *pop.Connection, assertions Assertions) models.SignedCertification {
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

	ppmID := assertions.PersonallyProcuredMove.ID
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.PersonallyProcuredMove.ID) {
		ppmID = MakePPM(db, assertions).ID
	}

	certificationType := models.SignedCertificationTypeSHIPMENT
	signedCertification := models.SignedCertification{
		MoveID:                   moveID,
		SubmittingUserID:         userID,
		PersonallyProcuredMoveID: &ppmID,
		CertificationType:        &certificationType,
		CertificationText:        "LEGAL TEXT",
		Signature:                "SIGNATURE",
		Date:                     NextValidMoveDate,
	}

	// Overwrite values with those from assertions
	mergeModels(&signedCertification, assertions.SignedCertification)

	mustCreate(db, &signedCertification, assertions.Stub)

	return signedCertification
}

// MakeDefaultSignedCertification returns a MoveDocument with default values
func MakeDefaultSignedCertification(db *pop.Connection) models.SignedCertification {
	return MakeSignedCertification(db, Assertions{})
}
