package dbtools

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/pop/v5"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
)

// NewTableFromSliceCreator is the public constructor for a TableFromSliceCreator using Pop
func NewTableFromSliceCreator(db *pop.Connection, logger Logger, isTemp bool, dropIfExists bool) services.TableFromSliceCreator {
	return &tableFromSliceCreator{
		db:           db,
		logger:       logger,
		isTemp:       isTemp,
		dropIfExists: dropIfExists,
	}
}

// tableFromSliceCreator is a service object to create/populate a table from a slice
type tableFromSliceCreator struct {
	db           *pop.Connection
	logger       Logger
	isTemp       bool
	dropIfExists bool
}

// CreateTableFromSlice creates and populates a table from a slice of structs
func (c tableFromSliceCreator) CreateTableFromSlice(slice interface{}) error {
	// Ensure we've got a slice or an array.
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Slice {
		return errors.New(fmt.Sprintf("Parameter must be slice or array, but got %v", sliceType))
	}

	// Ensure that the element type of the slice is a struct.
	elementType := sliceType.Elem()
	if elementType.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("Elements of slice must be type struct, but got %v", elementType))
	}

	// Create the table based on the struct
	var builder strings.Builder
	builder.WriteString("CREATE ")
	if c.isTemp {
		builder.WriteString("TEMP ")
	}
	builder.WriteString("TABLE ")
	tableName := name.Tableize(elementType.Name())
	builder.WriteString(tableName)
	builder.WriteString("(")
	numFields := elementType.NumField()
	fieldDbNames := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		field := elementType.Field(i)
		if field.Type.Kind() != reflect.String {
			return errors.New(fmt.Sprintf("All fields of struct must be string, but field %v is %v", field.Name, field.Type))
		}
		fieldDbNames[i] = field.Tag.Get("db")
		if len(fieldDbNames[i]) == 0 {
			fieldDbNames[i] = strings.ToLower(field.Name)
		}
		builder.WriteString(fieldDbNames[i])
		builder.WriteString(" text NOT NULL")
		if i < numFields-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(");")
	createTableQuery := builder.String()

	// check to see if we are already in a transaction
	if c.db.TX != nil {
		return c.executeSQL(c.db, tableName, createTableQuery, fieldDbNames, numFields, slice)
	}
	return c.db.Transaction(func(tx *pop.Connection) error {
		return c.executeSQL(tx, tableName, createTableQuery, fieldDbNames, numFields, slice)
	})
}

func (c tableFromSliceCreator) executeSQL(tx *pop.Connection, tableName string, createTableQuery string, fieldDbNames []string, numFields int, slice interface{}) error {
	if c.dropIfExists {
		err := tx.RawQuery("DROP TABLE IF EXISTS " + tableName).Exec()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error dropping table: '%s'", tableName))
		}
	}
	err := tx.RawQuery(createTableQuery).Exec()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating table: '%s'", tableName))
	}

	// Put data into the table
	stmt, err := tx.TX.Prepare(pq.CopyIn(tableName, fieldDbNames...))
	if err != nil {
		return errors.Wrap(err, "Error preparing CopyIn statement")
	}

	sliceValue := reflect.ValueOf(slice)
	for i := 0; i < sliceValue.Len(); i++ {
		structValue := sliceValue.Index(i)

		fieldValues := make([]interface{}, numFields)
		for j := 0; j < numFields; j++ {
			fieldValue := structValue.Field(j)
			fieldValues[j] = fieldValue.Interface()
		}
		_, execErr := stmt.Exec(fieldValues...)
		if execErr != nil {
			return errors.Wrapf(execErr, "Error executing CopyIn statement with values %q", fieldValues)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return errors.Wrap(err, "Error flushing CopyIn statement")
	}

	if err := stmt.Close(); err != nil {
		return errors.Wrap(err, "Error closing CopyIn statement")
	}

	return nil
}
