package dbfmt

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gobuffalo/uuid"
)

// TODO
// Print just the object graph

func recursivePrettyStringWithPadding(model interface{}, padding string) string {
	prettyString := ""

	modelValue := reflect.Indirect(reflect.ValueOf(model))

	if !modelValue.IsValid() {
		return "<nil>"
	}

	modelType := modelValue.Type()

	switch modelType {
	case reflect.TypeOf(time.Time{}):
		layout := "2006-01-02 15:04:05"
		aTime := model.(time.Time)
		return aTime.Format(layout)
	}

	switch modelValue.Kind() {
	case reflect.Slice:
		if modelValue.Len() == 0 {
			prettyString += "[]"
		} else {
			prettyString += "[ "
			for i := 0; i < modelValue.Len(); i++ {
				sub := modelValue.Index(i)
				prettyString += recursivePrettyStringWithPadding(sub.Interface(), padding+"    ")
				if i < modelValue.Len()-1 {
					prettyString += ",\n"
				}
			}
			prettyString += fmt.Sprintf("\n%s%20s]", padding, "")
		}
	case reflect.Struct:
		// check to see if it is one of our models. We don't want to be recursing down otherwise.
		idField := modelValue.FieldByName("ID")
		if idField == (reflect.Value{}) {
			return fmt.Sprintf("%v", modelValue.Interface())
		}

		// If we are a struct with a field named ID and that field is the default
		// value, then this model hasn't been loaded and we won't display it.
		id := idField.Interface().(uuid.UUID)
		emptyUUID := uuid.UUID{}
		if id == emptyUUID {
			return "<<not loaded>>"
		}

		prettyString += "\n"

		for i := 0; i < modelValue.NumField(); i++ {
			valField := modelValue.Field(i)
			typeField := modelType.Field(i)

			var fieldRep interface{}
			fieldRep = recursivePrettyStringWithPadding(valField.Interface(), padding+"    ")

			prettyString += fmt.Sprintf("%s%-20s = %v\n", padding, typeField.Name, fieldRep)
		}

		// gotta remove the last \n
		prettyString = strings.TrimSuffix(prettyString, "\n")

	default:
		return fmt.Sprintf("%v", modelValue.Interface())
	}

	return prettyString
}

// PrettyString returns a cleaned up model diagram
func PrettyString(model interface{}) string {
	prettyString := recursivePrettyStringWithPadding(model, "")
	return strings.TrimSpace(prettyString)
}
