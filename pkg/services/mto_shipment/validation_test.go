package mtoshipment

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
					context.Background(),
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
		now := time.Now()
		hide := false
		primeShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
		})
		nonPrimeShipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		hiddenPrimeShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
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
					suite.IsType(err, services.NotFoundError{})
					suite.Contains(err.Error(), nonPrimeShipment.ID.String())
				},
			},
			"disabled move": {
				hiddenPrimeShipment.ID,
				func(err error) {
					suite.Require().Error(err)
					suite.IsType(err, services.NotFoundError{})
					suite.Contains(err.Error(), hiddenPrimeShipment.ID.String())
				},
			},
			"does not exist": {
				badUUID,
				func(err error) {
					suite.Require().Error(err)
					suite.IsType(err, services.NotFoundError{})
					suite.Contains(err.Error(), badUUID.String())
				},
			},
		}

		for name, tc := range testCases {
			suite.Run(name, func() {
				checker := checkAvailToPrime(suite.DB())
				err := checker.Validate(context.Background(), &models.MTOShipment{ID: tc.id}, nil)
				tc.verf(err)
			})
		}
	})
}
