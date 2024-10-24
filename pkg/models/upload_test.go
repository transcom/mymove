package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *ModelSuite) Test_ValidateUpload() {
	upload := models.Upload{
		ID:          uuid.Must(uuid.NewV4()),
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: uploader.FileTypePDF,
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}

	var expErrors = map[string][]string{}
	suite.verifyValidationErrors(&upload, expErrors)
}

func (suite *ModelSuite) Test_UploadValidationErrors() {
	upload := &models.Upload{}

	var expErrors = map[string][]string{
		"checksum":     {"Checksum can not be blank."},
		"bytes":        {"Bytes can not be blank."},
		"filename":     {"Filename can not be blank."},
		"content_type": {"ContentType can not be blank."},
		"upload_type":  {"UploadType is not in the list [USER, PRIME, OFFICE, APP]."},
	}

	suite.verifyValidationErrors(upload, expErrors)
}
