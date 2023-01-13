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

	suite.Run("checkBaseRequiredFields", func() {
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

				err := checkBaseRequiredFields().Validate(suite.AppContextForTest(),
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
				err := checkBaseRequiredFields().Validate(suite.AppContextForTest(),
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

	suite.Run("checkAdditionalRequiredFields", func() {
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

				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
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

				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
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

				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
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

				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(), updatedMovingExpense, existingMovingExpense)
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
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
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
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
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
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
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
		})
	})

	suite.Run("verifyReasonAndStatusAreConstant", func() {
		docApprovedStatus := models.PPMDocumentStatusApproved
		docRejectedStatus := models.PPMDocumentStatusRejected

		suite.Run("Success", func() {
			constantMovingExpensesTestCases := map[string]struct {
				newMovingExpense models.MovingExpense
				oldMovingExpense models.MovingExpense
			}{
				"Status is nil for both": {
					newMovingExpense: models.MovingExpense{Status: nil},
					oldMovingExpense: models.MovingExpense{Status: nil},
				},
				"Status is rejected for both": {
					newMovingExpense: models.MovingExpense{Status: &docRejectedStatus},
					oldMovingExpense: models.MovingExpense{Status: &docRejectedStatus},
				},
				"Reason is nil for both": {
					newMovingExpense: models.MovingExpense{Reason: nil},
					oldMovingExpense: models.MovingExpense{Reason: nil},
				},
				"Reason is filled for both": {
					newMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
					oldMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
				},
			}

			for name, constantMovingExpenses := range constantMovingExpensesTestCases {
				name := name
				constantMovingExpenses := constantMovingExpenses

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&constantMovingExpenses.newMovingExpense,
						&constantMovingExpenses.oldMovingExpense,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedMovingExpenseTestCases := map[string]struct {
				newMovingExpense models.MovingExpense
				oldMovingExpense models.MovingExpense
				expectedErrorKey string
				expectedErrorMsg string
			}{
				"Status changed from nil to Approved": {
					newMovingExpense: models.MovingExpense{Status: nil},
					oldMovingExpense: models.MovingExpense{Status: &docApprovedStatus},
					expectedErrorKey: "Status",
					expectedErrorMsg: "status cannot be modified",
				},
				"Status changed from Rejected to nil": {
					newMovingExpense: models.MovingExpense{Status: &docRejectedStatus},
					oldMovingExpense: models.MovingExpense{Status: nil},
					expectedErrorKey: "Status",
					expectedErrorMsg: "status cannot be modified",
				},
				"Status is changed from Approved to Rejected": {
					newMovingExpense: models.MovingExpense{Status: &docRejectedStatus},
					oldMovingExpense: models.MovingExpense{Status: &docApprovedStatus},
					expectedErrorKey: "Status",
					expectedErrorMsg: "status cannot be modified",
				},
				"Reason is changed from nil to something": {
					newMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					oldMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed from something to nil": {
					newMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed": {
					newMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
			}

			for name, changedMovingExpensesTestCase := range changedMovingExpenseTestCases {
				name := name
				changedMovingExpensesTestCase := changedMovingExpensesTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&changedMovingExpensesTestCase.newMovingExpense,
						&changedMovingExpensesTestCase.oldMovingExpense,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedMovingExpensesTestCase.expectedErrorKey)

					suite.Contains(
						verrs.Get(changedMovingExpensesTestCase.expectedErrorKey),
						changedMovingExpensesTestCase.expectedErrorMsg,
					)
				})
			}
		})
	})

	suite.Run("verifyReasonAndStatusAreValid", func() {
		docApprovedStatus := models.PPMDocumentStatusApproved
		docExcludedStatus := models.PPMDocumentStatusExcluded
		docRejectedStatus := models.PPMDocumentStatusRejected

		suite.Run("Success", func() {
			validMovingExpenseTestCases := map[string]models.MovingExpense{
				"Status is Approved with a nil reason": {
					Status: &docApprovedStatus,
					Reason: nil,
				},
				"Status is Excluded with a reason": {
					Status: &docExcludedStatus,
					Reason: models.StringPointer("not a valid expense."),
				},
				"Status is Rejected with a reason": {
					Status: &docRejectedStatus,
					Reason: models.StringPointer("bad document!"),
				},
			}

			for name, validMovingExpense := range validMovingExpenseTestCases {
				name := name
				validMovingExpense := validMovingExpense

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&validMovingExpense,
						nil,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedMovingExpenseTestCases := map[string]struct {
				newMovingExpense models.MovingExpense
				expectedErrorKey string
				expectedErrorMsg string
			}{
				"Reason exists without a status": {
					newMovingExpense: models.MovingExpense{
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason should not be set if the status is not set",
				},
				"Status is Approved and a blank reason is provided": {
					newMovingExpense: models.MovingExpense{
						Status: &docApprovedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Approved and a reason is provided": {
					newMovingExpense: models.MovingExpense{
						Status: &docApprovedStatus,
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Excluded and reason is nil": {
					newMovingExpense: models.MovingExpense{
						Status: &docExcludedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Excluded and reason is blank": {
					newMovingExpense: models.MovingExpense{
						Status: &docExcludedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is nil": {
					newMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is blank": {
					newMovingExpense: models.MovingExpense{
						Status: &docRejectedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
			}

			for name, changedMovingExpensesTestCase := range changedMovingExpenseTestCases {
				name := name
				changedMovingExpensesTestCase := changedMovingExpensesTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&changedMovingExpensesTestCase.newMovingExpense,
						nil,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedMovingExpensesTestCase.expectedErrorKey)

					suite.Contains(verrs.Get(
						changedMovingExpensesTestCase.expectedErrorKey),
						changedMovingExpensesTestCase.expectedErrorMsg,
					)
				})
			}
		})
	})
}
