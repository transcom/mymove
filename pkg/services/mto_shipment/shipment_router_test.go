package mtoshipment

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type routerApproveSubtestData struct {
	appContext      appcontext.AppContext
	shipmentRouter  services.ShipmentRouter
	unsavedShipment models.MTOShipment
}

func (suite *MTOShipmentServiceSuite) createRouterApproveSubtestData() (subtestData *routerApproveSubtestData) {
	subtestData = &routerApproveSubtestData{}

	subtestData.shipmentRouter = NewShipmentRouter()

	subtestData.unsavedShipment = testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		Stub: true,
	})

	subtestData.appContext = suite.AppContextForTest()

	return subtestData
}

func (suite *MTOShipmentServiceSuite) TestApprove() {
	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Submitted", models.MTOShipmentStatusSubmitted},
		{"Diversion Requested", models.MTOShipmentStatusDiversionRequested},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			subtestData := suite.createRouterApproveSubtestData()

			shipment := subtestData.unsavedShipment

			shipment.Status = validStatus.status
			// special case for diversion requested
			shipment.Diversion = true

			err := subtestData.shipmentRouter.Approve(subtestData.appContext, &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
			suite.NotNil(shipment.ApprovedDate)
		})
	}

	invalidStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusApproved},
		{"Draft", models.MTOShipmentStatusDraft},
		{"Canceled", models.MTOShipmentStatusCanceled},
		{"Rejected", models.MTOShipmentStatusRejected},
		{"Cancellation Requested", models.MTOShipmentStatusCancellationRequested},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			subtestData := suite.createRouterApproveSubtestData()

			shipment := subtestData.unsavedShipment

			shipment.Status = invalidStatus.status

			err := subtestData.shipmentRouter.Approve(subtestData.appContext, &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status 'APPROVED' from [\"SUBMITTED\" \"DIVERSION_REQUESTED\"]", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}

	suite.Run("does not approve a shipment if the move is not Approved or Approvals Requested", func() {
		subtestData := suite.createRouterApproveSubtestData()

		submittedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{Stub: true})

		err := subtestData.shipmentRouter.Approve(subtestData.appContext, &submittedShipment)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "Cannot approve a shipment if the move isn't approved")
	})
}

func (suite *MTOShipmentServiceSuite) TestSubmit() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.Submit(suite.AppContextForTest(), &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusSubmitted, shipment.Status)
		})
	}

	invalidStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Canceled", models.MTOShipmentStatusCanceled},
		{"Rejected", models.MTOShipmentStatusRejected},
		{"Cancellation Requested", models.MTOShipmentStatusCancellationRequested},
		{"Diversion Requested", models.MTOShipmentStatusDiversionRequested},
		{"Approved", models.MTOShipmentStatusApproved},
		{"Submitted", models.MTOShipmentStatusSubmitted},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			shipment.Status = invalidStatus.status

			err := shipmentRouter.Submit(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status 'SUBMITTED' from [\"DRAFT\"]", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestCancel() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Cancellation Requested", models.MTOShipmentStatusCancellationRequested},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.Cancel(suite.AppContextForTest(), &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusCanceled, shipment.Status)
		})
	}

	invalidStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Canceled", models.MTOShipmentStatusCanceled},
		{"Rejected", models.MTOShipmentStatusRejected},
		{"Diversion Requested", models.MTOShipmentStatusDiversionRequested},
		{"Approved", models.MTOShipmentStatusApproved},
		{"Submitted", models.MTOShipmentStatusSubmitted},
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			shipment.Status = invalidStatus.status

			err := shipmentRouter.Cancel(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestReject() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter()
	rejectionReason := "reason"

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Submitted", models.MTOShipmentStatusSubmitted},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.Reject(suite.AppContextForTest(), &shipment, &rejectionReason)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusRejected, shipment.Status)
			suite.Equal(&rejectionReason, shipment.RejectionReason)
		})
	}

	invalidStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Canceled", models.MTOShipmentStatusCanceled},
		{"Rejected", models.MTOShipmentStatusRejected},
		{"Diversion Requested", models.MTOShipmentStatusDiversionRequested},
		{"Approved", models.MTOShipmentStatusApproved},
		{"Cancellation Requested", models.MTOShipmentStatusCancellationRequested},
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			shipment.Status = invalidStatus.status

			err := shipmentRouter.Reject(suite.AppContextForTest(), &shipment, &rejectionReason)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestRequestDiversion() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusApproved},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.RequestDiversion(suite.AppContextForTest(), &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusDiversionRequested, shipment.Status)
		})
	}

	invalidStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Canceled", models.MTOShipmentStatusCanceled},
		{"CANCELLATION_REQUESTED", models.MTOShipmentStatusCancellationRequested},
		{"Rejected", models.MTOShipmentStatusRejected},
		{"Diversion Requested", models.MTOShipmentStatusDiversionRequested},
		{"Submitted", models.MTOShipmentStatusSubmitted},
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			shipment.Status = invalidStatus.status

			err := shipmentRouter.RequestDiversion(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestApproveDiversion() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter()

	suite.Run("fails when the Diversion field is false", func() {
		err := shipmentRouter.ApproveDiversion(suite.AppContextForTest(), &shipment)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "Cannot approve the diversion")
	})

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusSubmitted},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status
			shipment.Diversion = true

			err := shipmentRouter.ApproveDiversion(suite.AppContextForTest(), &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
		})
	}

	invalidStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Canceled", models.MTOShipmentStatusCanceled},
		{"CANCELLATION_REQUESTED", models.MTOShipmentStatusCancellationRequested},
		{"Rejected", models.MTOShipmentStatusRejected},
		{"Diversion Requested", models.MTOShipmentStatusApproved},
		{"Submitted", models.MTOShipmentStatusDiversionRequested},
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			shipment.Status = invalidStatus.status
			shipment.Diversion = true

			err := shipmentRouter.ApproveDiversion(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestRequestCancellation() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusApproved},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.RequestCancellation(suite.AppContextForTest(), &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusCancellationRequested, shipment.Status)
		})
	}

	invalidStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Canceled", models.MTOShipmentStatusCanceled},
		{"CANCELLATION_REQUESTED", models.MTOShipmentStatusCancellationRequested},
		{"Rejected", models.MTOShipmentStatusRejected},
		{"Diversion Requested", models.MTOShipmentStatusDiversionRequested},
		{"Submitted", models.MTOShipmentStatusSubmitted},
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			shipment.Status = invalidStatus.status

			err := shipmentRouter.RequestCancellation(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}
