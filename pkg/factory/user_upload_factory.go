package factory

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

type UserUploadExtendedParams struct {
	UserUploader *uploader.UserUploader
	UploaderID   uuid.UUID
	File         afero.File
	AppContext   appcontext.AppContext
}

// BuildUserUpload creates an UserUpload.
//
// The customization for BuildUserUpload allows dev to provide an
// UserUploadExtendedParams object. This extended mode uses an
// Uploader object to create the upload vs. using the model. If an
// Uploader is provided, the model customizations are ignored in favor
// of the actual file provided for upload. In addition, an AppContext
// must be provided.
//
// Params:
//   - customs is a slice that will be modified by the factory
//   - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUserUpload(db *pop.Connection, customs []Customization, traits []Trait) models.UserUpload {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cUserUpload models.UserUpload
	var cUserUploadParams *UserUploadExtendedParams
	if result := findValidCustomization(customs, UserUpload); result != nil {
		cUserUpload = result.Model.(models.UserUpload)

		if result.LinkOnly {
			return cUserUpload
		}

		// If extendedParams were provided, extract them
		typedResult, ok := result.ExtendedParams.(*UserUploadExtendedParams)
		if result.ExtendedParams != nil && !ok {
			log.Panic("To create UserUpload model, ExtendedParams must be nil or a pointer to UserUploadExtendedParams")
		}
		cUserUploadParams = typedResult

	}

	// Find/create the Document model
	document := BuildDocument(db, customs, traits)

	// UPLOADER MODE
	//
	// The user upload customization has an extended parameter
	// struct that includes a UserUploader interface and a file.
	// If the UserUploader is passed in, models.UserUpload
	// assertions are ignored in favor of the UserUploader. The
	// Document and UploaderID customizations are still used if
	// provided. The UserUploader functionality is used to add the
	// file. This creates the Upload model.
	if db != nil && cUserUploadParams != nil && cUserUploadParams.UserUploader != nil {
		// Appcontext required if uploader mode used.
		if cUserUploadParams.AppContext == nil {
			log.Panic("If UserUploader is provided, AppContext must also be provided.")
		}

		// Get file object
		var file afero.File
		if cUserUploadParams.File != nil {
			file = cUserUploadParams.File
		} else {
			file = FixtureOpen("test.pdf")
		}

		var uploaderID uuid.UUID
		if !cUserUploadParams.UploaderID.IsNil() {
			uploaderID = cUserUploadParams.UploaderID
		} else {
			uploaderID = document.ServiceMember.UserID
		}

		// Create file userUpload
		userUpload, verrs, err := cUserUploadParams.UserUploader.CreateUserUploadForDocument(
			cUserUploadParams.AppContext,
			&document.ID,
			uploaderID,
			uploader.File{File: file},
			uploader.AllowedTypesPPMDocuments,
		)

		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving user upload %v, %v", verrs, err))
		}
		// CreateUserUploadForDocument does not assign Document (just
		// DocumentID), so do it manually to be consistent with when
		// not using an uploader
		userUpload.Document = document
		return *userUpload
	}

	// Find/create the Upload model
	tempUploadCustoms := customs
	result := findValidCustomization(customs, Uploads.UploadTypeUser)
	if result != nil {
		tempUploadCustoms = convertCustomizationInList(tempUploadCustoms, Uploads.UploadTypeUser, Upload)
	}
	// Treat any Upload customizations that aren't specified as Uploads.UploadTypePrime as user uploads
	upload := BuildUpload(db, tempUploadCustoms, traits)

	// Ensure the UserUpload has the correct UploadType
	if upload.UploadType != models.UploadTypeUSER {
		log.Panic("UserUpload must have UploadTypeUSER")
	}

	uploaderID := document.ServiceMember.UserID
	// create upload
	userUpload := models.UserUpload{
		DocumentID: &document.ID,
		Document:   document,
		UploaderID: uploaderID,
		Upload:     upload,
		UploadID:   upload.ID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&userUpload, cUserUpload)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &userUpload)
	}

	return userUpload
}

// sometimes have to create one out of nowhere for uploads
func uploaderAppContext(db *pop.Connection) appcontext.AppContext {
	// *sigh*, use global zap logger
	return appcontext.NewAppContext(db, zap.L(), nil, nil)
}

// buildDocumentWithUploads builds a document and creates an upload associated with the document. Returns the document at the end.
//
// Usage example:
//
//	emptyDocument := buildDocumentWithUploads(db, userUploader)
func buildDocumentWithUploads(db *pop.Connection, userUploader *uploader.UserUploader, serviceMember models.ServiceMember, file afero.File) models.Document {

	customs := []Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}

	if db != nil {
		customs = append(customs, Customization{
			Model: models.UserUpload{},
			ExtendedParams: &UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   uploaderAppContext(db),
				File:         file,
			},
		})
	}

	upload := BuildUserUpload(db, customs, nil)

	doc := upload.Document

	doc.UserUploads = append(doc.UserUploads, upload)

	return doc
}

func GetTraitTimestampedUserUpload() []Customization {
	return []Customization{
		{
			Model: models.UserUpload{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}
