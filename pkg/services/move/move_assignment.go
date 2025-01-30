package move

import (
	"slices"

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

	transactionErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// track moves we've already assigned for payment requests
		var paymentRequestMoveList []uuid.UUID

		for _, move := range movesToAssign {
			for _, officeUser := range officeUserData {
				if officeUser != nil && officeUser.MoveAssignments > 0 {
					officeUserId := uuid.FromStringOrNil(officeUser.ID.String())

					switch queueType {
					case string(models.QueueTypeCounseling):
						move.SCAssignedID = &officeUserId
					case string(models.QueueTypeCloseout):
						move.SCAssignedID = &officeUserId
					case string(models.QueueTypeTaskOrder):
						move.TOOAssignedID = &officeUserId
					case string(models.QueueTypePaymentRequest):
						if !slices.Contains(paymentRequestMoveList, move.ID) {
							move.TIOAssignedID = &officeUserId

							// add move id to list so we can ignore them for the rest of the loop
							paymentRequestMoveList = append(paymentRequestMoveList, move.ID)
						}
					}

					officeUser.MoveAssignments -= 1

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
