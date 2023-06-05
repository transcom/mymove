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

type ServiceRequestDocumentUploadExtendedParams struct {
	ServiceRequestDocumentUploader *uploader.ServiceRequestUploader
	UploaderID                     uuid.UUID
	File                           afero.File
	AppContext                     appcontext.AppContext
}

// BuildServiceRequestDocumentUpload creates a ServiceRequestDocumentUpload.
//
// The customization for BuildServiceRequestDocumentUpload allows dev to provide an
// ServiceRequestDocumentUploadExtendedParams object. This extended mode uses an
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
func BuildServiceRequestDocumentUpload(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceRequestDocumentUpload {
	// Make sure that any uploads created for ServiceRequestDocumentUpload have UploadType: models.UploadTypePRIME
	traits = append(traits, GetTraitUploadTypePrimeServiceRequest)
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cServiceRequestDocumentUpload models.ServiceRequestDocumentUpload
	var cServiceRequestDocumentUploadParams *ServiceRequestDocumentUploadExtendedParams
	if result := findValidCustomization(customs, ServiceRequestDocumentUpload); result != nil {
		cServiceRequestDocumentUpload = result.Model.(models.ServiceRequestDocumentUpload)

		if result.LinkOnly {
			return cServiceRequestDocumentUpload
		}

		// If extendedParams were provided, extract them
		typedResult, ok := result.ExtendedParams.(*ServiceRequestDocumentUploadExtendedParams)
		if result.ExtendedParams != nil && !ok {
			log.Panic("To create ServiceRequestDocumentUpload model, ExtendedParams must be nil or a pointer to ServiceRequestDocumentUploadExtendedParams")
		}
		cServiceRequestDocumentUploadParams = typedResult

	}

	contractor := FetchOrBuildDefaultContractor(db, customs, traits)
	serviceRequestDocument := BuildServiceRequestDocument(db, customs, traits)

	// UPLOADER MODE
	//
	// The prime upload customization has an extended parameter
	// struct that includes a ServiceRequestDocumentUploader interface and a file.
	// If the ServiceRequestDocumentUploader is passed in, models.ServiceRequestDocumentUpload
	// assertions are ignored in favor of the ServiceRequestDocumentUploader. The
	// ServiceRequestDocument and ContractorID customizations are still used if
	// provided. The ServiceRequestDocumentUploader functionality is used to add the
	// file. This creates the Upload model.
	if db != nil && cServiceRequestDocumentUploadParams != nil && cServiceRequestDocumentUploadParams.ServiceRequestDocumentUploader != nil {
		// Appcontext required if uploader mode used.
		if cServiceRequestDocumentUploadParams.AppContext == nil {
			log.Panic("If ServiceRequestDocumentUploader is provided, AppContext must also be provided.")
		}

		// Get file object
		var file afero.File
		if cServiceRequestDocumentUploadParams.File != nil {
			file = cServiceRequestDocumentUploadParams.File
		} else {
			file = FixtureOpen("test.pdf")
		}

		// Create file serviceRequestUpload
		serviceRequestUpload, verrs, err := cServiceRequestDocumentUploadParams.ServiceRequestDocumentUploader.CreateServiceRequestUploadForDocument(
			cServiceRequestDocumentUploadParams.AppContext,
			&serviceRequestDocument.ID,
			contractor.ID,
			uploader.File{File: file},
			uploader.AllowedTypesServiceMember,
		)

		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving prime upload %v, %v", verrs, err))
		}
		// CreateServiceRequestDocumentUploadForDocument does not assign ServiceRequestDocument or Contractor (just
		// ServiceRequestDocumentID and ContractorID), so do it manually to be consistent with when
		// not using an uploader
		serviceRequestUpload.ServiceRequestDocument = serviceRequestDocument
		serviceRequestUpload.Contractor = contractor
		return *serviceRequestUpload
	}

	// Find/create the Upload model with type models.UploadTypePRIME
	// GetTraitUploadTypePrimeServiceRequest was appended to traits at the beginning of this function
	tempUploadCustoms := customs
	tempUploadCustoms = convertCustomizationInList(tempUploadCustoms, Uploads.UploadTypePrime, Upload)
	upload := BuildUpload(db, tempUploadCustoms, traits)

	// Ensure the ServiceRequestDocumentUpload has the correct UploadType
	if upload.UploadType != models.UploadTypePRIME {
		log.Panic("ServiceRequestDocumentUpload must have UploadTypePRIME")
	}

	// create upload
	serviceRequestUpload := models.ServiceRequestDocumentUpload{
		ServiceRequestDocumentID: serviceRequestDocument.ID,
		ServiceRequestDocument:   serviceRequestDocument,
		ContractorID:             contractor.ID,
		Contractor:               contractor,
		Upload:                   upload,
		UploadID:                 upload.ID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&serviceRequestUpload, cServiceRequestDocumentUpload)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &serviceRequestUpload)
	}

	return serviceRequestUpload
}

func GetTraitUploadTypePrimeServiceRequest() []Customization {
	return []Customization{
		{
			Model: models.Upload{
				UploadType: models.UploadTypePRIME,
			},
			Type: &Uploads.UploadTypePrime,
		},
	}
}

func GetTraitServiceRequestDocumentUploadDeleted() []Customization {
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
