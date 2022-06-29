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

				suite.Error(err)
				suite.IsType(&validate.Errors{}, err)
				suite.Contains(err.Error(), "new WeightTicket ID must match original WeightTicket ID")
			})
		})
	})

	suite.Run("checkRequiredFields", func() {
		weightTicketID := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		emptyDocID := uuid.Must(uuid.NewV4())
		fullDocID := uuid.Must(uuid.NewV4())
		proofOfOwnershipID := uuid.Must(uuid.NewV4())

		existingWeightTicket := &models.WeightTicket{
			ID:                                weightTicketID,
			PPMShipmentID:                     ppmShipmentID,
			EmptyDocumentID:                   emptyDocID,
			FullDocumentID:                    fullDocID,
			ProofOfTrailerOwnershipDocumentID: proofOfOwnershipID,
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

			// TODO: Add test for required docs?
		})

		suite.Run("Failure", func() {
			suite.Run("Update WeightTicket - missing fields", func() {
				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.WeightTicket{
						ID:                       weightTicketID,
						VehicleDescription:       models.StringPointer("1994 Mazda MX-5 Miata"),
						EmptyWeight:              models.PoundPointer(2500),
						MissingEmptyWeightTicket: models.BoolPointer(true),
						FullWeight:               models.PoundPointer(3300),
						MissingFullWeightTicket:  models.BoolPointer(true),
					},
					existingWeightTicket,
				)
				suite.Error(err)
			})
		})
	})
}
