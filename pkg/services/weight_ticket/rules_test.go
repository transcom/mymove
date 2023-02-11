package weightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *WeightTicketSuite) TestValidationRules() {
	// setup some shared data
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
						AdjustedNetWeight:        models.PoundPointer(800),
						NetWeightRemarks:         models.StringPointer("Net weight has been adjusted"),
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
					suite.Equal(len(verr.Keys()), 15)
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
					suite.Contains(verr.Keys(), "AdjustedNetWeight")
					suite.Contains(verr.Keys(), "NetWeightRemarks")
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
						AdjustedNetWeight:        models.PoundPointer(3000),
						NetWeightRemarks:         models.StringPointer("Weight was adjusted"),
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
				err := verifyProofOfTrailerOwnershipDocument().Validate(suite.AppContextForTest(),
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

	suite.Run("verifyReasonAndStatusAreConstant", func() {
		docApprovedStatus := models.PPMDocumentStatusApproved
		docRejectedStatus := models.PPMDocumentStatusRejected

		suite.Run("Success", func() {
			constantWeightTicketTestCases := map[string]struct {
				newWeightTicket models.WeightTicket
				oldWeightTicket models.WeightTicket
			}{
				"Status is nil for both": {
					newWeightTicket: models.WeightTicket{Status: nil},
					oldWeightTicket: models.WeightTicket{Status: nil},
				},
				"Status is rejected for both": {
					newWeightTicket: models.WeightTicket{Status: &docRejectedStatus},
					oldWeightTicket: models.WeightTicket{Status: &docRejectedStatus},
				},
				"Reason is nil for both": {
					newWeightTicket: models.WeightTicket{Reason: nil},
					oldWeightTicket: models.WeightTicket{Reason: nil},
				},
				"Reason is filled for both": {
					newWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
					oldWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document"),
					},
				},
			}

			for name, constantWeightTicket := range constantWeightTicketTestCases {
				name := name
				constantWeightTicket := constantWeightTicket

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&constantWeightTicket.newWeightTicket,
						&constantWeightTicket.oldWeightTicket,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedWeightTicketTestCases := map[string]struct {
				newWeightTicket  models.WeightTicket
				oldWeightTicket  models.WeightTicket
				expectedErrorKey string
				expectedErrorMsg string
			}{
				"Status changed from nil to Approved": {
					newWeightTicket:  models.WeightTicket{Status: nil},
					oldWeightTicket:  models.WeightTicket{Status: &docApprovedStatus},
					expectedErrorKey: "Status",
					expectedErrorMsg: "status cannot be modified",
				},
				"Status changed from Rejected to nil": {
					newWeightTicket:  models.WeightTicket{Status: &docRejectedStatus},
					oldWeightTicket:  models.WeightTicket{Status: nil},
					expectedErrorKey: "Status",
					expectedErrorMsg: "status cannot be modified",
				},
				"Status is changed from Approved to Rejected": {
					newWeightTicket:  models.WeightTicket{Status: &docRejectedStatus},
					oldWeightTicket:  models.WeightTicket{Status: &docApprovedStatus},
					expectedErrorKey: "Status",
					expectedErrorMsg: "status cannot be modified",
				},
				"Reason is changed from nil to something": {
					newWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					oldWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed from something to nil": {
					newWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
				"Reason is changed": {
					newWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("bad document!"),
					},
					oldWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer("document is ok!"),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason cannot be modified",
				},
			}

			for name, changedWeightTicketTestCase := range changedWeightTicketTestCases {
				name := name
				changedWeightTicketTestCase := changedWeightTicketTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreConstant().Validate(
						suite.AppContextForTest(),
						&changedWeightTicketTestCase.newWeightTicket,
						&changedWeightTicketTestCase.oldWeightTicket,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedWeightTicketTestCase.expectedErrorKey)

					suite.Contains(
						verrs.Get(changedWeightTicketTestCase.expectedErrorKey),
						changedWeightTicketTestCase.expectedErrorMsg,
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
			validWeightTicketTestCases := map[string]models.WeightTicket{
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

			for name, validWeightTicket := range validWeightTicketTestCases {
				name := name
				validWeightTicket := validWeightTicket

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&validWeightTicket,
						nil,
					)

					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("Failure", func() {
			changedWeightTicketTestCases := map[string]struct {
				newWeightTicket  models.WeightTicket
				expectedErrorKey string
				expectedErrorMsg string
			}{
				"Reason exists without a status": {
					newWeightTicket: models.WeightTicket{
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason should not be set if the status is not set",
				},
				"Status is Approved and a blank reason is provided": {
					newWeightTicket: models.WeightTicket{
						Status: &docApprovedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Approved and a reason is provided": {
					newWeightTicket: models.WeightTicket{
						Status: &docApprovedStatus,
						Reason: models.StringPointer("interesting document..."),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason must not be set if the status is Approved",
				},
				"Status is Excluded and reason is nil": {
					newWeightTicket: models.WeightTicket{
						Status: &docExcludedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Excluded and reason is blank": {
					newWeightTicket: models.WeightTicket{
						Status: &docExcludedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is nil": {
					newWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: nil,
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
				"Status is Rejected and reason is blank": {
					newWeightTicket: models.WeightTicket{
						Status: &docRejectedStatus,
						Reason: models.StringPointer(""),
					},
					expectedErrorKey: "Reason",
					expectedErrorMsg: "reason is mandatory if the status is Excluded or Rejected",
				},
			}

			for name, changedWeightTicketTestCase := range changedWeightTicketTestCases {
				name := name
				changedWeightTicketTestCase := changedWeightTicketTestCase

				suite.Run(name, func() {
					err := verifyReasonAndStatusAreValid().Validate(
						suite.AppContextForTest(),
						&changedWeightTicketTestCase.newWeightTicket,
						nil,
					)

					suite.Error(err)

					suite.IsType(&validate.Errors{}, err)
					verrs := err.(*validate.Errors)

					suite.Len(verrs.Errors, 1)

					suite.Contains(verrs.Keys(), changedWeightTicketTestCase.expectedErrorKey)

					suite.Contains(verrs.Get(
						changedWeightTicketTestCase.expectedErrorKey),
						changedWeightTicketTestCase.expectedErrorMsg,
					)
				})
			}
		})
	})
}
