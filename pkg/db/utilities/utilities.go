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
		// model and "delete" it and its associations
		// either do a raw query setting the deleted_at or use db.Save()
		// use reflect to get associations
		//
		modelValue := reflect.ValueOf(model)
		modelType := reflect.TypeOf(model)

		for field := 0; field < modelValue.NumField(); field++ {
			modelField := modelType.Field(field)
			fieldValue := reflect.ValueOf(modelField)
			name := modelField.Name

			if name == deletedAt && fieldValue.CanSet() {
				now := time.Now()
				reflectTime := reflect.ValueOf(now)
				fieldValue.Set(reflectTime)
				verrs, err = db.ValidateAndSave(model)

				if err != nil || verrs.HasAny() {
					return errors.New("error updating model")
				}
			}
			return nil
		}
		return errors.New("Rollback The transaction")
	})

	if transactionError != nil || verrs.HasAny() {
		return transactionError
	}
	return nil
}
