package signedcertification

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *SignedCertificationSuite) TestCheckSignedCertificationID() {
	suite.Run("Success", func() {
		suite.Run("Create a signed certification without setting an ID", func() {
			err := checkSignedCertificationID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				nil,
			)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("Update a signed certification with a matching ID", func() {
			id := uuid.Must(uuid.NewV4())

			err := checkSignedCertificationID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{ID: id},
				&models.SignedCertification{ID: id},
			)

			suite.NilOrNoVerrs(err)
		})
	})

	suite.Run("Failure", func() {
		suite.Run("Try to create a signed certification with an ID", func() {
			err := checkSignedCertificationID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{ID: uuid.Must(uuid.NewV4())},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "cannot manually set a new Signed Certification's UUID")
		})

		suite.Run("Try to update a signed certification with a different ID", func() {
			err := checkSignedCertificationID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{ID: uuid.Must(uuid.NewV4())},
				&models.SignedCertification{ID: uuid.Must(uuid.NewV4())},
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "cannot change a Signed Certification's UUID")
		})
	})
}

func (suite *SignedCertificationSuite) TestCheckSubmittingUserID() {
	suite.Run("Success", func() {
		suite.Run("Create a signed certification with a SubmittingUserID", func() {
			err := checkSubmittingUserID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{SubmittingUserID: uuid.Must(uuid.NewV4())},
				nil,
			)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("Update a signed certification with a matching SubmittingUserID", func() {
			id := uuid.Must(uuid.NewV4())

			err := checkSubmittingUserID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{SubmittingUserID: id},
				&models.SignedCertification{SubmittingUserID: id},
			)

			suite.NilOrNoVerrs(err)
		})
	})

	suite.Run("Failure", func() {
		suite.Run("Try to create a signed certification without a SubmittingUserID", func() {
			err := checkSubmittingUserID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "SubmittingUserID is required")
		})

		suite.Run("Try to update a signed certification with a different SubmittingUserID", func() {
			err := checkSubmittingUserID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{SubmittingUserID: uuid.Must(uuid.NewV4())},
				&models.SignedCertification{SubmittingUserID: uuid.Must(uuid.NewV4())},
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "SubmittingUserID cannot be changed")
		})
	})
}

func (suite *SignedCertificationSuite) TestCheckMoveID() {
	suite.Run("Success", func() {
		suite.Run("Create a signed certification with a MoveID", func() {
			signedCertification := models.SignedCertification{MoveID: uuid.Must(uuid.NewV4())}

			err := checkMoveID().Validate(
				suite.AppContextForTest(),
				signedCertification,
				nil,
			)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("Update a signed certification with a matching MoveID", func() {
			id := uuid.Must(uuid.NewV4())

			err := checkMoveID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{MoveID: id},
				&models.SignedCertification{MoveID: id},
			)

			suite.NilOrNoVerrs(err)
		})
	})

	suite.Run("Failure", func() {
		suite.Run("Try to create a signed certification without a MoveID", func() {
			err := checkMoveID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "MoveID is required")
		})

		suite.Run("Try to update a signed certification with a different MoveID", func() {
			err := checkMoveID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{MoveID: uuid.Must(uuid.NewV4())},
				&models.SignedCertification{MoveID: uuid.Must(uuid.NewV4())},
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "MoveID cannot be changed")
		})
	})
}

