package progearweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ProgearWeightTicketSuite) TestValidationRules() {
	// setup some shared data
	progearWeightTicketID := uuid.Must(uuid.NewV4())
	ppmShipmentID := uuid.Must(uuid.NewV4())
	docID := uuid.Must(uuid.NewV4())

	uploads := models.UserUploads{}
	uploads = append(uploads, models.UserUpload{
		DocumentID: &docID,
	})

	existingProgearWeightTicket := &models.ProgearWeightTicket{
		ID:            progearWeightTicketID,
		PPMShipmentID: ppmShipmentID,
		DocumentID:    docID,
		Document: models.Document{
			UserUploads: uploads,
		},
	}

	suite.Run("checkID", func() {
		suite.Run("Success", func() {
			suite.Run("Create ProgearWeightTicket without an ID", func() {
				err := checkID().Validate(suite.AppContextForTest(), nil, nil)

				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update ProgearWeightTicket with matching ID", func() {
				id := uuid.Must(uuid.NewV4())

				err := checkID().Validate(suite.AppContextForTest(), &models.ProgearWeightTicket{ID: id}, &models.ProgearWeightTicket{ID: id})

				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Return an error if the IDs don't match", func() {
				err := checkID().Validate(suite.AppContextForTest(), &models.ProgearWeightTicket{ID: uuid.Must(uuid.NewV4())}, &models.ProgearWeightTicket{ID: uuid.Must(uuid.NewV4())})

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
		suite.Run("Success", func() {
			suite.Run("Update ProgearWeightTicket - all fields", func() {

				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
					},
					existingProgearWeightTicket,
				)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Update ProgearWeightTicket - missing fields", func() {
				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID: progearWeightTicketID,
					},
					&models.ProgearWeightTicket{
						ID: progearWeightTicketID,
					},
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 5)
					suite.Contains(verr.Keys(), "Document")
					suite.Contains(verr.Keys(), "Description")
					suite.Contains(verr.Keys(), "Weight")
					suite.Contains(verr.Keys(), "HasWeightTickets")
					suite.Contains(verr.Keys(), "BelongsToSelf")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})

			suite.Run("Update ProgearWeightTicket - invalid weight", func() {
				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           nil,
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
					},
					existingProgearWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "Weight")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})

			docLessProgearWeightTicket := &models.ProgearWeightTicket{
				ID:            progearWeightTicketID,
				PPMShipmentID: ppmShipmentID,
			}
			suite.Run("Update ProgearWeightTicket - documents required", func() {
				err := checkRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
					},
					docLessProgearWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "Document")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}

			})
		})
	})

	suite.Run("verifyReasonAndStatusAreConstant", func() {
		suite.Run("Success", func() {
			err := verifyReasonAndStatusAreConstant().Validate(suite.AppContextForTest(),
				&models.ProgearWeightTicket{
					ID:               progearWeightTicketID,
					Description:      models.StringPointer("self progear"),
					Weight:           models.PoundPointer(2500),
					HasWeightTickets: models.BoolPointer(true),
					BelongsToSelf:    models.BoolPointer(true),
				},
				existingProgearWeightTicket,
			)

			suite.NilOrNoVerrs(err)
		})

		suite.Run("Failure", func() {
			status := models.PPMDocumentStatusRejected
			err := verifyReasonAndStatusAreConstant().Validate(suite.AppContextForTest(),
				&models.ProgearWeightTicket{
					ID:               progearWeightTicketID,
					Description:      models.StringPointer("self progear"),
					Weight:           models.PoundPointer(2500),
					HasWeightTickets: models.BoolPointer(true),
					BelongsToSelf:    models.BoolPointer(true),
					Status:           &status,
					Reason:           models.StringPointer("bad data"),
				},
				existingProgearWeightTicket,
			)

			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Equal(len(verr.Keys()), 2)
				suite.Contains(verr.Keys(), "Reason")
				suite.Contains(verr.Keys(), "Status")
			default:
				suite.Failf("expected *validate.Errors", "%t - %v", err, err)
			}
		})
	})

	suite.Run("verifyReasonAndStatusAreValid", func() {
		suite.Run("Success", func() {
			suite.Run("no status or reason", func() {
				err := verifyReasonAndStatusAreValid().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
					},
					existingProgearWeightTicket,
				)
				suite.NilOrNoVerrs(err)
			})

			suite.Run("excluded with reason", func() {
				status := models.PPMDocumentStatusExcluded
				err := verifyReasonAndStatusAreValid().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
						Status:           &status,
						Reason:           models.StringPointer("bad data"),
					},
					existingProgearWeightTicket,
				)
				suite.NilOrNoVerrs(err)
			})

			suite.Run("approved with no reason", func() {
				status := models.PPMDocumentStatusApproved
				err := verifyReasonAndStatusAreValid().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
						Status:           &status,
						Reason:           nil,
					},
					existingProgearWeightTicket,
				)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Reason cannot be provided without status", func() {
				err := verifyReasonAndStatusAreValid().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
						Reason:           models.StringPointer("reason without status"),
					},
					existingProgearWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "Reason")
					suite.Equal("reason should be empty", err.Error())
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})

			suite.Run("Reason must be empty when status is approved", func() {
				status := models.PPMDocumentStatusApproved
				err := verifyReasonAndStatusAreValid().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
						Status:           &status,
						Reason:           models.StringPointer("bad data"),
					},
					existingProgearWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "Reason")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})

			suite.Run("Reason must be populated when status is excluded", func() {
				status := models.PPMDocumentStatusExcluded
				err := verifyReasonAndStatusAreValid().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
						Status:           &status,
						Reason:           models.StringPointer(""),
					},
					existingProgearWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "Reason")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}

			})

			suite.Run("Reason must be populated when status is rejected", func() {
				status := models.PPMDocumentStatusRejected
				err := verifyReasonAndStatusAreValid().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:               progearWeightTicketID,
						Description:      models.StringPointer("self progear"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
						BelongsToSelf:    models.BoolPointer(true),
						Status:           &status,
						Reason:           models.StringPointer(""),
					},
					existingProgearWeightTicket,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 1)
					suite.Contains(verr.Keys(), "Reason")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
		})
	})
}
