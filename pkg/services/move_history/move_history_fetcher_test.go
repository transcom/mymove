package movehistory

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
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

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(20)}
		moveHistory, totalCount, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.FatalNoError(err)

		suite.Equal(totalCount, int64(6), "total count should be 6")

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
						oldData := removeEscapeJSON(h.OldData)
						if oldData["city"] == oldAddress.City && oldData["state"] == oldAddress.State && oldData["postal_code"] == oldAddress.PostalCode {
							verifyOldPickupAddress = true
						}
					}
					if h.ChangedData != nil {
						changedData := removeEscapeJSON(h.ChangedData)
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
				if h.OldData != nil {
					oldData := removeEscapeJSON(h.OldData)
					if len(oldData["tio_remarks"]) == 0 {
						verifyOldTIORemarks = true
					}
				}
				if *h.ObjectID == approvedMove.ID {
					if h.ChangedData != nil {
						changedData := removeEscapeJSON(h.ChangedData)
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

		params := services.FetchMoveHistoryParams{Locator: "QX97UY", Page: swag.Int64(1), PerPage: swag.Int64(20)}
		_, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

}

func removeEscapeJSON(data *string) map[string]string {
	var result map[string]string
	if data == nil || *data == "" {
		return result
	}
	var byteData = []byte(*data)

	_ = json.Unmarshal(byteData, &result)
	return result
}

func (suite *MoveHistoryServiceSuite) TestMoveFetcherWithFakeData() {
	moveHistoryFetcher := NewMoveHistoryFetcher()

	suite.T().Run("returns Audit History with session information", func(t *testing.T) {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		fakeRole := testdatagen.MakeTOORole(suite.DB())
		fakeUser := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{})
		_ = testdatagen.MakeUsersRoles(suite.DB(), testdatagen.Assertions{
			User: fakeUser,
			UsersRoles: models.UsersRoles{
				RoleID: fakeRole.ID,
			},
		})
		_ = testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				User:   fakeUser,
				UserID: &fakeUser.ID,
			},
		})

		_ = testdatagen.MakeAuditHistory(suite.DB(), testdatagen.Assertions{
			User: fakeUser,
			Move: models.Move{
				ID: approvedMove.ID,
			},
		})

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(20)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		suite.NotEmpty(moveHistoryData.AuditHistories, "AuditHistories should not be empty")
		suite.NotEmpty(moveHistoryData.AuditHistories[0].SessionUserID, "AuditHistories contains an AuditHistory with a SessionUserID")
		suite.NotEmpty(moveHistoryData.AuditHistories[0].SessionUserFirstName, "AuditHistories contains an AuditHistory with a SessionUserFirstName")
		suite.NotEmpty(moveHistoryData.AuditHistories[0].SessionUserLastName, "AuditHistories contains an AuditHistory with a SessionUserLastName")
		suite.NotEmpty(moveHistoryData.AuditHistories[0].SessionUserEmail, "AuditHistories contains an AuditHistory with a SessionUserEmail")
		suite.NotEmpty(moveHistoryData.AuditHistories[0].SessionUserTelephone, "AuditHistories contains an AuditHistory with a SessionUserTelephone")
	})

	suite.T().Run("filters shipments from different move ", func(t *testing.T) {

		auditHistoryContains := func(auditHistories models.AuditHistories, keyword string) func() (success bool) {
			return func() (success bool) {
				for _, record := range auditHistories {
					if strings.Contains(*record.ChangedData, keyword) {
						return true
					}
				}
				return false
			}
		}

		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		approvedShipment := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &approvedMove, testdatagen.Assertions{})

		approvedMoveToFilter := testdatagen.MakeAvailableMove(suite.DB())
		approvedShipmentToFilter := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &approvedMoveToFilter, testdatagen.Assertions{})

		customerRemarks := "fragile"
		approvedShipment.CustomerRemarks = &customerRemarks
		suite.MustSave(&approvedShipment)

		customerRemarksFilter := "sturdy"
		approvedShipmentToFilter.CustomerRemarks = &customerRemarksFilter
		suite.MustSave(&approvedShipmentToFilter)

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(20)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		suite.Equal(5, len(moveHistoryData.AuditHistories), "should not have more than 5")
		suite.Condition(auditHistoryContains(moveHistoryData.AuditHistories, "fragile"), "should contain fragile")
		containsSturdy := auditHistoryContains(moveHistoryData.AuditHistories, "sturdy")()
		suite.False(containsSturdy, "should not contain sturdy")

	})

	suite.T().Run("has context", func(t *testing.T) {
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()

		updater := mtoserviceitem.NewMTOServiceItemUpdater(builder, moveRouter)
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		rejectionReason := swag.String("")

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		params := services.FetchMoveHistoryParams{Locator: move.Locator, Page: swag.Int64(1), PerPage: swag.Int64(20)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		suite.NotEmpty(moveHistoryData.AuditHistories, "AuditHistories should not be empty")
		verifyServiceItemStatusContext := false
		for _, h := range moveHistoryData.AuditHistories {
			if h.TableName == "mto_service_items" {
				if *h.ObjectID == updatedServiceItem.ID {
					if h.Context != nil {
						context := removeEscapeJSON(h.Context)
						if context != nil && context["name"] == serviceItem.ReService.Name && context["shipment_type"] == string(serviceItem.MTOShipment.ShipmentType) {
							verifyServiceItemStatusContext = true
						}
					}
				}
			}
		}
		suite.True(verifyServiceItemStatusContext, "AuditHistories contains an AuditHistory with a Context when a service item is approved")
	})

	suite.T().Run("has paginated results", func(t *testing.T) {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())

		// update move
		tioRemarks := "updating TIO remarks for test"
		approvedMove.TIORemarks = &tioRemarks
		suite.MustSave(&approvedMove)

		// update move
		tioRemarks = "updating TIO remarks for test AGAIN"
		approvedMove.TIORemarks = &tioRemarks
		suite.MustSave(&approvedMove)

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(2)}
		moveHistoryData, totalCount, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)
		suite.Greater(totalCount, int64(2), "total count should be 5")
		suite.Equal(2, len(moveHistoryData.AuditHistories), "should have 2 rows due to pagination")

	})
}
