package move

import (
	"database/sql"
	"io"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

type additionalDocumentsUploader struct {
	uploadCreator services.UploadCreator
	checks        []validator
}

// NewMoveAdditionalDocumentsUploader returns a new additionalDocumentsUploader
func NewMoveAdditionalDocumentsUploader(uploadCreator services.UploadCreator) services.MoveAdditionalDocumentsUploader {
	return &additionalDocumentsUploader{uploadCreator, basicChecks()}
}

// CreateAdditionalDocumentsUpload uploads an additional document and updates the move with the new upload info
func (u *additionalDocumentsUploader) CreateAdditionalDocumentsUpload(
	appCtx appcontext.AppContext,
	userID uuid.UUID,
	moveID uuid.UUID,
	file io.ReadCloser,
	filename string,
	storer storage.FileStorer,
	uploadType models.UploadType,
) (models.Upload, string, *validate.Errors, error) {
	moveToUpdate, findErr := u.findMoveWithAdditionalDocuments(appCtx, moveID)
	if findErr != nil {
		return models.Upload{}, "", nil, findErr
	}

	serviceMemberID, memberFindErr := findServiceMemberIDWithOrderID(appCtx, moveToUpdate.OrdersID)
	if memberFindErr != nil {
		return models.Upload{}, "", nil, memberFindErr
	}

	userUpload, url, verrs, err := u.additionalDoc(appCtx, userID, *serviceMemberID, *moveToUpdate, file, filename, storer, uploadType)
	if verrs.HasAny() || err != nil {
		return models.Upload{}, "", verrs, err
	}

	return userUpload.Upload, url, nil, nil
}

func (u *additionalDocumentsUploader) findMoveWithAdditionalDocuments(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.Move, error) {
	var move models.Move

	query := appCtx.DB().Q().EagerPreload("AdditionalDocuments").Where("moves.id = ?", moveID)

	err := query.Find(&move, moveID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "while looking for move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	return &move, nil
}

func (u *additionalDocumentsUploader) additionalDoc(appCtx appcontext.AppContext, userID uuid.UUID, serviceMemberID uuid.UUID, move models.Move, file io.ReadCloser, filename string, storer storage.FileStorer, uploadType models.UploadType) (models.UserUpload, string, *validate.Errors, error) {
	// If move does not have a Document for additional document uploads, then create a new one
	var err error
	savedAdditionalDoc := move.AdditionalDocuments
	if move.AdditionalDocuments == nil {
		additionalDocument := &models.Document{
			ServiceMemberID: serviceMemberID,
		}
		savedAdditionalDoc, err = u.saveAdditionalDocumentForMove(appCtx, additionalDocument)
		if err != nil {
			return models.UserUpload{}, "", nil, err
		}

		// save new AdditionalDocumentID (document ID) to move
		move.AdditionalDocuments = savedAdditionalDoc
		move.AdditionalDocumentsID = &savedAdditionalDoc.ID
		_, _, err = u.updateMove(appCtx, move)
		if err != nil {
			return models.UserUpload{}, "", nil, err
		}
	}

	// Create new user upload for additional document
	var userUpload *models.UserUpload
	var verrs *validate.Errors
	var url string
	userUpload, url, verrs, err = uploader.CreateUserUploadForDocumentWrapper(
		appCtx,
		userID,
		storer,
		file,
		filename,
		uploader.MaxCustomerUserUploadFileSizeLimit,
		uploader.AllowedTypesServiceMember,
		&savedAdditionalDoc.ID,
		uploadType,
	)

	if verrs.HasAny() || err != nil {
		return models.UserUpload{}, "", verrs, err
	}

	move.AdditionalDocuments.UserUploads = append(move.AdditionalDocuments.UserUploads, *userUpload)

	return *userUpload, url, nil, nil
}

func (u *additionalDocumentsUploader) saveAdditionalDocumentForMove(appCtx appcontext.AppContext, doc *models.Document) (*models.Document, error) {
	var docID uuid.UUID
	if doc != nil {
		docID = doc.ID
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var verrs *validate.Errors
		var err error
		verrs, err = txnAppCtx.DB().ValidateAndSave(doc)
		return handleError(docID, verrs, err)
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return doc, nil
}

func (f *additionalDocumentsUploader) updateMove(appCtx appcontext.AppContext, move models.Move) (*models.Move, uuid.UUID, error) {
	var returnedMove *models.Move
	var err error

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		returnedMove, err = updateMoveInTx(txnAppCtx, move)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, uuid.Nil, transactionError
	}

	return returnedMove, returnedMove.ID, nil
}

func updateMoveInTx(appCtx appcontext.AppContext, move models.Move) (*models.Move, error) {
	var verrs *validate.Errors
	var err error

	verrs, err = appCtx.DB().ValidateAndUpdate(&move)
	if e := handleError(move.ID, verrs, err); e != nil {
		return nil, e
	}

	return &move, nil
}

func findServiceMemberIDWithOrderID(appCtx appcontext.AppContext, orderID uuid.UUID) (*uuid.UUID, error) {
	var order models.Order

	err := appCtx.DB().Q().EagerPreload("ServiceMember.User",
		"ServiceMemberID").
		Find(&order, orderID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(orderID, "while looking for order")
		default:
			return nil, apperror.NewQueryError("Order", err, "")
		}
	}

	return &order.ServiceMemberID, nil
}
