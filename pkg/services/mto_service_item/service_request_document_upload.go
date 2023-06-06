package mtoserviceitem

import (
	"database/sql"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

const (
	// VersionTimeFormat is the Go time format for creating a version number.
	VersionTimeFormat string = "20060102150405"
)

type serviceRequestDocumentUploadCreator struct {
	fileStorer    storage.FileStorer
	fileSizeLimit uploader.ByteSize
}

// NewServiceRequestDocumentUploadCreator returns a new payment request upload creator
func NewServiceRequestDocumentUploadCreator(fileStorer storage.FileStorer) services.ServiceRequestDocumentUploadCreator {
	return &serviceRequestDocumentUploadCreator{fileStorer, uploader.MaxFileSizeLimit}
}

func (p *serviceRequestDocumentUploadCreator) assembleUploadFilePathName(appCtx appcontext.AppContext, mtoServiceItemID uuid.UUID, filename string) (string, error) {
	var mtoServiceItem models.MTOServiceItem
	err := appCtx.DB().Where("id=$1", mtoServiceItemID).First(&mtoServiceItem)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", apperror.NewNotFoundError(mtoServiceItemID, "")
		default:
			return "", apperror.NewQueryError("MTOServiceItem", err, "")
		}
	}

	newfilename := time.Now().Format(VersionTimeFormat) + "-" + filename
	uploadFilePath := fmt.Sprintf("/payment-request-uploads/mto-%s/payment-request-%s", mtoServiceItem.MoveTaskOrderID, mtoServiceItem.ID)
	uploadFileName := path.Join(uploadFilePath, newfilename)

	return uploadFileName, err
}

func (p *serviceRequestDocumentUploadCreator) CreateUpload(appCtx appcontext.AppContext, file io.ReadCloser, mtoServiceItemID uuid.UUID, contractorID uuid.UUID, uploadFilename string) (*models.Upload, error) {
	var upload *models.Upload
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		newUploader, err := uploader.NewServiceRequestUploader(p.fileStorer, p.fileSizeLimit)
		if err != nil {
			if err == uploader.ErrFileSizeLimitExceedsMax {
				return apperror.NewBadDataError(err.Error())
			}
			return err
		}

		fileName, err := p.assembleUploadFilePathName(txnAppCtx, mtoServiceItemID, uploadFilename)
		if err != nil {
			return err
		}

		aFile, err := newUploader.PrepareFileForUpload(txnAppCtx, file, fileName)
		if err != nil {
			return err
		}

		newUploader.SetUploadStorageKey(fileName)

		var mtoServiceItem models.MTOServiceItem
		err = txnAppCtx.DB().Find(&mtoServiceItem, mtoServiceItemID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(mtoServiceItemID, "")
			default:
				return apperror.NewQueryError("MTOServiceItem", err, "")
			}
		}
		// create proof of service doc
		serviceRequestDocument := models.ServiceRequestDocument{
			MTOServiceItemID: mtoServiceItemID,
			MTOServiceItem:   mtoServiceItem,
		}
		verrs, err := txnAppCtx.DB().ValidateAndCreate(&serviceRequestDocument)
		if err != nil {
			return fmt.Errorf("failure creating proof of service doc: %w", err) // server err
		}
		if verrs.HasAny() {
			return apperror.NewInvalidCreateInputError(verrs, "validation error with creating proof of service doc")
		}

		srdID := &serviceRequestDocument.ID
		serviceRequestUpload, verrs, err := newUploader.CreateServiceRequestUploadForDocument(txnAppCtx, srdID, contractorID, uploader.File{File: aFile}, uploader.AllowedTypesServiceRequest)
		if verrs.HasAny() {
			return apperror.NewInvalidCreateInputError(verrs, "validation error with creating payment request")
		}
		if err != nil {
			return fmt.Errorf("failure creating service request document upload: %w", err)
		}
		upload = &serviceRequestUpload.Upload
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return upload, nil
}
