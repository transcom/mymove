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

func (suite *MTOShipmentServiceSuite) TestMTOShipmentHasValidRequestedPickupDate() {
	uuidTest, _ := uuid.NewV4()
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	futureDate := models.TimePointer(tomorrow)
	pastDate := models.TimePointer(today.Add(-24 * time.Hour))
	zeroTime := time.Time{}
	requiredDateError := "RequestedPickupDate is required to create an HHG shipment"
	invalidDateError := "RequestedPickupDate must be greater than or equal to tomorrow's date"

	testCases := []struct {
		name          string
		newer         *models.MTOShipment
		older         *models.MTOShipment
		expectedError bool
		errorMessage  string
	}{
		{
			name: "Zero RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				RequestedPickupDate: &zeroTime,
				ShipmentType:        models.MTOShipmentTypeHHG,
			},
			expectedError: true,
			errorMessage:  requiredDateError,
		},
		{
			name: "RequestedPickupDate today",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				RequestedPickupDate: &today,
				ShipmentType:        models.MTOShipmentTypeHHG,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "RequestedPickupDate in future",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				RequestedPickupDate: futureDate,
				ShipmentType:        models.MTOShipmentTypeHHG,
			},
			expectedError: false,
		},
		{
			name: "Nil RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:           uuidTest,
				ShipmentType: models.MTOShipmentTypeHHG,
			},
			expectedError: true,
			errorMessage:  requiredDateError,
		},
		{
			name: "PPM shipment with nil RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:           uuidTest,
				ShipmentType: models.MTOShipmentTypePPM,
			},
			expectedError: false,
		},
		{
			name: "HHG shipment with past RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				RequestedPickupDate: pastDate,
				ShipmentType:        models.MTOShipmentTypeHHG,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "Boat shipment with nil RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:           uuidTest,
				ShipmentType: models.MTOShipmentTypeBoatHaulAway,
			},
			expectedError: false,
		},
		{
			name: "Boat shipment with past RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeBoatHaulAway,
				RequestedPickupDate: pastDate,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "Boat shipment with todays date for RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeBoatHaulAway,
				RequestedPickupDate: &today,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "Mobile home shipment with nil RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:           uuidTest,
				ShipmentType: models.MTOShipmentTypeMobileHome,
			},
			expectedError: false,
		},
		{
			name: "Mobile home shipment with past RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeMobileHome,
				RequestedPickupDate: pastDate,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "Mobile home shipment with todays date for RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeMobileHome,
				RequestedPickupDate: &today,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "HHG Out of NTS shipment with nil RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:           uuidTest,
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTS,
			},
			expectedError: false,
		},
		{
			name: "HHG Out of NTS shipment with past RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeHHGOutOfNTS,
				RequestedPickupDate: pastDate,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "HHG Out of NTS shipment with todays date for RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeHHGOutOfNTS,
				RequestedPickupDate: &today,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "HHG Into NTS shipment with nil RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:           uuidTest,
				ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
			},
			expectedError: true,
			errorMessage:  "RequestedPickupDate is required to create an HHG_INTO_NTS shipment",
		},
		{
			name: "HHG Into NTS shipment with past RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeHHGIntoNTS,
				RequestedPickupDate: pastDate,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "HHG Into NTS shipment with todays date for RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeHHGIntoNTS,
				RequestedPickupDate: &today,
			},
			expectedError: true,
			errorMessage:  invalidDateError,
		},
		{
			name: "HHG Into NTS shipment with future date for RequestedPickupDate",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				ShipmentType:        models.MTOShipmentTypeHHGIntoNTS,
				RequestedPickupDate: futureDate,
			},
			expectedError: false,
		},
		{
			name: "Update from valid date to nil",
			newer: &models.MTOShipment{
				ID:           uuidTest,
				ShipmentType: models.MTOShipmentTypeHHG,
			},
			older: &models.MTOShipment{
				ID:                  uuidTest,
				RequestedPickupDate: &tomorrow,
				ShipmentType:        models.MTOShipmentTypeHHG,
			},
			expectedError: true,
			errorMessage:  requiredDateError,
		},
		{
			name: "Update from valid date to new valid date",
			newer: &models.MTOShipment{
				ID:                  uuidTest,
				RequestedPickupDate: models.TimePointer(tomorrow.Add(24 * time.Hour)),
				ShipmentType:        models.MTOShipmentTypeHHG,
			},
			older: &models.MTOShipment{
				ID:                  uuidTest,
				RequestedPickupDate: &tomorrow,
				ShipmentType:        models.MTOShipmentTypeHHG,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			validator := MTOShipmentHasValidRequestedPickupDate()
			err := validator.Validate(suite.AppContextForTest(), tc.newer, tc.older)
			if tc.expectedError {
				suite.Error(err)
				suite.Contains(err.Error(), tc.errorMessage)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestGetAorAnWithShipmentType() {
	suite.Run("GetAorAnWithShipmentType", func() {
		testCases := map[models.MTOShipmentType]string{
			models.MTOShipmentTypeHHG:                  "an",
			models.MTOShipmentTypeHHGIntoNTS:           "an",
			models.MTOShipmentTypeHHGOutOfNTS:          "an",
			models.MTOShipmentTypeUnaccompaniedBaggage: "an",
			models.MTOShipmentTypeMobileHome:           "a",
			models.MTOShipmentTypeBoatHaulAway:         "a",
			models.MTOShipmentTypeBoatTowAway:          "a",
			models.MTOShipmentTypePPM:                  "a",
			"UnknownType":                              "a",
		}

		for shipmentType, expected := range testCases {
			suite.Run(string(shipmentType), func() {
				result := GetAorAnByShipmentType(shipmentType)
				suite.Equal(expected, result)
			})
		}
	})
}
