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
	fmt.Println(reflect.TypeOf(model))
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

		if fieldValue.CanInterface() {
			association := fieldValue.Interface()

			if association != nil {
				hasOneTag := modelType.Field(pos).Tag.Get("has_one")
				hasManyTag := modelType.Field(pos).Tag.Get("has_many")

				if hasOneTag != "" {
					foreignKeyAssociations = append(foreignKeyAssociations, GetHasOneForeignKeyAssociation(association))
				}

				if hasManyTag != "" {
					foreignKeyAssociations = append(foreignKeyAssociations, GetHasManyForeignKeyAssociations(association)...)
				}
			}
		}
	}
	return foreignKeyAssociations
}

// GetHasOneForeignKeyAssociation fetches the "has_one" foreign key association if not an empty model
func GetHasOneForeignKeyAssociation(model interface{}) interface{} {
	modelValue := reflect.ValueOf(model).Elem()
	idField := modelValue.FieldByName("ID")

	if idField.CanInterface() && idField.Interface() != uuid.Nil {
		return model
	}
	return nil
}

// GetHasManyForeignKeyAssociations fetches the "has_many" foreing key association if not an empty model
func GetHasManyForeignKeyAssociations(model interface{}) []interface{} {
	var hasManyForeignKeyAssociations []interface{}
	associations := reflect.ValueOf(model)
	for pos := 0; pos < associations.Len(); pos++ {
		association := associations.Index(pos)
		idField := association.FieldByName("ID")

		if idField.CanInterface() && idField.Interface() != uuid.Nil {
			associationPtr := association.Addr().Interface()
			hasManyForeignKeyAssociations = append(hasManyForeignKeyAssociations, associationPtr)
		}
	}
	return hasManyForeignKeyAssociations
}
