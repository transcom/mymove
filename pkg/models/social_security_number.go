package models

import (
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

var ssnFormatRegex = regexp.MustCompile(`^\d{3}-\d{2}-\d{4}$`)

// BuildSocialSecurityNumber returns an *unsaved* SSN that has the ssn hash set based on the passed in raw ssn
func BuildSocialSecurityNumber(unencryptedSSN string) (*SocialSecurityNumber, *validate.Errors, error) {
	verrs := validate.NewErrors()
	if !ssnFormatRegex.MatchString(unencryptedSSN) {
		verrs.Add("social_secuirty_number", "SSN must be in 123-45-6789 format")
		return nil, verrs, nil
	}

	byteHash, err := bcrypt.GenerateFromPassword([]byte(unencryptedSSN), -1) // -1 chooses the default cost
	if err != nil {
		return nil, verrs, err
	}
	hash := string(byteHash)
	ssn := SocialSecurityNumber{
		EncryptedHash: hash,
	}
	return &ssn, verrs, nil
}

// Matches returns true if the encrypted_hahs matches the unencryptedSSN
func (s SocialSecurityNumber) Matches(unencryptedSSN string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(s.EncryptedHash), []byte(unencryptedSSN))
	return err == nil
}
