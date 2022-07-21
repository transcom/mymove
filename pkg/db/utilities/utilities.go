package utilities

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

const deletedAt = "DeletedAt"
const modelsPkgPath = "github.com/transcom/mymove/pkg/models"

// SoftDestroy soft deletes a record and all foreign key associations from the database
func SoftDestroy(c *pop.Connection, model interface{}) error {
	var verrs *validate.Errors
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

	associations, err := GetForeignKeyAssociations(c, model)

	if err != nil {
		return err
	}

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
func GetForeignKeyAssociations(c *pop.Connection, model interface{}) ([]interface{}, error) {
	var foreignKeyAssociations []interface{}

	err := c.Load(model)

	if err != nil {
		return nil, err
	}

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
	return foreignKeyAssociations, nil
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

// ExcludeDeletedScope is a chainable way to remove soft deleted models from your queries using Pop's Scope() command
// If you are filtering the same model that is being returned you can call this with no arguments and it should give
// you the desired behavior: db.Scope(utilities.ExcludedDeletedScope()).All(&mtoShipments).
//
// If you are filtering on a join table(s) or need to disambiguate multiple tables with deleted_at columns, then you
// will need to pass in the model(s) so we can derive the real table names:
// db.Scope(utilities.ExcludeDeletedScope(models.MTOShipment{})).
//     Join("mto_shipments", "mto_shipments.move_id = moves.id").
//     All(&moves)
//
// You won't be able to use this if you have given your table name(s) an alias, so just fall back to a normal where.
// You also cannot combine Scopes with RawQuery, which disregards any Join, Where, Scope, or Eager that may have been
// chained to the query.
//
// https://gobuffalo.io/documentation/database/scoping/
func ExcludeDeletedScope(models ...interface{}) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(models) == 0 {
			q.Where("deleted_at IS NULL")
		}
		for _, model := range models {
			var tableName string
			// Some models override the table name instead of the snake case pluralizing of the model name
			if modelTable, ok := model.(pop.TableNameAble); ok {
				tableName = modelTable.TableName()
			} else {
				modelType := reflect.TypeOf(model)
				tableName = flect.Pluralize(flect.Underscore(modelType.Name()))
			}
			q.Where(fmt.Sprintf("%s.deleted_at IS NULL", tableName))
		}
		return q
	}
}
