package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// BackupContact is a model representing a backup contact for a service member
type BackupContact struct {
	ID              uuid.UUID                                `json:"id" db:"id"`
	CreatedAt       time.Time                                `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                                `json:"updated_at" db:"updated_at"`
	ServiceMemberID uuid.UUID                                `json:"service_member_id" db:"service_member_id"`
	ServiceMember   ServiceMember                            `belongs_to:"service_member"`
	Permission      internalmessages.BackupContactPermission `json:"permission" db:"permission"`
	Name            string                                   `json:"name" db:"name"`
	Email           string                                   `json:"email" db:"email"`
	Phone           *string                                  `json:"phone" db:"phone"`
}

// String is not required by pop and may be deleted
func (b BackupContact) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BackupContacts is not required by pop and may be deleted
type BackupContacts []BackupContact

// String is not required by pop and may be deleted
func (b BackupContacts) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BackupContact) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: b.Name, Name: "Name"},
		&validators.StringIsPresent{Field: b.Email, Name: "Email"},
		&BackupContactPermissionIsPresent{Field: b.Permission, Name: "Permission"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *BackupContact) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *BackupContact) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchBackupContact returns a specific backup contact model
func FetchBackupContact(db *pop.Connection, authUser User, id uuid.UUID) (BackupContact, error) {
	var contact BackupContact
	eagerConn := *db
	err := eagerConn.Eager().Find(&contact, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return BackupContact{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return BackupContact{}, err
	}
	// TODO: Handle case where more than one user is authorized to modify contact
	if contact.ServiceMember.UserID != authUser.ID {
		return BackupContact{}, ErrFetchForbidden
	}

	return contact, nil
}
