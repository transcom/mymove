package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// MakePrimeUpload creates a single PrimeUpload.
func MakePrimeUpload(db *pop.Connection, assertions Assertions) models.PrimeUpload {
	posDoc := assertions.PrimeUpload.ProofOfServiceDoc
	if isZeroUUID(assertions.PrimeUpload.ProofOfServiceDocID) {
		if isZeroUUID(assertions.ProofOfServiceDoc.ID) {
			posDoc = MakeProofOfServiceDoc(db, assertions)
		} else {
			posDoc = assertions.ProofOfServiceDoc
		}
	}

	contractor := assertions.PrimeUpload.Contractor
	if isZeroUUID(assertions.PrimeUpload.ContractorID) {
		if isZeroUUID(assertions.Contractor.ID) {
			contractor = MakeContractor(db, assertions)
		} else {
			contractor = assertions.Contractor
		}
	}

	// Users can either assert a PrimeUploader (and a real file is used), or can optionally assert fields
	var primeUpload *models.PrimeUpload
	if assertions.PrimeUploader != nil {
		// If an PrimeUploader is passed in, PrimeUpload assertions are ignored
		var err error
		var verrs *validate.Errors
		file := Fixture("test.pdf")

		if assertions.File != nil {
			file = assertions.File
		}
		// Ugh. Use the global logger. All testdatagen methods should
		// take a logger
		appCtx := appcontext.NewAppContext(db, zap.L(), nil)
		primeUpload, verrs, err = assertions.PrimeUploader.CreatePrimeUploadForDocument(appCtx, &posDoc.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesServiceMember)
		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving prime upload %v, %v", verrs, err))
		}

		primeUpload.ProofOfServiceDocID = posDoc.ID
		primeUpload.ProofOfServiceDoc = posDoc
	} else {
		// If no PrimeUploader is being stored, use asserted fields

		if assertions.PrimeUpload.Upload.ID != uuid.Nil {
			assertions.Upload = assertions.PrimeUpload.Upload
		}
		assertions.Upload.UploadType = models.UploadTypePRIME
		upload := MakeUpload(db, assertions)

		primeUpload = &models.PrimeUpload{
			ProofOfServiceDocID: posDoc.ID,
			ProofOfServiceDoc:   posDoc,
			ContractorID:        contractor.ID,
			Upload:              upload,
			UploadID:            upload.ID,
		}

		mergeModels(primeUpload, assertions.PrimeUpload)

		mustCreate(db, primeUpload, assertions.Stub)
	}

	return *primeUpload
}

// MakeDefaultPrimeUpload makes an PrimeUpload with default values
func MakeDefaultPrimeUpload(db *pop.Connection) models.PrimeUpload {
	return MakePrimeUpload(db, Assertions{})
}
