package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// ServiceMember is a user of type service member
type ServiceMember struct {
	ID                        uuid.UUID                           `json:"id" db:"id"`
	CreatedAt                 time.Time                           `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time                           `json:"updated_at" db:"updated_at"`
	UserID                    uuid.UUID                           `json:"user_id" db:"user_id"`
	User                      User                                `belongs_to:"user"`
	Edipi                     *string                             `json:"edipi" db:"edipi"`
	Branch                    *internalmessages.MilitaryBranch    `json:"branch" db:"branch"`
	Rank                      *internalmessages.ServiceMemberRank `json:"rank" db:"rank"`
	FirstName                 *string                             `json:"first_name" db:"first_name"`
	MiddleInitial             *string                             `json:"middle_initial" db:"middle_initial"`
	LastName                  *string                             `json:"last_name" db:"last_name"`
	Suffix                    *string                             `json:"suffix" db:"suffix"`
	Telephone                 *string                             `json:"telephone" db:"telephone"`
	SecondaryTelephone        *string                             `json:"secondary_telephone" db:"secondary_telephone"`
	PersonalEmail             *string                             `json:"personal_email" db:"personal_email"`
	PhoneIsPreferred          *bool                               `json:"phone_is_preferred" db:"phone_is_preferred"`
	SecondaryPhoneIsPreferred *bool                               `json:"secondary_phone_is_preferred" db:"secondary_phone_is_preferred"`
	EmailIsPreferred          *bool                               `json:"email_is_preferred" db:"email_is_preferred"`
	ResidentialAddressID      *uuid.UUID                          `json:"residential_address_id" db:"residential_address_id"`
	ResidentialAddress        *Address                            `belongs_to:"address"`
	BackupMailingAddressID    *uuid.UUID                          `json:"backup_mailing_address_id" db:"backup_mailing_address_id"`
	BackupMailingAddress      *Address                            `belongs_to:"address"`
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

// ServiceMemberResult is returned by GetServiceMemberForUser and encapsulates whether the call succeeded and why it failed.
type ServiceMemberResult struct {
	valid         bool
	errorCode     FetchError
	serviceMember ServiceMember
}

// IsValid indicates whether the ServiceMemberResult is valid.
func (m ServiceMemberResult) IsValid() bool {
	return m.valid
}

// ServiceMember returns the serviceMember if and only if the serviceMember was correctly fetched
func (m ServiceMemberResult) ServiceMember() ServiceMember {
	if !m.valid {
		zap.L().Fatal("Check if this isValid before accessing the ServiceMember()!")
	}
	return m.serviceMember
}

// ErrorCode returns the error if and only if the serviceMember was not correctly fetched
func (m ServiceMemberResult) ErrorCode() FetchError {
	if m.valid {
		zap.L().Fatal("Check that this !isValid before accessing the ErrorCode()!")
	}
	return m.errorCode
}

// NewInvalidServiceMemberResult creates an invalid ServiceMemberResult
func NewInvalidServiceMemberResult(errorCode FetchError) ServiceMemberResult {
	return ServiceMemberResult{
		errorCode: errorCode,
	}
}

// NewValidServiceMemberResult creates a valid ServiceMemberResult
func NewValidServiceMemberResult(serviceMember ServiceMember) ServiceMemberResult {
	return ServiceMemberResult{
		valid:         true,
		serviceMember: serviceMember,
	}
}

// GetServiceMemberForUser returns a serviceMember only if it is allowed for the given user to access that serviceMember.
// If the user is not authorized to access that serviceMember, it behaves as if no such serviceMember exists.
func GetServiceMemberForUser(db *pop.Connection, userID uuid.UUID, id uuid.UUID) (ServiceMemberResult, error) {
	var result ServiceMemberResult
	var serviceMember ServiceMember
	err := db.Eager().Find(&serviceMember, id)
	if err != nil {
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			result = NewInvalidServiceMemberResult(FetchErrorNotFound)
			err = nil
		}
		// Otherwise, it's an unexpected err so we return that.
	} else {
		if serviceMember.UserID != userID {
			result = NewInvalidServiceMemberResult(FetchErrorForbidden)
		} else {
			result = NewValidServiceMemberResult(serviceMember)
		}
	}

	return result, err
}

// CreateOrUpdateServiceMemberWithAddresses takes a serviceMember with Address structs and coordinates saving it all in a transaction
func CreateOrUpdateServiceMemberWithAddresses(dbConnection *pop.Connection, serviceMember *ServiceMember) (*validate.Errors, error) {
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

		if transactionError == nil {
			serviceMember.ResidentialAddressID = GetAddressID(serviceMember.ResidentialAddress)
			serviceMember.BackupMailingAddressID = GetAddressID(serviceMember.BackupMailingAddress)

			if verrs, err := dbConnection.ValidateAndCreate(serviceMember); verrs.HasAny() || err != nil {
				transactionError = errors.New("Rollback The transaction")
				responseVErrors = verrs
				responseError = err
			}
		}

		return transactionError

	})

	return responseVErrors, responseError

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
		if payload.MiddleInitial != nil {
			s.MiddleInitial = payload.MiddleInitial
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
			s.PersonalEmail = payload.PersonalEmail
		}
		if payload.PhoneIsPreferred != nil {
			s.PhoneIsPreferred = payload.PhoneIsPreferred
		}
		if payload.SecondaryPhoneIsPreferred != nil {
			s.SecondaryPhoneIsPreferred = payload.SecondaryPhoneIsPreferred
		}
		if payload.EmailIsPreferred != nil {
			s.EmailIsPreferred = payload.EmailIsPreferred
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
	if s.PhoneIsPreferred == nil && s.SecondaryPhoneIsPreferred == nil && s.EmailIsPreferred == nil {
		return false
	}
	if s.ResidentialAddress == nil {
		return false
	}
	// TODO: add check for station, SSN, and backup contacts
	// All required fields have a set value
	return true
}
