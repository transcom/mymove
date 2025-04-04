package move

import (
	"container/list"
	"fmt"

	"github.com/go-openapi/strfmt"
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
	case string(models.QueueTypeCounseling), string(models.QueueTypeCloseout):
		assign = func(move *models.Move, userID uuid.UUID) { move.SCAssignedID = &userID }
	case string(models.QueueTypeTaskOrder):
		assign = func(move *models.Move, userID uuid.UUID) { move.TOOAssignedID = &userID }
	case string(models.QueueTypePaymentRequest):
		assign = func(move *models.Move, userID uuid.UUID) { move.TIOAssignedID = &userID }
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
			fmt.Printf("Processing user: %s with current assignment count: %d\n", userID, moveAssignments[userID])
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
			fmt.Printf("User %s assignment count after decrement: %d\n", userID, moveAssignments[userID])
			moveIndex++

			// If user still has remaining assignments, re-queue them
			if moveAssignments[userID] > 0 {
				queue.PushBack(userID)
				fmt.Printf("User %s requeued with remaining assignments: %d\n", userID, moveAssignments[userID])
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

// get the incoming office users, and the user being re-assigned from
// get the moves related the the reassigned user
// find all moves related to the user being reassigned
// assign the moves to the other users
func (a moveAssigner) BulkMoveReAssignment(appCtx appcontext.AppContext, queueType string, officeUserData []*ghcmessages.BulkAssignmentForUser, reAssignFrom strfmt.UUID) (*models.Moves, error) {

	var movesToReAssignUsersFrom models.Moves

	whereClause, err := a.getBulkReAssigeeMoveQuery(queueType, reAssignFrom)
	if err != nil {
		return nil, apperror.NewBadDataError("Unable to generate move retrieval query")
	}

	queryError := appCtx.DB().Q().Where(whereClause).All(&movesToReAssignUsersFrom)
	if queryError != nil {
		return nil, apperror.NewBadDataError("unable to fetch re-assignment moves!")
	}

	moveAssignments, userQueue := a.getUserAssignmentCounts(officeUserData)

	// keep track of the updatedMovesForBatchSave to batch save
	updatedMovesForBatchSave := a.bulkAssignUsers(movesToReAssignUsersFrom, userQueue, moveAssignments, queueType)

	if len(updatedMovesForBatchSave) > 0 {
		verrs, err := appCtx.DB().ValidateAndUpdate(updatedMovesForBatchSave) // Bulk update
		if err != nil || verrs.HasAny() {
			return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Bulk assignment failed")
		}
	}

	return nil, nil
}

// Assign/ReAssign users
func (a moveAssigner) bulkAssignUsers(movesToReAssignUsersFrom models.Moves, userQueue *list.List, moveAssignments map[uuid.UUID]int, queueType string) models.Moves {
	moveIndex := 0
	updatedMoves := make([]models.Move, 0, len(movesToReAssignUsersFrom))
	// Re-Assign Users
	for i, move := range movesToReAssignUsersFrom {
		fmt.Printf("Iteration %d, queue length: %d\n", i, userQueue.Len())

		user := userQueue.Front()
		userID := user.Value.(uuid.UUID)
		fmt.Printf("Processing user: %s with current assignment count: %d\n", userID, moveAssignments[userID])
		userQueue.Remove(user)

		ordersType := move.Orders.OrdersType
		if ordersType != internalmessages.OrdersTypeSAFETY && ordersType != internalmessages.OrdersTypeBLUEBARK && ordersType != internalmessages.OrdersTypeWOUNDEDWARRIOR {
			updatedMove, err := a.commitMoveAssignmentToMove(queueType, &move, userID)
			if err == nil {
				updatedMoves = append(updatedMoves, *updatedMove)
			}
		}

		// decrement the user's assignment count
		moveAssignments[userID]--
		fmt.Printf("User %s assignment count after decrement: %d\n", userID, moveAssignments[userID])
		moveIndex++

		// If user still has remaining assignments, re-queue them
		if moveAssignments[userID] > 0 {
			userQueue.PushBack(userID)
			fmt.Printf("User %s requeued with remaining assignments: %d\n", userID, moveAssignments[userID])
		}

	}

	return updatedMoves
}

func (a moveAssigner) getBulkReAssigeeMoveQuery(queueType string, reAssignFrom strfmt.UUID) (string, error) {
	var queryFieldTerm string
	switch queueType {
	case string(models.QueueTypeCounseling), string(models.QueueTypeCloseout):
		queryFieldTerm = "sc_assigned_id"
	case string(models.QueueTypeTaskOrder):
		queryFieldTerm = "too_assigned_id"
	case string(models.QueueTypePaymentRequest):
		queryFieldTerm = "tio_assigned_id"
	default:
		return "", apperror.NewBadDataError("Invalid queue type")
	}

	whereClause := fmt.Sprintf(`%s = '%s'`, queryFieldTerm, reAssignFrom)

	return whereClause, nil
}

func (a moveAssigner) commitMoveAssignmentToMove(queueType string, input_Move *models.Move, uuid uuid.UUID) (*models.Move, error) {

	switch queueType {
	case string(models.QueueTypeCounseling), string(models.QueueTypeCloseout):
		input_Move.SCAssignedID = &uuid
	case string(models.QueueTypeTaskOrder):
		input_Move.TOOAssignedID = &uuid
	case string(models.QueueTypePaymentRequest):
		input_Move.TIOAssignedID = &uuid
	default:
		return nil, apperror.NewBadDataError("Invalid queue type")
	}
	return input_Move, nil
}

func (a moveAssigner) getUserAssignmentCounts(officeUserData []*ghcmessages.BulkAssignmentForUser) (map[uuid.UUID]int, *list.List) {

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
	return moveAssignments, queue
}
