package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ServiceMember is a user of type service member
type ServiceMember struct {
	ID                        uuid.UUID  `json:"id" db:"id"`
	CreatedAt                 time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
	UserID                    uuid.UUID  `json:"user_id" db:"user_id"`
	User                      User       `belongs_to:"user"`
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
	ResidentialAddress        *Address   `belongs_to:"address"`
	BackupMailingAddressID    *uuid.UUID `json:"backup_mailing_address_id" db:"backup_mailing_address_id"`
	BackupMailingAddress      *Address   `belongs_to:"address"`
}

// TODO add func to evaluate whether profile is complete - add call to payload struct in handler

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
	err := db.Find(&serviceMember, id)
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

// ValidateServiceMemberOwnership validates that a user has a serviceMember that exists
func ValidateServiceMemberOwnership(db *pop.Connection, userID uuid.UUID, id uuid.UUID) (bool, bool) {
	exists := false
	userOwns := false
	var serviceMember ServiceMember
	err := db.Find(&serviceMember, id)
	if err == nil {
		exists = true
		// TODO: Handle case where more than one user is authorized to modify serviceMember
		if uuid.Equal(serviceMember.UserID, userID) {
			userOwns = true
		}
	}

	return exists, userOwns
}

// CreateServiceMemberWithAddresses takes a serviceMember with Address structs and coordinates saving it all in a transaction
func CreateServiceMemberWithAddresses(dbConnection *pop.Connection, serviceMember *ServiceMember) (*validate.Errors, error) {
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

// IsProfileComplete checks if the profile has been completely filled out
func (s *ServiceMember) IsProfileComplete() bool {
	fmt.Println("profile complete hit")
	// TODO: check if every field is not 0 value and return true if so
	return false
}
