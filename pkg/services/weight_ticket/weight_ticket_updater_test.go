package weightticket

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *WeightTicketSuite) TestUpdateWeightTicket() {
	ppmShipmentUpdater := mocks.PPMShipmentUpdater{}

	setupForTest := func(overrides *models.WeightTicket, hasEmptyFiles bool, hasFullFiles bool, hasProofFiles bool) *models.WeightTicket {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		emptyDocument := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)
		fullDocument := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)
		proofOfOwnership := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		now := time.Now()
		if hasEmptyFiles {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time

				markAsDeleted := i == 1

				if markAsDeleted {
					deletedAt = &now
				}

				userUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
					{
						Model:    emptyDocument,
						LinkOnly: true,
					},
					{
						Model: models.UserUpload{
							DeletedAt: deletedAt,
						},
					},
				}, nil)

				if !markAsDeleted {
					emptyDocument.UserUploads = append(emptyDocument.UserUploads, userUpload)
				}
			}
		}

		if hasFullFiles {
			for i := 0; i < 2; i++ {
				userUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
					{
						Model:    fullDocument,
						LinkOnly: true,
					},
				}, nil)

				fullDocument.UserUploads = append(fullDocument.UserUploads, userUpload)
			}
		}

		if hasProofFiles {
			for i := 0; i < 2; i++ {
				userUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
					{
						Model:    proofOfOwnership,
						LinkOnly: true,
					},
				}, nil)

				proofOfOwnership.UserUploads = append(proofOfOwnership.UserUploads, userUpload)
			}
		}

		originalWeightTicket := models.WeightTicket{
			EmptyDocumentID:                   emptyDocument.ID,
			EmptyDocument:                     emptyDocument,
			FullDocumentID:                    fullDocument.ID,
			FullDocument:                      fullDocument,
			ProofOfTrailerOwnershipDocumentID: proofOfOwnership.ID,
			ProofOfTrailerOwnershipDocument:   proofOfOwnership,
			PPMShipmentID:                     ppmShipment.ID,
			PPMShipment:                       ppmShipment,
		}

		if overrides != nil {
			testdatagen.MergeModels(&originalWeightTicket, overrides)
		}

		verrs, err := suite.DB().ValidateAndCreate(&originalWeightTicket)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalWeightTicket.ID)

		return &originalWeightTicket
	}

	setUpFetcher := func(returnValue ...interface{}) services.WeightTicketFetcher {
		mockFetcher := &mocks.WeightTicketFetcher{}

		mockFetcher.On(
			"GetWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(returnValue...)

		return mockFetcher
	}

	suite.Run("Returns an error if the original doesn't exist", func() {
		badWeightTicket := models.WeightTicket{
			ID: uuid.Must(uuid.NewV4()),
		}

		notFoundErr := apperror.NewNotFoundError(badWeightTicket.ID, "while looking for weight ticket")

		updater := NewCustomerWeightTicketUpdater(setUpFetcher(nil, notFoundErr), &ppmShipmentUpdater)

		updatedWeightTicket, err := updater.UpdateWeightTicket(suite.AppContextForTest(), badWeightTicket, "")

		suite.Nil(updatedWeightTicket)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				notFoundErr.Error(),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		originalWeightTicket := setupForTest(nil, false, false, false)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		updater := NewCustomerWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)

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
		override := models.WeightTicket{
			EmptyWeight: models.PoundPointer(3000),
			FullWeight:  models.PoundPointer(4200),
		}

		originalWeightTicket := setupForTest(&override, true, true, false)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		updater := NewCustomerWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(3000),
			MissingEmptyWeightTicket: models.BoolPointer(true),
			FullWeight:               models.PoundPointer(4200),
			MissingFullWeightTicket:  models.BoolPointer(true),
			OwnsTrailer:              models.BoolPointer(false),
			TrailerMeetsCriteria:     models.BoolPointer(false),
			AdjustedNetWeight:        models.PoundPointer(1200),
			AllowableWeight:          models.PoundPointer(1200),
			NetWeightRemarks:         models.StringPointer("Weight has been adjusted"),
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
		suite.Equal(*desiredWeightTicket.AdjustedNetWeight, *updatedWeightTicket.AdjustedNetWeight)
		suite.Equal(*desiredWeightTicket.NetWeightRemarks, *updatedWeightTicket.NetWeightRemarks)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.SubmittedEmptyWeight)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.SubmittedFullWeight)
	})

	suite.Run("Succesfully updates when files are required", func() {
		override := models.WeightTicket{
			EmptyWeight: models.PoundPointer(3000),
			FullWeight:  models.PoundPointer(4200),
		}
		originalWeightTicket := setupForTest(&override, true, true, true)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		updater := NewCustomerWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(3000),
			MissingEmptyWeightTicket: models.BoolPointer(false),
			FullWeight:               models.PoundPointer(4200),
			MissingFullWeightTicket:  models.BoolPointer(false),
			OwnsTrailer:              models.BoolPointer(true),
			TrailerMeetsCriteria:     models.BoolPointer(true),
			AdjustedNetWeight:        models.PoundPointer(1200),
			NetWeightRemarks:         models.StringPointer("Weight has been adjusted"),
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
		suite.Equal(*desiredWeightTicket.AdjustedNetWeight, *updatedWeightTicket.AdjustedNetWeight)
		suite.Equal(*desiredWeightTicket.NetWeightRemarks, *updatedWeightTicket.NetWeightRemarks)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.SubmittedEmptyWeight)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.SubmittedFullWeight)
		suite.Equal(1, len(updatedWeightTicket.EmptyDocument.UserUploads))
		suite.Equal(2, len(updatedWeightTicket.FullDocument.UserUploads))
		suite.Equal(2, len(updatedWeightTicket.ProofOfTrailerOwnershipDocument.UserUploads))
	})

	suite.Run("Successfully updates and calls the ppmShipmentUpdater when weights are updated", func() {
		originalWeightTicket := setupForTest(nil, true, true, false)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		updater := NewOfficeWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)
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
			AdjustedNetWeight:        models.PoundPointer(1200),
			NetWeightRemarks:         models.StringPointer("Weight has been adjusted"),
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
		suite.Equal(*desiredWeightTicket.AdjustedNetWeight, *updatedWeightTicket.AdjustedNetWeight)
		suite.Equal(*desiredWeightTicket.NetWeightRemarks, *updatedWeightTicket.NetWeightRemarks)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.SubmittedEmptyWeight)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.SubmittedFullWeight)
	})

	suite.Run("Successfully updates and does not call ppmShipmentUpdater when total weight is unchanged", func() {
		override := models.WeightTicket{
			EmptyWeight: models.PoundPointer(3000),
			FullWeight:  models.PoundPointer(4200),
		}
		originalWeightTicket := setupForTest(&override, true, true, false)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		updater := NewOfficeWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(1000),
			MissingEmptyWeightTicket: models.BoolPointer(true),
			FullWeight:               models.PoundPointer(2200),
			MissingFullWeightTicket:  models.BoolPointer(true),
			OwnsTrailer:              models.BoolPointer(false),
			TrailerMeetsCriteria:     models.BoolPointer(false),
			AdjustedNetWeight:        models.PoundPointer(1200),
			NetWeightRemarks:         models.StringPointer("Weight has been adjusted"),
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
		suite.Equal(*desiredWeightTicket.AdjustedNetWeight, *updatedWeightTicket.AdjustedNetWeight)
		suite.Equal(*desiredWeightTicket.NetWeightRemarks, *updatedWeightTicket.NetWeightRemarks)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.SubmittedEmptyWeight)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.SubmittedFullWeight)
	})

	suite.Run("Successfully updates when total weight is changed - taking adjustedNetWeight into account", func() {
		override := models.WeightTicket{
			EmptyWeight:       models.PoundPointer(3000),
			FullWeight:        models.PoundPointer(4200),
			AdjustedNetWeight: models.PoundPointer(1200),
			NetWeightRemarks:  models.StringPointer("Weight has been adjusted"),
		}
		originalWeightTicket := setupForTest(&override, true, true, false)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		updater := NewOfficeWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(1000),
			MissingEmptyWeightTicket: models.BoolPointer(true),
			FullWeight:               models.PoundPointer(2200),
			MissingFullWeightTicket:  models.BoolPointer(true),
			OwnsTrailer:              models.BoolPointer(false),
			TrailerMeetsCriteria:     models.BoolPointer(false),
			AdjustedNetWeight:        models.PoundPointer(1000),
			NetWeightRemarks:         models.StringPointer("Weight has been adjusted"),
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
		suite.Equal(*desiredWeightTicket.AdjustedNetWeight, *updatedWeightTicket.AdjustedNetWeight)
		suite.Equal(*desiredWeightTicket.NetWeightRemarks, *updatedWeightTicket.NetWeightRemarks)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.SubmittedEmptyWeight)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.SubmittedFullWeight)
	})

	suite.Run("Fails to update when files are missing", func() {
		originalWeightTicket := setupForTest(nil, false, false, false)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		updater := NewCustomerWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(3000),
			MissingEmptyWeightTicket: models.BoolPointer(false),
			FullWeight:               models.PoundPointer(4200),
			MissingFullWeightTicket:  models.BoolPointer(false),
			OwnsTrailer:              models.BoolPointer(true),
			TrailerMeetsCriteria:     models.BoolPointer(true),
			AdjustedNetWeight:        models.PoundPointer(1200),
			NetWeightRemarks:         models.StringPointer("Weight has been adjusted"),
		}

		updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

		suite.Nil(updatedWeightTicket)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.Equal("Invalid input found while validating the weight ticket.", updateErr.Error())
	})

	suite.Run("hasTotalWeightChanged function", func() {
		suite.Run("should return true if there's a change in total weight based off adjusted net weight for both tickets", func() {
			//Default net weight of 4,000 - full weight of 18500 - empty weight of 14500
			oldAdjustedNetWeight := unit.Pound(3999)
			originalWeightTicket := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
				{
					Model: models.WeightTicket{
						AdjustedNetWeight: &oldAdjustedNetWeight,
					},
				},
			}, nil)

			newWeightTicket := originalWeightTicket
			newAdjustedNetWeight := unit.Pound(3000)
			newWeightTicket.AdjustedNetWeight = &newAdjustedNetWeight

			totalWeightHasChanged := hasTotalWeightChanged(originalWeightTicket, newWeightTicket)
			suite.Equal(true, totalWeightHasChanged)
		})

		suite.Run("should return true if there's a change in total weight when only one ticket has an adjusted net weight value", func() {
			//Default net weight of 4,000 - full weight of 18500 - empty weight of 14500
			originalWeightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

			newWeightTicket := originalWeightTicket
			newAdjustedNetWeight := unit.Pound(3500)
			newWeightTicket.AdjustedNetWeight = &newAdjustedNetWeight

			totalWeightHasChanged := hasTotalWeightChanged(originalWeightTicket, newWeightTicket)
			suite.Equal(true, totalWeightHasChanged)
		})

		suite.Run("should return true if there's a change in total weight based off full and empty weight", func() {
			//Default net weight of 4,000 - full weight of 18500 - empty weight of 14500
			originalWeightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

			newWeightTicket := originalWeightTicket
			newFullWeight := unit.Pound(15000)
			newEmptyWeight := unit.Pound(10000)
			newWeightTicket.FullWeight = &newFullWeight
			newWeightTicket.EmptyWeight = &newEmptyWeight

			totalWeightHasChanged := hasTotalWeightChanged(originalWeightTicket, newWeightTicket)
			suite.Equal(true, totalWeightHasChanged)
		})

		suite.Run("should return false when the total weight is the same", func() {
			//Default net weight of 4,000 - full weight of 18500 - empty weight of 14500
			originalWeightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

			newWeightTicket := originalWeightTicket
			newFullWeight := unit.Pound(16500)
			newEmptyWeight := unit.Pound(12500)
			newWeightTicket.FullWeight = &newFullWeight
			newWeightTicket.EmptyWeight = &newEmptyWeight

			totalWeightHasChanged := hasTotalWeightChanged(originalWeightTicket, newWeightTicket)
			suite.Equal(false, totalWeightHasChanged)
		})

		suite.Run("should return false when there's different values but the total weight remains the same", func() {
			//Default net weight of 4,000 - full weight of 18500 - empty weight of 14500
			originalWeightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

			newWeightTicket := originalWeightTicket
			newAdjustedNetWeight := unit.Pound(4000)
			newWeightTicket.AdjustedNetWeight = &newAdjustedNetWeight

			totalWeightHasChanged := hasTotalWeightChanged(originalWeightTicket, newWeightTicket)
			suite.Equal(false, totalWeightHasChanged)
		})
	})

	suite.Run("Status and reason related", func() {
		suite.Run("successfully", func() {

			suite.Run("changes status and reason", func() {
				originalWeightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.MilApp,
					ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
				})

				updater := NewOfficeWeightTicketUpdater(setUpFetcher(&originalWeightTicket, nil), &ppmShipmentUpdater)

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
				status := models.PPMDocumentStatusExcluded
				originalWeightTicket := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.MilApp,
					ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
				})

				updater := NewOfficeWeightTicketUpdater(setUpFetcher(&originalWeightTicket, nil), &ppmShipmentUpdater)

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
				status := models.PPMDocumentStatusExcluded
				originalWeightTicket := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.MilApp,
					ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
				})

				updater := NewOfficeWeightTicketUpdater(setUpFetcher(&originalWeightTicket, nil), &ppmShipmentUpdater)

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
				originalWeightTicket := setupForTest(nil, true, true, false)

				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.MilApp,
					ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
				})

				updater := NewCustomerWeightTicketUpdater(setUpFetcher(originalWeightTicket, nil), &ppmShipmentUpdater)

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
					AdjustedNetWeight:        models.PoundPointer(1000),
					NetWeightRemarks:         models.StringPointer("Weight has been adjusted"),
				}

				updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

				suite.Nil(updatedWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("Invalid input found while validating the weight ticket.", updateErr.Error())
			})

			suite.Run("to update status if reason is also set when approving", func() {
				status := models.PPMDocumentStatusExcluded
				originalWeightTicket := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.WeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.MilApp,
					ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
				})

				updater := NewOfficeWeightTicketUpdater(setUpFetcher(&originalWeightTicket, nil), &ppmShipmentUpdater)

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
				originalWeightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.MilApp,
					ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
				})

				updater := NewOfficeWeightTicketUpdater(setUpFetcher(&originalWeightTicket, nil), &ppmShipmentUpdater)

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