func (suite *SignedCertificationSuite) TestCheckPersonallyProcuredMoveID() {
	successCases := map[string]*uuid.UUID{
		"nil":   nil,
		"valid": models.UUIDPointer(uuid.Must(uuid.NewV4())),
	}

	for name, id := range successCases {
		name := name
		id := id

		suite.Run(fmt.Sprintf("Success creating when PersonallyProcuredMoveID is %s", name), func() {

			err := checkPersonallyProcuredMoveID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{PersonallyProcuredMoveID: id},
				nil,
			)

			suite.NilOrNoVerrs(err)
		})

		suite.Run(fmt.Sprintf("Success updating when PersonallyProcuredMoveID is %s", name), func() {
			originalPersonallyProcuredMoveID := id
			newPersonallyProcuredMoveID := id

			if id != nil {
				// Copying the value to make sure we're comparing values rather than pointers
				originalPersonallyProcuredMoveID = models.UUIDPointer(*id)
				newPersonallyProcuredMoveID = models.UUIDPointer(*id)
			}

			err := checkPersonallyProcuredMoveID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{PersonallyProcuredMoveID: newPersonallyProcuredMoveID},
				&models.SignedCertification{PersonallyProcuredMoveID: originalPersonallyProcuredMoveID},
			)

			suite.NilOrNoVerrs(err)
		})
	}

	suite.Run("Failure", func() {
		suite.Run("Try to create a signed certification with an invalid PersonallyProcuredMoveID", func() {
			err := checkPersonallyProcuredMoveID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{PersonallyProcuredMoveID: &uuid.Nil},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "PersonallyProcuredMoveID is not a valid UUID")
		})

		updateFailureCases := map[string]*uuid.UUID{
			"an invalid UUID":  &uuid.Nil,
			"a different UUID": models.UUIDPointer(uuid.Must(uuid.NewV4())),
		}

		for name, id := range updateFailureCases {
			name := name
			id := id

			suite.Run(fmt.Sprintf("Try to update a signed certification with %s checkPersonallyProcuredMoveID", name), func() {
				err := checkPersonallyProcuredMoveID().Validate(
					suite.AppContextForTest(),
					models.SignedCertification{PersonallyProcuredMoveID: id},
					&models.SignedCertification{PersonallyProcuredMoveID: models.UUIDPointer(uuid.Must(uuid.NewV4()))},
				)

				suite.NotNil(err)
				suite.Contains(err.Error(), "PersonallyProcuredMoveID cannot be changed")
			})
		}
	})
}

func (suite *SignedCertificationSuite) TestCheckPpmID() {
	successCases := map[string]*uuid.UUID{
		"nil":   nil,
		"valid": models.UUIDPointer(uuid.Must(uuid.NewV4())),
	}

	for name, id := range successCases {
		name := name
		id := id

		suite.Run(fmt.Sprintf("Success creating when PpmID is %s", name), func() {

			err := checkPpmID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{PpmID: id},
				nil,
			)

			suite.NilOrNoVerrs(err)
		})

		suite.Run(fmt.Sprintf("Success updating when PpmID is %s", name), func() {
			originalPpmID := id
			newPpmID := id

			if id != nil {
				// Copying the value to make sure we're comparing values rather than pointers
				originalPpmID = models.UUIDPointer(*id)
				newPpmID = models.UUIDPointer(*id)
			}

			err := checkPpmID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{PpmID: newPpmID},
				&models.SignedCertification{PpmID: originalPpmID},
			)

			suite.NilOrNoVerrs(err)
		})
	}

	suite.Run("Failure", func() {
		suite.Run("Try to create a signed certification with an invalid PpmID", func() {
			err := checkPpmID().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{PpmID: &uuid.Nil},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "PpmID is not a valid UUID")
		})

		updateFailureCases := map[string]*uuid.UUID{
			"an invalid UUID":  &uuid.Nil,
			"a different UUID": models.UUIDPointer(uuid.Must(uuid.NewV4())),
		}

		for name, id := range updateFailureCases {
			name := name
			id := id

			suite.Run(fmt.Sprintf("Try to update a signed certification with %s PpmID", name), func() {
				err := checkPpmID().Validate(
					suite.AppContextForTest(),
					models.SignedCertification{PpmID: id},
					&models.SignedCertification{PpmID: models.UUIDPointer(uuid.Must(uuid.NewV4()))},
				)

				suite.NotNil(err)
				suite.Contains(err.Error(), "PpmID cannot be changed")
			})
		}
	})
}

func (suite *SignedCertificationSuite) TestCheckCertificationType() {
	successCases := map[string]models.SignedCertificationType{
		"Shipment":    models.SignedCertificationTypeSHIPMENT,
		"PPM Payment": models.SignedCertificationTypePPMPAYMENT,
	}

	for name, certType := range successCases {
		name := name
		certType := certType

		suite.Run(fmt.Sprintf("Success creating when type is %s", name), func() {
			err := checkCertificationType().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{CertificationType: &certType},
				nil,
			)

			suite.NilOrNoVerrs(err)
		})

		suite.Run(fmt.Sprintf("Success updating when type is %s", name), func() {
			// Copying the value to make sure we're comparing values rather than pointers
			originalCertType := certType
			newCertType := certType

			err := checkCertificationType().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{CertificationType: &newCertType},
				&models.SignedCertification{CertificationType: &originalCertType},
			)

			suite.NilOrNoVerrs(err)
		})
	}

	suite.Run("Creation failure", func() {

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

			suite.Run(fmt.Sprintf("when type is %s", name), func() {
				err := checkCertificationType().Validate(
					suite.AppContextForTest(),
					models.SignedCertification{CertificationType: failureCase.certType},
					nil,
				)

				suite.NotNil(err)
				suite.Contains(err.Error(), failureCase.error)
			})
		}
	})

	suite.Run("Update failure if certification type changes", func() {
		shipmentCertType := models.SignedCertificationTypeSHIPMENT
		ppmCertType := models.SignedCertificationTypePPMPAYMENT

		err := checkCertificationType().Validate(
			suite.AppContextForTest(),
			models.SignedCertification{CertificationType: &ppmCertType},
			&models.SignedCertification{CertificationType: &shipmentCertType},
		)

		suite.NotNil(err)
		suite.Contains(err.Error(), "CertificationType cannot be changed")
	})
}

