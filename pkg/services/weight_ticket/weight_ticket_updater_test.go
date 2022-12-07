package weightticket

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite WeightTicketSuite) TestUpdateWeightTicket() {
	ppmShipmentUpdater := mocks.PPMShipmentUpdater{}

	setupForTest := func(appCtx appcontext.AppContext, overrides *models.WeightTicket, hasEmptyFiles bool, hasFullFiles bool, hasProofFiles bool) *models.WeightTicket {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		ppmShipment := testdatagen.MakeMinimalDefaultPPMShipment(suite.DB())

		baseDocumentAssertions := testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		}

		emptyDocument := testdatagen.MakeDocument(appCtx.DB(), baseDocumentAssertions)
		fullDocument := testdatagen.MakeDocument(appCtx.DB(), baseDocumentAssertions)
		proofOfOwnership := testdatagen.MakeDocument(appCtx.DB(), baseDocumentAssertions)

		now := time.Now()
		if hasEmptyFiles {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = &now
				}
				testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
					UserUpload: models.UserUpload{
						UploaderID: serviceMember.UserID,
						DocumentID: &emptyDocument.ID,
						Document:   emptyDocument,
						DeletedAt:  deletedAt,
					},
				})
			}
		}

		if hasFullFiles {
			for i := 0; i < 2; i++ {
				testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
					UserUpload: models.UserUpload{
						UploaderID: serviceMember.UserID,
						DocumentID: &fullDocument.ID,
						Document:   fullDocument,
					},
				})
			}
		}

		if hasProofFiles {
			for i := 0; i < 2; i++ {
				testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
					UserUpload: models.UserUpload{
						UploaderID: serviceMember.UserID,
						DocumentID: &proofOfOwnership.ID,
						Document:   proofOfOwnership,
					},
				})
			}
		}

		originalWeightTicket := models.WeightTicket{
			EmptyDocumentID:                   emptyDocument.ID,
			FullDocumentID:                    fullDocument.ID,
			ProofOfTrailerOwnershipDocumentID: proofOfOwnership.ID,
			PPMShipmentID:                     ppmShipment.ID,
			PPMShipment:                       ppmShipment,
		}

		if overrides != nil {
			testdatagen.MergeModels(&originalWeightTicket, overrides)
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&originalWeightTicket)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalWeightTicket.ID)

		return &originalWeightTicket
	}

	suite.Run("Returns an error if the original doesn't exist", func() {
		badWeightTicket := models.WeightTicket{
			ID: uuid.Must(uuid.NewV4()),
		}

		updater := NewCustomerWeightTicketUpdater(&ppmShipmentUpdater)

		updatedWeightTicket, err := updater.UpdateWeightTicket(suite.AppContextForTest(), badWeightTicket, "")

		suite.Nil(updatedWeightTicket)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for WeightTicket", badWeightTicket.ID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		appCtx := suite.AppContextForTest()

		originalWeightTicket := setupForTest(appCtx, nil, false, false, false)

		updater := NewCustomerWeightTicketUpdater(&ppmShipmentUpdater)

		updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *originalWeightTicket, "")

		suite.Nil(updatedWeightTicket)

		if suite.Error(updateErr) {
			suite.IsType(apperror.PreconditionFailedError{}, updateErr)

			suite.Equal(
				fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", originalWeightTicket.ID.String()),
				updateErr.Error(),
			)
		}
	})

	suite.Run("Successfully updates", func() {
		appCtx := suite.AppContextForTest()

		originalWeightTicket := setupForTest(appCtx, nil, true, true, false)

		updater := NewCustomerWeightTicketUpdater(&ppmShipmentUpdater)
		ppmShipmentUpdater.
			On(
				"UpdatePPMShipmentWithDefaultCheck",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.PPMShipment"),
				mock.AnythingOfType("uuid.UUID"),
			).
			Return(nil, nil)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(3000),
			MissingEmptyWeightTicket: models.BoolPointer(true),
			FullWeight:               models.PoundPointer(4200),
			MissingFullWeightTicket:  models.BoolPointer(true),
			OwnsTrailer:              models.BoolPointer(false),
			TrailerMeetsCriteria:     models.BoolPointer(false),
		}

		updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalWeightTicket.ID, updatedWeightTicket.ID)
		suite.Equal(originalWeightTicket.EmptyDocumentID, updatedWeightTicket.EmptyDocumentID)
		suite.Equal(originalWeightTicket.FullDocumentID, updatedWeightTicket.FullDocumentID)
		suite.Equal(originalWeightTicket.ProofOfTrailerOwnershipDocumentID, updatedWeightTicket.ProofOfTrailerOwnershipDocumentID)
		suite.Equal(*desiredWeightTicket.VehicleDescription, *updatedWeightTicket.VehicleDescription)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.EmptyWeight)
		suite.Equal(*desiredWeightTicket.MissingEmptyWeightTicket, *updatedWeightTicket.MissingEmptyWeightTicket)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.FullWeight)
		suite.Equal(*desiredWeightTicket.MissingFullWeightTicket, *updatedWeightTicket.MissingFullWeightTicket)
		suite.Equal(*desiredWeightTicket.OwnsTrailer, *updatedWeightTicket.OwnsTrailer)
		suite.Equal(*desiredWeightTicket.TrailerMeetsCriteria, *updatedWeightTicket.TrailerMeetsCriteria)
	})

	suite.Run("Succesfully updates when files are required", func() {
		appCtx := suite.AppContextForTest()

		originalWeightTicket := setupForTest(appCtx, nil, true, true, true)

		updater := NewCustomerWeightTicketUpdater(&ppmShipmentUpdater)
		ppmShipmentUpdater.
			On(
				"UpdatePPMShipmentWithDefaultCheck",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.PPMShipment"),
				mock.AnythingOfType("uuid.UUID"),
			).
			Return(nil, nil)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(3000),
			MissingEmptyWeightTicket: models.BoolPointer(false),
			FullWeight:               models.PoundPointer(4200),
			MissingFullWeightTicket:  models.BoolPointer(false),
			OwnsTrailer:              models.BoolPointer(true),
			TrailerMeetsCriteria:     models.BoolPointer(true),
		}

		updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalWeightTicket.ID, updatedWeightTicket.ID)
		suite.Equal(originalWeightTicket.EmptyDocumentID, updatedWeightTicket.EmptyDocumentID)
		suite.Equal(originalWeightTicket.FullDocumentID, updatedWeightTicket.FullDocumentID)
		suite.Equal(originalWeightTicket.ProofOfTrailerOwnershipDocumentID, updatedWeightTicket.ProofOfTrailerOwnershipDocumentID)
		suite.Equal(*desiredWeightTicket.VehicleDescription, *updatedWeightTicket.VehicleDescription)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.EmptyWeight)
		suite.Equal(*desiredWeightTicket.MissingEmptyWeightTicket, *updatedWeightTicket.MissingEmptyWeightTicket)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.FullWeight)
		suite.Equal(*desiredWeightTicket.MissingFullWeightTicket, *updatedWeightTicket.MissingFullWeightTicket)
		suite.Equal(*desiredWeightTicket.OwnsTrailer, *updatedWeightTicket.OwnsTrailer)
		suite.Equal(*desiredWeightTicket.TrailerMeetsCriteria, *updatedWeightTicket.TrailerMeetsCriteria)
		suite.Equal(1, len(updatedWeightTicket.EmptyDocument.UserUploads))
		suite.Equal(2, len(updatedWeightTicket.FullDocument.UserUploads))
		suite.Equal(2, len(updatedWeightTicket.ProofOfTrailerOwnershipDocument.UserUploads))
	})

	suite.Run("Fails to update when files are missing", func() {
		appCtx := suite.AppContextForTest()

		originalWeightTicket := setupForTest(appCtx, nil, false, false, false)

		updater := NewCustomerWeightTicketUpdater(&ppmShipmentUpdater)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(3000),
			MissingEmptyWeightTicket: models.BoolPointer(false),
			FullWeight:               models.PoundPointer(4200),
			MissingFullWeightTicket:  models.BoolPointer(false),
			OwnsTrailer:              models.BoolPointer(true),
			TrailerMeetsCriteria:     models.BoolPointer(true),
		}

		updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

		suite.Nil(updatedWeightTicket)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.Equal("Invalid input found while validating the weight ticket.", updateErr.Error())
	})

	suite.Run("Status and reason related", func() {
		suite.Run("successfully", func() {

			suite.Run("changes status and reason", func() {
				appCtx := suite.AppContextForTest()

				originalWeightTicket := testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{})

				updater := NewOfficeWeightTicketUpdater(&ppmShipmentUpdater)

				status := models.PPMDocumentStatusExcluded

				desiredWeightTicket := &models.WeightTicket{
					ID:     originalWeightTicket.ID,
					Status: &status,
					Reason: models.StringPointer("bad data"),
				}

				updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedWeightTicket)
				suite.Equal(*desiredWeightTicket.Status, *updatedWeightTicket.Status)
				suite.Equal(*desiredWeightTicket.Reason, *updatedWeightTicket.Reason)
			})

			suite.Run("changes reason", func() {
				appCtx := suite.AppContextForTest()

				status := models.PPMDocumentStatusExcluded
				originalWeightTicket := testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
					WeightTicket: models.WeightTicket{
						Status: &status,
						Reason: models.StringPointer("some temporary reason"),
					},
				})

				updater := NewOfficeWeightTicketUpdater(&ppmShipmentUpdater)

				desiredWeightTicket := &models.WeightTicket{
					ID:     originalWeightTicket.ID,
					Reason: models.StringPointer("bad data"),
				}

				updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedWeightTicket)
				suite.Equal(status, *updatedWeightTicket.Status)
				suite.Equal(*desiredWeightTicket.Reason, *updatedWeightTicket.Reason)
			})

			suite.Run("changes reason from rejected to approved", func() {
				appCtx := suite.AppContextForTest()

				status := models.PPMDocumentStatusExcluded
				originalWeightTicket := testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
					WeightTicket: models.WeightTicket{
						Status: &status,
						Reason: models.StringPointer("some temporary reason"),
					},
				})

				updater := NewOfficeWeightTicketUpdater(&ppmShipmentUpdater)

				desiredStatus := models.PPMDocumentStatusApproved
				desiredWeightTicket := &models.WeightTicket{
					ID:     originalWeightTicket.ID,
					Status: &desiredStatus,
					Reason: models.StringPointer(""),
				}

				updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedWeightTicket)
				suite.Equal(desiredStatus, *updatedWeightTicket.Status)
				suite.Equal((*string)(nil), updatedWeightTicket.Reason)
			})
		})

		suite.Run("fails", func() {
			suite.Run("to update when status or reason are changed", func() {
				appCtx := suite.AppContextForTest()

				originalWeightTicket := setupForTest(appCtx, nil, true, true, false)

				updater := NewCustomerWeightTicketUpdater(&ppmShipmentUpdater)

				status := models.PPMDocumentStatusExcluded

				desiredWeightTicket := &models.WeightTicket{
					ID:                       originalWeightTicket.ID,
					VehicleDescription:       models.StringPointer("2004 Ford Fiesta"),
					EmptyWeight:              models.PoundPointer(3000),
					MissingEmptyWeightTicket: models.BoolPointer(false),
					FullWeight:               models.PoundPointer(4000),
					MissingFullWeightTicket:  models.BoolPointer(false),
					OwnsTrailer:              models.BoolPointer(false),
					TrailerMeetsCriteria:     models.BoolPointer(false),
					Status:                   &status,
					Reason:                   models.StringPointer("bad data"),
				}

				updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

				suite.Nil(updatedWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("Invalid input found while validating the weight ticket.", updateErr.Error())
			})

			suite.Run("to update status", func() {
				appCtx := suite.AppContextForTest()

				status := models.PPMDocumentStatusExcluded
				originalWeightTicket := testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
					WeightTicket: models.WeightTicket{
						Status: &status,
						Reason: models.StringPointer("some temporary reason"),
					},
				})

				updater := NewOfficeWeightTicketUpdater(&ppmShipmentUpdater)

				desiredStatus := models.PPMDocumentStatusApproved
				desiredWeightTicket := &models.WeightTicket{
					ID:     originalWeightTicket.ID,
					Status: &desiredStatus,
					Reason: models.StringPointer("bad data"),
				}

				updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

				suite.Nil(updatedWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("Invalid input found while validating the weight ticket.", updateErr.Error())
			})

			suite.Run("to update because of invalid status", func() {
				appCtx := suite.AppContextForTest()

				originalWeightTicket := testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{})

				updater := NewOfficeWeightTicketUpdater(&ppmShipmentUpdater)

				status := models.PPMDocumentStatus("invalid status")
				desiredWeightTicket := &models.WeightTicket{
					ID:     originalWeightTicket.ID,
					Status: &status,
					Reason: models.StringPointer("bad data"),
				}

				updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

				suite.Nil(updatedWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("invalid input found while updating the WeightTicket", updateErr.Error())
			})
		})
	})
}
