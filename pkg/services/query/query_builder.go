package query

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/services"
)

// allowed comparators for this query builder implementation
const equals = "="
const greaterThan = ">"

// Error message constants
const fetchManyReflectionMessage = "Model should be pointer to slice of structs"
const fetchOneReflectionMessage = "Model should be pointer to struct"

// Builder is a wrapper around pop
// with more flexible query patterns to MilMove
type Builder struct {
	db *pop.Connection
}

// NewQueryBuilder returns a new query builder implemented with pop
// constructor is for Dependency Injection frameworks requiring a function instead of struct
func NewQueryBuilder(db *pop.Connection) *Builder {
	return &Builder{db}
}

// Lookup to check if a specific string is inside the db field tags of the type
func getDBColumn(t reflect.Type, field string) (string, bool) {
	for i := 0; i < t.NumField(); i++ {
		dbTag, ok := t.Field(i).Tag.Lookup("db")
		if ok && dbTag == field {
			return dbTag, true
		}
	}
	return "", false
}

// check that we have a valid comparator
func getComparator(comparator string) (string, bool) {
	switch comparator {
	case equals:
		return equals, true
	case greaterThan:
		return greaterThan, true
	default:
		return "", false
	}
}

func filteredQuery(query *pop.Query, filters []services.QueryFilter, t reflect.Type) (*pop.Query, error) {
	invalidFields := make([]string, 0)
	for _, f := range filters {
		column, ok := getDBColumn(t, f.Column())
		if !ok {
			invalidFields = append(
				invalidFields,
				fmt.Sprintf("%s %s", f.Column(), f.Comparator()),
			)
		}
		comparator, ok := getComparator(f.Comparator())
		if !ok {
			invalidFields = append(
				invalidFields,
				fmt.Sprintf("%s %s", f.Column(), f.Comparator()),
			)
			continue
		}
		// Column lookup should always adhere to SQL injection input validations
		// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.md#defense-option-3-whitelist-input-validation
		columnQuery := fmt.Sprintf("%s %s ?", column, comparator)
		query = query.Where(columnQuery, f.Value())
	}
	if len(invalidFields) != 0 {
		return query, fmt.Errorf("%v is not valid input", invalidFields)
	}
	return query, nil
}

// FetchOne fetches a single model record using pop's First method
// Will return error if model is not pointer to struct
func (p *Builder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New(fetchOneReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New(fetchOneReflectionMessage)
	}
	query := p.db.Q()
	query, err := filteredQuery(query, filters, t)
	if err != nil {
		return err
	}
	return query.First(model)
}

// FetchMany fetches multiple model records using pop's All method
// Will return error if model is not pointer to slice of structs
func (p *Builder) FetchMany(model interface{}, filters []services.QueryFilter) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New(fetchManyReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return errors.New(fetchManyReflectionMessage)
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New(fetchManyReflectionMessage)
	}
	query := p.db.Q()
	query, err := filteredQuery(query, filters, t)
	if err != nil {
		return err
	}
	return query.All(model)
}
