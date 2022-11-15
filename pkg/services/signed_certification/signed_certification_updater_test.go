package signedcertification

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *SignedCertificationSuite) TestUpdateSignedCertification() {
	var originalSignedCertification models.SignedCertification

	suite.PreloadData(func() {
		move := testdatagen.MakeDefaultMove(suite.DB())
		serviceMember := move.Orders.ServiceMember

		shipmentCertType := models.SignedCertificationTypeSHIPMENT

		originalSignedCertification = models.SignedCertification{
			SubmittingUserID:  serviceMember.User.ID,
			MoveID:            move.ID,
			CertificationType: &shipmentCertType,
			CertificationText: "I certify that the information I have provided is true and complete to the best of my knowledge.",
			Signature:         fmt.Sprintf("%s %s", *serviceMember.FirstName, *serviceMember.LastName),
			Date:              testdatagen.NextValidMoveDate,
		}

		verrs, err := suite.DB().ValidateAndCreate(&originalSignedCertification)

		suite.FatalNoVerrs(verrs)
		suite.FatalNoError(err)
		suite.FatalNotNil(originalSignedCertification.ID)
	})

	suite.Run("Returns an error if the original signed certification is not found", func() {
		newSignedCertification := models.SignedCertification{
			ID: uuid.Must(uuid.NewV4()),
		}

		updater := NewSignedCertificationUpdater()

		updatedSignedCertification, err := updater.UpdateSignedCertification(
			suite.AppContextForTest(),
			newSignedCertification,
			"",
		)

		suite.Nil(updatedSignedCertification)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for SignedCertification", newSignedCertification.ID),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		updater := NewSignedCertificationUpdater()

		updatedSignedCertification, err := updater.UpdateSignedCertification(
			suite.AppContextForTest(),
			originalSignedCertification,
			"",
		)

		suite.Nil(updatedSignedCertification)

		if suite.Error(err) {
			suite.IsType(apperror.PreconditionFailedError{}, err)

			suite.Equal(
				fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", originalSignedCertification.ID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Merges new changes before validating", func() {
		updater := NewSignedCertificationUpdater()

		badSignedCertification := models.SignedCertification{
			ID:        originalSignedCertification.ID,
			Signature: "",
		}

		updatedSignedCertification, err := updater.UpdateSignedCertification(
			suite.AppContextForTest(),
			badSignedCertification,
			etag.GenerateEtag(originalSignedCertification.UpdatedAt),
		)

		// If we didn't merge before validating, then the signature would be empty and we'd get an error.
		// This also makes it hard to test that validation is being run though because the merge takes care of basically
		// every error.
		suite.NoError(err)

		suite.NotNil(updatedSignedCertification)
	})

	suite.Run("Can successfully update a signed certification", func() {
		updater := NewSignedCertificationUpdater()

		newSignedCertification := originalSignedCertification
		newSignedCertification.Signature = "New Signature"

		updatedSignedCertification, err := updater.UpdateSignedCertification(
			suite.AppContextForTest(),
			newSignedCertification,
			etag.GenerateEtag(originalSignedCertification.UpdatedAt),
		)

		suite.NoError(err)

		if suite.NotNil(updatedSignedCertification) {
			suite.Equal(newSignedCertification.Signature, updatedSignedCertification.Signature)
		}
	})
}

func (suite *SignedCertificationSuite) TestMergeSignedCertification() {
	suite.Run("Can change certification text, signature, and date", func() {
		shipmentCertType := models.SignedCertificationTypeSHIPMENT
		today := time.Now()

		originalSignedCertification := models.SignedCertification{
			ID:                       uuid.Must(uuid.NewV4()),
			SubmittingUserID:         uuid.Must(uuid.NewV4()),
			MoveID:                   uuid.Must(uuid.NewV4()),
			PersonallyProcuredMoveID: models.UUIDPointer(uuid.Must(uuid.NewV4())),
			PpmID:                    models.UUIDPointer(uuid.Must(uuid.NewV4())),
			CertificationType:        &shipmentCertType,
			CertificationText:        "Original Certification Text",
			Signature:                "Original Signature",
			Date:                     today,
		}

		newSignedCertification := models.SignedCertification{
			ID:                       uuid.Must(uuid.NewV4()),
			SubmittingUserID:         uuid.Must(uuid.NewV4()),
			MoveID:                   uuid.Must(uuid.NewV4()),
			PersonallyProcuredMoveID: models.UUIDPointer(uuid.Must(uuid.NewV4())),
			PpmID:                    models.UUIDPointer(uuid.Must(uuid.NewV4())),
			CertificationType:        &shipmentCertType,
			CertificationText:        "New Certification Text",
			Signature:                "New Signature",
			Date:                     today.AddDate(0, 0, 1),
		}

		mergedSignedCertification := mergeSignedCertification(newSignedCertification, &originalSignedCertification)

		// fields that should be unchanged
		suite.Equal(originalSignedCertification.ID, mergedSignedCertification.ID)
		suite.Equal(originalSignedCertification.SubmittingUserID, mergedSignedCertification.SubmittingUserID)
		suite.Equal(originalSignedCertification.MoveID, mergedSignedCertification.MoveID)
		suite.Equal(originalSignedCertification.PersonallyProcuredMoveID, mergedSignedCertification.PersonallyProcuredMoveID)
		suite.Equal(originalSignedCertification.PpmID, mergedSignedCertification.PpmID)
		suite.Equal(originalSignedCertification.CertificationType, mergedSignedCertification.CertificationType)

		// fields that should be changed
		suite.NotEqual(originalSignedCertification.CertificationText, mergedSignedCertification.CertificationText)
		suite.Equal(newSignedCertification.CertificationText, mergedSignedCertification.CertificationText)
		suite.NotEqual(originalSignedCertification.Signature, mergedSignedCertification.Signature)
		suite.Equal(newSignedCertification.Signature, mergedSignedCertification.Signature)
		suite.NotEqual(originalSignedCertification.Date, mergedSignedCertification.Date)
		suite.True(newSignedCertification.Date.Equal(mergedSignedCertification.Date), "new and merged dates should be equal")
	})

	suite.Run("Does not change certification text, signature, and date if they are empty", func() {
		shipmentCertType := models.SignedCertificationTypeSHIPMENT
		today := time.Now()

		originalSignedCertification := models.SignedCertification{
			ID:                       uuid.Must(uuid.NewV4()),
			SubmittingUserID:         uuid.Must(uuid.NewV4()),
			MoveID:                   uuid.Must(uuid.NewV4()),
			PersonallyProcuredMoveID: models.UUIDPointer(uuid.Must(uuid.NewV4())),
			PpmID:                    models.UUIDPointer(uuid.Must(uuid.NewV4())),
			CertificationType:        &shipmentCertType,
			CertificationText:        "Original Certification Text",
			Signature:                "Original Signature",
			Date:                     today,
		}

		newSignedCertification := models.SignedCertification{
			CertificationText: "",
			Signature:         "",
			Date:              time.Time{},
		}

		mergedSignedCertification := mergeSignedCertification(newSignedCertification, &originalSignedCertification)

		// fields that should be unchanged
		suite.Equal(originalSignedCertification.ID, mergedSignedCertification.ID)
		suite.Equal(originalSignedCertification.SubmittingUserID, mergedSignedCertification.SubmittingUserID)
		suite.Equal(originalSignedCertification.MoveID, mergedSignedCertification.MoveID)
		suite.Equal(originalSignedCertification.PersonallyProcuredMoveID, mergedSignedCertification.PersonallyProcuredMoveID)
		suite.Equal(originalSignedCertification.PpmID, mergedSignedCertification.PpmID)
		suite.Equal(originalSignedCertification.CertificationType, mergedSignedCertification.CertificationType)
		suite.Equal(originalSignedCertification.CertificationText, mergedSignedCertification.CertificationText)
		suite.Equal(originalSignedCertification.Signature, mergedSignedCertification.Signature)
		suite.True(originalSignedCertification.Date.Equal(mergedSignedCertification.Date), "original and merged dates should be equal")
	})
}
