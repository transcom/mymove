package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestServiceRequestDocumentUploadValidation() {
	suite.Run("test valid ServiceRequestDocumentUpload", func() {
		validServiceRequestDocumentUpload := models.ServiceRequestDocumentUpload{
			ContractorID:             uuid.Must(uuid.NewV4()),
			ServiceRequestDocumentID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validServiceRequestDocumentUpload, expErrors, nil)
	})

	suite.Run("test empty ServiceRequestDocumentUpload", func() {
		invalidServiceRequestDocumentUpload := models.ServiceRequestDocumentUpload{}

		expErrors := map[string][]string{
			"contractor_id":               {"ContractorID can not be blank."},
			"service_request_document_id": {"ServiceRequestDocumentID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidServiceRequestDocumentUpload, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchDeletedServiceRequestUpload() {
	serviceRequestUpload := factory.BuildServiceRequestDocumentUpload(suite.DB(), nil, nil)

	err := models.DeleteServiceRequestDocumentUpload(suite.DB(), &serviceRequestUpload)

	suite.Nil(err)
}
