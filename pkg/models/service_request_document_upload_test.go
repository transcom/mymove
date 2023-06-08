package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *ModelSuite) TestServiceRequestDocumentUploadValidation() {
	suite.Run("test valid ServiceRequestDocumentUpload", func() {
		validServiceRequestDocumentUpload := models.ServiceRequestDocumentUpload{
			ContractorID:             uuid.Must(uuid.NewV4()),
			ServiceRequestDocumentID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validServiceRequestDocumentUpload, expErrors)
	})

	suite.Run("test empty ServiceRequestDocumentUpload", func() {
		invalidServiceRequestDocumentUpload := models.ServiceRequestDocumentUpload{}

		expErrors := map[string][]string{
			"contractor_id":               {"ContractorID can not be blank."},
			"service_request_document_id": {"ServiceRequestDocumentID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidServiceRequestDocumentUpload, expErrors)
	})
}

func (suite *ModelSuite) TestDeletedServiceRequestUpload() {
	t := suite.T()

	srDoc := factory.BuildServiceRequestDocument(suite.DB(), nil, nil)
	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: uploader.FileTypePDF,
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypePRIME,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)
	if err != nil {
		t.Errorf("could not save Upload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect Upload validation errors: %v", verrs)
	}

	serviceRequestUpload := models.ServiceRequestDocumentUpload{
		ServiceRequestDocumentID: srDoc.ID,
		ContractorID:             contractor.ID,
		UploadID:                 upload.ID,
		Upload:                   upload,
	}

	verrs, err = suite.DB().ValidateAndSave(&serviceRequestUpload)
	if err != nil {
		t.Errorf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}

	err = models.DeleteServiceRequestDocumentUpload(suite.DB(), &serviceRequestUpload)
	suite.Nil(err)
	srUp, err := models.FetchPrimeUpload(suite.DB(), contractor.ID, serviceRequestUpload.ID)
	suite.Equal("error fetching prime_uploads: FETCH_NOT_FOUND", err.Error())

	// fetches a nil primeupload
	suite.Equal(srUp.ID, uuid.Nil)
}
