package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
)

// A Document represents a physical artifact such as a multipage form that was
// filled out by hand. A Document can have many associated Uploads, which allows
// for handling multiple files that belong to the same document.
type Document struct {
	ID              uuid.UUID     `db:"id"`
	ServiceMemberID uuid.UUID     `db:"service_member_id"`
	ServiceMember   ServiceMember `belongs_to:"service_members" fk_id:"service_member_id"`
	CreatedAt       time.Time     `db:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"`
	DeletedAt       *time.Time    `db:"deleted_at"`
	UserUploads     UserUploads   `has_many:"user_uploads" fk_id:"document_id" order_by:"created_at asc"`
}

// TableName overrides the table name used by Pop.
func (d Document) TableName() string {
	return "documents"
}

type Documents []Document

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (d *Document) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: d.ServiceMemberID, Name: "ServiceMemberID"},
	), nil
}

// FetchDocument returns a document if the user has access to that document
func FetchDocument(db *pop.Connection, session *auth.Session, id uuid.UUID) (Document, error) {
	return fetchDocumentWithAccessibilityCheck(db, session, id, true)
}

// FetchDocument returns a document regardless if user has access to that document
func FetchDocumentWithNoRestrictions(db *pop.Connection, session *auth.Session, id uuid.UUID) (Document, error) {
	return fetchDocumentWithAccessibilityCheck(db, session, id, false)
}

// FetchDocument returns a document if the user has access to that document
func fetchDocumentWithAccessibilityCheck(db *pop.Connection, session *auth.Session, id uuid.UUID, checkUserAccessiability bool) (Document, error) {
	var document Document
	var uploads []Upload
	query := db.Q()
	// Giving the cursors names in which they will be defined as after opened in the database function.
	// Doing so we can reference the specific cursor we want by the defined name as opposed to <unnamed portal 1>, <unnamed portal 2>
	// which causes syntax errors when used in the FETCH ALL IN query.
	documentCursor := "documentcursor"
	userUploadCursor := "useruploadcursor"
	uploadCursor := "uploadcursor"

	documentsQuery := `SELECT fetch_documents(?, ?, ?, ?);`

	err := query.RawQuery(documentsQuery, documentCursor, userUploadCursor, uploadCursor, id).Exec()

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Document{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Document{}, err
	}

	// Since we know the name of the cursor we can fetch the specific one we are interested in
	// using FETCH ALL IN and populate the appropriate model
	fetchDocument := `FETCH ALL IN ` + documentCursor + `;`
	fetchUserUploads := `FETCH ALL IN ` + userUploadCursor + `;`
	fetchUploads := `FETCH ALL IN ` + uploadCursor + `;`

	err = query.RawQuery(fetchDocument).First(&document)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Document{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Document{}, err
	}

	err = query.RawQuery(fetchUserUploads).All(&document.UserUploads)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Document{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Document{}, err
	}

	err = query.RawQuery(fetchUploads).All(&uploads)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Document{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Document{}, err
	}

	// We have an array of UserUploads inside Document model, to populate that Upload model we need to loop and apply
	// the resulting uploads into the appropriate UserUpload.Upload model by matching the upload ids
	for i := 0; i < len(document.UserUploads); i++ {
		for j := 0; j < len(uploads); j++ {
			if document.UserUploads[i].UploadID == uploads[j].ID {
				document.UserUploads[i].Upload = uploads[j]
			}
		}
	}

	// We close all the cursors we opened during the fetch_documents call
	closeDocCursor := `CLOSE ` + documentCursor + `;`
	closeUserCursor := `CLOSE ` + userUploadCursor + `;`
	closeUploadCursor := `CLOSE ` + uploadCursor + `;`

	closeErr := query.RawQuery(closeDocCursor).Exec()

	if closeErr != nil {
		return Document{}, fmt.Errorf("error closing documents cursor: %w", closeErr)
	}

	closeErr = query.RawQuery(closeUserCursor).Exec()

	if closeErr != nil {
		return Document{}, fmt.Errorf("error closing user uploads cursor: %w", closeErr)
	}

	closeErr = query.RawQuery(closeUploadCursor).Exec()

	if closeErr != nil {
		return Document{}, fmt.Errorf("error closing uploads cursor: %w", closeErr)
	}

	if checkUserAccessiability {
		_, smErr := FetchServiceMemberForUser(db, session, document.ServiceMemberID)
		if smErr != nil {
			return Document{}, smErr
		}
	}

	return document, nil
}
