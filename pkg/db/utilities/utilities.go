package utilities

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/gofrs/uuid"
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
					fmt.Println(modelValue.Type())
					fmt.Println(err)
					fmt.Println(verrs)
					return errors.New("error updating model")
				}
			} else {
				return errors.New("can not soft delete this model")
			}
		} else {
			return errors.New("this model does not have deleted_at property")
		}

		associations := GetForeignKeyAssociations(c, model)
		for _, association := range associations {
			err = SoftDestroy(c, association)
			if err != nil {
				return err
			}
		}
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
func GetForeignKeyAssociations(c *pop.Connection, model interface{}) []interface{} {
	var foreignKeyAssociations []interface{}
	c.Load(model)

	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()

	for pos := 0; pos < modelValue.NumField(); pos++ {
		fieldValue := modelValue.Field(pos)

		// check if association is not nil and if they have an assigned ID
		// why is this working when passed through the SoftDestroy method but not here?
		// specifically reflect.ValueOf(association).Elem()
		// it gets angry about it here but not when passed through that method
		// is association being transformed somehow???
		// may be same issue when getting the address of the model within function instead of passing it along

		if fieldValue.CanInterface() {
			association := fieldValue.Interface()

			if association != nil {
				hasOneTag := modelType.Field(pos).Tag.Get("has_one")
				hasManyTag := modelType.Field(pos).Tag.Get("has_many")

				if hasOneTag != "" {
					associationValue := reflect.ValueOf(association).Elem()
					associationIDField := associationValue.FieldByName("ID")
					if associationIDField.CanInterface() && associationIDField.Interface() != uuid.Nil {
						foreignKeyAssociations = append(foreignKeyAssociations, fieldValue.Interface())
					}
				}

				if hasManyTag != "" {
					association := fieldValue.Interface()
					fmt.Println(association)
					// for object in objects
					foreignKeyAssociations = append(foreignKeyAssociations, fieldValue.Interface())
				}
			}
		}
	}
	fmt.Println(foreignKeyAssociations)
	return foreignKeyAssociations
}
