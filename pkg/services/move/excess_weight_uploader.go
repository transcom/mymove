package move

import (
	"fmt"
	"io"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
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
		return nil, services.NewNotFoundError(moveID, "while looking for move")
	}

	// Run the (read-only) validations
	if verr := validateMove(appCtx, *move, nil, u.checks...); verr != nil {
		return nil, verr
	}

	// Open transaction to create upload and update the move
	txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		excessWeightUpload, err := u.uploadCreator.CreateUpload(txnAppCtx, file, uploadFilename, uploadType)
		if err != nil {
			return err
		}

		now := time.Now()
		move.ExcessWeightQualifiedAt = &now
		move.ExcessWeightUploadID = &excessWeightUpload.ID
		move.ExcessWeightUpload = excessWeightUpload

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(move)
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidCreateInputError(
				verrs, "Validation errors found while updating excess weight info on move")
		} else if err != nil {
			return fmt.Errorf("Failure to update excess weight info on move: %v", err)
		}

		return nil
	})
	if txnErr != nil {
		return nil, txnErr
	}

	return move, nil
}
