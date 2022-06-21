package mtoshipment

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/auth"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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

	suite.Run("checkAvailToPrime", func() {
		appCtx := suite.AppContextForTest()

		now := time.Now()
		hide := false
		availableToPrimeMove := testdatagen.MakeAvailableMove(appCtx.DB())
		primeShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: availableToPrimeMove,
		})
		nonPrimeShipment := testdatagen.MakeDefaultMTOShipmentMinimal(appCtx.DB())
		externalShipment := testdatagen.MakeMTOShipmentMinimal(appCtx.DB(), testdatagen.Assertions{
			Move: availableToPrimeMove,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				UsesExternalVendor: true,
			},
		})
		hiddenPrimeShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
				Show:               &hide,
			},
		})
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
		servicesCounselor := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{})
		servicesCounselorSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *servicesCounselor.UserID,
			OfficeUserID:    servicesCounselor.ID,
		}
		servicesCounselorSession.Roles = append(servicesCounselorSession.Roles, servicesCounselor.User.Roles...)

		too := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{})
		tooSession := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *too.UserID,
			OfficeUserID:    too.ID,
		}
		tooSession.Roles = append(tooSession.Roles, too.User.Roles...)

		tio := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})
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
					models.MTOShipmentStatusApproved:              false,
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
				appCtx := suite.AppContextWithSessionForTest(&tc.session)

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
	suite.Run("checkDeleteAllowed", func() {
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
				shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
					Move: models.Move{
						Status: status,
					},
				})

				err := checkDeleteAllowed().Validate(
					suite.AppContextForTest(),
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
			models.MTOShipmentTypeHHG:              false,
			models.MTOShipmentTypeInternationalHHG: false,
			models.MTOShipmentTypeInternationalUB:  false,
			models.MTOShipmentTypeHHGLongHaulDom:   false,
			models.MTOShipmentTypeHHGShortHaulDom:  false,
			models.MTOShipmentTypeHHGIntoNTSDom:    false,
			models.MTOShipmentTypeHHGOutOfNTSDom:   false,
			models.MTOShipmentTypeMotorhome:        false,
			models.MTOShipmentTypeBoatHaulAway:     false,
			models.MTOShipmentTypeBoatTowAway:      false,
		}

		for shipmentType, allowed := range testCases {
			suite.Run("Shipment type "+string(shipmentType), func() {
				shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
					MTOShipment: models.MTOShipment{
						ShipmentType: shipmentType,
					},
					Stub: true,
				})

				err := checkPrimeDeleteAllowed().Validate(
					suite.AppContextForTest(),
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

	suite.Run("checkPrimeDeleteAllowed based on PPM status", func() {
		testCases := map[models.PPMShipmentStatus]bool{
			models.PPMShipmentStatusDraft:                true,
			models.PPMShipmentStatusSubmitted:            true,
			models.PPMShipmentStatusWaitingOnCustomer:    false,
			models.PPMShipmentStatusNeedsAdvanceApproval: true,
			models.PPMShipmentStatusNeedsPaymentApproval: true,
			models.PPMShipmentStatusPaymentApproved:      true,
		}

		for status, allowed := range testCases {
			now := time.Now()
			suite.Run("PPM status "+string(status), func() {
				ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
					PPMShipment: models.PPMShipment{
						Status: status,
					},
					Move: models.Move{
						AvailableToPrimeAt: &now,
					},
				})

				err := checkPrimeDeleteAllowed().Validate(
					suite.AppContextForTest(),
					nil,
					&ppmShipment.Shipment,
				)

				if allowed {
					suite.NoError(err)
				} else {
					suite.Error(err)
				}
			})
		}
	})
}
