package mtoshipment

import (
	"fmt"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestApprove() {
	shipmentRouter := NewShipmentRouter()

	setUpTestData := func(overrides testdatagen.Assertions) models.MTOShipment {
		fullAssertions := testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
			Stub: true,
		}

		// Merge the overrides into the base assertions
		testdatagen.MergeModels(&fullAssertions, overrides)

		return testdatagen.MakeMTOShipment(suite.DB(), fullAssertions)
	}

	validStatuses := []models.MTOShipmentStatus{
		models.MTOShipmentStatusSubmitted,
		models.MTOShipmentStatusDiversionRequested,
	}

	for _, validStatus := range validStatuses {
		validStatus := validStatus

		suite.Run("from valid status: "+string(validStatus), func() {
			overrides := testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					Status: validStatus,
				},
			}

			// special case for diversion requested
			if validStatus == models.MTOShipmentStatusDiversionRequested {
				overrides.MTOShipment.Diversion = true
			}

			shipment := setUpTestData(overrides)

			err := shipmentRouter.Approve(suite.AppContextForTest(), &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
			suite.NotNil(shipment.ApprovedDate)
		})
	}

	invalidStatuses := []models.MTOShipmentStatus{
		models.MTOShipmentStatusApproved,
		models.MTOShipmentStatusDraft,
		models.MTOShipmentStatusCanceled,
		models.MTOShipmentStatusRejected,
		models.MTOShipmentStatusCancellationRequested,
	}
	for _, invalidStatus := range invalidStatuses {
		invalidStatus := invalidStatus

		suite.Run("from invalid status: "+string(invalidStatus), func() {
			shipment := setUpTestData(testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					Status: invalidStatus,
				},
			})

			err := shipmentRouter.Approve(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status 'APPROVED' from [\"SUBMITTED\" \"DIVERSION_REQUESTED\"]", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus))
		})
	}

	invalidMoveStatuses := []models.MoveStatus{
		models.MoveStatusDRAFT,
		models.MoveStatusSUBMITTED,
		models.MoveStatusCANCELED,
		models.MoveStatusNeedsServiceCounseling,
		models.MoveStatusServiceCounselingCompleted,
	}

	for _, invalidMoveStatus := range invalidMoveStatuses {
		invalidMoveStatus := invalidMoveStatus

		suite.Run(fmt.Sprintf("Doesn't approve a shipment if the move status is %s", invalidMoveStatus), func() {
			move := testdatagen.MakeStubbedMoveWithStatus(suite.DB(), invalidMoveStatus)

			overrides := testdatagen.Assertions{
				Move: move,
				MTOShipment: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			}

			shipment := setUpTestData(overrides)

			err := shipmentRouter.Approve(suite.AppContextForTest(), &shipment)

			if suite.Error(err) {
				suite.IsType(apperror.ConflictError{}, err)
				suite.Contains(
					err.Error(),
					fmt.Sprintf(
						"Cannot approve a shipment if the move status isn't %s or %s, or if it isn't a PPM shipment with a move status of %s. The current status for the move with ID %s is %s",
						models.MoveStatusAPPROVED,
						models.MoveStatusAPPROVALSREQUESTED,
						models.MoveStatusNeedsServiceCounseling,
						move.ID,
						move.Status,
					),
				)
			}
		})
	}

	suite.Run(fmt.Sprintf("can approve a shipment if it is a PPM shipment and the move status is %s", models.MoveStatusNeedsServiceCounseling), func() {
		move := testdatagen.MakeStubbedMoveWithStatus(suite.DB(), models.MoveStatusNeedsServiceCounseling)

		overrides := testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
			Stub: true,
		}

		ppmShipment := testdatagen.MakePPMShipment(suite.DB(), overrides)

		err := shipmentRouter.Approve(suite.AppContextForTest(), &ppmShipment.Shipment)

		if suite.NoError(err) {
			suite.Equal(models.MTOShipmentStatusApproved, ppmShipment.Shipment.Status)
			suite.NotNil(ppmShipment.Shipment.ApprovedDate)
		}
	})

	suite.Run("does not approve a shipment if the shipment uses an external vendor", func() {
		shipment := setUpTestData(testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				UsesExternalVendor: true,
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				Status:             models.MTOShipmentStatusSubmitted,
			},
		})

		err := shipmentRouter.Approve(suite.AppContextForTest(), &shipment)

		suite.Contains(err.Error(), "cannot approve a shipment if it uses an external vendor")
		suite.Equal(models.MTOShipmentStatusSubmitted, shipment.Status)
		suite.Nil(shipment.ApprovedDate)
	})
}

func (suite *MTOShipmentServiceSuite) TestSubmit() {

	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Blank", models.MTOShipmentStatus("")},
		{"Draft", models.MTOShipmentStatusDraft},
	}
	for _, validStatus := range validStatuses {
		validStatus := validStatus

		suite.Run("from valid status: "+string(validStatus.desc), func() {
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
			shipment.Status = validStatus.status

			err := shipmentRouter.Submit(suite.AppContextForTest(), &shipment)

			suite.NoError(err)
			suite.Equal(models.MTOShipmentStatusSubmitted, shipment.Status)
		})
	}

	invalidStatuses := []models.MTOShipmentStatus{
		models.MTOShipmentStatusCanceled,
		models.MTOShipmentStatusRejected,
		models.MTOShipmentStatusCancellationRequested,
		models.MTOShipmentStatusDiversionRequested,
		models.MTOShipmentStatusApproved,
		models.MTOShipmentStatusSubmitted,
	}
	for _, invalidStatus := range invalidStatuses {
		invalidStatus := invalidStatus

		suite.Run("from invalid status: "+string(invalidStatus), func() {
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
			shipment.Status = invalidStatus

			err := shipmentRouter.Submit(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status 'SUBMITTED' from [\"DRAFT\"]", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus))
		})
	}
}

func (suite *MTOShipmentServiceSuite) TestCancel() {

	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Cancellation Requested", models.MTOShipmentStatusCancellationRequested},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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

	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusApproved},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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

	shipmentRouter := NewShipmentRouter()

	suite.Run("fails when the Diversion field is false", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
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
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
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

func (suite *MTOShipmentServiceSuite) TestApproveDiversionUsesExternal() {

	shipmentRouter := NewShipmentRouter()

	suite.Run("fails when the UsesExternal field is true", func() {

		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		shipment.UsesExternalVendor = true
		shipment.Diversion = true
		err := shipmentRouter.ApproveDiversion(suite.AppContextForTest(), &shipment)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "has the UsesExternalVendor field set to true")
	})
}

func (suite *MTOShipmentServiceSuite) TestRequestCancellation() {

	shipmentRouter := NewShipmentRouter()

	validStatuses := []struct {
		desc   string
		status models.MTOShipmentStatus
	}{
		{"Approved", models.MTOShipmentStatusApproved},
	}
	for _, validStatus := range validStatuses {
		suite.Run("from valid status: "+string(validStatus.status), func() {
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
			shipment.Status = validStatus.status
			shipment.UsesExternalVendor = true
			shipment.Diversion = true

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
			shipment := testdatagen.MakeStubbedShipment(suite.DB())
			shipment.Status = invalidStatus.status

			err := shipmentRouter.RequestCancellation(suite.AppContextForTest(), &shipment)

			suite.Error(err)
			suite.IsType(ConflictStatusError{}, err)
			suite.Contains(err.Error(), fmt.Sprintf("Shipment with id '%s' can only transition to status", shipment.ID))
			suite.Contains(err.Error(), fmt.Sprintf("but its current status is '%s'", invalidStatus.status))
		})
	}
}
