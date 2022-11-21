package movehistory

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/reweigh"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MoveHistoryServiceSuite) TestMoveFetcher() {
	moveHistoryFetcher := NewMoveHistoryFetcher()

	suite.Run("successfully returns submitted move history available to prime", func() {

		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		now := time.Now()
		pickupDate := now.AddDate(0, 0, 10)
		secondaryPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
		approvedShipment := testdatagen.MakeMTOShipmentWithMove(suite.DB(), &approvedMove, testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        &now,
				ScheduledPickupDate: &pickupDate,
			},
			Move:                   approvedMove,
			SecondaryPickupAddress: secondaryPickupAddress,
		})

		testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				FirstName:    swag.String("Test1"),
				LastName:     swag.String("Agent"),
				Email:        swag.String("test@test.email.com"),
				MTOAgentType: models.MTOAgentReceiving,
			},
			MTOShipment: approvedShipment,
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

		// update Secondary Pickup Address
		oldSecondaryPickupAddress := *approvedShipment.SecondaryPickupAddress
		updateSecondaryPickupAddress := approvedShipment.SecondaryPickupAddress
		updateSecondaryPickupAddress.City = "Hampton"
		updateSecondaryPickupAddress.State = "VA"
		updateSecondaryPickupAddress.PostalCode = "23661"
		suite.MustSave(updateSecondaryPickupAddress)

		// update move
		tioRemarks := "updating TIO remarks for test"
		approvedMove.TIORemarks = &tioRemarks
		suite.MustSave(&approvedMove)

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistory, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.FatalNoError(err)

		// address update
		verifyOldPickupAddress := false
		verifyNewPickupAddress := false
		verifyOldSecondaryPickupAddress := false
		verifyNewSecondaryPickupAddress := false
		// agent update
		verifyNewAgent := false
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
				} else if *h.ObjectID == updateSecondaryPickupAddress.ID {
					if h.OldData != nil {
						oldData := removeEscapeJSONtoObject(h.OldData)
						if oldData["city"] == oldSecondaryPickupAddress.City && oldData["state"] == oldSecondaryPickupAddress.State && oldData["postal_code"] == oldSecondaryPickupAddress.PostalCode {
							verifyOldSecondaryPickupAddress = true
						}
					}
					if h.ChangedData != nil {
						changedData := removeEscapeJSONtoObject(h.ChangedData)
						if changedData["city"] == updateSecondaryPickupAddress.City && changedData["state"] == updateSecondaryPickupAddress.State && changedData["postal_code"] == updateSecondaryPickupAddress.PostalCode {
							verifyNewSecondaryPickupAddress = true
						}
					}
				}
			} else if h.TableName == "orders" {
				if *h.ObjectID == approvedMove.Orders.ID {
					if h.OldData != nil {
						oldData := removeEscapeJSONtoObject(h.OldData)
						if sac, ok := oldData["sac"]; !ok || sac == nil {
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
			} else if h.TableName == "mto_agents" {
				if h.ChangedData != nil {
					changedData := removeEscapeJSONtoObject(h.ChangedData)
					if changedData["agent_type"] == string(models.MTOAgentReceiving) {
						verifyNewAgent = true
					}
				}
			} else if h.TableName == "entitlements" {
				if h.ChangedData != nil {
					oldData := removeEscapeJSONtoObject(h.OldData)
					if authorizedWeight, ok := oldData["authorized_weight"]; !ok || authorizedWeight == nil {
						verifyDBAuthorizedWeight = true
					}
				}
			} else if h.TableName == "moves" {
				if h.OldData != nil {
					oldData := removeEscapeJSONtoObject(h.OldData)
					if tioRemarks, ok := oldData["tio_remarks"]; !ok || tioRemarks == nil {
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
		// secondary address update
		suite.True(verifyOldSecondaryPickupAddress, "verifyOldSecondaryPickupAddress")
		suite.True(verifyNewSecondaryPickupAddress, "verifyNewSecondaryPickupAddress")
		// agent update
		suite.True(verifyNewAgent, "verifyNewAgent")
		// orders update
		suite.True(verifyOldSAC, "verifyOldSAC")
		suite.True(verifyNewSAC, "verifyNewSAC")
		// move update
		suite.True(verifyOldTIORemarks, "verifyOldTIORemarks")
		suite.True(verifyTIORemarks, "verifyTIORemarks")

		suite.True(verifyDBAuthorizedWeight, "verifyDBAuthorizedWeight")
	})

	suite.Run("returns not found error for unknown locator", func() {
		_ = testdatagen.MakeAvailableMove(suite.DB())

		params := services.FetchMoveHistoryParams{Locator: "QX97UY", Page: swag.Int64(1), PerPage: swag.Int64(100)}
		_, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("returns Orders fields and context", func() {
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
			PerPage: swag.Int64(100),
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

				//changedData["amended_orders_acknowledged_at"] is being converted from a string to a formatted time.Time type
				layout := "2006-01-02T15:04:05.000000"
				changedDataTimeStamp, _ := time.Parse(layout, changedData["amended_orders_acknowledged_at"].(string))
				//CircleCi seems to add on nanoseconds to the tested time stamps so this is being used with Truncate to shave those nanoseconds off
				d := 1000 * time.Nanosecond
				//We assert if it falls within a range starting at the original order.AmendedOrdersAcknowledgedAt time and ending with a added 2000 microsecond buffer
				suite.WithinRange(changedDataTimeStamp, order.AmendedOrdersAcknowledgedAt.Truncate(d), order.AmendedOrdersAcknowledgedAt.Add(2000*time.Microsecond).Truncate(d))

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
	suite.Run("returns user uploads fields and context", func() {
		// Make an approved move and get the associated orders, service member, uploaded orders and related document
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		orders := approvedMove.Orders
		serviceMember := orders.ServiceMember
		uploadedOrdersDocument := orders.UploadedOrders
		userUploadedOrders := uploadedOrdersDocument.UserUploads[0]

		// Create an amended orders that is associated with the service member
		userUploadedAmendedOrders := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMember:   serviceMember,
				ServiceMemberID: serviceMember.ID,
			},
		})

		// Update the orders with the amended orders
		orders.UploadedAmendedOrdersID = &userUploadedAmendedOrders.Document.ID
		orders.UploadedAmendedOrders = &userUploadedAmendedOrders.Document
		suite.MustSave(&orders)

		parameters := services.FetchMoveHistoryParams{
			Locator: approvedMove.Locator,
			Page:    swag.Int64(1),
			PerPage: swag.Int64(100),
		}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &parameters)
		suite.FatalNoError(err)

		foundUserUploadOrdersRecord := false
		foundUserUploadAmendedOrdersRecord := false
		for _, historyRecord := range moveHistoryData.AuditHistories {
			if *historyRecord.ObjectID == userUploadedOrders.ID && historyRecord.Action == "INSERT" {
				context := removeEscapeJSONtoArray(historyRecord.Context)[0]
				suite.Equal(userUploadedOrders.Upload.Filename, context["filename"])
				suite.Equal("orders", context["upload_type"])

				foundUserUploadOrdersRecord = true
			} else if *historyRecord.ObjectID == userUploadedAmendedOrders.ID && historyRecord.Action == "INSERT" {
				context := removeEscapeJSONtoArray(historyRecord.Context)[0]
				suite.Equal(userUploadedAmendedOrders.Upload.Filename, context["filename"])
				suite.Equal("amendedOrders", context["upload_type"])

				foundUserUploadAmendedOrdersRecord = true
			}
		}
		// double check that we found the records we're looking for
		suite.True(foundUserUploadOrdersRecord, "foundUserUploadOrdersRecord")
		suite.True(foundUserUploadAmendedOrdersRecord, "foundUserUploadAmendedOrdersRecord")

	})

}

func removeEscapeJSONtoObject(data *string) map[string]interface{} {
	var result map[string]interface{}
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

	suite.Run("returns Audit History with session information", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		fakeRole, _ := testdatagen.LookupOrMakeRoleByRoleType(suite.DB(), roles.RoleTypeTOO)
		fakeUser := factory.BuildUser(suite.DB(), nil, nil)
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

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
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

	suite.Run("filters shipments and service items from different move ", func() {

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

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
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

	suite.Run("has context", func() {
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

		params := services.FetchMoveHistoryParams{Locator: move.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
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

	suite.Run("has paginated results", func() {
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
		suite.Equal(2, len(moveHistoryData.AuditHistories))
	})

	suite.Run("approved payment request shows up", func() {
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		cents := unit.Cents(1000)
		approvedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusPending,
			},
			Move: approvedMove,
		})

		testServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: approvedMove,
		})

		paymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Name: "Test",
			},
			PaymentServiceItem: models.PaymentServiceItem{
				Status:     models.PaymentServiceItemStatusRequested,
				PriceCents: &cents,
			},
			PaymentRequest: approvedPaymentRequest,
			MTOServiceItem: testServiceItem,
		})

		approvedPaymentRequest.Status = models.PaymentRequestStatusReviewed
		suite.MustSave(&approvedPaymentRequest)
		paymentServiceItem.Status = models.PaymentServiceItemStatusApproved
		suite.MustSave(&paymentServiceItem)

		params := services.FetchMoveHistoryParams{Locator: approvedMove.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		verifyPaymentRequestHistoryFound := false
		for _, h := range moveHistoryData.AuditHistories {
			if h.TableName == "payment_requests" && *h.ObjectID == approvedPaymentRequest.ID {
				if h.Context != nil {
					context := removeEscapeJSONtoArray(h.Context)
					if context != nil {
						suite.Contains(context[0]["status"], paymentServiceItem.Status)
						verifyPaymentRequestHistoryFound = true
					}
				}
				break
			}
		}
		suite.True(verifyPaymentRequestHistoryFound, "AuditHistories contains an AuditHistory with an approved payment request")
	})

	suite.Run("has audit history records for reweighs", func() {
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

		params := services.FetchMoveHistoryParams{Locator: shipment.MoveTaskOrder.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
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

	suite.Run("has audit history records for service item dimensions", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)

		dimension := models.MTOServiceItemDimension{
			Type:      models.DimensionTypeItem,
			Length:    12000,
			Height:    12000,
			Width:     12000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		reServiceDDFSIT := testdatagen.MakeDDFSITReService(suite.DB())

		serviceItem := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService:       reServiceDDFSIT,
			MTOShipmentID:   &shipment.ID,
			MTOShipment:     shipment,
			Dimensions:      models.MTOServiceItemDimensions{dimension},
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItem)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		params := services.FetchMoveHistoryParams{Locator: shipment.MoveTaskOrder.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		verifyServiceItemDimensionsHistoryFound := false

		for _, h := range moveHistoryData.AuditHistories {
			if h.TableName == "mto_service_item_dimensions" {
				if h.ChangedData != nil {
					changedData := removeEscapeJSONtoObject(h.ChangedData)
					if changedData["type"] == "ITEM" {
						verifyServiceItemDimensionsHistoryFound = true
						break
					}
				}
			}
		}
		suite.True(verifyServiceItemDimensionsHistoryFound, "AuditHistories contains an AuditHistory with a service item dimensions creation")
	})

	suite.Run("has audit history records for service item customer contacts", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)

		reService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		})

		sitEntryDate := time.Now()
		contact1 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			FirstAvailableDeliveryDate: sitEntryDate,
			TimeMilitary:               "0815Z",
		}
		contact2 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeSecond,
			FirstAvailableDeliveryDate: sitEntryDate,
			TimeMilitary:               "0815Z",
		}
		var contacts models.MTOServiceItemCustomerContacts
		contacts = append(contacts, contact1, contact2)

		serviceItem := models.MTOServiceItem{
			MoveTaskOrderID:  move.ID,
			MoveTaskOrder:    move,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			CustomerContacts: contacts,
			ReService:        reService,
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItem)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		params := services.FetchMoveHistoryParams{Locator: shipment.MoveTaskOrder.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		verifyServiceItemDimensionsHistoryFound := false

		for _, h := range moveHistoryData.AuditHistories {
			if h.TableName == "mto_service_item_customer_contacts" {
				if h.ChangedData != nil {
					changedData := removeEscapeJSONtoObject(h.ChangedData)
					if changedData["time_military"] == "0815Z" {
						verifyServiceItemDimensionsHistoryFound = true
						break
					}
				}
			}
		}
		suite.True(verifyServiceItemDimensionsHistoryFound, "AuditHistories contains an AuditHistory with a service item customer contacts creation")
	})

	suite.Run("has audit history records for mto_agents", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		mtoAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		// Make two audit history entries, one with an event name we should find
		// and another with the eventName we are intentionally not returning in our query.
		eventNameToFind := "updateShipment"
		eventNameToNotFind := "deleteShipment"
		tableName := "mto_agents"
		testdatagen.MakeAuditHistory(suite.DB(), testdatagen.Assertions{
			TestDataAuditHistory: testdatagen.TestDataAuditHistory{
				EventName:   &eventNameToFind,
				TableNameDB: tableName,
				ObjectID:    &mtoAgent.ID,
			},
			Move: models.Move{
				ID: move.ID,
			},
		})
		testdatagen.MakeAuditHistory(suite.DB(), testdatagen.Assertions{
			TestDataAuditHistory: testdatagen.TestDataAuditHistory{
				EventName:   &eventNameToNotFind,
				TableNameDB: tableName,
				ObjectID:    &mtoAgent.ID,
			},
			Move: models.Move{
				ID: move.ID,
			},
		})
		params := services.FetchMoveHistoryParams{Locator: move.Locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistoryData, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.NotNil(moveHistoryData)
		suite.NoError(err)

		verifyEventNameFound := false
		verifyEventNameNotFound := false

		for _, h := range moveHistoryData.AuditHistories {
			if h.TableName == "mto_agents" {
				if *h.EventName == eventNameToFind {
					verifyEventNameFound = true
				}
				if *h.EventName == eventNameToNotFind {
					verifyEventNameNotFound = true
				}
			}
		}
		suite.True(verifyEventNameFound, "MTO Agent event name to find.")
		suite.False(verifyEventNameNotFound, "MTO Agent event name to NOT find.")
	})
}

