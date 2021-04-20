package move

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveRouter struct {
	db *pop.Connection
}

// NewMoveRouter creates a new moveRouter service
func NewMoveRouter(db *pop.Connection) services.MoveRouter {
	return &moveRouter{db}
}

// Submit is called when the customer submits their move. It determines whether
// to send the move to Service Counseling or directly to the TOO. If it goes to
// Service Counseling, its status becomes "Needs Service Counseling", otherwise,
// "Submitted".
func (router moveRouter) Submit(move *models.Move) error {
	var err error

	if router.needsServiceCounseling() {
		err = router.sendToServiceCounselor(move)
	} else {
		err = router.sendNewMoveToOfficeUser(move)
	}
	if err != nil {
		return err
	}
	return nil
}

// TODO: Replace the code in this function to determine whether or not the move
// needs service counseling based on the service member's origin duty station.
// Then remove all code related to the service counseling feature flag here and
// in pkg/cli/featureflag.go, and remove any references to
// `FEATURE_FLAG_SERVICE_COUNSELING` from the entire project.
// You'll need to update the test setup in TestSubmitMoveForServiceCounselingHandler
// so that the move's origin duty station will trigger service counseling.
func (router moveRouter) needsServiceCounseling() bool {
	logger := zap.NewNop()
	flag := pflag.CommandLine
	v := viper.New()
	err := v.BindPFlags(flag)
	if err != nil {
		logger.Fatal("could not bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return v.GetBool(cli.FeatureFlagServiceCounseling)
}

// sendToServiceCounselor makes the move available for a Service Counselor to review
func (router moveRouter) sendToServiceCounselor(move *models.Move) error {
	if move.Status != models.MoveStatusDRAFT {
		return errors.Wrap(
			models.ErrInvalidTransition, fmt.Sprintf(
				"Cannot move to NeedsServiceCounseling state when the Move is not in Draft status. Its current status is %s",
				move.Status,
			),
		)
	}
	move.Status = models.MoveStatusNeedsServiceCounseling
	now := time.Now()
	move.SubmittedAt = &now

	return nil
}

// sendNewMoveToOfficeUser makes the move available for a TOO to review
// The Submitted status indicates to the TOO that this is a new move.
func (router moveRouter) sendNewMoveToOfficeUser(move *models.Move) error {
	if move.Status != models.MoveStatusDRAFT {
		return errors.Wrap(models.ErrInvalidTransition, "Submit")
	}
	move.Status = models.MoveStatusSUBMITTED
	now := time.Now()
	move.SubmittedAt = &now

	// Update PPM status too
	for i := range move.PersonallyProcuredMoves {
		ppm := &move.PersonallyProcuredMoves[i]
		err := ppm.Submit(now)
		if err != nil {
			return err
		}
	}

	for _, ppm := range move.PersonallyProcuredMoves {
		if ppm.Advance != nil {
			err := ppm.Advance.Request()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Approve makes the Move available to the Prime. The Prime cannot create
// Service Items unless the Move is approved.
func (router moveRouter) Approve(move *models.Move) error {
	if router.approvable(move) {
		move.Status = models.MoveStatusAPPROVED
		return nil
	}
	if router.alreadyApproved(move) {
		return nil
	}
	return errors.Wrap(
		models.ErrInvalidTransition, fmt.Sprintf(
			"A move can only be approved if it's in one of these states: %q. However, its current status is: %s",
			validStatusesBeforeApproval, move.Status,
		),
	)
}

func (router moveRouter) alreadyApproved(move *models.Move) bool {
	return move.Status == models.MoveStatusAPPROVED
}

func (router moveRouter) approvable(move *models.Move) bool {
	return statusSliceContains(validStatusesBeforeApproval, move.Status)
}

func statusSliceContains(statusSlice []models.MoveStatus, status models.MoveStatus) bool {
	for _, validStatus := range statusSlice {
		if status == validStatus {
			return true
		}
	}
	return false
}

var validStatusesBeforeApproval = []models.MoveStatus{
	models.MoveStatusSUBMITTED,
	models.MoveStatusAPPROVALSREQUESTED,
	models.MoveStatusServiceCounselingCompleted,
}

// SendToOfficeUserToReviewNewServiceItems sets the moves status to
// "Approvals Requested", which indicates to the TOO that they have new
// service items to review.
func (router moveRouter) SendToOfficeUserToReviewNewServiceItems(move *models.Move) error {
	// Do nothing if it's already in the desired state
	if move.Status == models.MoveStatusAPPROVALSREQUESTED {
		return nil
	}
	if move.Status != models.MoveStatusAPPROVED {
		return errors.Wrap(models.ErrInvalidTransition, fmt.Sprintf("The status for the Move with ID %s can only be set to 'Approvals Requested' from the 'Approved' status, but its current status is %s.", move.ID, move.Status))
	}
	move.Status = models.MoveStatusAPPROVALSREQUESTED
	return nil
}

// Cancel cancels the Move and its associated PPMs
func (router moveRouter) Cancel(reason string, move *models.Move) error {
	// We can cancel any move that isn't already complete.
	if move.Status == models.MoveStatusCANCELED {
		return errors.Wrap(models.ErrInvalidTransition, "Cancel")
	}

	move.Status = models.MoveStatusCANCELED

	// If a reason was submitted, add it to the move record.
	if reason != "" {
		move.CancelReason = &reason
	}

	// This will work only if you use the PPM in question rather than a var representing it
	// i.e. you can't use _, ppm := range PPMs, has to be PPMS[i] as below
	for i := range move.PersonallyProcuredMoves {
		err := move.PersonallyProcuredMoves[i].Cancel()
		if err != nil {
			return err
		}
	}

	// TODO: Orders can exist after related moves are canceled
	err := move.Orders.Cancel()
	if err != nil {
		return err
	}

	return nil

}

// CompleteServiceCounseling sets the move status to "Service Counseling Completed",
// which makes the move available to the TOO. This gets called when the Service
// Counselor is done reviewing the move and submits it.
func (router moveRouter) CompleteServiceCounseling(move *models.Move) error {
	if move.Status != models.MoveStatusNeedsServiceCounseling {
		return errors.Wrap(
			models.ErrInvalidTransition,
			fmt.Sprintf("The status for the Move with ID %s can only be set to 'Service Counseling Completed' from the 'Needs Service Counseling' status, but its current status is %s.",
				move.ID, move.Status,
			),
		)
	}

	now := time.Now()
	move.ServiceCounselingCompletedAt = &now
	move.Status = models.MoveStatusServiceCounselingCompleted

	return nil
}
