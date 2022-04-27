package movehistory

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/reweigh"

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
		// update authorized weight
		updateDBAuthorizedWeight := 500
		approvedMove.Orders.Entitlement.DBAuthorizedWeight = &updateDBAuthorizedWeight
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
		moveHistory, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.FatalNoError(err)

		// address update
		verifyOldPickupAddress := false
		verifyNewPickupAddress := false
		// orders update
		verifyOldSAC := false
		verifyNewSAC := false
		// move update
		verifyOldTIORemarks := false
		verifyTIORemarks := false
		verifyDBAuthorizedWeight := false

		for _, h := range moveHistory.AuditHistories {

			if h.TableName == "addresses" {
				if *h.ObjectID == updateAddress.ID {
					if h.OldData != nil {
						oldData := removeEscapeJSONtoObject(h.OldData)
						if oldData["city"] == oldAddress.City && oldData["state"] == oldAddress.State && oldData["postal_code"] == oldAddress.PostalCode {
							verifyOldPickupAddress = true
						}
					}
					if h.ChangedData != nil {
						changedData := removeEscapeJSONtoObject(h.ChangedData)
						if changedData["city"] == updateAddress.City && changedData["state"] == updateAddress.State && changedData["postal_code"] == updateAddress.PostalCode {
							verifyNewPickupAddress = true
						}
					}
				}
			} else if h.TableName == "orders" {
				if *h.ObjectID == approvedMove.Orders.ID {
					if h.OldData != nil {
						oldData := removeEscapeJSONtoObject(h.OldData)
						if len(oldData["sac"]) == 0 {
							verifyOldSAC = true
						}
					}
					if h.ChangedData != nil {
						changedData := removeEscapeJSONtoObject(h.ChangedData)
						if changedData["sac"] == updateSAC {
							verifyNewSAC = true
						}
					}
				}
			} else if h.TableName == "entitlements" {
				if h.ChangedData != nil {
					oldData := removeEscapeJSONtoObject(h.OldData)
					if len(oldData["authorized_weight"]) == 0 {
						verifyDBAuthorizedWeight = true
					}
				}
			} else if h.TableName == "moves" {
				if h.OldData != nil {
					oldData := removeEscapeJSONtoObject(h.OldData)
					if len(oldData["tio_remarks"]) == 0 {
						verifyOldTIORemarks = true
					}
				}
				if *h.ObjectID == approvedMove.ID {
					if h.ChangedData != nil {
						changedData := removeEscapeJSONtoObject(h.ChangedData)
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
		suite.True(verifyOldSAC, "verifyOldSAC")
		suite.True(verifyNewSAC, "verifyNewSAC")
		// move update
		suite.True(verifyOldTIORemarks, "verifyOldTIORemarks")
		suite.True(verifyTIORemarks, "verifyTIORemarks")

		suite.True(verifyDBAuthorizedWeight, "verifyDBAuthorizedWeight")
	})

	suite.T().Run("returns not found error for unknown locator", func(t *testing.T) {
		_ = testdatagen.MakeAvailableMove(suite.DB())

		params := services.FetchMoveHistoryParams{Locator: "QX97UY", Page: swag.Int64(1), PerPage: swag.Int64(20)}
		_, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("returns Orders fields and context", func(t *testing.T) {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		order := approvedMove.Orders
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		testdatagen.MakeMTOShipmentWithMove(suite.DB(), &approvedMove, testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move: approvedMove,
		})

		changeOldDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		changeNewDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())

		// Make sure we're testing for all the things that we can update on the Orders page
		// README: This list of properties below here is taken from
		// swagger-def/ghc.yaml#UpdateOrderPayload
		// README: issueDate, reportByDate, ordersType, ordersTypeDetail,
		// originDutyLocationId, newDutyLocationId, ordersNumber, tac, sac,
		// ntsTac, ntsSac, departmentIndicator, ordersAcknowledgement
		orderNumber := "030-00362"
		tac := "1234"
		sac := "2345"
		ntsTac := "3456"
		ntsSac := "4567"

		order.IssueDate = now.AddDate(0, 0, 20)
		order.ReportByDate = now.AddDate(0, 0, 25)
		order.OrdersType = internalmessages.OrdersTypeRETIREMENT
		order.OrdersTypeDetail = internalmessages.NewOrdersTypeDetail(internalmessages.OrdersTypeDetailDELAYEDAPPROVAL)
		order.OriginDutyLocationID = &changeOldDutyLocation.ID
		order.OriginDutyLocation = &changeOldDutyLocation
		order.NewDutyLocationID = changeNewDutyLocation.ID
		order.NewDutyLocation = changeNewDutyLocation
		order.OrdersNumber = &orderNumber
		order.TAC = &tac
		order.SAC = &sac
		order.NtsTAC = &ntsTac
		order.NtsSAC = &ntsSac
		order.DepartmentIndicator = (*string)(internalmessages.NewDeptIndicator(internalmessages.DeptIndicatorARMY))
		order.AmendedOrdersAcknowledgedAt = &now
		// this is gathered on the customer flow
		rank := string(models.ServiceMemberRankE9SPECIALSENIORENLISTED)
		order.Grade = &rank

		suite.MustSave(&order)

		parameters := services.FetchMoveHistoryParams{
			Locator: approvedMove.Locator,
			Page:    swag.Int64(1),
			PerPage: swag.Int64(20),
		}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &parameters)
		suite.FatalNoError(err)

		foundUpdateOrderRecord := false
		for _, historyRecord := range moveHistoryData.AuditHistories {
			if *historyRecord.ObjectID == order.ID && historyRecord.Action == "UPDATE" {
				changedData := removeEscapeJSONtoObject(historyRecord.ChangedData)
				// Date format here: https://go.dev/src/time/format.go
				suite.Equal(order.IssueDate.Format("2006-01-02"), changedData["issue_date"])
				suite.Equal(order.ReportByDate.Format("2006-01-02"), changedData["report_by_date"])
				suite.Equal(string(order.OrdersType), changedData["orders_type"])
				suite.Equal((string)(*order.OrdersTypeDetail), changedData["orders_type_detail"])
				suite.Equal(order.OriginDutyLocationID.String(), changedData["origin_duty_location_id"])
				suite.Equal(order.NewDutyLocationID.String(), changedData["new_duty_location_id"])
				suite.Equal(*order.OrdersNumber, changedData["orders_number"])
				suite.Equal(*order.TAC, changedData["tac"])
				suite.Equal(*order.SAC, changedData["sac"])
				suite.Equal(*order.NtsTAC, changedData["nts_tac"])
				suite.Equal(*order.NtsSAC, changedData["nts_sac"])
				suite.Equal(*order.DepartmentIndicator, changedData["department_indicator"])
				//TODO: make a better comparison test for time
				//suite.Equal(order.AmendedOrdersAcknowledgedAt.Format("2006-01-02T15:04:05.000000"), changedData["amended_orders_acknowledged_at"])

				// rank/grade is also on orders
				suite.Equal(*order.Grade, changedData["grade"])

				// test context as well
				context := removeEscapeJSONtoArray(historyRecord.Context)[0]
				suite.Equal(order.OriginDutyLocation.Name, context["origin_duty_location_name"])
				suite.Equal(order.NewDutyLocation.Name, context["new_duty_location_name"])

				foundUpdateOrderRecord = true
				break
			}
		}

		// double check that we found the record we're looking for
		suite.True(foundUpdateOrderRecord)
	})

}

