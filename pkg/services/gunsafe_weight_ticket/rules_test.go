package gunsafeweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GunSafeWeightTicketSuite) TestValidationRules() {
	// setup some shared data
	gunSafeWeightTicketID := uuid.Must(uuid.NewV4())
	ppmShipmentID := uuid.Must(uuid.NewV4())
	docID := uuid.Must(uuid.NewV4())

	uploads := models.UserUploads{}
	uploads = append(uploads, models.UserUpload{
		DocumentID: &docID,
	})

	existingGunSafeWeightTicket := &models.GunSafeWeightTicket{
		ID:            gunSafeWeightTicketID,
		PPMShipmentID: ppmShipmentID,
		DocumentID:    docID,
		Document: models.Document{
			UserUploads: uploads,
		},
	}

	suite.Run("checkID", func() {
		suite.Run("Success", func() {
			suite.Run("Create GunSafeWeightTicket without an ID", func() {
				err := checkID().Validate(suite.AppContextForTest(), nil, nil)

				suite.NilOrNoVerrs(err)
			})

			suite.Run("Update GunSafeWeightTicket with matching ID", func() {
				id := uuid.Must(uuid.NewV4())

				err := checkID().Validate(suite.AppContextForTest(), &models.GunSafeWeightTicket{ID: id}, &models.GunSafeWeightTicket{ID: id})

				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Return an error if the IDs don't match", func() {
				err := checkID().Validate(suite.AppContextForTest(), &models.GunSafeWeightTicket{ID: uuid.Must(uuid.NewV4())}, &models.GunSafeWeightTicket{ID: uuid.Must(uuid.NewV4())})

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
			suite.Run("Create GunSafeWeightTIcket", func() {

				err := checkBaseRequiredFields().Validate(suite.AppContextForTest(),
					&models.GunSafeWeightTicket{
						ID:            gunSafeWeightTicketID,
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
			suite.Run("Create GunSafeWeightTicket - missing fields", func() {
				err := checkBaseRequiredFields().Validate(suite.AppContextForTest(),
					&models.GunSafeWeightTicket{},
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
			suite.Run("Update GunSafeWeightTicket - all fields", func() {

				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
					&models.GunSafeWeightTicket{
						ID:               gunSafeWeightTicketID,
						Description:      models.StringPointer("self gunSafe"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
					},
					existingGunSafeWeightTicket,
				)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			suite.Run("Update GunSafeWeightTicket - missing fields", func() {
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
					&models.GunSafeWeightTicket{
						ID: gunSafeWeightTicketID,
					},
					&models.GunSafeWeightTicket{
						ID: gunSafeWeightTicketID,
					},
				)

				switch verr := err.(type) {
				case *validate.Errors:
					suite.True(verr.HasAny())
					suite.Equal(len(verr.Keys()), 4)
					suite.Contains(verr.Keys(), "Document")
					suite.Contains(verr.Keys(), "Description")
					suite.Contains(verr.Keys(), "Weight")
					suite.Contains(verr.Keys(), "HasWeightTickets")
				default:
					suite.Failf("expected *validate.Errors", "%t - %v", err, err)
				}
			})

			suite.Run("Update GunSafeWeightTicket - invalid weight", func() {
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
					&models.GunSafeWeightTicket{
						ID:               gunSafeWeightTicketID,
						Description:      models.StringPointer("self gunSafe"),
						Weight:           nil,
						HasWeightTickets: models.BoolPointer(true),
					},
					existingGunSafeWeightTicket,
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

			docLessGunSafeWeightTicket := &models.GunSafeWeightTicket{
				ID:            gunSafeWeightTicketID,
				PPMShipmentID: ppmShipmentID,
			}
			suite.Run("Update GunSafeWeightTicket - documents required", func() {
				err := checkAdditionalRequiredFields().Validate(suite.AppContextForTest(),
					&models.GunSafeWeightTicket{
						ID:               gunSafeWeightTicketID,
						Description:      models.StringPointer("self gunSafe"),
						Weight:           models.PoundPointer(2500),
						HasWeightTickets: models.BoolPointer(true),
					},
					docLessGunSafeWeightTicket,
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
			constantGunSafeWeightTicketTestCases := map[string]struct {
				newGunSafeWeightTicket models.GunSafeWeightTicket
				oldGunSafeWeightTicket models.GunSafeWeightTicket
			}{
				"Status is nil for both": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{Status: nil},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{Status: nil},
				},
				"Status is rejected for both": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{Status: &docRejectedStatus},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{Status: &docRejectedStatus},
				},
				"Reason is nil for both": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{Reason: nil},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{Reason: nil},
				},
				"Reason is filled for both": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
				},
			}

			for name, constantGunSafeWeightTickets := range constantGunSafeWeightTicketTestCases {
				name := name
				constantGunSafeWeightTickets := constantGunSafeWeightTickets

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&constantGunSafeWeightTickets.newGunSafeWeightTicket,
						&constantGunSafeWeightTickets.oldGunSafeWeightTicket,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedGunSafeWeightTicketTestCases := map[string]struct {
				newGunSafeWeightTicket models.GunSafeWeightTicket
				oldGunSafeWeightTicket models.GunSafeWeightTicket
				expectedErrorKey       string
				expectedErrorMsg       string
			}{
				"Status changed from nil to Approved": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{Status: nil},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{Status: &docApprovedStatus},
					expectedErrorKey:       "Status",
					expectedErrorMsg:       "status cannot be modified",
				},
				"Status changed from Rejected to nil": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{Status: &docRejectedStatus},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{Status: nil},
					expectedErrorKey:       "Status",
					expectedErrorMsg:       "status cannot be modified",
				},
				"Status is changed from Approved to Rejected": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{Status: &docRejectedStatus},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{Status: &docApprovedStatus},
					expectedErrorKey:       "Status",
					expectedErrorMsg:       "status cannot be modified",
				},
				"Reason is changed from nil to something": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed from something to nil": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
			}

			for name, changedGunSafeWeightTicketTestCase := range changedGunSafeWeightTicketTestCases {
				name := name
				changedGunSafeWeightTicketTestCase := changedGunSafeWeightTicketTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&changedGunSafeWeightTicketTestCase.newGunSafeWeightTicket,
						&changedGunSafeWeightTicketTestCase.oldGunSafeWeightTicket,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedGunSafeWeightTicketTestCase.expectedErrorKey)

					suite.Contains(
						verrs.Get(changedGunSafeWeightTicketTestCase.expectedErrorKey),
						changedGunSafeWeightTicketTestCase.expectedErrorMsg,
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
			validGunSafeWeightTicketTestCases := map[string]models.GunSafeWeightTicket{
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

			for name, validGunSafeWeightTicket := range validGunSafeWeightTicketTestCases {
				name := name
				validGunSafeWeightTicket := validGunSafeWeightTicket

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&validGunSafeWeightTicket,
						nil,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedGunSafeWeightTicketTestCases := map[string]struct {
				newGunSafeWeightTicket models.GunSafeWeightTicket
				expectedErrorKey       string
				expectedErrorMsg       string
			}{
				"Reason exists without a status": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason should not be set if the status is not set",
				},
				"Status is Approved and a blank reason is provided": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docApprovedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Approved and a reason is provided": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docApprovedStatus,
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Excluded and reason is nil": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docExcludedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Excluded and reason is blank": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docExcludedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is nil": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is blank": {
					newGunSafeWeightTicket: models.GunSafeWeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
			}

			for name, changedGunSafeWeightTicketTestCase := range changedGunSafeWeightTicketTestCases {
				name := name
				changedGunSafeWeightTicketTestCase := changedGunSafeWeightTicketTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&changedGunSafeWeightTicketTestCase.newGunSafeWeightTicket,
						nil,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedGunSafeWeightTicketTestCase.expectedErrorKey)

					suite.Contains(verrs.Get(
						changedGunSafeWeightTicketTestCase.expectedErrorKey),
						changedGunSafeWeightTicketTestCase.expectedErrorMsg,
					)
				})
			}
		})
	})
}
