package utilities

import (
	"errors"
	"reflect"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

const deletedAt = "DeletedAt"
const modelsPkgPath = "github.com/transcom/mymove/pkg/models"

// SoftDestroy deletes a record and all foreign key associations from the database
func SoftDestroy(c *pop.Connection, model interface{}) error {
	verrs := validate.NewErrors()
	var err error

	if !IsModel(model) {
		return errors.New("can only soft delete type model")
	}

	transactionError := c.Transaction(func(db *pop.Connection) error {
		modelValue := reflect.ValueOf(model).Elem()
		deletedAtField := modelValue.FieldByName(deletedAt)

		if deletedAtField.IsValid() {
			if deletedAtField.CanSet() {
				now := time.Now()
				reflectTime := reflect.ValueOf(&now)
				deletedAtField.Set(reflectTime)
				verrs, err = db.ValidateAndSave(model)

				if err != nil || verrs.HasAny() {
					return errors.New("error updating model")
				}
			} else {
				return errors.New("can not soft delete this model")
			}
		} else {
			return errors.New("this model does not have deleted_at property")
		}

		// TODO get associations and do this all over again recursively????
		return nil
	})

	if transactionError != nil || verrs.HasAny() {
		return transactionError
	}
	return nil
}

// IsModel verifies if the given interface is a model
func IsModel(model interface{}) bool {
	pkgPath := reflect.TypeOf(model).Elem().PkgPath()
	return pkgPath == modelsPkgPath
}

// GetForeignKeyAssociations fetches all the foreign key associations the model has
func GetForeignKeyAssociations(model interface{}) []interface{} {
	var foreignKeyAssociations []interface{}
	return foreignKeyAssociations
	// go through all the fields
	// fetch those of tag 'has_many'
}