func removeEscapeJSONtoObject(data *string) map[string]string {
	var result map[string]string
	if data == nil || *data == "" {
		return result
	}
	var byteData = []byte(*data)

	_ = json.Unmarshal(byteData, &result)
	return result
}

func removeEscapeJSONtoArray(data *string) []map[string]string {
	var result []map[string]string
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

	suite.T().Run("filters shipments and service items from different move ", func(t *testing.T) {

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
		serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: approvedMove,
		})

		approvedMoveToFilter := testdatagen.MakeAvailableMove(suite.DB())
		approvedShipmentToFilter := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &approvedMoveToFilter, testdatagen.Assertions{})
		serviceItemToFilter := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: approvedMoveToFilter,
		})

		reason := "heavy"
		serviceItem.Reason = &reason
		suite.MustSave(&serviceItem)

		reasonFilter := "light"
		serviceItemToFilter.Reason = &reasonFilter
		suite.MustSave(&serviceItemToFilter)

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

		suite.Condition(auditHistoryContains(moveHistoryData.AuditHistories, "fragile"), "should contain fragile")
		containsSturdy := auditHistoryContains(moveHistoryData.AuditHistories, "sturdy")()
		suite.False(containsSturdy, "should not contain sturdy")

		suite.Condition(auditHistoryContains(moveHistoryData.AuditHistories, "heavy"), "should contain heavy")
		containsLight := auditHistoryContains(moveHistoryData.AuditHistories, "light")()
		suite.False(containsLight, "should not contain light")

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
						context := removeEscapeJSONtoArray(h.Context)
						if context != nil && context[0]["name"] == serviceItem.ReService.Name && context[0]["shipment_type"] == string(serviceItem.MTOShipment.ShipmentType) {
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
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

	})

	suite.T().Run("has audit history records for reweighs", func(t *testing.T) {
		shipment := testdatagen.MakeMTOShipmentWithMove(suite.DB(), nil, testdatagen.Assertions{})
		// Create a valid reweigh for the move
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterTOO,
			Shipment:    shipment,
			ShipmentID:  shipment.ID,
		}
		reweighCreator := reweigh.NewReweighCreator()
		createdReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)
		suite.NoError(err)

		params := services.FetchMoveHistoryParams{Locator: shipment.MoveTaskOrder.Locator, Page: swag.Int64(1), PerPage: swag.Int64(5)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		verifyReweighHistoryFound := false
		verifyReweighContext := false

		for _, h := range moveHistoryData.AuditHistories {
			if h.TableName == "reweighs" && *h.ObjectID == createdReweigh.ID {
				verifyReweighHistoryFound = true
				if h.Context != nil {
					context := removeEscapeJSONtoArray(h.Context)
					if context != nil && context[0]["shipment_type"] == string(shipment.ShipmentType) {
						verifyReweighContext = true
					}
				}
				break
			}
		}
		suite.True(verifyReweighHistoryFound, "AuditHistories contains an AuditHistory with a Reweigh creation")
		suite.True(verifyReweighContext, "Reweigh creation AuditHistory contains a context with the appropriate shipment type")
	})
}
