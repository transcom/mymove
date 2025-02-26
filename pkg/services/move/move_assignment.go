package move

import (
	"container/list"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
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

	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// while we have a queue...
		for moveIndex < len(movesToAssign) && queue.Len() > 0 {
			// grab that ID off the front
			user := queue.Front()
			userID := user.Value.(uuid.UUID)
			queue.Remove(user)

			// do our assignment logic
			move := movesToAssign[moveIndex]
			switch queueType {
			case string(models.QueueTypeCounseling):
				move.SCAssignedID = &userID
			case string(models.QueueTypeCloseout):
				move.SCAssignedID = &userID
			case string(models.QueueTypeTaskOrder):
				move.TOOAssignedID = &userID
			case string(models.QueueTypePaymentRequest):
				move.TIOAssignedID = &userID
			}

			verrs, err := appCtx.DB().ValidateAndUpdate(&move)
			if err != nil || verrs.HasAny() {
				return apperror.NewInvalidInputError(move.ID, err, verrs, "")
			}

			// decrement the users assignment count
			moveAssignments[userID]--
			// increment our index
			moveIndex++

			// If user still has remaining assignments, re-queue them
			if moveAssignments[userID] > 0 {
				queue.PushBack(userID)
			}
		}

		return nil
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return nil, nil
}
