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

		// update move
		tioRemarks := "updating TIO remarks for test"
		approvedMove.TIORemarks = &tioRemarks
		suite.MustSave(&approvedMove)

		moveHistory, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), approvedMove.Locator)
		suite.FatalNoError(err)

		// address update
		verifyOldPickupAddress := false
		verifyNewPickupAddress := false
		// orders update
		// verifyOldSAC := false
		// verifyNewSAC := false
		// move update
		verifyOldTIORemarks := false
		verifyTIORemarks := false

		for _, h := range moveHistory.AuditHistories {

			if h.TableName == "addresses" {
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

func (suite *MoveHistoryServiceSuite) TestMoveFetcherWithFakeData() {
	moveHistoryFetcher := NewMoveHistoryFetcher()

	suite.T().Run("returns Audit History with session information", func(t *testing.T) {

		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		fakeRole := testdatagen.MakeRole(suite.DB(), testdatagen.Assertions{})
		fakeUser := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{})
		_ = testdatagen.MakeUsersRoles(suite.DB(), testdatagen.Assertions{
			User: fakeUser,
			UsersRoles: models.UsersRoles{
				RoleID: fakeRole.ID,
			},
		})

		_ = testdatagen.MakeAuditHistory(suite.DB(), testdatagen.Assertions{
			User: fakeUser,
			Move: models.Move{
				ID: approvedMove.ID,
			},
		})

		moveHistoryData, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), approvedMove.Locator)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

	})

}
