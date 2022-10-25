package signedcertification

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *SignedCertificationSuite) TestSignedCertificationRules() {
	suite.Run("checkSignedCertificationID", func() {
		suite.Run("success", func() {
			err := checkSignedCertificationID().Validate(suite.AppContextForTest(), models.SignedCertification{})

			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			signedCertification := models.SignedCertification{ID: uuid.Must(uuid.NewV4())}

			err := checkSignedCertificationID().Validate(suite.AppContextForTest(), signedCertification)

			suite.NotNil(err)
			suite.Contains(err.Error(), "cannot manually set a new Signed Certification's UUID")
		})
	})

	suite.Run("checkSubmittingUserID", func() {
		suite.Run("success", func() {
			signedCertification := models.SignedCertification{SubmittingUserID: uuid.Must(uuid.NewV4())}

			err := checkSubmittingUserID().Validate(suite.AppContextForTest(), signedCertification)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			err := checkSubmittingUserID().Validate(suite.AppContextForTest(), models.SignedCertification{})

			suite.NotNil(err)
			suite.Contains(err.Error(), "SubmittingUserID is required")
		})
	})

	suite.Run("checkMoveID", func() {
		suite.Run("success", func() {
			signedCertification := models.SignedCertification{MoveID: uuid.Must(uuid.NewV4())}

			err := checkMoveID().Validate(suite.AppContextForTest(), signedCertification)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			err := checkMoveID().Validate(suite.AppContextForTest(), models.SignedCertification{})

			suite.NotNil(err)
			suite.Contains(err.Error(), "MoveID is required")
		})
	})

	suite.Run("checkPersonallyProcuredMoveID", func() {
		successCases := map[string]*uuid.UUID{
			"nil":   nil,
			"valid": models.UUIDPointer(uuid.Must(uuid.NewV4())),
		}

		for name, id := range successCases {
			name := name
			id := id

			suite.Run(fmt.Sprintf("success when ID is %s", name), func() {
				signedCertification := models.SignedCertification{PersonallyProcuredMoveID: id}

				err := checkPersonallyProcuredMoveID().Validate(suite.AppContextForTest(), signedCertification)

				suite.NilOrNoVerrs(err)
			})
		}

		suite.Run("failure when ID is an invalid UUID", func() {
			signedCertification := models.SignedCertification{PersonallyProcuredMoveID: &uuid.Nil}

			err := checkPersonallyProcuredMoveID().Validate(suite.AppContextForTest(), signedCertification)

			suite.NotNil(err)
			suite.Contains(err.Error(), "PersonallyProcuredMoveID is not a valid UUID")
		})
	})

	suite.Run("checkPpmID", func() {
		successCases := map[string]*uuid.UUID{
			"nil":   nil,
			"valid": models.UUIDPointer(uuid.Must(uuid.NewV4())),
		}

		for name, id := range successCases {
			name := name
			id := id

			suite.Run(fmt.Sprintf("success when ID is %s", name), func() {
				signedCertification := models.SignedCertification{PpmID: id}

				err := checkPpmID().Validate(suite.AppContextForTest(), signedCertification)

				suite.NilOrNoVerrs(err)
			})
		}

		suite.Run("failure when ID is an invalid UUID", func() {
			signedCertification := models.SignedCertification{PpmID: &uuid.Nil}

			err := checkPpmID().Validate(suite.AppContextForTest(), signedCertification)

			suite.NotNil(err)
			suite.Contains(err.Error(), "PpmID is not a valid UUID")
		})
	})

	suite.Run("checkCertificationType", func() {
		successCases := map[string]models.SignedCertificationType{
			"Shipment":    models.SignedCertificationTypeSHIPMENT,
			"PPM Payment": models.SignedCertificationTypePPMPAYMENT,
		}

		for name, certType := range successCases {
			name := name
			certType := certType

			suite.Run(fmt.Sprintf("success when type is %s", name), func() {
				signedCertification := models.SignedCertification{CertificationType: &certType}

				err := checkCertificationType().Validate(suite.AppContextForTest(), signedCertification)

				suite.NilOrNoVerrs(err)
			})
		}

		blankCertType := models.SignedCertificationType("")

		failureCases := map[string]struct {
			certType *models.SignedCertificationType
			error    string
		}{
			"nil": {
				certType: nil,
				error:    "CertificationType is required",
			},
			"blank": {
				certType: &blankCertType,
				error:    "CertificationType is not a valid option",
			},
		}

		for name, failureCase := range failureCases {
			name := name
			failureCase := failureCase

			suite.Run(fmt.Sprintf("failure when type is %s", name), func() {
				signedCertification := models.SignedCertification{CertificationType: failureCase.certType}

				err := checkCertificationType().Validate(suite.AppContextForTest(), signedCertification)

				suite.NotNil(err)
				suite.Contains(err.Error(), failureCase.error)
			})
		}
	})

	suite.Run("checkCertificationText", func() {
		suite.Run("success", func() {
			signedCertification := models.SignedCertification{CertificationText: "I certify that I have read and understand the information above."}

			err := checkCertificationText().Validate(suite.AppContextForTest(), signedCertification)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			err := checkCertificationText().Validate(suite.AppContextForTest(), models.SignedCertification{})

			suite.NotNil(err)
			suite.Contains(err.Error(), "CertificationText is required")
		})
	})

	suite.Run("checkSignature", func() {
		suite.Run("success", func() {
			signedCertification := models.SignedCertification{Signature: "Best Customer"}

			err := checkSignature().Validate(suite.AppContextForTest(), signedCertification)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			err := checkSignature().Validate(suite.AppContextForTest(), models.SignedCertification{})

			suite.NotNil(err)
			suite.Contains(err.Error(), "Signature is required")
		})
	})

	suite.Run("checkDate", func() {
		suite.Run("success", func() {
			signedCertification := models.SignedCertification{Date: time.Now()}

			err := checkDate().Validate(suite.AppContextForTest(), signedCertification)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			err := checkDate().Validate(suite.AppContextForTest(), models.SignedCertification{})

			suite.NotNil(err)
			suite.Contains(err.Error(), "Date is required")
		})
	})
}
