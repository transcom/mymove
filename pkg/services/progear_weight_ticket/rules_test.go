package progearweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ProgearWeightTicketSuite) TestValidationRules() {
	suite.Run("checkID", func() {
		suite.Run("Success", func() {
			suite.Run("Create ProgearWeightTicket without an ID", func() {
				progearWeightTicket := &models.ProgearWeightTicket{
					PPMShipmentID: uuid.Must(uuid.NewV4()),
					Document: models.Document{
						ServiceMemberID: uuid.Must(uuid.NewV4()),
					},
				}
				err := checkID().Validate(suite.AppContextForTest(), progearWeightTicket, nil)

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

	suite.Run("checkCreateRequiredFields", func() {
		progearWeightTicketID := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		docID := uuid.Must(uuid.NewV4())
		serviceMemberID := uuid.Must(uuid.NewV4())

		documentUploads := models.UserUploads{}
		documentUploads = append(documentUploads, models.UserUpload{
			DocumentID: &docID,
		})

		suite.Run("Success", func() {
			suite.Run("Create ProgearWeightTIcket", func() {

				err := checkCreateRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:            progearWeightTicketID,
						PPMShipmentID: ppmShipmentID,
						Document: models.Document{
							ServiceMemberID: serviceMemberID,
							UserUploads:     documentUploads,
						},
					},
					nil,
				)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Create ProgearWeightTicket - missing fields", func() {
				err := checkCreateRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{},
					nil,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 2)
					suite.Contains(verr.Keys(), "PPMShipmentID")
					suite.Contains(verr.Keys(), "ServiceMemberID")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
		})
	})

	suite.Run("checkUpdateRequiredFields", func() {
		progearWeightTicketID := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		documentID := uuid.Must(uuid.NewV4())
		serviceMemberID := uuid.Must(uuid.NewV4())

		documentUploads := models.UserUploads{}
		documentUploads = append(documentUploads, models.UserUpload{
			DocumentID: &documentID,
		})

		existingProgearWeightTicket := &models.ProgearWeightTicket{
			ID:            progearWeightTicketID,
			PPMShipmentID: ppmShipmentID,
			DocumentID:    documentID,
			Document: models.Document{
				ID:              documentID,
				ServiceMemberID: serviceMemberID,
				UserUploads:     documentUploads,
			},
		}
		suite.Run("Success", func() {
			suite.Run("Update ProgearWeightTicket - Description", func() {
				updatedProgearWeightTicket := &models.ProgearWeightTicket{
					ID:            existingProgearWeightTicket.ID,
					PPMShipmentID: existingProgearWeightTicket.PPMShipmentID,
					DocumentID:    existingProgearWeightTicket.DocumentID,
					Document: models.Document{
						ID:              existingProgearWeightTicket.DocumentID,
						ServiceMemberID: existingProgearWeightTicket.Document.ServiceMemberID,
						UserUploads: models.UserUploads{
							{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
					Description: models.StringPointer("Self Progear"),
				}

				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(), updatedProgearWeightTicket, existingProgearWeightTicket)
				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update ProgearWeightTicket - approved status", func() {
				approvedStatus := models.PPMDocumentStatusApproved
				updatedProgearWeightTicket := &models.ProgearWeightTicket{
					ID:            existingProgearWeightTicket.ID,
					PPMShipmentID: existingProgearWeightTicket.PPMShipmentID,
					DocumentID:    existingProgearWeightTicket.DocumentID,
					Document: models.Document{
						ID:              existingProgearWeightTicket.DocumentID,
						ServiceMemberID: existingProgearWeightTicket.Document.ServiceMemberID,
						UserUploads: models.UserUploads{
							{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
					BelongsToSelf:    models.BoolPointer(true),
					Description:      models.StringPointer("Self Progear"),
					HasWeightTickets: models.BoolPointer(true),
					Weight:           models.PoundPointer(unit.Pound(100)),
					Status:           &approvedStatus,
				}

				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(), updatedProgearWeightTicket, existingProgearWeightTicket)
				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update ProgearWeightTicket - not approved status with reason", func() {
				rejectedStatus := models.PPMDocumentStatusRejected
				updatedProgearWeightTicket := &models.ProgearWeightTicket{
					ID:            existingProgearWeightTicket.ID,
					PPMShipmentID: existingProgearWeightTicket.PPMShipmentID,
					DocumentID:    existingProgearWeightTicket.DocumentID,
					Document: models.Document{
						ID:              existingProgearWeightTicket.DocumentID,
						ServiceMemberID: existingProgearWeightTicket.Document.ServiceMemberID,
						UserUploads: models.UserUploads{
							{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
					Description:      models.StringPointer("Self Progear"),
					HasWeightTickets: models.BoolPointer(false),
					Weight:           models.PoundPointer(unit.Pound(100)),
					Status:           &rejectedStatus,
					Reason:           models.StringPointer("Missing a progear weight ticket"),
				}

				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(), updatedProgearWeightTicket, existingProgearWeightTicket)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Update ProgearWeightTicket - missing required fields", func() {
				progearWeightTicketMissingUploads := &models.ProgearWeightTicket{
					ID:            existingProgearWeightTicket.ID,
					PPMShipmentID: existingProgearWeightTicket.PPMShipmentID,
					DocumentID:    existingProgearWeightTicket.DocumentID,
					//Document:      models.Document{},
				}
				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{
						ID:            progearWeightTicketMissingUploads.ID,
						PPMShipmentID: progearWeightTicketMissingUploads.PPMShipmentID,
						DocumentID:    progearWeightTicketMissingUploads.DocumentID,
						Document:      progearWeightTicketMissingUploads.Document,
					},
					progearWeightTicketMissingUploads,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 6)
					suite.Contains(verr.Keys(), "Description")
					suite.Contains(verr.Keys(), "HasWeightTickets")
					suite.Contains(verr.Keys(), "Weight")
					suite.Contains(verr.Keys(), "Document")
					suite.Contains(verr.Keys(), "Status")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
		})
	})
}
