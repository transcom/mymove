package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// ServiceMember is a user of type service member
type ServiceMember struct {
	ID                     uuid.UUID                           `json:"id" db:"id"`
	CreatedAt              time.Time                           `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time                           `json:"updated_at" db:"updated_at"`
	UserID                 uuid.UUID                           `json:"user_id" db:"user_id"`
	User                   User                                `belongs_to:"user"`
	Edipi                  *string                             `json:"edipi" db:"edipi"`
	Branch                 *internalmessages.MilitaryBranch    `json:"branch" db:"branch"`
	Rank                   *internalmessages.ServiceMemberRank `json:"rank" db:"rank"`
	FirstName              *string                             `json:"first_name" db:"first_name"`
	MiddleName             *string                             `json:"middle_name" db:"middle_name"`
	LastName               *string                             `json:"last_name" db:"last_name"`
	Suffix                 *string                             `json:"suffix" db:"suffix"`
	Telephone              *string                             `json:"telephone" db:"telephone"`
	SecondaryTelephone     *string                             `json:"secondary_telephone" db:"secondary_telephone"`
	PersonalEmail          *string                             `json:"personal_email" db:"personal_email"`
	PhoneIsPreferred       *bool                               `json:"phone_is_preferred" db:"phone_is_preferred"`
	TextMessageIsPreferred *bool                               `json:"text_message_is_preferred" db:"text_message_is_preferred"`
	EmailIsPreferred       *bool                               `json:"email_is_preferred" db:"email_is_preferred"`
	ResidentialAddressID   *uuid.UUID                          `json:"residential_address_id" db:"residential_address_id"`
	ResidentialAddress     *Address                            `belongs_to:"address"`
	BackupMailingAddressID *uuid.UUID                          `json:"backup_mailing_address_id" db:"backup_mailing_address_id"`
	BackupMailingAddress   *Address                            `belongs_to:"address"`
	SocialSecurityNumberID *uuid.UUID                          `json:"social_security_number_id" db:"social_security_number_id"`
	SocialSecurityNumber   *SocialSecurityNumber               `belongs_to:"address"`
	BackupContacts         *BackupContacts                     `has_many:"backup_contacts"`
}

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

// FetchServiceMember returns a service member only if it is allowed for the given user to access that service member.
func FetchServiceMember(db *pop.Connection, user User, id uuid.UUID) (ServiceMember, error) {
	var serviceMember ServiceMember
	err := db.Eager().Find(&serviceMember, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return ServiceMember{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return ServiceMember{}, err
	}
	// TODO: Handle case where more than one user is authorized to modify serviceMember
	if serviceMember.UserID != user.ID {
		return ServiceMember{}, ErrFetchForbidden
	}

	return serviceMember, nil
}

// CreateServiceMember takes a serviceMember with Address structs and coordinates saving it all in a transaction
func CreateServiceMember(dbConnection *pop.Connection, serviceMember *ServiceMember) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	// If the passed in function returns an error, the transaction is rolled back
	dbConnection.Transaction(func(dbConnection *pop.Connection) error {

		var transactionError error
		addressModels := []*Address{
			serviceMember.ResidentialAddress,
			serviceMember.BackupMailingAddress,
		}

		for _, model := range addressModels {
			if model == nil {
				continue
			} else if verrs, err := dbConnection.ValidateAndCreate(model); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				transactionError = errors.New("Rollback The transaction")
				// Halt what we're doing if we get a database error
				if err != nil {
					responseError = err
					break
				}
			}
		}

		if transactionError == nil && serviceMember.SocialSecurityNumber != nil {
			verrs, err := dbConnection.ValidateAndCreate(serviceMember.SocialSecurityNumber)
			if err != nil || verrs.HasAny() {
				responseVErrors.Append(verrs)
				transactionError = errors.New("Rollback The Transaction")
				if err != nil {
					responseError = err
				}
			}
		}

		if transactionError == nil {
			if serviceMember.SocialSecurityNumber != nil {
				serviceMember.SocialSecurityNumberID = &serviceMember.SocialSecurityNumber.ID
			}
			serviceMember.ResidentialAddressID = GetAddressID(serviceMember.ResidentialAddress)
			serviceMember.BackupMailingAddressID = GetAddressID(serviceMember.BackupMailingAddress)

			if verrs, err := dbConnection.ValidateAndCreate(serviceMember); verrs.HasAny() || err != nil {
				// Return a reasonable error if someone tries to create a second SM when one already exists for this user
				if strings.HasPrefix(errors.Cause(err).Error(), UniqueConstraintViolationErrorPrefix) {
					responseError = ErrCreateViolatesUniqueConstraint
				} else {
					responseError = err
				}

				transactionError = errors.New("Rollback The transaction")
				responseVErrors = verrs

			}
		}

		return transactionError

	})

	return responseVErrors, responseError

}

// CreateBackupContact creates a backup contact model tied to the service member
func (s ServiceMember) CreateBackupContact(db *pop.Connection, name string, email string, phone *string, permission internalmessages.BackupContactPermission) (BackupContact, *validate.Errors, error) {
	newContact := BackupContact{
		ServiceMemberID: s.ID,
		ServiceMember:   s,
		Name:            name,
		Email:           email,
		Phone:           phone,
		Permission:      permission,
	}

	verrs, err := db.ValidateAndCreate(&newContact)
	if err != nil || verrs.HasAny() {
		newContact = BackupContact{}
	}
	return newContact, verrs, err
}

// PatchServiceMemberWithPayload patches service member with payload
func (s *ServiceMember) PatchServiceMemberWithPayload(db *pop.Connection, payload *internalmessages.PatchServiceMemberPayload) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error
	db.Transaction(func(dbConnection *pop.Connection) error {
		var transactionError error

		if payload.Edipi != nil {
			s.Edipi = payload.Edipi
		}
		if payload.Branch != nil {
			s.Branch = payload.Branch
		}
		if payload.Rank != nil {
			s.Rank = payload.Rank
		}
		if payload.FirstName != nil {
			s.FirstName = payload.FirstName
		}
		if payload.MiddleName != nil {
			s.MiddleName = payload.MiddleName
		}
		if payload.LastName != nil {
			s.LastName = payload.LastName
		}
		if payload.Suffix != nil {
			s.Suffix = payload.Suffix
		}
		if payload.Telephone != nil {
			s.Telephone = payload.Telephone
		}
		if payload.SecondaryTelephone != nil {
			s.SecondaryTelephone = payload.SecondaryTelephone
		}
		if payload.PersonalEmail != nil {
			s.PersonalEmail = swag.String(payload.PersonalEmail.String())
		}
		if payload.PhoneIsPreferred != nil {
			s.PhoneIsPreferred = payload.PhoneIsPreferred
		}
		if payload.TextMessageIsPreferred != nil {
			s.TextMessageIsPreferred = payload.TextMessageIsPreferred
		}
		if payload.EmailIsPreferred != nil {
			s.EmailIsPreferred = payload.EmailIsPreferred
		}
		if payload.SocialSecurityNumber != nil {
			if s.SocialSecurityNumber != nil {
				// If SSN model exists
				ssn := s.SocialSecurityNumber
				if verrs, err := ssn.SetEncryptedHash(payload.SocialSecurityNumber.String()); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = err
					return errors.New("New Transaction Error")
				}

				if verrs, err := dbConnection.ValidateAndUpdate(ssn); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = err
					return errors.New("New Transaction Error")
				}
			} else {
				// Else create an SSN model
				newSSN := SocialSecurityNumber{}
				if verrs, err := newSSN.SetEncryptedHash(payload.SocialSecurityNumber.String()); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = err
					return errors.New("New Transaction Error")
				}

				if verrs, err := dbConnection.ValidateAndCreate(&newSSN); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = err
					return errors.New("New Transaction Error")
				}
				s.SocialSecurityNumber = &newSSN
				s.SocialSecurityNumberID = &newSSN.ID
			}

		}

		if payload.ResidentialAddress != nil {
			residentialAddress := AddressModelFromPayload(payload.ResidentialAddress)
			if verrs, err := dbConnection.ValidateAndCreate(residentialAddress); verrs.HasAny() || err != nil {
				if verrs.HasAny() {
					responseVErrors = verrs
				} else {
					responseError = err
				}
				transactionError = errors.New("Failed saving residential address model")
				return transactionError
			}
			s.ResidentialAddressID = &residentialAddress.ID
		}
		if payload.BackupMailingAddress != nil {
			backupMailingAddress := AddressModelFromPayload(payload.BackupMailingAddress)
			if verrs, err := dbConnection.ValidateAndCreate(backupMailingAddress); verrs.HasAny() || err != nil {
				if verrs.HasAny() {
					responseVErrors = verrs
				} else {
					responseError = err
				}
				transactionError = errors.New("Failed saving backup mailing address model")
				return transactionError
			}
			s.BackupMailingAddressID = &backupMailingAddress.ID
		}

		if verrs, err := dbConnection.ValidateAndUpdate(s); verrs.HasAny() || err != nil {
			if verrs.HasAny() {
				responseVErrors = verrs
			} else {
				responseError = err
			}
			transactionError = errors.New("Failed saving service member model")
		}

		return transactionError
	})
	return responseVErrors, responseError
}

// IsProfileComplete checks if the profile has been completely filled out
func (s *ServiceMember) IsProfileComplete() bool {

	// The following fields are required to be set for a profile to be complete
	if s.Edipi == nil {
		return false
	}
	if s.Branch == nil {
		return false
	}
	if s.Rank == nil {
		return false
	}
	if s.FirstName == nil {
		return false
	}
	if s.LastName == nil {
		return false
	}
	if s.Telephone == nil {
		return false
	}
	if s.PersonalEmail == nil {
		return false
	}
	if s.PhoneIsPreferred == nil && s.TextMessageIsPreferred == nil && s.EmailIsPreferred == nil {
		return false
	}
	if s.ResidentialAddress == nil {
		return false
	}
	// TODO: add check for station, SSN, and backup contacts
	// All required fields have a set value
	return true
}
