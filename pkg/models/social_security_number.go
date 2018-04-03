package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"golang.org/x/crypto/bcrypt"
)

// SocialSecurityNumber represents an SSN in the database. It stores the SSN securely by hashing it.
type SocialSecurityNumber struct {
	ID            uuid.UUID `db:"id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	UserID        uuid.UUID `db:"user_id"`
	EncryptedHash string    `db:"encrypted_hash"`
}

// SocialSecurityNumbers is a list of SSNs
type SocialSecurityNumbers []SocialSecurityNumber

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s *SocialSecurityNumber) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.EncryptedHash, Name: "EncryptedHash"},
		&StringDoesNotContainSSN{Field: s.EncryptedHash, Name: "EncryptedHash"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (s *SocialSecurityNumber) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (s *SocialSecurityNumber) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ErrSSNBadFormat is returned if you attempted to hash an SSN not in the format '123-12-1234'
var ErrSSNBadFormat = errors.New("SSNs must be in the format '123-12-1234'")
var ssnFormatValidator = regexp.MustCompile(`^\d{3}-\d{2}-\d{4}$`)

// BuildSocialSecurityNumber returns an *unsaved* SSN that has the ssn hash set based on the passed in raw ssn
func BuildSocialSecurityNumber(userID uuid.UUID, unencryptedSSN string) (SocialSecurityNumber, error) {
	if !ssnFormatValidator.Match([]byte(unencryptedSSN)) {
		return SocialSecurityNumber{}, ErrSSNBadFormat
	}

	byteHash, err := bcrypt.GenerateFromPassword([]byte(unencryptedSSN), -1) // -1 chooses the default cost
	if err != nil {
		return SocialSecurityNumber{}, err
	}
	hash := string(byteHash)
	ssn := SocialSecurityNumber{
		UserID:        userID,
		EncryptedHash: hash,
	}
	return ssn, nil
}

// Matches returns true if the encrypted_hahs matches the unencryptedSSN
func (s SocialSecurityNumber) Matches(unencryptedSSN string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(s.EncryptedHash), []byte(unencryptedSSN))
	return err == nil
}
