package models

import (
	"context"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
)

// A Document represents a physical artifact such as a multipage form that was
// filled out by hand. A Document can have many associated Uploads, which allows
// for handling multiple files that belong to the same document.
type Document struct {
	ID              uuid.UUID     `db:"id"`
	ServiceMemberID uuid.UUID     `db:"service_member_id"`
	ServiceMember   ServiceMember `belongs_to:"service_members"`
	CreatedAt       time.Time     `db:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"`
	Uploads         Uploads       `has_many:"uploads" order_by:"created_at asc"`
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
func FetchDocument(ctx context.Context, db *pop.Connection, session *auth.Session, id uuid.UUID) (Document, error) {

	ctx, span := beeline.StartSpan(ctx, "FetchDocument")
	defer span.Send()

	var document Document
	err := db.Q().Eager().Find(&document, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Document{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Document{}, err
	}

	_, smErr := FetchServiceMemberForUser(ctx, db, session, document.ServiceMemberID)
	if smErr != nil {
		return Document{}, smErr
	}
	return document, nil
}
