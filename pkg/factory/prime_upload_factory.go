package factory

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

type PrimeUploadExtendedParams struct {
	PrimeUploader *uploader.PrimeUploader
	UploaderID    uuid.UUID
	File          afero.File
	AppContext    appcontext.AppContext
}

// BuildPrimeUpload creates a PrimeUpload.
//
// The customization for BuildPrimeUpload allows dev to provide an
// PrimeUploadExtendedParams object. This extended mode uses an
// Uploader object to create the upload vs. using the model. If an
// Uploader is provided, the model customizations are ignored in favor
// of the actual file provided for upload. In addition, an AppContext
// must be provided.
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
//
// Notes:
// If you want to customize the Contractor and are using LinkOnly customizations,
// make sure the Contractor customizations are set when BuildMove is called
// This function uses FetchOrBuildDefaultContractor and the Contractor is first created
func BuildPrimeUpload(db *pop.Connection, customs []Customization, traits []Trait) models.PrimeUpload {
	// Make sure that any uploads created for PrimeUpload have UploadType: models.UploadTypePRIME
	traits = append(traits, GetTraitUploadTypePrime)
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cPrimeUpload models.PrimeUpload
	var cPrimeUploadParams *PrimeUploadExtendedParams
	if result := findValidCustomization(customs, PrimeUpload); result != nil {
		cPrimeUpload = result.Model.(models.PrimeUpload)

		if result.LinkOnly {
			return cPrimeUpload
		}

		// If extendedParams were provided, extract them
		typedResult, ok := result.ExtendedParams.(*PrimeUploadExtendedParams)
		if result.ExtendedParams != nil && !ok {
			log.Panic("To create PrimeUpload model, ExtendedParams must be nil or a pointer to PrimeUploadExtendedParams")
		}
		cPrimeUploadParams = typedResult

	}

	contractor := FetchOrBuildDefaultContractor(db, customs, traits)
	proofOfServiceDoc := BuildProofOfServiceDoc(db, customs, traits)

	// UPLOADER MODE
	//
	// The prime upload customization has an extended parameter
	// struct that includes a PrimeUploader interface and a file.
	// If the PrimeUploader is passed in, models.PrimeUpload
	// assertions are ignored in favor of the PrimeUploader. The
	// ProofOfServiceDoc and ContractorID customizations are still used if
	// provided. The PrimeUploader functionality is used to add the
	// file. This creates the Upload model.
	if db != nil && cPrimeUploadParams != nil && cPrimeUploadParams.PrimeUploader != nil {
		// Appcontext required if uploader mode used.
		if cPrimeUploadParams.AppContext == nil {
			log.Panic("If PrimeUploader is provided, AppContext must also be provided.")
		}

		// Get file object
		var file afero.File
		if cPrimeUploadParams.File != nil {
			file = cPrimeUploadParams.File
		} else {
			file = FixtureOpen("test.pdf")
		}

		// Create file primeUpload
		primeUpload, verrs, err := cPrimeUploadParams.PrimeUploader.CreatePrimeUploadForDocument(
			cPrimeUploadParams.AppContext,
			&proofOfServiceDoc.ID,
			contractor.ID,
			uploader.File{File: file},
			uploader.AllowedTypesServiceMember,
		)

		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving prime upload %v, %v", verrs, err))
		}
		// CreatePrimeUploadForDocument does not assign ProofOfServiceDoc or Contractor (just
		// ProofOfServiceDocID and ContractorID), so do it manually to be consistent with when
		// not using an uploader
		primeUpload.ProofOfServiceDoc = proofOfServiceDoc
		primeUpload.Contractor = contractor
		return *primeUpload
	}

	// Find/create the Upload model with type models.UploadTypePRIME
	// GetTraitUploadTypePrime was appended to traits at the beginning of this function
	tempUploadCustoms := customs
	tempUploadCustoms = convertCustomizationInList(tempUploadCustoms, Uploads.UploadTypePrime, Upload)
	upload := BuildUpload(db, tempUploadCustoms, traits)

	// Ensure the PrimeUpload has the correct UploadType
	if upload.UploadType != models.UploadTypePRIME {
		log.Panic("PrimeUpload must have UploadTypePRIME")
	}

	// create upload
	primeUpload := models.PrimeUpload{
		ProofOfServiceDocID: proofOfServiceDoc.ID,
		ProofOfServiceDoc:   proofOfServiceDoc,
		ContractorID:        contractor.ID,
		Contractor:          contractor,
		Upload:              upload,
		UploadID:            upload.ID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&primeUpload, cPrimeUpload)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &primeUpload)
	}

	return primeUpload
}

func GetTraitUploadTypePrime() []Customization {
	return []Customization{
		{
			Model: models.Upload{
				UploadType: models.UploadTypePRIME,
			},
			Type: &Uploads.UploadTypePrime,
		},
	}
}

func GetTraitPrimeUploadDeleted() []Customization {
	return []Customization{
		{
			Model: models.Upload{
				UploadType: models.UploadTypePRIME,
				DeletedAt:  models.TimePointer(time.Now()),
			},
			Type: &Uploads.UploadTypePrime,
		},
	}
}
