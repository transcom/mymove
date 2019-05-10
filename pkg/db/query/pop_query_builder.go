package query

import (
	"fmt"
	"reflect"

	"github.com/gobuffalo/pop"
)

type popQueryBuilder struct {
	db *pop.Connection
}

func NewPopQueryBuilder(db *pop.Connection) popQueryBuilder {
	return popQueryBuilder{db}
}

type InvalidInputError struct {
	input []string
}

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

func (p popQueryBuilder) FetchOne(model interface{}, field string, value interface{}) error {
	// todo: pointer check on type
	column, ok := getDBColumn(reflect.TypeOf(model).Elem(), field)
	if !ok {
		return &InvalidInputError{[]string{field}}
	}
	columnQuery := fmt.Sprintf("%s = ?", column)
	query := p.db.Where(columnQuery,value)
	return query.First(model)
}

func (p popQueryBuilder) FetchMany(model interface{}, filters map[string]interface{}) error {
	query := p.db.Q()
	invalidFields := make([]string, 0)
	t := reflect.TypeOf(model).Elem().Elem() // todo: add slice check
	for field, value := range filters {
		column, ok := getDBColumn(t, field)
		if !ok {
			invalidFields = append(invalidFields, field)
		}
		columnQuery := fmt.Sprintf("%s = ?", column)
		query = query.Where(columnQuery,value)
	}
	if len(invalidFields) != 0 {
		return &InvalidInputError{invalidFields}
	}
	return query.All(model)
}

