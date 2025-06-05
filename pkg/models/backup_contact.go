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

// BackupContactPermission represents the permissions granted to a backup contact
type BackupContactPermission string

const (
	// BackupContactPermissionNONE captures enum value "NONE"
	BackupContactPermissionNONE BackupContactPermission = "NONE"
	// BackupContactPermissionVIEW captures enum value "VIEW"
	BackupContactPermissionVIEW BackupContactPermission = "VIEW"
	// BackupContactPermissionEDIT captures enum value "EDIT"
	BackupContactPermissionEDIT BackupContactPermission = "EDIT"
)

// BackupContact is a model representing a backup contact for a service member
type BackupContact struct {
	ID              uuid.UUID               `json:"id" db:"id"`
	CreatedAt       time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at" db:"updated_at"`
	ServiceMemberID uuid.UUID               `json:"service_member_id" db:"service_member_id"`
	ServiceMember   ServiceMember           `belongs_to:"service_member" fk_id:"service_member_id"`
	Permission      BackupContactPermission `json:"permission" db:"permission"`
	FirstName       string                  `json:"first_name" db:"first_name"`
	LastName        string                  `json:"last_name" db:"last_name"`
	Email           string                  `json:"email" db:"email"`
	Phone           string                  `json:"phone" db:"phone"`
}

// TableName overrides the table name used by Pop.
func (b BackupContact) TableName() string {
	return "backup_contacts"
}

type BackupContacts []BackupContact

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (b *BackupContact) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: b.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: b.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: b.Email, Name: "Email"},
		&validators.StringIsPresent{Field: string(b.Permission), Name: "Permission"},
	), nil
}

// FetchBackupContact returns a specific backup contact model
func FetchBackupContact(db *pop.Connection, session *auth.Session, id uuid.UUID) (BackupContact, error) {
	var contact BackupContact
	err := db.Q().Eager().Find(&contact, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return BackupContact{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return BackupContact{}, err
	}
	// TODO: Handle case where more than one user is authorized to modify contact
	if session.IsMilApp() && contact.ServiceMember.ID != session.ServiceMemberID {
		return BackupContact{}, ErrFetchForbidden
	}
	return contact, nil
}
