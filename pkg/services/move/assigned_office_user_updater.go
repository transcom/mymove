package move

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type AssignedOfficeUserUpdater struct {
	services.MoveFetcher
}

func NewAssignedOfficeUserUpdater(moveFetcher services.MoveFetcher) services.MoveAssignedOfficeUserUpdater {
	return &AssignedOfficeUserUpdater{moveFetcher}
}

// arguments and return here correspond to what is setup in services/moves.go type MoveAssignedOfficeUserUpdater interface
func (s AssignedOfficeUserUpdater) UpdateAssignedOfficeUser(appCtx appcontext.AppContext, moveID uuid.UUID, officeUser *models.OfficeUser, queueType models.QueueType) (*models.Move, error) {
	var move models.Move
	err := appCtx.DB().Q().Find(&move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "while looking for move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	switch queueType {
	case models.QueueTypeCounseling:
		move.SCAssignedID = &officeUser.ID
		move.SCAssignedUser = officeUser
	case models.QueueTypeCloseout:
		move.SCAssignedID = &officeUser.ID
		move.SCAssignedUser = officeUser
	case models.QueueTypeTaskOrder:
		move.TOOAssignedID = &officeUser.ID
		move.TOOAssignedUser = officeUser
	case models.QueueTypeDestinationRequest:
		move.TOODestinationAssignedID = &officeUser.ID
		move.TOODestinationAssignedUser = officeUser
	case models.QueueTypePaymentRequest:
		move.TIOAssignedID = &officeUser.ID
		move.TIOAssignedUser = officeUser
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if err != nil || verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "")
	}

	return &move, nil
}

func (s AssignedOfficeUserUpdater) DeleteAssignedOfficeUser(appCtx appcontext.AppContext, moveID uuid.UUID, queueType models.QueueType) (*models.Move, error) {
	var move models.Move
	err := appCtx.DB().Q().Find(&move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "while looking for move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	switch queueType {
	case models.QueueTypeCounseling:
		move.SCAssignedID = nil
		move.SCAssignedUser = nil
	case models.QueueTypeCloseout:
		move.SCAssignedID = nil
		move.SCAssignedUser = nil
	case models.QueueTypeTaskOrder:
		move.TOOAssignedID = nil
		move.TOOAssignedUser = nil
	case models.QueueTypeDestinationRequest:
		move.TOODestinationAssignedID = nil
		move.TOODestinationAssignedUser = nil
	case models.QueueTypePaymentRequest:
		move.TIOAssignedID = nil
		move.TIOAssignedUser = nil
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if err != nil || verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "")
	}

	return &move, nil
}
