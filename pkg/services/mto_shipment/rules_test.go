package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *MTOShipmentServiceSuite) TestUpdateValidations() {
	suite.Run("checkStatus", func() {
		testCases := map[models.MTOShipmentStatus]bool{
			"":                                            true,
			models.MTOShipmentStatusDraft:                 true,
			models.MTOShipmentStatusSubmitted:             true,
			"random_junk_status":                          false,
			models.MTOShipmentStatusApproved:              false,
			models.MTOShipmentStatusRejected:              false,
			models.MTOShipmentStatusCancellationRequested: false,
			models.MTOShipmentStatusCanceled:              false,
			models.MTOShipmentStatusDiversionRequested:    false,
		}
		for status, allowed := range testCases {
			suite.Run("status "+string(status), func() {
				err := checkStatus().Validate(
					suite.AppContextForTest(),
					&models.MTOShipment{Status: status},
					nil,
				)
				if allowed {
					suite.Empty(err.Error())
				} else {
					suite.NotEmpty(err.Error())
				}
			})
		}
	})

	suite.Run("MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate Invalid add tertiary address without secondary", func() {
		tertiaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, nil)

		minimalMove := models.MTOShipment{
			TertiaryDeliveryAddress: &tertiaryDeliveryAddress,
		}

		mtoShipment_ThNScndP_address_Move := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		checker := MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate()
		err := checker.Validate(suite.AppContextForTest(), &minimalMove, &mtoShipment_ThNScndP_address_Move.MTOShipments[0])
		suite.Error(err)
	})

	suite.Run("MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate Valid add secondary address", func() {
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		minimalMove := models.MTOShipment{
			SecondaryPickupAddress: &secondaryPickupAddress,
		}

		mtoShipment_ThNScndP_address_Move := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		checker := MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate()
		err := checker.Validate(suite.AppContextForTest(), &mtoShipment_ThNScndP_address_Move.MTOShipments[0], &minimalMove)
		suite.NoError(err)
	})

	suite.Run("MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate Valid remove secondary address", func() {
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		oldMove := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		oldMove.MTOShipments[0].SecondaryPickupAddress = &secondaryPickupAddress

		newMove := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		checker := MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate()
		err := checker.Validate(suite.AppContextForTest(), &newMove.MTOShipments[0], &oldMove.MTOShipments[0])
		suite.NoError(err)
	})

	suite.Run("MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate Valid", func() {
		tertiaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		minimalMove := models.MTOShipment{
			TertiaryPickupAddress: &tertiaryPickupAddress,
		}

		mtoShipment_ThNScndP_address_Move := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		checker := MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate()
		err := checker.Validate(suite.AppContextForTest(), &mtoShipment_ThNScndP_address_Move.MTOShipments[0], &minimalMove)
		suite.NoError(err)
	})

	suite.Run("MTOShipmentHasTertiaryAddressWithNoSecondaryAddressCreate No Secondary Address With Tertiary Invalid", func() {
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		TertiaryDestinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		tertiaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		mtoShipment_Valid_address := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		mtoShipment_Valid_address.MTOShipments[0].SecondaryPickupAddress = &secondaryPickupAddress
		mtoShipment_Valid_address.MTOShipments[0].TertiaryPickupAddress = &tertiaryPickupAddress
		mtoShipment_Valid_address.MTOShipments[0].TertiaryDeliveryAddress = &TertiaryDestinationAddress

		checker := MTOShipmentHasTertiaryAddressWithNoSecondaryAddressCreate()
		err := checker.Validate(suite.AppContextForTest(), &mtoShipment_Valid_address.MTOShipments[0], nil)
		suite.Error(err)
	})

	suite.Run("MTOShipmentHasTertiaryAddressWithNoSecondaryAddressCreate with Secondary Address Valid", func() {
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		TertiaryDestinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		tertiaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		mtoShipment_Valid_address := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		mtoShipment_Valid_address.MTOShipments[0].SecondaryPickupAddress = &secondaryPickupAddress
		mtoShipment_Valid_address.MTOShipments[0].TertiaryPickupAddress = &tertiaryPickupAddress
		mtoShipment_Valid_address.MTOShipments[0].TertiaryDeliveryAddress = &TertiaryDestinationAddress

		checker := MTOShipmentHasTertiaryAddressWithNoSecondaryAddressCreate()
		err := checker.Validate(suite.AppContextForTest(), &mtoShipment_Valid_address.MTOShipments[0], nil)
		suite.Error(err)
	})

	suite.Run("MTOShipmentHasTertiaryAddressWithNoSecondaryAddressCreate Valid", func() {
		SecondaryDestinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		TertiaryDestinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		tertiaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		mtoShipment_Valid_address := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		mtoShipment_Valid_address.MTOShipments[0].SecondaryPickupAddress = &secondaryPickupAddress
		mtoShipment_Valid_address.MTOShipments[0].SecondaryDeliveryAddress = &SecondaryDestinationAddress
		mtoShipment_Valid_address.MTOShipments[0].TertiaryPickupAddress = &tertiaryPickupAddress
		mtoShipment_Valid_address.MTOShipments[0].TertiaryDeliveryAddress = &TertiaryDestinationAddress

		checker := MTOShipmentHasTertiaryAddressWithNoSecondaryAddressCreate()
		err := checker.Validate(suite.AppContextForTest(), &mtoShipment_Valid_address.MTOShipments[0], nil)
		suite.NoError(err)
	})

	suite.Run("checkAvailToPrime", func() {
		appCtx := suite.AppContextForTest()

		now := time.Now()
		hide := false
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		primeShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
		}, nil)
		nonPrimeShipment := factory.BuildMTOShipmentMinimal(appCtx.DB(), nil, nil)
		externalShipment := factory.BuildMTOShipmentMinimal(appCtx.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
					UsesExternalVendor: true,
				},
			},
		}, nil)
		hiddenPrimeShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Show:               &hide,
				},
			},
		}, nil)
		badUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

		testCases := map[string]struct {
			id   uuid.UUID
			verf func(error)
		}{
			"happy path": {
				primeShipment.ID,
				func(err error) {
					suite.Require().NoError(err)
				},
			},
			"exists unavailable": {
				nonPrimeShipment.ID,
				func(err error) {
					suite.Require().Error(err)
					suite.IsType(apperror.NotFoundError{}, err)
					suite.Contains(err.Error(), nonPrimeShipment.ID.String())
				},
			},
			"external vendor": {
				externalShipment.ID,
				func(err error) {
					suite.Require().Error(err)
					suite.IsType(apperror.NotFoundError{}, err)
					suite.Contains(err.Error(), externalShipment.ID.String())
				},
			},
			"disabled move": {
				hiddenPrimeShipment.ID,
				func(err error) {
					suite.Require().Error(err)
					suite.IsType(apperror.NotFoundError{}, err)
					suite.Contains(err.Error(), hiddenPrimeShipment.ID.String())
				},
			},
			"does not exist": {
				badUUID,
				func(err error) {
					suite.Require().Error(err)
					suite.IsType(apperror.NotFoundError{}, err)
					suite.Contains(err.Error(), badUUID.String())
				},
			},
		}

		for name, tc := range testCases {
			suite.Run(name, func() {
				checker := checkAvailToPrime()
				err := checker.Validate(appCtx, &models.MTOShipment{ID: tc.id}, nil)
				tc.verf(err)
			})
		}
	})

	suite.Run("checkUpdateAllowed", func() {
		servicesCounselor := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		servicesCounselorSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *servicesCounselor.UserID,
			OfficeUserID:    servicesCounselor.ID,
		}
		servicesCounselorSession.Roles = append(servicesCounselorSession.Roles, servicesCounselor.User.Roles...)

		too := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		tooSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		tooSession.Roles = append(tooSession.Roles, too.User.Roles...)

		tio := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})
		tioSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *tio.UserID,
			OfficeUserID:    tio.ID,
		}
		tioSession.Roles = append(tioSession.Roles, tio.User.Roles...)

		testCases := map[string]struct {
			session auth.Session
			tests   map[models.MTOShipmentStatus]bool
		}{
			"Service Counselor": {
				servicesCounselorSession,
				map[models.MTOShipmentStatus]bool{
					models.MTOShipmentStatusSubmitted:             true,
					models.MTOShipmentStatusApproved:              true,
					models.MTOShipmentStatusCancellationRequested: false,
					models.MTOShipmentStatusCanceled:              false,
					models.MTOShipmentStatusDiversionRequested:    false,
					models.MTOShipmentStatusDraft:                 true,
				},
			},
			"TOO": {
				tooSession,
				map[models.MTOShipmentStatus]bool{
					models.MTOShipmentStatusSubmitted:             true,
					models.MTOShipmentStatusApproved:              true,
					models.MTOShipmentStatusCancellationRequested: true,
					models.MTOShipmentStatusCanceled:              true,
					models.MTOShipmentStatusDiversionRequested:    true,
					models.MTOShipmentStatusDraft:                 false,
				},
			},
			"TIO": {
				tioSession,
				map[models.MTOShipmentStatus]bool{
					models.MTOShipmentStatusSubmitted:             false,
					models.MTOShipmentStatusApproved:              true,
					models.MTOShipmentStatusCancellationRequested: false,
					models.MTOShipmentStatusCanceled:              false,
					models.MTOShipmentStatusDiversionRequested:    false,
					models.MTOShipmentStatusDraft:                 false,
				},
			},
			"Non-office user": {
				auth.Session{},
				map[models.MTOShipmentStatus]bool{
					models.MTOShipmentStatusSubmitted: false,
				},
			},
		}

		for name, tc := range testCases {
			for status, canUpdate := range tc.tests {
				appCtx := suite.AppContextWithSessionForTest(&tc.session) //#nosec G601

				suite.Run(fmt.Sprintf("User:%v Shipment Status:%v", name, status), func() {
					checker := checkUpdateAllowed()
					err := checker.Validate(appCtx, nil, &models.MTOShipment{Status: status})
					if canUpdate {
						suite.NoError(err)
					} else {
						suite.Error(err)
					}
				})
			}
		}
	})

}

