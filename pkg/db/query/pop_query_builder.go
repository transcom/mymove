package query

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

const equals = "="

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

// check that we have a valid comparator
func getComparator(comparator string) (string, bool) {
	switch comparator {
	case equals:
		return equals, true
	default:
		return "", false
	}
}

func filteredQuery(query *pop.Query, filters []filter, t reflect.Type) (*pop.Query, error) {
	invalidFields := make([]string, 0)
	for _, f := range filters {
		column, ok := getDBColumn(t, f.Column())
		if !ok {
			invalidFields = append(invalidFields, f.Column())
		}
		comparator, ok := getComparator(f.Comparator())
		if !ok {
			invalidFields = append(invalidFields, f.Column())
		}
		columnQuery := fmt.Sprintf("%s %s ?", column, comparator)
		query = query.Where(columnQuery, f.Value())
	}
	if len(invalidFields) != 0 {
		return query, &InvalidInputError{invalidFields}
	}
	return query, nil
}

// FetchOne fetches a single model record using pop's First method
// Will return error if model is not pointer to struct
func (p *PopQueryBuilder) FetchOne(model interface{}, filters ...filter) error {
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Ptr {
		return errors.New("Model should be pointer to struct")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("Model should be pointer to struct")
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
func (p *PopQueryBuilder) FetchMany(model interface{}, filters ...filter) error {
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
	query, err := filteredQuery(query, filters, t)
	if err != nil {
		return err
	}
	return query.All(model)
}

// EqualsUUID creates a PopFilter based on the given column and UUID
func (p *PopQueryBuilder) EqualsUUID(column string, value uuid.UUID) PopFilter {
	return PopFilter{
		column:     column,
		comparator: equals,
		value:      value.String(),
	}
}

type filter interface {
	Column() string
	Comparator() string
	Value() string
}

// PopFilter is a filter used by the pop query builder
type PopFilter struct {
	column     string
	comparator string
	value      string
}

// Column returns the filter's column as a string
func (f PopFilter) Column() string {
	return f.column
}

// Comparator returns the filter's comparator as a string
func (f PopFilter) Comparator() string {
	return f.comparator
}

// Value returns the filter's value as a string
func (f PopFilter) Value() string {
	return f.value
}
