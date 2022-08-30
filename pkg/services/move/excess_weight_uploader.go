package move

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type excessWeightUploader struct {
	uploadCreator services.UploadCreator
	checks        []validator
}

// NewMoveExcessWeightUploader returns a new excessWeightUploader
func NewMoveExcessWeightUploader(uploadCreator services.UploadCreator) services.MoveExcessWeightUploader {
	return &excessWeightUploader{uploadCreator, basicChecks()}
}

// NewPrimeMoveExcessWeightUploader returns a new excessWeightUploader
func NewPrimeMoveExcessWeightUploader(uploadCreator services.UploadCreator) services.MoveExcessWeightUploader {
	return &excessWeightUploader{uploadCreator, primeChecks()}
}

// CreateExcessWeightUpload uploads an excess weight document and updates the move with the new upload info
func (u *excessWeightUploader) CreateExcessWeightUpload(
	appCtx appcontext.AppContext,
	moveID uuid.UUID,
	file io.ReadCloser,
	uploadFilename string,
	uploadType models.UploadType,
) (*models.Move, error) {
	// Get existing move
	move := &models.Move{}
	err := appCtx.DB().Find(move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "while looking for move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	// Run the (read-only) validations
	if verr := validateMove(appCtx, *move, nil, u.checks...); verr != nil {
		return nil, verr
	}

	// Open transaction to create upload and update the move
	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		excessWeightUpload, err := u.uploadCreator.CreateUpload(
			txnAppCtx, file, fmt.Sprintf("move/%s/%s", move.ID, uploadFilename), uploadType)
		if err != nil {
			return err
		}

		move.ExcessWeightUploadID = &excessWeightUpload.ID
		move.ExcessWeightUpload = excessWeightUpload

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(move)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(
				move.ID, err, verrs, "Validation errors found while updating excess weight info on move")
		} else if err != nil {
			return apperror.NewQueryError("Move", err, "Failed to update excess weight info on move")
		}

		return nil
	})
	if txnErr != nil {
		return nil, txnErr
	}

	return move, nil
}
