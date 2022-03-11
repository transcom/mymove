package movehistory

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveHistoryServiceSuite) TestMoveFetcher() {
	moveHistoryFetcher := NewMoveHistoryFetcher()

	suite.T().Run("successfully returns submitted move history available to prime", func(t *testing.T) {

		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		approvedShipment := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &approvedMove, testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: approvedMove,
		})

		// update HHG SAC
		updateSAC := "23456"
		approvedMove.Orders.SAC = &updateSAC
		suite.MustSave(&approvedMove.Orders)

		// update Pickup Address
		oldAddress := *approvedShipment.PickupAddress
		updateAddress := approvedShipment.PickupAddress
		updateAddress.City = "Norfolk"
		updateAddress.State = "VA"
		updateAddress.PostalCode = "23503"
		suite.MustSave(updateAddress)

		//update shipment weight
		//oldWeight := *approvedShipment.PrimeActualWeight
		//updateWeight := *approvedShipment.PrimeActualWeight + unit.Pound(1000)
		//*approvedShipment.PrimeActualWeight = updateWeight
		//suite.MustSave(&approvedShipment)

		// update move
		tioRemarks := "updating TIO remarks for test"
		approvedMove.TIORemarks = &tioRemarks
		suite.MustSave(&approvedMove)

		moveHistory, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), approvedMove.Locator)
		suite.FatalNoError(err)

		// address update
		verifyOldPickupAddress := false
		verifyNewPickupAddress := false
		// shipment update
		//verifyOldWeight := false
		//verifyNewWeight := false
		// orders update
		// verifyOldSAC := false
		// verifyNewSAC := false
		// move update
		verifyOldTIORemarks := false
		verifyTIORemarks := false

		for _, h := range moveHistory.AuditHistories {

			if h.TableName == "mto_shipments" {
				/*
					if *h.ObjectID == approvedShipment.ID {
						if h.OldData != nil {
							oldData := *h.OldData
							if oldData["prime_actual_weight"] == oldWeight {
								verifyOldWeight = true
							}
						}
						if h.ChangedData != nil {
							changedData := *h.ChangedData
							fmt.Printf("+%v\n", changedData["prime_actual_weight"].(float64))
							fmt.Printf("+%v\n", updateWeight)
							weight, ok := changedData["prime_actual_weight"].(float64)
							if ok {
								// w, _ := strconv.ParseFloat(weight, 64)
								if weight == updateWeight.Float64() {
									verifyNewWeight = true
								}
							}

						}
					}
				*/
			} else if h.TableName == "addresses" {
				if *h.ObjectID == updateAddress.ID {
					if h.OldData != nil {
						oldData := *h.OldData
						if oldData["city"] == oldAddress.City && oldData["state"] == oldAddress.State && oldData["postal_code"] == oldAddress.PostalCode {
							verifyOldPickupAddress = true
						}
					}
					if h.ChangedData != nil {
						changedData := *h.ChangedData
						if changedData["city"] == updateAddress.City && changedData["state"] == updateAddress.State && changedData["postal_code"] == updateAddress.PostalCode {
							verifyNewPickupAddress = true
						}
					}
				}
				/*} else if h.TableName == "orders" {
				if *h.ObjectID == approvedMove.Orders.ID {
					if h.OldData != nil {
						oldData := *h.OldData
						if oldData["sac"] == nil {
							verifyOldSAC = true
						}
					}
					if h.ChangedData != nil {
						changedData := *h.ChangedData
						if changedData["sac"] == updateSAC {
							verifyNewSAC = true
						}
					}
				}*/
			} else if h.TableName == "moves" {
				if *h.ObjectID == approvedMove.ID {
					if h.OldData != nil {
						oldData := *h.OldData
						if oldData["tio_remarks"] == nil {
							verifyOldTIORemarks = true
						}
					}
					if h.ChangedData != nil {
						changedData := *h.ChangedData
						if changedData["tio_remarks"] == tioRemarks {
							verifyTIORemarks = true
						}
					}
				}
			}

		}

		suite.Equal(approvedMove.ID, moveHistory.ID)
		suite.Equal(approvedMove.Locator, moveHistory.Locator)
		suite.Equal(approvedMove.ReferenceID, moveHistory.ReferenceID)

		// address update
		suite.True(verifyOldPickupAddress, "verifyOldPickupAddress")
		suite.True(verifyNewPickupAddress, "verifyNewPickupAddress")
		// shipment update
		//suite.True(verifyOldWeight, "verifyOldWeight")
		//suite.True(verifyNewWeight, "verifyNewWeight")
		// orders update
		// suite.True(verifyOldSAC, "verifyOldSAC")
		// suite.True(verifyNewSAC, "verifyNewSAC")
		// move update
		suite.True(verifyOldTIORemarks, "verifyOldTIORemarks")
		suite.True(verifyTIORemarks, "verifyTIORemarks")
	})

	suite.T().Run("returns not found error for unknown locator", func(t *testing.T) {
		_ = testdatagen.MakeAvailableMove(suite.DB())

		_, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), "QX97UY")
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

}