func (suite *MTOShipmentServiceSuite) TestCheckAddressUpdateAllowed() {
	suite.Run("checkStatusAllowsAddressUpdates", func() {
		v4ID := uuid.Must(uuid.NewV4())
		bannedErrMsgPartial := "does not allow address updates"
		hhgIntoNtsErrMsgPartial := "cannot update the destination address of an NTS shipment directly"
		hhgOutOfNtsErrMsgPartial := "cannot update the pickup address of an NTS-Release shipment directly"
		approvedDestinationAddressErrMsgPartial := "please use the updateShipmentDestinationAddress endpoint / ShipmentAddressUpdateRequester service"
		testCases := map[string]struct {
			status           models.MTOShipmentStatus
			sType            models.MTOShipmentType
			canUpdate        bool
			applyIds         bool
			errorMsgIncludes string
		}{
			"Draft is not banned": {
				status:    models.MTOShipmentStatusDraft,
				canUpdate: true,
			},
			"Submitted is not banned": {
				status:    models.MTOShipmentStatusSubmitted,
				canUpdate: true,
			},
			"CancellationRequested is not banned": {
				status:    models.MTOShipmentStatusCancellationRequested,
				canUpdate: true,
			},
			"DiversionRequested is not banned": {
				status:    models.MTOShipmentStatusDiversionRequested,
				canUpdate: true,
			},
			"Approved is not banned": {
				status:    models.MTOShipmentStatusApproved,
				canUpdate: true,
			},
			"Rejected is banned": {
				status:           models.MTOShipmentStatusRejected,
				canUpdate:        false,
				errorMsgIncludes: bannedErrMsgPartial,
			},
			"Canceled is banned": {
				status:           models.MTOShipmentStatusCanceled,
				canUpdate:        false,
				errorMsgIncludes: bannedErrMsgPartial,
			},
			"TerminatedForCause is banned": {
				status:           models.MTOShipmentStatusTerminatedForCause,
				canUpdate:        false,
				errorMsgIncludes: bannedErrMsgPartial,
			},
			"HHG into NTS can't update dest address directly": {
				sType:            models.MTOShipmentTypeHHGIntoNTS,
				canUpdate:        false,
				applyIds:         true,
				errorMsgIncludes: hhgIntoNtsErrMsgPartial,
			},
			"HHG out of NTS can't update pickup address directly": {
				sType:            models.MTOShipmentTypeHHGOutOfNTS,
				canUpdate:        false,
				applyIds:         true,
				errorMsgIncludes: hhgOutOfNtsErrMsgPartial,
			},
			"Approved cannot have its destination address changed from this service, it must use the ShipmentAddressUpdateRequester service": {
				status:           models.MTOShipmentStatusApproved,
				canUpdate:        false,
				applyIds:         true,
				errorMsgIncludes: approvedDestinationAddressErrMsgPartial,
			},
		}
		// !IMPORANT!
		// Update this count on every new test case that isn't related to the status check of checkStatusNotBannedFromAddressUpdates
		var countOfNonStatusTestCases = 3

		appCtx := suite.AppContextForTest()

		// Check that we have a test case for all counts of possible shipment statuses
		type statusRow struct {
			Status string `db:"status"`
		}
		var rows []statusRow
		err := appCtx.DB().
			RawQuery(`SELECT unnest(enum_range(NULL::public.mto_shipment_status)) AS status`).
			All(&rows)
		suite.FatalNoError(err)
		suite.Require().Equal(len(testCases)-countOfNonStatusTestCases, len(rows), "The count of shipment status test cases do not match the amount pulled from the database enum")

		checker := checkAddressUpdateAllowed()

		for name, tc := range testCases {
			suite.Run(name, func() {
				address := models.Address{}
				shipment := models.MTOShipment{
					Status:       tc.status,
					ShipmentType: tc.sType,
				}
				if tc.applyIds {
					address.ID = v4ID
					shipment.PickupAddressID = &v4ID
					shipment.DestinationAddressID = &v4ID
				}
				err := checker.Validate(appCtx, &address, &shipment)
				if tc.canUpdate {
					suite.NoError(err, "expected no error for status %s", tc.status)
				} else {
					suite.Error(err, "expected error for status %s", tc.status)
					suite.ErrorContains(err, tc.errorMsgIncludes, "expected error to match the test case partial err msg case")
				}
			})
		}
	})

	suite.Run("Check error if shipment or address is nil", func() {
		testCases := map[string]struct {
			address  *models.Address
			shipment *models.MTOShipment
		}{
			"shipment should error if nil": {
				address: &models.Address{},
			},
			"address should error if nil": {
				shipment: &models.MTOShipment{},
			},
		}
		appCtx := suite.AppContextForTest()
		checker := checkAddressUpdateAllowed()
		for name, tc := range testCases {
			suite.Run(name, func() {
				err := checker.Validate(appCtx, tc.address, tc.shipment)
				suite.Error(err)
				suite.ErrorContains(err, "shipment address updater is not passing needed validator values")
			})
		}
	})
}

