package query

import (
	"fmt"
	"reflect"

	"github.com/gobuffalo/pop"
)

type PopQueryBuilder struct {
	db *pop.Connection
}

func NewPopQueryBuilder(db *pop.Connection) PopQueryBuilder {
	return PopQueryBuilder{db}
}
func (p PopQueryBuilder) FetchOne(model interface{}, field string, value string) error {
	columnNames := make([]string, 0)
	val := reflect.ValueOf(model)
	for i := 0; i < val.Type().NumField(); i++ {
		columnNames = append(dbNames, val.Type().Field(i).Tag.Get("db"))
	}
	columnNames.contains(field)
	column := fmt.Sprintf("%s=", columnNames[i])
	query := p.db.Where(column, " = $2", field, value)
	return query.First(model)
}
