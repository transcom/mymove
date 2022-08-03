package movingexpense

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MovingExpenseSuite) TestValidationRules() {
	suite.Run("checkID", func() {
		suite.Run("Success", func() {
			suite.Run("Create MovingExpense without an ID", func() {
				movingExpense := &models.MovingExpense{
					PPMShipmentID: uuid.Must(uuid.NewV4()),
					Document: models.Document{
						ServiceMemberID: uuid.Must(uuid.NewV4()),
					},
				}
				err := checkID().Validate(suite.AppContextForTest(), movingExpense, nil)

				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update WeightTicket with matching ID", func() {
				id := uuid.Must(uuid.NewV4())

				err := checkID().Validate(suite.AppContextForTest(), &models.MovingExpense{ID: id}, &models.MovingExpense{ID: id})

				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Return an error if the IDs don't match", func() {
				err := checkID().Validate(suite.AppContextForTest(), &models.MovingExpense{ID: uuid.Must(uuid.NewV4())}, &models.MovingExpense{ID: uuid.Must(uuid.NewV4())})

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
		movingExpenseID := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		docID := uuid.Must(uuid.NewV4())
		serviceMemberID := uuid.Must(uuid.NewV4())

		docUploads := models.UserUploads{}
		docUploads = append(docUploads, models.UserUpload{
			DocumentID: &docID,
		})

		suite.Run("Success", func() {
			suite.Run("Create MovingExpense", func() {

				err := checkCreateRequiredFields().Validate(suite.AppContextForTest(),
					&models.MovingExpense{
						ID:            movingExpenseID,
						PPMShipmentID: ppmShipmentID,
						Document: models.Document{
							ServiceMemberID: serviceMemberID,
							UserUploads:     docUploads,
						},
					},
					nil,
				)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Create MovingExpense - missing fields", func() {
				err := checkCreateRequiredFields().Validate(suite.AppContextForTest(),
					&models.MovingExpense{},
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
		movingExpenseID := uuid.Must(uuid.NewV4())
		ppmShipmentID := uuid.Must(uuid.NewV4())
		docID := uuid.Must(uuid.NewV4())
		serviceMemberID := uuid.Must(uuid.NewV4())

		docUploads := models.UserUploads{}
		docUploads = append(docUploads, models.UserUpload{
			DocumentID: &docID,
		})

		existingMovingExpense := &models.MovingExpense{
			ID:            movingExpenseID,
			PPMShipmentID: ppmShipmentID,
			DocumentID:    docID,
			Document: models.Document{
				ID:              docID,
				ServiceMemberID: serviceMemberID,
				UserUploads:     docUploads,
			},
		}
		suite.Run("Success", func() {
			suite.Run("Update MovingExpense - non-storage type", func() {
				tollsExpenseType := models.MovingExpenseReceiptTypeTolls
				updatedMovingExpense := &models.MovingExpense{
					ID:            existingMovingExpense.ID,
					PPMShipmentID: existingMovingExpense.PPMShipmentID,
					DocumentID:    existingMovingExpense.DocumentID,
					Document: models.Document{
						ID:              existingMovingExpense.DocumentID,
						ServiceMemberID: existingMovingExpense.Document.ServiceMemberID,
						UserUploads: models.UserUploads{
							{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
					MovingExpenseType: &tollsExpenseType,
					Description:       models.StringPointer("Pennsylvania Turnpike Oakmont to King of Prussia"),
					Amount:            models.CentPointer(unit.Cents(7500)),
					MissingReceipt:    models.BoolPointer(false),
					PaidWithGTCC:      models.BoolPointer(false),
				}

				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update MovingExpense - storage type", func() {
				storageExpenseType := models.MovingExpenseReceiptTypeStorage
				updatedMovingExpense := &models.MovingExpense{
					ID:            existingMovingExpense.ID,
					PPMShipmentID: existingMovingExpense.PPMShipmentID,
					DocumentID:    existingMovingExpense.DocumentID,
					Document: models.Document{
						ID:              existingMovingExpense.DocumentID,
						ServiceMemberID: existingMovingExpense.Document.ServiceMemberID,
						UserUploads: models.UserUploads{
							{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
					MovingExpenseType: &storageExpenseType,
					Description:       models.StringPointer("UHaul storage unit rental"),
					Amount:            models.CentPointer(unit.Cents(88800)),
					MissingReceipt:    models.BoolPointer(true),
					PaidWithGTCC:      models.BoolPointer(true),
					SITStartDate:      models.TimePointer(time.Now()),
					SITEndDate:        models.TimePointer(time.Now().Add(30 * time.Hour * 24)),
				}

				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update MovingExpense - approved status", func() {
				tollsExpenseType := models.MovingExpenseReceiptTypeTolls
				approvedStatus := models.PPMDocumentStatusApproved
				updatedMovingExpense := &models.MovingExpense{
					ID:            existingMovingExpense.ID,
					PPMShipmentID: existingMovingExpense.PPMShipmentID,
					DocumentID:    existingMovingExpense.DocumentID,
					Document: models.Document{
						ID:              existingMovingExpense.DocumentID,
						ServiceMemberID: existingMovingExpense.Document.ServiceMemberID,
						UserUploads: models.UserUploads{
							{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
					MovingExpenseType: &tollsExpenseType,
					Description:       models.StringPointer("Pennsylvania Turnpike Oakmont to King of Prussia"),
					Amount:            models.CentPointer(unit.Cents(7500)),
					MissingReceipt:    models.BoolPointer(false),
					PaidWithGTCC:      models.BoolPointer(false),
					Status:            &approvedStatus,
				}

				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update MovingExpense - not approved status with reason", func() {
				tollsExpenseType := models.MovingExpenseReceiptTypeTolls
				excludedStatus := models.PPMDocumentStatusExcluded
				updatedMovingExpense := &models.MovingExpense{
					ID:            existingMovingExpense.ID,
					PPMShipmentID: existingMovingExpense.PPMShipmentID,
					DocumentID:    existingMovingExpense.DocumentID,
					Document: models.Document{
						ID:              existingMovingExpense.DocumentID,
						ServiceMemberID: existingMovingExpense.Document.ServiceMemberID,
						UserUploads: models.UserUploads{
							{
								ID: uuid.Must(uuid.NewV4()),
							},
						},
					},
					MovingExpenseType: &tollsExpenseType,
					Description:       models.StringPointer("Pennsylvania Turnpike Oakmont to King of Prussia"),
					Amount:            models.CentPointer(unit.Cents(7500)),
					MissingReceipt:    models.BoolPointer(false),
					PaidWithGTCC:      models.BoolPointer(false),
					Status:            &excludedStatus,
					Reason:            models.StringPointer("Duplicate expense"),
				}

				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Update MovingExpense - missing required fields", func() {
				movingExpenseMissingUploads := &models.MovingExpense{
					ID:            existingMovingExpense.ID,
					PPMShipmentID: existingMovingExpense.PPMShipmentID,
					DocumentID:    existingMovingExpense.DocumentID,
					Document:      models.Document{},
				}
				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(),
					&models.MovingExpense{
						ID:            movingExpenseMissingUploads.ID,
						PPMShipmentID: movingExpenseMissingUploads.PPMShipmentID,
						DocumentID:    movingExpenseMissingUploads.DocumentID,
						Document:      movingExpenseMissingUploads.Document,
					},
					movingExpenseMissingUploads,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 6)
					suite.Contains(verr.Keys(), "MovingExpenseType")
					suite.Contains(verr.Keys(), "Description")
					suite.Contains(verr.Keys(), "PaidWithGTCC")
					suite.Contains(verr.Keys(), "Amount")
					suite.Contains(verr.Keys(), "MissingReceipt")
					suite.Contains(verr.Keys(), "Document")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
			suite.Run("Update WeightTicket - storage type dates missing", func() {
				storageExpenseType := models.MovingExpenseReceiptTypeStorage
				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(),
					&models.MovingExpense{
						ID:                existingMovingExpense.ID,
						PPMShipmentID:     existingMovingExpense.PPMShipmentID,
						DocumentID:        existingMovingExpense.DocumentID,
						Document:          existingMovingExpense.Document,
						MovingExpenseType: &storageExpenseType,
					},
					existingMovingExpense,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Contains(verr.Keys(), "SITStartDate")
					suite.Contains(verr.Keys(), "SITEndDate")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
			suite.Run("Update WeightTicket - storage SITStartDate not before SITEndDate", func() {
				storageExpenseType := models.MovingExpenseReceiptTypeStorage
				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(),
					&models.MovingExpense{
						ID:                existingMovingExpense.ID,
						PPMShipmentID:     existingMovingExpense.PPMShipmentID,
						DocumentID:        existingMovingExpense.DocumentID,
						Document:          existingMovingExpense.Document,
						MovingExpenseType: &storageExpenseType,
						SITStartDate:      models.TimePointer(time.Now().Add(1 * time.Hour * 24)),
						SITEndDate:        models.TimePointer(time.Now()),
					},
					existingMovingExpense,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Contains(verr.Keys(), "SITStartDate")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
			suite.Run("Update WeightTicket - unapproved status missing reason", func() {
				storageExpenseType := models.MovingExpenseReceiptTypeStorage
				rejectedStatus := models.PPMDocumentStatusRejected
				err := checkUpdateRequiredFields().Validate(suite.AppContextForTest(),
					&models.MovingExpense{
						ID:                existingMovingExpense.ID,
						PPMShipmentID:     existingMovingExpense.PPMShipmentID,
						DocumentID:        existingMovingExpense.DocumentID,
						Document:          existingMovingExpense.Document,
						MovingExpenseType: &storageExpenseType,
						SITStartDate:      models.TimePointer(time.Now().Add(1 * time.Hour * 24)),
						SITEndDate:        models.TimePointer(time.Now()),
						Status:            &rejectedStatus,
					},
					existingMovingExpense,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Contains(verr.Keys(), "Reason")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
		})
	})
}
