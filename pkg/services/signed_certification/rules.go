package signedcertification

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkSignedCertificationID checks that the ID has not been set if creating, otherwise checks that it hasn't changed
func checkSignedCertificationID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if originalSignedCertification == nil {
			if !newSignedCertification.ID.IsNil() {
				verrs.Add("ID", "cannot manually set a new Signed Certification's UUID")
			}
		} else {
			if newSignedCertification.ID != originalSignedCertification.ID {
				verrs.Add("ID", "cannot change a Signed Certification's UUID")
			}
		}

		return verrs
	})
}

// checkSubmittingUserID check that the SubmittingUserID is not nil if creating, otherwise checks that it hasn't changed
func checkSubmittingUserID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.SubmittingUserID.IsNil() {
			verrs.Add("SubmittingUserID", "SubmittingUserID is required")
		}

		if originalSignedCertification != nil && newSignedCertification.SubmittingUserID != originalSignedCertification.SubmittingUserID {
			verrs.Add("SubmittingUserID", "SubmittingUserID cannot be changed")
		}

		return verrs
	})
}

// checkMoveID check that the MoveID is not nil if creating, otherwise checks that it hasn't changed
func checkMoveID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.MoveID.IsNil() {
			verrs.Add("MoveID", "MoveID is required")
		}

		if originalSignedCertification != nil && newSignedCertification.MoveID != originalSignedCertification.MoveID {
			verrs.Add("MoveID", "MoveID cannot be changed")
		}

		return verrs
	})
}

// checkPersonallyProcuredMoveID check that the PersonallyProcuredMoveID is either nil or a valid UUID if creating,
// otherwise checks that it is valid and hasn't changed.
func checkPersonallyProcuredMoveID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.PersonallyProcuredMoveID != nil && newSignedCertification.PersonallyProcuredMoveID.IsNil() {
			verrs.Add("PersonallyProcuredMoveID", "PersonallyProcuredMoveID is not a valid UUID")
		}

		if originalSignedCertification != nil {
			if (newSignedCertification.PersonallyProcuredMoveID != nil && originalSignedCertification.PersonallyProcuredMoveID == nil) ||
				(newSignedCertification.PersonallyProcuredMoveID == nil && originalSignedCertification.PersonallyProcuredMoveID != nil) ||
				(newSignedCertification.PersonallyProcuredMoveID != nil && originalSignedCertification.PersonallyProcuredMoveID != nil && *newSignedCertification.PersonallyProcuredMoveID != *originalSignedCertification.PersonallyProcuredMoveID) {
				verrs.Add("PersonallyProcuredMoveID", "PersonallyProcuredMoveID cannot be changed")
			}
		}

		return verrs
	})
}

// checkPpmID check that the PpmID is either nil or a valid UUID if creating, otherwise checks that it is valid and
// :hasn't changed.
func checkPpmID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.PpmID != nil && newSignedCertification.PpmID.IsNil() {
			verrs.Add("PpmID", "PpmID is not a valid UUID")
		}

		if originalSignedCertification != nil {
			if (newSignedCertification.PpmID == nil && originalSignedCertification.PpmID != nil) ||
				(newSignedCertification.PpmID != nil && originalSignedCertification.PpmID == nil) ||
				(newSignedCertification.PpmID != nil && originalSignedCertification.PpmID != nil && *newSignedCertification.PpmID != *originalSignedCertification.PpmID) {
				verrs.Add("PpmID", "PpmID cannot be changed")
			}
		}

		return verrs
	})
}

// checkCertificationType check that the CertificationType is a valid option
func checkCertificationType() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.CertificationType == nil {
			verrs.Add("CertificationType", "CertificationType is required")

			return verrs
		}

		if originalSignedCertification == nil {
			for _, validCertificationType := range models.AllowedSignedCertificationTypes {
				if string(*newSignedCertification.CertificationType) == validCertificationType {
					return nil
				}
			}

			verrs.Add("CertificationType", "CertificationType is not a valid option")
		} else {
			if *newSignedCertification.CertificationType != *originalSignedCertification.CertificationType {
				verrs.Add("CertificationType", "CertificationType cannot be changed")
			}
		}

		return verrs
	})
}

// checkCertificationText check that the CertificationText is not empty
func checkCertificationText() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.CertificationText == "" {
			verrs.Add("CertificationText", "CertificationText is required")
		}

		return verrs
	})
}

// checkSignature check that the Signature is not empty
func checkSignature() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.Signature == "" {
			verrs.Add("Signature", "Signature is required")
		}

		return verrs
	})
}

// checkDate checks that the Date is valid
func checkDate() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification, originalSignedCertification *models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.Date.IsZero() {
			verrs.Add("Date", "Date is required")
		}

		return verrs
	})
}

// basicSignedCertificationChecks returns the rules that should run for any SignedCertification validation
func basicSignedCertificationChecks() []signedCertificationValidator {
	return []signedCertificationValidator{
		checkSignedCertificationID(),
		checkSubmittingUserID(),
		checkMoveID(),
		checkPersonallyProcuredMoveID(),
		checkPpmID(),
		checkCertificationType(),
		checkCertificationText(),
		checkSignature(),
		checkDate(),
	}
}
