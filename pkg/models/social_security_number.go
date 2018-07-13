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

// SetEncryptedHash correctly sets the encrypted hash for the given SSN
func (s *SocialSecurityNumber) SetEncryptedHash(unencryptedSSN string) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	if !ssnFormatRegex.MatchString(unencryptedSSN) {
		verrs.Add("social_security_number", "SSN must be in 123-45-6789 format")
		return verrs, nil
	}

	workFactor := 13 // This was timed on the staging infrastructure to be ~1 second per hash. ( .76 sph actual)
	// # of iterations is 2 ^ workFactor
	// ~1 second per hash was picked based on this description:
	// https://security.stackexchange.com/questions/3959/recommended-of-iterations-when-using-pkbdf2-sha256/3993
	byteHash, err := bcrypt.GenerateFromPassword([]byte(unencryptedSSN), workFactor)
	if err != nil {
		return verrs, err
	}
	hash := string(byteHash)
	s.EncryptedHash = hash

	return verrs, nil
}

// BuildSocialSecurityNumber returns an *unsaved* SSN that has the ssn hash set based on the passed in raw ssn
func BuildSocialSecurityNumber(unencryptedSSN string) (*SocialSecurityNumber, *validate.Errors, error) {
	ssn := SocialSecurityNumber{}
	verrs, err := ssn.SetEncryptedHash(unencryptedSSN)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}
	return &ssn, verrs, nil
}

// Matches returns true if the encrypted_hash matches the unencryptedSSN
func (s SocialSecurityNumber) Matches(unencryptedSSN string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(s.EncryptedHash), []byte(unencryptedSSN))
	return err == nil
}
