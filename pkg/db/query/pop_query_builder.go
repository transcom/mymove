package query

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gobuffalo/pop"
)

// PopQueryBuilder is a wrapper aroudn pop
// with more flexible query patterns to MilMove
type PopQueryBuilder struct {
	db *pop.Connection
}

// NewPopQueryBuilder returns a new query builder implemented with pop
// constructor is for Dependency Injection frameworks requiring a function instead of struct
func NewPopQueryBuilder(db *pop.Connection) *PopQueryBuilder {
	return &PopQueryBuilder{db}
}

// InvalidInputError is an error for when query inputs are incorrect
type InvalidInputError struct {
	input []string
}

// Error returns the error message
func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("%v is not valid input", e.input)
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

// FetchOne fetches a single model record using pop's First method
// Will return error if model is not pointer to struct
func (p *PopQueryBuilder) FetchOne(model interface{}, filters map[string]interface{}) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New("Model should be pointer to struct")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("Model should be pointer to struct")
	}
	query := p.db.Q()
	invalidFields := make([]string, 0)
	for field, value := range filters {
		column, ok := getDBColumn(t, field)
		if !ok {
			invalidFields = append(invalidFields, field)
		}
		columnQuery := fmt.Sprintf("%s = ?", column)
		query = query.Where(columnQuery, value)
	}
	if len(invalidFields) != 0 {
		return &InvalidInputError{invalidFields}
	}
	return query.First(model)
}

// FetchMany fetches multiple model records using pop's All method
// Will return error if model is not pointer to slice of structs
func (p *PopQueryBuilder) FetchMany(model interface{}, filters map[string]interface{}) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New("Model should be pointer to slice of structs")
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return errors.New("Model should be pointer to slice of structs")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("Model should be pointer to slice of structs")
	}
	query := p.db.Q()
	invalidFields := make([]string, 0)
	for field, value := range filters {
		column, ok := getDBColumn(t, field)
		if !ok {
			invalidFields = append(invalidFields, field)
		}
		columnQuery := fmt.Sprintf("%s = ?", column)
		query = query.Where(columnQuery, value)
	}
	if len(invalidFields) != 0 {
		return &InvalidInputError{invalidFields}
	}
	return query.All(model)
}
