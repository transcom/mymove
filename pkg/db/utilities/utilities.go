package utilities

import (
	"errors"
	"reflect"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"
)

const deletedAt = "DeletedAt"
const modelsPkgPath = "github.com/transcom/mymove/pkg/models"

// SoftDestroy soft deletes a record and all foreign key associations from the database
func SoftDestroy(c *pop.Connection, model interface{}) error {
	verrs := validate.NewErrors()
	var err error

	if !IsModel(model) {
		return errors.New("can only soft delete type model")
	}

	modelValue := reflect.ValueOf(model).Elem()
	deletedAtField := modelValue.FieldByName(deletedAt)

	if deletedAtField.IsValid() {
		if deletedAtField.CanSet() {
			now := time.Now()
			reflectTime := reflect.ValueOf(&now)
			deletedAtField.Set(reflectTime)
			verrs, err = c.ValidateAndSave(model)

			if err != nil || verrs.HasAny() {
				return errors.New("error updating model")
			}
		} else {
			return errors.New("can not soft delete this model")
		}
	} else {
		return errors.New("this model does not have deleted_at field")
	}

	associations := GetForeignKeyAssociations(c, model)
	if len(associations) > 0 {
		for _, association := range associations {
			err = SoftDestroy(c, association)
			if err != nil {
				return err
			}
		}
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
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is used later on
	//RA: Given the assigned variable is being used in a different line and the functions that are flagged by the linter are being used to assign variables
	//RA: that lead to error handling, there is no risk
	//RA Developer Status: False Positive
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	c.Load(model) // nolint:errcheck

	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()

	for pos := 0; pos < modelValue.NumField(); pos++ {
		fieldValue := modelValue.Field(pos)

		if fieldValue.CanInterface() {
			association := fieldValue.Interface()

			if association != nil {
				hasOneTag := modelType.Field(pos).Tag.Get("has_one")
				hasManyTag := modelType.Field(pos).Tag.Get("has_many")

				if hasOneTag != "" && GetHasOneForeignKeyAssociation(association) != nil {
					foreignKeyAssociations = append(foreignKeyAssociations, association)
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
