package models_test

import (
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"
)

func (suite *ModelSuite) Test_PrimeUploadCreate() {
	t := suite.T()

	posDoc := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypePRIME,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)
	if err != nil {
		t.Fatalf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect PrimeUpload validation errors: %v", verrs)
	}

	primeUpload := models.PrimeUpload{
		ProofOfServiceDocID: posDoc.ID,
		ContractorID:        contractor.ID,
		Upload:              upload,
	}

	verrs, err = suite.DB().ValidateAndSave(&primeUpload)

	if err != nil {
		t.Fatalf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect PrimeUpload validation errors: %v", verrs)
	}
}

func (suite *ModelSuite) Test_PrimeUploadCreateWithID() {
	t := suite.T()

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypePRIME,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)
	if err != nil {
		t.Fatalf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect PrimeUpload validation errors: %v", verrs)
	}

	posDoc := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	id := uuid.Must(uuid.NewV4())
	primeUpload := models.PrimeUpload{
		ID:                  id,
		ProofOfServiceDocID: posDoc.ID,
		ContractorID:        contractor.ID,
		Upload:              upload,
	}

	verrs, err = suite.DB().ValidateAndSave(&primeUpload)

	if err != nil {
		t.Fatalf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect PrimeUpload validation errors: %v", verrs)
	}

	if primeUpload.ID.String() != id.String() {
		t.Errorf("wrong uuid for PrimeUpload: expected %s, got %s", id.String(), primeUpload.ID.String())
	}
}

func (suite *ModelSuite) Test_PrimeUploadValidations() {
	primeUpload := &models.PrimeUpload{}

	var expErrors = map[string][]string{
		"proof_of_service_doc_id": {"ProofOfServiceDocID can not be blank."},
		"contractor_id":           {"ContractorID can not be blank."},
	}

	suite.verifyValidationErrors(primeUpload, expErrors)
}

func (suite *ModelSuite) TestFetchPrimeUploadWithNoUpload() {
	posDoc := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	primeUpload := models.PrimeUpload{
		ProofOfServiceDocID: posDoc.ID,
		ContractorID:        contractor.ID,
	}

	_, err := suite.DB().ValidateAndSave(&primeUpload)
	suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.ForeignKeyViolation, "prime_uploads_uploads_id_fkey"), "expected primeupload error")
}

func (suite *ModelSuite) TestFetchPrimeUpload() {
	t := suite.T()

	posDoc := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypePRIME,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)
	if err != nil {
		t.Fatalf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect PrimeUpload validation errors: %v", verrs)
	}
	primeUpload := models.PrimeUpload{
		ProofOfServiceDocID: posDoc.ID,
		ContractorID:        contractor.ID,
		Upload:              upload,
	}

	verrs, err = suite.DB().ValidateAndSave(&primeUpload)
	if err != nil {
		t.Fatalf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect PrimeUpload validation errors: %v", verrs)
	}

	primeUp, _ := models.FetchPrimeUpload(suite.DB(), contractor.ID, primeUpload.ID)
	suite.Equal(primeUp.ID, primeUpload.ID)
	suite.Equal(upload.ID, primeUpload.Upload.ID)
	suite.Equal(upload.ID, primeUpload.UploadID)
}

func (suite *ModelSuite) TestFetchDeletedPrimeUpload() {
	t := suite.T()

	posDoc := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())
	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypePRIME,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)
	if err != nil {
		t.Fatalf("could not save Upload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect Upload validation errors: %v", verrs)
	}

	primeUpload := models.PrimeUpload{
		ProofOfServiceDocID: posDoc.ID,
		ContractorID:        contractor.ID,
		UploadID:            upload.ID,
		Upload:              upload,
	}

	verrs, err = suite.DB().ValidateAndSave(&primeUpload)
	if err != nil {
		t.Fatalf("could not save PrimeUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}

	err = models.DeletePrimeUpload(suite.DB(), &primeUpload)
	suite.Nil(err)
	primeUp, err := models.FetchPrimeUpload(suite.DB(), contractor.ID, primeUpload.ID)
	suite.Equal("error fetching prime_uploads: FETCH_NOT_FOUND", err.Error())

	// fetches a nil primeupload
	suite.Equal(primeUp.ID, uuid.Nil)
}
