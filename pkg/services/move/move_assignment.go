package move

import (
	"container/list"

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

func (a moveAssigner) getUserAssignmentCounts(officeUserData []*ghcmessages.BulkAssignmentForUser) (map[uuid.UUID]int, list.List) {

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
	return moveAssignments, *queue
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

// get the incoming office users, and the user being re-assigned from
// get the moves related the the reassigned user
// find all moves related to the user being reassigned
// assign the moves to the other users
func (a moveAssigner) BulkMoveReAssignment(appCtx appcontext.AppContext, queueType string, officeUserData *ghcmessages.BulkReAssignmentTakingWork, reAssignFrom strfmt.UUID) (*models.Moves, error) {

	var movesToReAssignUsersFrom *models.Moves

	//build Bulk Re-Assignment Query
	//queryStub := "SELECT * FROM moves WHERE"
	//var queryFieldTerm string

	// set the field to search for ID with in the query
	/*
		switch queueType {
		case string(models.QueueTypeCounseling), string(models.QueueTypeCloseout):
			queryFieldTerm = "sc_assigned_id"
		case string(models.QueueTypeTaskOrder):
			queryFieldTerm = "too_assigned_id"
		case string(models.QueueTypePaymentRequest):
			queryFieldTerm = "tio_assigned_id"
		default:
			return nil, apperror.NewBadDataError("Invalid queue type")
		}
	*/

	// format the query
	//fullQueryString := fmt.Sprintf(`%s %s=?`, queryStub, queryFieldTerm)
	//fullQueryString := fmt.Sprintf(`%s %s='%s'`, queryStub, queryFieldTerm, reAssignFrom)
	//fullQueryString := fmt.Sprintf(`%s =$1`, queryFieldTerm)

	err := appCtx.DB().Q().
		Eager("Moves").
		Where("sc_assigned_id = ?", reAssignFrom).All(&movesToReAssignUsersFrom)
	// direct injection to get moves to re-assign from
	//err := appCtx.DB().RawQuery(fullQueryString).All(&movesToReAssignUsersFrom)
	if err != nil {
		return nil, apperror.NewBadDataError("Invalid queue type")
	}

	moveAssignments, userQueue := a.getUserAssignmentCounts(nil)

	// keep track of the updatedMovesForBatchSave to batch save
	updatedMovesForBatchSave := make([]models.Move, 0, len(*movesToReAssignUsersFrom))

	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// while we have a queue...

		updatedMovesForBatchSave = a.bulkAssignUsers(*movesToReAssignUsersFrom, userQueue, moveAssignments, queueType)

		return nil
	})

	if len(updatedMovesForBatchSave) > 0 {
		verrs, err := appCtx.DB().ValidateAndUpdate(updatedMovesForBatchSave) // Bulk update
		if err != nil || verrs.HasAny() {
			return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Bulk assignment failed")
		}
	}

	if transactionErr != nil {
		return nil, transactionErr
	}

	return nil, nil
}

// Assign/ReAssign users
func (a moveAssigner) bulkAssignUsers(movesToReAssignUsersFrom models.Moves, userQueue list.List, moveAssignments map[uuid.UUID]int, queueType string) models.Moves {
	moveIndex := 0
	updatedMoves := make([]models.Move, 0, len(movesToReAssignUsersFrom))
	// Re-Assign Users
	for _, move := range movesToReAssignUsersFrom {
		user := userQueue.Front()
		userID := user.Value.(uuid.UUID)
		userQueue.Remove(user)

		ordersType := move.Orders.OrdersType
		if ordersType != internalmessages.OrdersTypeSAFETY && ordersType != internalmessages.OrdersTypeBLUEBARK && ordersType != internalmessages.OrdersTypeWOUNDEDWARRIOR {
			updatedMove, err := a.commitMoveAssignmentToMove(queueType, &move, userID)
			if err != nil {
				updatedMoves = append(updatedMoves, *updatedMove)
			}
		}

		// decrement the user's assignment count
		moveAssignments[userID]--
		moveIndex++

		// If user still has remaining assignments, re-queue them
		if moveAssignments[userID] > 0 {
			userQueue.PushBack(userID)
		}
	}

	return updatedMoves
}
