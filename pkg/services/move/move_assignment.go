package move

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveAssigner struct {
}

func NewMoveAssignerBulkAssignment() services.MoveAssigner {
	return &moveAssigner{}
}

func (a moveAssigner) BulkMoveAssignment(appCtx appcontext.AppContext, queueType string, officeUsers []models.OfficeUserWithWorkload, movesToAssign models.Moves) (*models.Moves, error) {
	if len(movesToAssign) == 0 {
		return nil, apperror.NewBadDataError("No moves to assign")
	}

	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		for _, move := range movesToAssign {
			for _, officeUser := range officeUsers {
				if officeUser.Workload > 0 {
					switch queueType {
					case string(models.QueueTypeCounseling):
						move.SCAssignedID = &officeUser.ID
					case string(models.QueueTypeCloseout):
						move.SCAssignedID = &officeUser.ID
					case string(models.QueueTypeTaskOrder):
						move.TOOAssignedID = &officeUser.ID
					}

					officeUser.Workload -= 1

					verrs, err := appCtx.DB().ValidateAndUpdate(&move)
					if err != nil || verrs.HasAny() {
						return apperror.NewInvalidInputError(move.ID, err, verrs, "")
					}

					break
				}
			}
		}

		return nil
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return nil, nil
}

// {
// 	officeUsers: [
// 		'1': 2,
// 		'2': 1,
//         '3': 3,
//         '4': 4,
//         '5': 5,
//         '6': 6,
// 	]
// 	moveIds: []
// 	queueType
// }
