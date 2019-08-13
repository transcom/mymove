package utilities

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

const deletedAt = "DeletedAt"

// SoftDestroy deletes a record and all foreign key associations from the database
func SoftDestroy(c *pop.Connection, model interface{}) error {
	verrs := validate.NewErrors()
	var err error

	fmt.Println(reflect.TypeOf(model))

	//TODO check if the model is a model
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
		return nil
	})

	if transactionError != nil || verrs.HasAny() {
		return transactionError
	}
	return nil
}
