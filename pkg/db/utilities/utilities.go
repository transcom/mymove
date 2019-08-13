package utilities

import (
	"errors"
	"reflect"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

const deletedAt = "DeletedAt"

// SoftDestroy deletes a record and all foreign key associations from the database
func SoftDestroy(c *pop.Connection, model interface{}) error {
	pop.Debug = true
	// destroy
	// use c for transaction
	verrs := validate.NewErrors()
	var err error

	//TODO check if the model is a model
	transactionError := c.Transaction(func(db *pop.Connection) error {
		modelValue := reflect.ValueOf(model)
		deletedAtField := modelValue.FieldByName(deletedAt)
		deletedAtValue := reflect.ValueOf(&deletedAtField).Elem()

		if deletedAtField.IsValid() && deletedAtValue.CanSet() {
			now := time.Now()
			reflectTime := reflect.ValueOf(now)
			deletedAtValue.Set(reflectTime)
			verrs, err = db.ValidateAndSave(model)

			if err != nil || verrs.HasAny() {
				return errors.New("error updating model")
			}
		} else {
			return errors.New("can not soft delete this model")
		}
		return nil
	})

	if transactionError != nil || verrs.HasAny() {
		return transactionError
	}
	return errors.New("error updating model")
}