func (suite *MTOShipmentServiceSuite) TestDeleteValidations() {
	suite.Run("checkDeleteAllowedTOO", func() {
		testCases := map[models.MoveStatus]bool{
			models.MoveStatusDRAFT:                      true,
			models.MoveStatusSUBMITTED:                  true,
			models.MoveStatusCANCELED:                   true,
			models.MoveStatusAPPROVALSREQUESTED:         true,
			models.MoveStatusNeedsServiceCounseling:     true,
			models.MoveStatusServiceCounselingCompleted: true,
		}

		for status, allowed := range testCases {
			suite.Run("Move status "+string(status), func() {
				shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
					{
						Model: models.Move{
							Status: status,
						},
					},
				}, nil)

				officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})

				appContext := suite.AppContextWithSessionForTest(&auth.Session{
					Roles:           officeUser.User.Roles,
					ApplicationName: auth.OfficeApp,
				})

				err := checkDeleteAllowed().Validate(
					appContext,
					nil,
					&shipment,
				)

				if allowed {
					suite.NoError(err)
				} else {
					suite.Error(err)
				}
			})
		}
	})

	suite.Run("checkDeleteAllowedTOO Approved MTO status", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})

		appContext := suite.AppContextWithSessionForTest(&auth.Session{
			Roles:           officeUser.User.Roles,
			ApplicationName: auth.OfficeApp,
		})

		err := checkDeleteAllowed().Validate(
			appContext,
			nil,
			&shipment,
		)

		if false {
			suite.NoError(err)
		} else {
			suite.Error(err)
		}
	})

	suite.Run("checkDeleteAllowedSC", func() {
		testCases := map[models.MoveStatus]bool{
			models.MoveStatusDRAFT:                      true,
			models.MoveStatusSUBMITTED:                  false,
			models.MoveStatusAPPROVED:                   false,
			models.MoveStatusCANCELED:                   false,
			models.MoveStatusAPPROVALSREQUESTED:         false,
			models.MoveStatusNeedsServiceCounseling:     true,
			models.MoveStatusServiceCounselingCompleted: false,
		}

		for status, allowed := range testCases {
			suite.Run("Move status "+string(status), func() {
				shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
					{
						Model: models.Move{
							Status: status,
						},
					},
				}, nil)

				officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

				appContext := suite.AppContextWithSessionForTest(&auth.Session{
					Roles:           officeUser.User.Roles,
					ApplicationName: auth.OfficeApp,
				})

				err := checkDeleteAllowed().Validate(
					appContext,
					nil,
					&shipment,
				)

				if allowed {
					suite.NoError(err)
				} else {
					suite.Error(err)
				}
			})
		}
	})

	suite.Run("checkPrimeDeleteAllowed for non-PPM shipments", func() {
		testCases := map[models.MTOShipmentType]bool{
			models.MTOShipmentTypeHHG:                  false,
			models.MTOShipmentTypeHHGIntoNTS:           false,
			models.MTOShipmentTypeHHGOutOfNTS:          false,
			models.MTOShipmentTypeMobileHome:           false,
			models.MTOShipmentTypeBoatHaulAway:         false,
			models.MTOShipmentTypeBoatTowAway:          false,
			models.MTOShipmentTypeUnaccompaniedBaggage: false,
		}

		for shipmentType, allowed := range testCases {
			suite.Run("Shipment type "+string(shipmentType), func() {
				now := time.Now()
				shipment := factory.BuildMTOShipment(nil, []factory.Customization{
					{
						Model: models.MTOShipment{
							ShipmentType: shipmentType,
						},
					},
					{
						Model: models.Move{
							AvailableToPrimeAt: &now,
							ApprovedAt:         &now,
						},
					},
				}, nil)

				err := checkPrimeDeleteAllowed().Validate(
					suite.AppContextForTest(),
					nil,
					&shipment,
				)

				if allowed {
					suite.NoError(err)
				} else {
					suite.Error(err)
					suite.Contains(err.Error(), "Prime can only delete PPM shipments")
				}
			})
		}
	})

	suite.Run("checkPrimeDeleteAllowed based on PPM status", func() {
		testCases := map[models.PPMShipmentStatus]bool{
			models.PPMShipmentStatusDraft:                true,
			models.PPMShipmentStatusSubmitted:            true,
			models.PPMShipmentStatusWaitingOnCustomer:    false,
			models.PPMShipmentStatusNeedsAdvanceApproval: true,
			models.PPMShipmentStatusNeedsCloseout:        true,
			models.PPMShipmentStatusCloseoutComplete:     true,
		}

		for status, allowed := range testCases {
			now := time.Now()
			suite.Run("PPM status "+string(status), func() {
				ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
					{
						Model: models.PPMShipment{
							Status: status,
						},
					},
					{
						Model: models.Move{
							AvailableToPrimeAt: &now,
							ApprovedAt:         &now,
						},
					},
				}, nil)
				err := checkPrimeDeleteAllowed().Validate(
					suite.AppContextForTest(),
					nil,
					&ppmShipment.Shipment,
				)

				if allowed {
					suite.NoError(err)
				} else {
					suite.Error(err)
					suite.Contains(err.Error(), "A PPM shipment with the status WAITING_ON_CUSTOMER cannot be deleted")
				}
			})
		}
	})

	suite.Run("checkPrimeDeleteAllowed for move not available to prime", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: nil,
					ApprovedAt:         nil,
				},
			},
		}, nil)

		err := checkPrimeDeleteAllowed().Validate(
			suite.AppContextForTest(),
			nil,
			&shipment,
		)
		suite.Error(err)
		suite.Contains(err.Error(), "not found for mtoShipment")
	})
}
