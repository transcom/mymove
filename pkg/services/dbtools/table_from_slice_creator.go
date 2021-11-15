package dbtools

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/lib/pq"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

// NewTableFromSliceCreator is the public constructor for a TableFromSliceCreator using Pop
func NewTableFromSliceCreator(isTemp bool, dropIfExists bool) services.TableFromSliceCreator {
	return &tableFromSliceCreator{
		isTemp:       isTemp,
		dropIfExists: dropIfExists,
	}
}

// tableFromSliceCreator is a service object to create/populate a table from a slice
type tableFromSliceCreator struct {
	isTemp       bool
	dropIfExists bool
}

// CreateTableFromSlice creates and populates a table from a slice of structs
func (c tableFromSliceCreator) CreateTableFromSlice(appCtx appcontext.AppContext, slice interface{}) error {
	// Ensure we've got a slice or an array.
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("Parameter must be slice or array, but got %v", sliceType)
	}

	// Ensure that the element type of the slice is a struct.
	elementType := sliceType.Elem()
	if elementType.Kind() != reflect.Struct {
		return fmt.Errorf("Elements of slice must be type struct, but got %v", elementType)
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
			return fmt.Errorf("All fields of struct must be string, but field %v is %v", field.Name, field.Type)
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

	return appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		return c.executeSQL(txnAppCtx, tableName, createTableQuery, fieldDbNames, numFields, slice)
	})
}

func (c tableFromSliceCreator) executeSQL(appCtx appcontext.AppContext, tableName string, createTableQuery string, fieldDbNames []string, numFields int, slice interface{}) error {
	if c.dropIfExists {
		err := appCtx.DB().RawQuery("DROP TABLE IF EXISTS " + tableName).Exec()
		if err != nil {
			return fmt.Errorf("Error dropping table: '%s': %w", tableName, err)
		}
	}
	err := appCtx.DB().RawQuery(createTableQuery).Exec()
	if err != nil {
		return fmt.Errorf("Error creating table: '%s': %w", tableName, err)
	}

	// Put data into the table
	stmt, err := appCtx.DB().TX.Prepare(pq.CopyIn(tableName, fieldDbNames...))
	if err != nil {
		return fmt.Errorf("Error preparing CopyIn statement: %w", err)
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
			return fmt.Errorf("Error executing CopyIn statement with values %q: %w", fieldValues, execErr)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("Error flushing CopyIn statement: %w", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("Error closing CopyIn statement: %w", err)
	}

	return nil
}
