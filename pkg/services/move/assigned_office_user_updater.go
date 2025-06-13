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
		move.SCCounselingAssignedID = &officeUser.ID
		move.SCCounselingAssignedUser = officeUser
	case models.QueueTypeCloseout:
		move.SCCloseoutAssignedID = &officeUser.ID
		move.SCCloseoutAssignedUser = officeUser
	case models.QueueTypeTaskOrder:
		move.TOOTaskOrderAssignedID = &officeUser.ID
		move.TOOTaskOrderAssignedUser = officeUser
	case models.QueueTypeDestinationRequest:
		move.TOODestinationAssignedID = &officeUser.ID
		move.TOODestinationAssignedUser = officeUser
	case models.QueueTypePaymentRequest:
		move.TIOPaymentRequestAssignedID = &officeUser.ID
		move.TIOPaymentRequestAssignedUser = officeUser
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
		move.SCCounselingAssignedID = nil
		move.SCCounselingAssignedUser = nil
	case models.QueueTypeCloseout:
		move.SCCloseoutAssignedID = nil
		move.SCCloseoutAssignedUser = nil
	case models.QueueTypeTaskOrder:
		move.TOOTaskOrderAssignedID = nil
		move.TOOTaskOrderAssignedUser = nil
	case models.QueueTypeDestinationRequest:
		move.TOODestinationAssignedID = nil
		move.TOODestinationAssignedUser = nil
	case models.QueueTypePaymentRequest:
		move.TIOPaymentRequestAssignedID = nil
		move.TIOPaymentRequestAssignedUser = nil
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if err != nil || verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "")
	}

	return &move, nil
}
