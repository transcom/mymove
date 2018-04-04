package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"time"
	// "github.com/transcom/mymove/pkg/gen/internalmessages"
)

// ServiceMember is a user of type service member
type ServiceMember struct {
	ID                        uuid.UUID  `json:"id" db:"id"`
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
	UserID                    uuid.UUID  `json:"user_id" db:"user_id"`
	User                      User       `json:"user" db:"user"`
	Edipi                     *string    `json:"edipi" db:"edipi"`
	FirstName                 *string    `json:"first_name" db:"first_name"`
	MiddleInitial             *string    `json:"middle_initial" db:"middle_initial"`
	LastName                  *string    `json:"last_name" db:"last_name"`
	Suffix                    *string    `json:"suffix" db:"suffix"`
	Telephone                 *string    `json:"telephone" db:"telephone"`
	SecondaryTelephone        *string    `json:"secondary_telephone" db:"secondary_telephone"`
	PersonalEmail             *string    `json:"personal_email" db:"personal_email"`
	PhoneIsPreferred          *bool      `json:"phone_is_preferred" db:"phone_is_preferred"`
	SecondaryPhoneIsPreferred *bool      `json:"secondary_phone_is_preferred" db:"secondary_phone_is_preferred"`
	EmailIsPreferred          *bool      `json:"email_is_preferred" db:"email_is_preferred"`
	ResidentialAddressID      *uuid.UUID `json:"residential_address_id" db:"residential_address_id"`
	ResidentialAddress        *Address   `json:"residential_address" db:"residential_address"`
	BackupMailingAddressID    *uuid.UUID `json:"backup_mailing_address_id" db:"backup_mailing_address_id"`
	BackupMailingAddress      *Address   `json:"backup_mailing_address" db:"backup_mailing_address"`
}

// todo add func to evaluate whether profile is complete - add call to payload struct in handler

// String is not required by pop and may be deleted
func (s ServiceMember) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ServiceMembers is not required by pop and may be deleted
type ServiceMembers []ServiceMember

// String is not required by pop and may be deleted
func (s ServiceMembers) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ServiceMember) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.UserID, Name: "UserID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ServiceMember) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ServiceMember) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
