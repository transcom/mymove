package move

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type AssignedOfficeUserUpdater struct {
	services.MoveFetcher
}

func NewAssignedOfficeUserUpdater(moveFetcher services.MoveFetcher) services.MoveAssignedOfficeUserUpdater {
	return &AssignedOfficeUserUpdater{moveFetcher}
}

// arguments and return here correspond to what is setup in services/moves.go type MoveAssignedOfficeUserUpdater interface
func (s AssignedOfficeUserUpdater) UpdateAssignedOfficeUser(appCtx appcontext.AppContext, moveID uuid.UUID, officeUser *models.OfficeUser, role roles.RoleType) (*models.Move, error) {
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

	switch role {
	case roles.RoleTypeServicesCounselor:
		move.SCAssignedID = &officeUser.ID
		move.SCAssignedUser = officeUser
	case roles.RoleTypeTOO:
		move.TOOAssignedID = &officeUser.ID
		move.TOOAssignedUser = officeUser
	case roles.RoleTypeTIO:
		move.TIOAssignedID = &officeUser.ID
		move.TIOAssignedUser = officeUser
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if err != nil || verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "")
	}

	return &move, nil
}

func (s AssignedOfficeUserUpdater) DeleteAssignedOfficeUser(appCtx appcontext.AppContext, moveID uuid.UUID, role roles.RoleType) (*models.Move, error) {
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

	switch role {
	case roles.RoleTypeServicesCounselor:
		move.SCAssignedID = nil
		move.SCAssignedUser = nil
	case roles.RoleTypeTOO:
		move.TOOAssignedID = nil
		move.TOOAssignedUser = nil
	case roles.RoleTypeTIO:
		move.TIOAssignedID = nil
		move.TIOAssignedUser = nil
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if err != nil || verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "")
	}

	return &move, nil
}