func (suite *MoveHistoryServiceSuite) TestMoveFetcherUserInfo() {
	moveHistoryFetcher := NewMoveHistoryFetcher()

	setupTestData := func(userID *uuid.UUID, userFirstName string, roleTypes []roles.RoleType, isOfficeUser bool) string {
		assertions := testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				FirstName: userFirstName,
			},
			User: models.User{
				ID: *userID,
			},
		}

		var user models.User
		if isOfficeUser {
			officeUser := testdatagen.MakeOfficeUserWithRoleTypes(suite.DB(), roleTypes, assertions)
			user = officeUser.User
		} else {
			user = testdatagen.MakeUserWithRoleTypes(suite.DB(), roleTypes, assertions)
		}
		approvedMove := testdatagen.MakeAvailableMove(suite.DB())
		testdatagen.MakeAuditHistory(suite.DB(), testdatagen.Assertions{
			User: user,
			Move: models.Move{
				ID: approvedMove.ID,
			},
		})
		return approvedMove.Locator
	}

	setupServiceMemberTestData := func(userFirstName string, fakeEventName string) (string, models.User) {
		assertions := testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				FirstName: &userFirstName,
			},
		}
		// Create an unsubmitted move with the service member attached to the orders.
		move := testdatagen.MakeMove(suite.DB(), assertions)
		user := move.Orders.ServiceMember.User
		testdatagen.MakeAuditHistory(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				ID: move.ID,
			},
			User: user,
			TestDataAuditHistory: testdatagen.TestDataAuditHistory{
				EventName: &fakeEventName,
			},
		})
		return move.Locator, user
	}

	suite.Run("Test with TOO user", func() {
		userID, _ := uuid.NewV4()
		userName := "TOO_user"
		locator := setupTestData(&userID, userName, []roles.RoleType{roles.RoleTypeTOO}, true)
		params := services.FetchMoveHistoryParams{Locator: locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistory, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Nil(err)
		auditHistoriesForUser := filterAuditHistoryByUserID(moveHistory.AuditHistories, userID)
		suite.Equal(1, len(auditHistoriesForUser))
		suite.Equal(userName, *auditHistoriesForUser[0].SessionUserFirstName)
	})

	suite.Run("Test with Prime user", func() {
		userID, _ := uuid.NewV4()
		userName := "Prime_user"
		locator := setupTestData(&userID, userName, []roles.RoleType{roles.RoleTypePrime}, false)
		params := services.FetchMoveHistoryParams{Locator: locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistory, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Nil(err)
		auditHistoriesForUser := filterAuditHistoryByUserID(moveHistory.AuditHistories, userID)
		suite.Equal(1, len(auditHistoriesForUser))
		suite.Equal("Prime", *auditHistoriesForUser[0].SessionUserFirstName)
	})

	suite.Run("Test with TOO and Prime Simulator user", func() {
		userID, _ := uuid.NewV4()
		userName := "TOO_and_prime_simulator_user"
		locator := setupTestData(&userID, userName, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypePrimeSimulator}, true)
		params := services.FetchMoveHistoryParams{Locator: locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistory, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Nil(err)
		auditHistoriesForUser := filterAuditHistoryByUserID(moveHistory.AuditHistories, userID)
		suite.Equal(1, len(auditHistoriesForUser))
		suite.Equal(userName, *auditHistoriesForUser[0].SessionUserFirstName)
	})

	suite.Run("Test with TOO and Customer user", func() {
		userID, _ := uuid.NewV4()
		userName := "TOO_and_customer_user"
		locator := setupTestData(&userID, userName, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeCustomer}, true)
		params := services.FetchMoveHistoryParams{Locator: locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistory, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Nil(err)
		auditHistoriesForUser := filterAuditHistoryByUserID(moveHistory.AuditHistories, userID)
		suite.Equal(1, len(auditHistoriesForUser))
		suite.Equal(userName, *auditHistoriesForUser[0].SessionUserFirstName)
	})

	suite.Run("Test with Service Member user", func() {
		userName := "service_member_creator"
		fakeEventName := "submitMoveForApproval"
		locator, user := setupServiceMemberTestData(userName, fakeEventName)
		params := services.FetchMoveHistoryParams{Locator: locator, Page: swag.Int64(1), PerPage: swag.Int64(100)}
		moveHistory, _, err := moveHistoryFetcher.FetchMoveHistory(suite.AppContextForTest(), &params)
		suite.Nil(err)
		auditHistoriesForUser := filterAuditHistoryByUserID(moveHistory.AuditHistories, user.ID)
		suite.Equal(1, len(auditHistoriesForUser))
		suite.Equal(userName, *auditHistoriesForUser[0].SessionUserFirstName)
		suite.Equal(fakeEventName, *auditHistoriesForUser[0].EventName)
	})
}

func filterAuditHistoryByUserID(auditHistories models.AuditHistories, userID uuid.UUID) models.AuditHistories {
	auditHistoriesForUser := models.AuditHistories{}
	for _, auditHistory := range auditHistories {
		if auditHistory.SessionUserID != nil && *auditHistory.SessionUserID == userID {
			auditHistoriesForUser = append(auditHistoriesForUser, auditHistory)
		}
	}
	return auditHistoriesForUser
}
