package move

import (
	"container/list"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveAssigner struct {
}

func NewMoveAssignerBulkAssignment() services.MoveAssigner {
	return &moveAssigner{}
}

func (a moveAssigner) BulkMoveAssignment(appCtx appcontext.AppContext, queueType string, officeUserData []*ghcmessages.BulkAssignmentForUser, movesToAssign models.Moves) (*models.Moves, error) {
	if len(movesToAssign) == 0 {
		return nil, apperror.NewBadDataError("No moves to assign")
	}

	var assign func(*models.Move, uuid.UUID)
	switch queueType {
	case string(models.QueueTypeCounseling):
		assign = func(move *models.Move, userID uuid.UUID) { move.SCCounselingAssignedID = &userID }
	case string(models.QueueTypeCloseout):
		assign = func(move *models.Move, userID uuid.UUID) { move.SCCloseoutAssignedID = &userID }
	case string(models.QueueTypeTaskOrder):
		assign = func(move *models.Move, userID uuid.UUID) { move.TOOTaskOrderAssignedID = &userID }
	case string(models.QueueTypeDestinationRequest):
		assign = func(move *models.Move, userID uuid.UUID) { move.TOODestinationAssignedID = &userID }
	case string(models.QueueTypePaymentRequest):
		assign = func(move *models.Move, userID uuid.UUID) { move.TIOPaymentRequestAssignedID = &userID }
	default:
		return nil, apperror.NewBadDataError("Invalid queue type")
	}

	// make a map to track users and their assignment counts
	// and a queue of userIDs
	moveAssignments := make(map[uuid.UUID]int)
	queue := list.New()
	for _, user := range officeUserData {
		if user != nil && user.MoveAssignments > 0 {
			userID := uuid.FromStringOrNil(user.ID.String())
			moveAssignments[userID] = int(user.MoveAssignments)
			queue.PushBack(userID)
		}
	}

	// point at the index in the movesToAssign set
	moveIndex := 0

	// keep track of the updatedMoves to batch save
	updatedMoves := make([]models.Move, 0, len(movesToAssign))

	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// while we have a queue...
		for moveIndex < len(movesToAssign) && queue.Len() > 0 {
			// grab that ID off the front
			user := queue.Front()
			userID := user.Value.(uuid.UUID)
			queue.Remove(user)

			// do our assignment logic
			move := movesToAssign[moveIndex]
			ordersType := move.Orders.OrdersType
			if ordersType != internalmessages.OrdersTypeSAFETY && ordersType != internalmessages.OrdersTypeBLUEBARK && ordersType != internalmessages.OrdersTypeWOUNDEDWARRIOR {
				assign(&move, userID)
				updatedMoves = append(updatedMoves, move)
			}

			// decrement the user's assignment count
			moveAssignments[userID]--
			moveIndex++

			// If user still has remaining assignments, re-queue them
			if moveAssignments[userID] > 0 {
				queue.PushBack(userID)
			}
		}

		return nil
	})

	if len(updatedMoves) > 0 {
		verrs, err := appCtx.DB().ValidateAndUpdate(updatedMoves) // Bulk update
		if err != nil || verrs.HasAny() {
			return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Bulk assignment failed")
		}
	}

	if transactionErr != nil {
		return nil, transactionErr
	}

	return nil, nil
}
