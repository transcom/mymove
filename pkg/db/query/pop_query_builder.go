package query

import (
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
func (p *PopQueryBuilder) FetchOne(model interface{}, field string, value interface{}) error {
	// todo: pointer check on type
	column, ok := getDBColumn(reflect.TypeOf(model).Elem(), field)
	if !ok {
		return &InvalidInputError{[]string{field}}
	}
	columnQuery := fmt.Sprintf("%s = ?", column)
	query := p.db.Where(columnQuery, value)
	return query.First(model)
}

// FetchMany fetches multiple model records using pop's All method
func (p *PopQueryBuilder) FetchMany(model interface{}, filters map[string]interface{}) error {
	query := p.db.Q()
	invalidFields := make([]string, 0)
	t := reflect.TypeOf(model).Elem().Elem() // todo: add slice check
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
