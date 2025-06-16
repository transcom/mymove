package sitextension

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

func (suite *SitExtensionServiceSuite) TestValidationRules() {

	suite.Run("checkDepartureDates", func() {
		shipmentId := uuid.Must(uuid.NewV4())
		today := time.Now()
		aFewDaysLater := today.Add(time.Hour * 72)

		reServiceDOASIT := models.ReService{
			Code: models.ReServiceCodeDOASIT,
		}
		reServiceDDASIT := models.ReService{
			Code: models.ReServiceCodeDDASIT,
		}
		sitWDOA := models.MTOServiceItem{
			MTOShipmentID:    &shipmentId,
			Status:           models.MTOServiceItemStatusApproved,
			SITEntryDate:     &today,
			SITDepartureDate: &aFewDaysLater,
			ReService:        reServiceDOASIT,
		}
		sitWDDA := models.MTOServiceItem{
			MTOShipmentID:    &shipmentId,
			Status:           models.MTOServiceItemStatusApproved,
			SITEntryDate:     &today,
			SITDepartureDate: &aFewDaysLater,
			ReService:        reServiceDDASIT,
		}

		suite.Run("success only DOASIT", func() {
			move := factory.BuildMove(suite.DB(), nil, nil)
			sitList := []models.MTOServiceItem{sitWDOA}
			shipment := models.MTOShipment{
				ID:                   shipmentId,
				OriginSITAuthEndDate: &today,
				MTOServiceItems:      sitList,
				MoveTaskOrder:        move,
				MoveTaskOrderID:      move.ID,
			}

			var emptySitExt models.SITDurationUpdate
			err := checkDepartureDate().Validate(suite.AppContextForTest(), emptySitExt, &shipment)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("success only DDASIT", func() {
			move := factory.BuildMove(suite.DB(), nil, nil)
			sitList := []models.MTOServiceItem{sitWDDA}
			shipment := models.MTOShipment{
				ID:                   shipmentId,
				OriginSITAuthEndDate: &today,
				MTOServiceItems:      sitList,
				MoveTaskOrder:        move,
				MoveTaskOrderID:      move.ID,
			}

			var emptySitExt models.SITDurationUpdate
			err := checkDepartureDate().Validate(suite.AppContextForTest(), emptySitExt, &shipment)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("success DDASIT and DOASIT", func() {
			sitList := []models.MTOServiceItem{sitWDOA, sitWDDA}

			move := factory.BuildMove(suite.DB(), nil, nil)
			shipment := models.MTOShipment{
				ID:                   shipmentId,
				OriginSITAuthEndDate: &today,
				MTOServiceItems:      sitList,
				MoveTaskOrder:        move,
				MoveTaskOrderID:      move.ID,
			}

			var emptySitExt models.SITDurationUpdate
			err := checkDepartureDate().Validate(suite.AppContextForTest(), emptySitExt, &shipment)
			suite.NilOrNoVerrs(err)
		})
	})

	suite.Run("checkDepartureDates", func() {
		suite.Run("failure", func() {
			shipmentId := uuid.Must(uuid.NewV4())
			today := time.Now()
			tomorrow := today.Add(time.Hour * 24)

			reService := models.ReService{
				Code: models.ReServiceCodeDDASIT,
			}
			var sitList []models.MTOServiceItem
			sit := models.MTOServiceItem{
				MTOShipmentID:    &shipmentId,
				Status:           models.MTOServiceItemStatusApproved,
				SITEntryDate:     &today,
				SITDepartureDate: &today,
				ReService:        reService,
			}
			sitList = append(sitList, sit)

			move := factory.BuildMove(suite.DB(), nil, nil)
			shipment := models.MTOShipment{
				ID:                        shipmentId,
				DestinationSITAuthEndDate: &today,
				ActualDeliveryDate:        &tomorrow,
				MTOServiceItems:           sitList,
				MoveTaskOrder:             move,
				MoveTaskOrderID:           move.ID,
			}

			var emptySitExt models.SITDurationUpdate
			err := checkDepartureDate().Validate(suite.AppContextForTest(), emptySitExt, &shipment)
			suite.Error(err)
			suite.Contains(err.Error(), "cannot be prior or equal to the SIT end date ")
		})
	})

	suite.Run("checkShipmentID", func() {
		suite.Run("success", func() {
			sit := models.SITDurationUpdate{MTOShipmentID: uuid.Must(uuid.NewV4())}
			err := checkShipmentID().Validate(suite.AppContextForTest(), sit, nil)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			var sit models.SITDurationUpdate
			err := checkShipmentID().Validate(suite.AppContextForTest(), sit, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Contains(verr.Keys(), "MTOShipmentID")
			default:
				suite.Failf("expected *validate.Errors", "%t - %v", err, err)
			}
		})
	})

	suite.Run("checkRequiredFields", func() {
		//takes an app context& sit extension
		//returns a verification error
		suite.Run("success", func() {
			shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
					LinkOnly: true,
				}, // Move status is automatically set to APPROVED
			}, nil)
			sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
				{
					Model:    shipment,
					LinkOnly: true,
				},
				{
					Model: models.SITDurationUpdate{
						RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
						Status:        models.SITExtensionStatusApproved,
						RequestedDays: 90,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})

			err := checkRequiredFields().Validate(suite.AppContextForTest(), sitExtension, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.NoVerrs(verr)
			default:
				suite.Failf("expected *validate.Errs", "%v", err)
			}
		})

		suite.Run("failure", func() {
			var sit models.SITDurationUpdate
			err := checkRequiredFields().Validate(suite.AppContextForTest(), sit, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Contains(verr.Keys(), "RequestedDays")
			default:
				suite.Failf("expected *validate.Errors", "%t - %v", err, err)
			}
		})
	})

	suite.Run("checkSITExtensionPending - Success", func() {
		// Testing: There is no new sit extension
		sit := models.SITDurationUpdate{MTOShipmentID: uuid.Must(uuid.NewV4())}
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			}, // Move status is automatically set to APPROVED
		}, nil)
		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NoError(err)
	})

	suite.Run("checkSITExtensionPending - Success after existing SIT is Approved", func() {
		// Testing: There is no new sit extension
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			}, // Move status is automatically set to APPROVED
		}, nil)

		// Approved Status SIT Extension
		// Changed Request Reason from the default
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
					Status:        models.SITExtensionStatusApproved,
					RequestedDays: 90,
				},
			},
		}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})

		sit := models.SITDurationUpdate{MTOShipmentID: uuid.Must(uuid.NewV4())}

		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NoError(err)
	})

	suite.Run("checkSITExtensionPending - Success after existing SIT is Denied", func() {
		// Testing: There is no new sit extension
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			}, // Move status is automatically set to APPROVED
		}, nil)

		// Denied SIT Extension
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
					Status:        models.SITExtensionStatusDenied,
					RequestedDays: 90,
				},
			},
		}, nil)
		sit := models.SITDurationUpdate{MTOShipmentID: uuid.Must(uuid.NewV4())}

		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NoError(err)
	})

	suite.Run("checkSITExtensionPending - Failure", func() {
		// Testing: There is a SIT extension and trying to be created
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			}, // Move status is automatically set to APPROVED
		}, nil)

		// Create SIT Extension #1 in DB
		// Change default status to Pending:
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
					Status:        models.SITExtensionStatusPending,
					RequestedDays: 90,
				},
			},
		}, nil)
		// Object we are trying to add to DB
		newSIT := models.SITDurationUpdate{MTOShipmentID: uuid.Must(uuid.NewV4()), Status: models.SITExtensionStatusPending, RequestedDays: 4}

		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), newSIT, &shipment)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("checkPrimeAvailability - Failure", func() {
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(suite.AppContextForTest(), models.SITDurationUpdate{}, nil)
		suite.NotNil(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal("Not found while looking for Prime-available Shipment", err.Error())
	})

	suite.Run("checkPrimeAvailability - Success", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			}, // Move status is automatically set to APPROVED
		}, nil)
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(suite.AppContextForTest(), models.SITDurationUpdate{}, &shipment)
		suite.NoError(err)
	})

	suite.Run("checkMinimumSITDuration - Success", func() {
		// Testing: There is a SIT duration of 5 days that can be reduced to 1 day
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			}, // Move status is automatically set to APPROVED
		}, nil)

		sitDaysAllowance := 5
		shipment.SITDaysAllowance = &sitDaysAllowance

		// New SIT Duration Update that decreases the SIT duration to 1 day
		approvedDays := -4
		sit := models.SITDurationUpdate{
			MTOShipmentID: shipment.ID,
			RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
			Status:        models.SITExtensionStatusApproved,
			RequestedDays: approvedDays,
			ApprovedDays:  &approvedDays,
		}

		err := checkMinimumSITDuration().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NoError(err)
	})

	suite.Run("checkMinimumSITDuration - Failure", func() {
		// Testing: There is a SIT duration of 5 days that cannot be reduced to 0 days
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			}, // Move status is automatically set to APPROVED
		}, nil)

		sitDaysAllowance := 5
		shipment.SITDaysAllowance = &sitDaysAllowance

		// New SIT Duration Update that decreases the SIT duration to 0 days
		approvedDays := -5
		sit := models.SITDurationUpdate{
			MTOShipmentID: shipment.ID,
			RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
			Status:        models.SITExtensionStatusApproved,
			RequestedDays: approvedDays,
			ApprovedDays:  &approvedDays,
		}

		err := checkMinimumSITDuration().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("can't reduce a SIT duration to less than one day", err.Error())
	})
}
