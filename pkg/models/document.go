package models

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
)

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: d.ServiceMemberID, Name: "ServiceMemberID"},
	), nil
}

type popDocumentDB struct {
	db *pop.Connection
}

// NewDocumentDB is the DI provider to create a pop based ServiceMemberDB
func NewDocumentDB(db *pop.Connection) DocumentDB {
	return &popDocumentDB{db}
}

// Fetch does a simple eager load of a Document using POP
func (pdb *popDocumentDB) Fetch(id uuid.UUID) (*Document, error) {
	var document Document
	err := pdb.db.Q().Eager().Find(&document, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	return &document, nil
}

// FetchUpload does a simple
