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

	suite.Run("checkBaseRequiredFields", func() {
		docID := uuid.Must(uuid.NewV4())
		serviceMemberID := uuid.Must(uuid.NewV4())

		documentUploads := models.UserUploads{}
		documentUploads = append(documentUploads, models.UserUpload{
			DocumentID: &docID,
		})

		suite.Run("Success", func() {
			suite.Run("Create ProgearWeightTIcket", func() {

				err := checkBaseRequiredFields().Validate(suite.AppContextForTest(),
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
				err := checkBaseRequiredFields().Validate(suite.AppContextForTest(),
					&models.ProgearWeightTicket{},
					nil,
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(2, len(verr.Keys()))
					suite.Contains(verr.Keys(), "PPMShipmentID")
					suite.Contains(verr.Keys(), "ServiceMemberID")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})
		})
	})

	suite.Run("checkAdditionalRequiredFields", func() {
		suite.Run("Success", func() {
			suite.Run("Update ProgearWeightTicket - all fields", func() {

				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
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
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
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
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
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
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
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
		docApprovedStatus := models.PPMDocumentStatusApproved
		docRejectedStatus := models.PPMDocumentStatusRejected

		suite.Run("Success", func() {
			constantProgearWeightTicketTestCases := map[string]struct {
				newProgearWeightTicket models.ProgearWeightTicket
				oldProgearWeightTicket models.ProgearWeightTicket
			}{
				"Status is nil for both": {
					newProgearWeightTicket: models.ProgearWeightTicket{Status: nil},
					oldProgearWeightTicket: models.ProgearWeightTicket{Status: nil},
				},
				"Status is rejected for both": {
					newProgearWeightTicket: models.ProgearWeightTicket{Status: &docRejectedStatus},
					oldProgearWeightTicket: models.ProgearWeightTicket{Status: &docRejectedStatus},
				},
				"Reason is nil for both": {
					newProgearWeightTicket: models.ProgearWeightTicket{Reason: nil},
					oldProgearWeightTicket: models.ProgearWeightTicket{Reason: nil},
				},
				"Reason is filled for both": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
					oldProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
				},
			}

			for name, constantProgearWeightTickets := range constantProgearWeightTicketTestCases {
				name := name
				constantProgearWeightTickets := constantProgearWeightTickets

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&constantProgearWeightTickets.newProgearWeightTicket,
						&constantProgearWeightTickets.oldProgearWeightTicket,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedProgearWeightTicketTestCases := map[string]struct {
				newProgearWeightTicket models.ProgearWeightTicket
				oldProgearWeightTicket models.ProgearWeightTicket
				expectedErrorKey       string
				expectedErrorMsg       string
			}{
				"Status changed from nil to Approved": {
					newProgearWeightTicket: models.ProgearWeightTicket{Status: nil},
					oldProgearWeightTicket: models.ProgearWeightTicket{Status: &docApprovedStatus},
					expectedErrorKey:       "Status",
					expectedErrorMsg:       "status cannot be modified",
				},
				"Status changed from Rejected to nil": {
					newProgearWeightTicket: models.ProgearWeightTicket{Status: &docRejectedStatus},
					oldProgearWeightTicket: models.ProgearWeightTicket{Status: nil},
					expectedErrorKey:       "Status",
					expectedErrorMsg:       "status cannot be modified",
				},
				"Status is changed from Approved to Rejected": {
					newProgearWeightTicket: models.ProgearWeightTicket{Status: &docRejectedStatus},
					oldProgearWeightTicket: models.ProgearWeightTicket{Status: &docApprovedStatus},
					expectedErrorKey:       "Status",
					expectedErrorMsg:       "status cannot be modified",
				},
				"Reason is changed from nil to something": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					oldProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed from something to nil": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
			}

			for name, changedProgearWeightTicketTestCase := range changedProgearWeightTicketTestCases {
				name := name
				changedProgearWeightTicketTestCase := changedProgearWeightTicketTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&changedProgearWeightTicketTestCase.newProgearWeightTicket,
						&changedProgearWeightTicketTestCase.oldProgearWeightTicket,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedProgearWeightTicketTestCase.expectedErrorKey)

					suite.Contains(
						verrs.Get(changedProgearWeightTicketTestCase.expectedErrorKey),
						changedProgearWeightTicketTestCase.expectedErrorMsg,
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
			validProgearWeightTicketTestCases := map[string]models.ProgearWeightTicket{
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

			for name, validProgearWeightTicket := range validProgearWeightTicketTestCases {
				name := name
				validProgearWeightTicket := validProgearWeightTicket

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&validProgearWeightTicket,
						nil,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedProgearWeightTicketTestCases := map[string]struct {
				newProgearWeightTicket models.ProgearWeightTicket
				expectedErrorKey       string
				expectedErrorMsg       string
			}{
				"Reason exists without a status": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason should not be set if the status is not set",
				},
				"Status is Approved and a blank reason is provided": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docApprovedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Approved and a reason is provided": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docApprovedStatus,
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Excluded and reason is nil": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docExcludedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Excluded and reason is blank": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docExcludedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is nil": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is blank": {
					newProgearWeightTicket: models.ProgearWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
			}

			for name, changedProgearWeightTicketTestCase := range changedProgearWeightTicketTestCases {
				name := name
				changedProgearWeightTicketTestCase := changedProgearWeightTicketTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&changedProgearWeightTicketTestCase.newProgearWeightTicket,
						nil,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedProgearWeightTicketTestCase.expectedErrorKey)

					suite.Contains(verrs.Get(
						changedProgearWeightTicketTestCase.expectedErrorKey),
						changedProgearWeightTicketTestCase.expectedErrorMsg,
					)
				})
			}
		})
	})
}
