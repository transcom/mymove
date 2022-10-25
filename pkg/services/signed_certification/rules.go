package signedcertification

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkSignedCertificationID checks that the ID has not been set
func checkSignedCertificationID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if !newSignedCertification.ID.IsNil() {
			verrs.Add("ID", "cannot manually set a new Signed Certification's UUID")
		}

		return verrs
	})
}

// checkSubmittingUserID check that the SubmittingUserID is not nil
func checkSubmittingUserID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.SubmittingUserID.IsNil() {
			verrs.Add("SubmittingUserID", "SubmittingUserID is required")
		}

		return verrs
	})
}

// checkMoveID check that the MoveID is not nil
func checkMoveID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.MoveID.IsNil() {
			verrs.Add("MoveID", "MoveID is required")
		}

		return verrs
	})
}

// checkPersonallyProcuredMoveID check that the PersonallyProcuredMoveID is either nil or a valid UUID
func checkPersonallyProcuredMoveID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.PersonallyProcuredMoveID != nil && newSignedCertification.PersonallyProcuredMoveID.IsNil() {
			verrs.Add("PersonallyProcuredMoveID", "PersonallyProcuredMoveID is not a valid UUID")
		}

		return verrs
	})
}

// checkPpmID check that the PpmID is either nil or a valid UUID
func checkPpmID() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.PpmID != nil && newSignedCertification.PpmID.IsNil() {
			verrs.Add("PpmID", "PpmID is not a valid UUID")
		}

		return verrs
	})
}

// checkCertificationType check that the CertificationType is a valid option
func checkCertificationType() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.CertificationType == nil {
			verrs.Add("CertificationType", "CertificationType is required")

			return verrs
		}

		for _, validCertificationType := range models.AllowedSignedCertificationTypes {
			if string(*newSignedCertification.CertificationType) == validCertificationType {
				return nil
			}
		}

		verrs.Add("CertificationType", "CertificationType is not a valid option")

		return verrs
	})
}

// checkCertificationText check that the CertificationText is not empty
func checkCertificationText() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.CertificationText == "" {
			verrs.Add("CertificationText", "CertificationText is required")
		}

		return verrs
	})
}

// checkSignature check that the Signature is not empty
func checkSignature() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
		verrs := validate.NewErrors()

		if newSignedCertification.Signature == "" {
			verrs.Add("Signature", "Signature is required")
		}

		return verrs
	})
}

// checkDate checks that the Date is valid
func checkDate() signedCertificationValidator {
	return signedCertificationValidatorFunc(func(_ appcontext.AppContext, newSignedCertification models.SignedCertification) error {
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
