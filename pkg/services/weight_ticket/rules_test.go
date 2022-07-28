package weightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *WeightTicketSuite) TestValidationRules() {
	suite.Run("checkID", func() {
		suite.Run("Success", func() {
			suite.Run("Create WeightTicket without an ID", func() {
				err := checkID().Validate(suite.AppContextForTest(), nil, nil)

				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update WeightTicket with matching ID", func() {
				id := uuid.Must(uuid.NewV4())

				err := checkID().Validate(suite.AppContextForTest(), &models.WeightTicket{ID: id}, &models.WeightTicket{ID: id})

				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Return an error if the IDs don't match", func() {
				err := checkID().Validate(suite.AppContextForTest(), &models.WeightTicket{ID: uuid.Must(uuid.NewV4())}, &models.WeightTicket{ID: uuid.Must(uuid.NewV4())})

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Contains(verr.Keys(), "ID")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
		})
	})

	suite.Run("checkRequiredFields", func() {
		weightTicketID := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		emptyDocID := uuid.Must(uuid.NewV4())
		fullDocID := uuid.Must(uuid.NewV4())
		proofOfOwnershipID := uuid.Must(uuid.NewV4())

		emptyUploads := models.UserUploads{}
		emptyUploads = append(emptyUploads, models.UserUpload{
			DocumentID: &emptyDocID,
		})

		fullUploads := models.UserUploads{}
		fullUploads = append(fullUploads, models.UserUpload{
			DocumentID: &fullDocID,
		})

		existingWeightTicket := &models.WeightTicket{
			ID:                                weightTicketID,
			PPMShipmentID:                     ppmShipmentID,
			EmptyDocumentID:                   emptyDocID,
			FullDocumentID:                    fullDocID,
			ProofOfTrailerOwnershipDocumentID: proofOfOwnershipID,
			EmptyDocument: models.Document{
				UserUploads: emptyUploads,
			},
			FullDocument: models.Document{
				UserUploads: fullUploads,
			},
		}

		suite.Run("Success", func() {
			suite.Run("Update WeightTicket - all fields", func() {

				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.WeightTicket{
						ID:                       weightTicketID,
						VehicleDescription:       models.StringPointer("1994 Mazda MX-5 Miata"),
						EmptyWeight:              models.PoundPointer(2500),
						MissingEmptyWeightTicket: models.BoolPointer(true),
						FullWeight:               models.PoundPointer(3300),
						MissingFullWeightTicket:  models.BoolPointer(true),
						OwnsTrailer:              models.BoolPointer(false),
						TrailerMeetsCriteria:     models.BoolPointer(false),
					},
					existingWeightTicket,
				)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Update WeightTicket - missing fields", func() {
				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.WeightTicket{
						ID: weightTicketID,
					},
					&models.WeightTicket{
						ID: weightTicketID,
					},
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 13)
					suite.Contains(verr.Keys(), "PPMShipmentID")
					suite.Contains(verr.Keys(), "EmptyDocumentID")
					suite.Contains(verr.Keys(), "FullDocumentID")
					suite.Contains(verr.Keys(), "ProofOfTrailerOwnershipDocumentID")
					suite.Contains(verr.Keys(), "VehicleDescription")
					suite.Contains(verr.Keys(), "EmptyWeight")
					suite.Contains(verr.Keys(), "MissingEmptyWeightTicket")
					suite.Contains(verr.Keys(), "FullWeight")
					suite.Contains(verr.Keys(), "MissingFullWeightTicket")
					suite.Contains(verr.Keys(), "EmptyWeightDocument")
					suite.Contains(verr.Keys(), "FullWeightDocument")
					suite.Contains(verr.Keys(), "OwnsTrailer")
					suite.Contains(verr.Keys(), "TrailerMeetsCriteria")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})

			suite.Run("Update WeightTicket - invalid weight", func() {
				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.WeightTicket{
						ID:                       weightTicketID,
						VehicleDescription:       models.StringPointer("1994 Mazda MX-5 Miata"),
						EmptyWeight:              models.PoundPointer(2500),
						MissingEmptyWeightTicket: models.BoolPointer(true),
						FullWeight:               models.PoundPointer(2400),
						MissingFullWeightTicket:  models.BoolPointer(true),
						OwnsTrailer:              models.BoolPointer(false),
						TrailerMeetsCriteria:     models.BoolPointer(false),
					},
					existingWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "FullWeight")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
			suite.Run("Update WeightTicket - documents required", func() {
				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.WeightTicket{
						ID:                       weightTicketID,
						VehicleDescription:       models.StringPointer("1994 Mazda MX-5 Miata"),
						EmptyWeight:              models.PoundPointer(2500),
						MissingEmptyWeightTicket: models.BoolPointer(false),
						FullWeight:               models.PoundPointer(3300),
						MissingFullWeightTicket:  models.BoolPointer(false),
						OwnsTrailer:              models.BoolPointer(true),
						TrailerMeetsCriteria:     models.BoolPointer(true),
					},
					existingWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "ProofOfTrailerOwnershipDocument")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}

			})
		})
	})
}
