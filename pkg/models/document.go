package models

import (
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

// Documents is not required by pop and may be deleted
type Documents []Document

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: d.ServiceMemberID, Name: "ServiceMemberID"},
	), nil
}

// FetchDocument returns a document if the user has access to that document
func FetchDocument(db *pop.Connection, session *auth.Session, id uuid.UUID, includeDeletedDocs bool) (Document, error) {
	var document Document
	query := db.Q()

	if !includeDeletedDocs {
		query = query.Where("documents.deleted_at is null and u.deleted_at is null")
	}

	err := query.Eager("UserUploads.Upload").
		LeftJoin("user_uploads as uu", "documents.id = uu.document_id").
		LeftJoin("uploads as u", "uu.upload_id = u.id").
		Find(&document, id)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Document{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Document{}, err
	}

	_, smErr := FetchServiceMemberForUser(db, session, document.ServiceMemberID)
	if smErr != nil {
		return Document{}, smErr
	}
	return document, nil
}
