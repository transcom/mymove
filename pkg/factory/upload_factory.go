package factory

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

type UploadExtendedParams struct {
	Uploader   *uploader.Uploader
	File       afero.File
	AppContext appcontext.AppContext
}

// BuildUpload creates an Upload.
//
// The customization for BuildUpload allows dev to provide an UploadExtendedParams object.
// This extended mode uses an Uploader object to create the upload vs. using the model.
// If an Uploader is provided, the model customizations are ignored in favor of the actual
// file provided for upload.
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUpload(db *pop.Connection, customs []Customization, traits []Trait) models.Upload {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cUpload models.Upload
	var cUploadParams *UploadExtendedParams
	if result := findValidCustomization(customs, Upload); result != nil {
		cUpload = result.Model.(models.Upload)

		// If extendedParams were provided, extract them
		typedResult, ok := result.ExtendedParams.(*UploadExtendedParams)
		if result.ExtendedParams != nil && !ok {
			log.Panic("To create Upload model, ExtendedParams must be nil or a pointer to UploadExtendedParams")
		}
		cUploadParams = typedResult

		if result.LinkOnly {
			return cUpload
		}
	}

	// UPLOADER MODE
	// The upload customization has an extended parameter struct that includes a Uploader interface and a file.
	// If the Uploader is passed in, models.Upload assertions are ignored in favor of the Uploader.
	// Instead we use the Uploader functionality to add the file. This creates the Upload model.
	if db != nil && cUploadParams != nil && cUploadParams.Uploader != nil {
		// Get file object
		var file afero.File
		if cUploadParams.File != nil {
			file = cUploadParams.File
		} else {
			file = FixtureOpen("test.pdf")
		}

		// Create file upload
		if cUploadParams.AppContext == nil {
			log.Panic("If Uploader is provided, AppContext must also be provided.")
		}
		upload, verrs, err := cUploadParams.Uploader.CreateUpload(cUploadParams.AppContext, uploader.File{File: file}, uploader.AllowedTypesServiceMember)
		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving upload %v, %v", verrs, err))
		}
		return *upload
	}

	// create upload
	upload := models.Upload{
		Filename:    "testFile.pdf",
		Bytes:       int64(2202009),
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&upload, cUpload)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &upload)
	}

	return upload
}

// BuildDefaultUpload returns an admin user with appropriate email
// Also creates
//   - User
func BuildDefaultUpload(db *pop.Connection) models.Upload {
	return BuildUpload(db, nil, nil)
}

func GetTraitTimestampedUpload() []Customization {
	return []Customization{
		{
			Model: models.Upload{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}