func (suite *SignedCertificationSuite) TestCheckCertificationText() {
	certText := "I certify that I have read and understand the information above."

	suite.Run("Create signed certification with CertificationText", func() {
		err := checkCertificationText().Validate(
			suite.AppContextForTest(),
			models.SignedCertification{CertificationText: certText},
			nil,
		)

		suite.NilOrNoVerrs(err)
	})

	updateSuccessCases := map[string]string{
		"the same": certText,
		"new":      "I certify that I have read and understand the information above. And I agree to it.",
	}

	for name, text := range updateSuccessCases {
		name := name
		text := text

		suite.Run(fmt.Sprintf("Update signed certification with %s CertificationText", name), func() {
			err := checkCertificationText().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{CertificationText: text},
				&models.SignedCertification{CertificationText: certText},
			)

			suite.NilOrNoVerrs(err)
		})
	}

	suite.Run("Failure", func() {
		suite.Run("Create signed certification without CertificationText", func() {
			err := checkCertificationText().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "CertificationText is required")
		})

		suite.Run("Update signed certification without CertificationText", func() {
			err := checkCertificationText().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				&models.SignedCertification{CertificationText: certText},
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "CertificationText is required")
		})
	})
}

func (suite *SignedCertificationSuite) TestCheckSignature() {
	signature := "Best Customer"

	suite.Run("Create signed certification with Signature", func() {
		err := checkSignature().Validate(
			suite.AppContextForTest(),
			models.SignedCertification{Signature: signature},
			nil,
		)

		suite.NilOrNoVerrs(err)
	})

	updateSuccessCases := map[string]string{
		"the same": signature,
		"a new":    "Best Customer Ever",
	}

	for name, newSignature := range updateSuccessCases {
		name := name
		newSignature := newSignature

		suite.Run(fmt.Sprintf("Update signed certification with %s Signature", name), func() {
			err := checkSignature().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{Signature: newSignature},
				&models.SignedCertification{Signature: signature},
			)

			suite.NilOrNoVerrs(err)
		})
	}

	suite.Run("Failure", func() {
		suite.Run("Create signed certification without Signature", func() {
			err := checkSignature().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "Signature is required")
		})

		suite.Run("Update signed certification without Signature", func() {
			err := checkSignature().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				&models.SignedCertification{Signature: signature},
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "Signature is required")
		})
	})
}

func (suite *SignedCertificationSuite) TestCheckDate() {
	date := time.Now()

	suite.Run("Create signed certification with Date", func() {
		err := checkDate().Validate(
			suite.AppContextForTest(),
			models.SignedCertification{Date: date},
			nil,
		)

		suite.NilOrNoVerrs(err)
	})

	updateSuccessCases := map[string]time.Time{
		"the same": date,
		"a new":    date.AddDate(0, 0, 1),
	}

	for name, newDate := range updateSuccessCases {
		name := name
		newDate := newDate

		suite.Run(fmt.Sprintf("Update signed certification with %s Date", name), func() {
			err := checkDate().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{Date: newDate},
				&models.SignedCertification{Date: date},
			)

			suite.NilOrNoVerrs(err)
		})
	}

	suite.Run("Failure", func() {
		suite.Run("Create signed certification without Date", func() {
			err := checkDate().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				nil,
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "Date is required")
		})

		suite.Run("Update signed certification without Date", func() {
			err := checkDate().Validate(
				suite.AppContextForTest(),
				models.SignedCertification{},
				&models.SignedCertification{Date: date},
			)

			suite.NotNil(err)
			suite.Contains(err.Error(), "Date is required")
		})
	})
}
