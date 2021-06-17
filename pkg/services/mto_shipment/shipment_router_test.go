package mtoshipment

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestApprove() {
	shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		Stub: true,
	})
	shipmentRouter := NewShipmentRouter(suite.DB())

	suite.Nil(shipment.ApprovedDate)

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Submitted", models.MTOShipmentStatusSubmitted},
		{"Diversion Requested", models.MTOShipmentStatusDiversionRequested},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status
			// special case for diversion requested
			shipment.Diversion = true

			err := shipmentRouter.Approve(&shipment)

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
			shipment.Status = invalidStatus.status

			err := shipmentRouter.Approve(&shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status 'APPROVED' from [\"SUBMITTED\" \"DIVERSION_REQUESTED\"]", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}

	suite.Run("does not approve a shipment if the move is not Approved or Approvals Requested", func() {
		submittedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{Stub: true})

		err := shipmentRouter.Approve(&submittedShipment)

		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
		suite.Contains(err.Error(), "Cannot approve a shipment if the move isn't approved")
	})
}

func (suite *MTOShipmentServiceSuite) TestSubmit() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter(suite.DB())

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.Submit(&shipment)

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

			err := shipmentRouter.Submit(&shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status 'SUBMITTED' from [\"DRAFT\"]", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestCancel() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter(suite.DB())

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Cancellation Requested", models.MTOShipmentStatusCancellationRequested},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.Cancel(&shipment)

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

			err := shipmentRouter.Cancel(&shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestReject() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter(suite.DB())
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

			err := shipmentRouter.Reject(&shipment, &rejectionReason)

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

			err := shipmentRouter.Reject(&shipment, &rejectionReason)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestRequestDiversion() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter(suite.DB())

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusApproved},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.RequestDiversion(&shipment)

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

			err := shipmentRouter.RequestDiversion(&shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestApproveDiversion() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter(suite.DB())

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusDiversionRequested},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.ApproveDiversion(&shipment)

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
		{"Submitted", models.MTOShipmentStatusSubmitted},
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, invalidStatus := range invalidStatuses {
		suite.Run("from invalid status: "+string(invalidStatus.status), func() {
			shipment.Status = invalidStatus.status

			err := shipmentRouter.ApproveDiversion(&shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestRequestCancellation() {
	shipment := testdatagen.MakeStubbedShipment(suite.DB())
	shipmentRouter := NewShipmentRouter(suite.DB())

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusApproved},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment.Status = validStatus.status

			err := shipmentRouter.RequestCancellation(&shipment)

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

			err := shipmentRouter.RequestCancellation(&shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}
