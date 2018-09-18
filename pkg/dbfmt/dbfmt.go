package dbfmt

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gobuffalo/uuid"
)

func recursivePrettyStringWithPadding(model interface{}, padding string) string {
	prettyString := ""

	modelValue := reflect.Indirect(reflect.ValueOf(model))

	if !modelValue.IsValid() {
		return "<nil>"
	}

	modelType := modelValue.Type()

	// Special cases where the default to_string is ugly
	switch modelType {
	case reflect.TypeOf(time.Time{}):
		layout := "2006-01-02 15:04:05"
		aTime := modelValue.Interface().(time.Time)
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
					prettyString += fmt.Sprintf("\n%s", padding+"    ")
				}
			}
			prettyString += fmt.Sprintf("\n%s]", padding)
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

		prettyString += "{\n"
		indentedPadding := padding + "    "

		for i := 0; i < modelValue.NumField(); i++ {
			valField := modelValue.Field(i)
			typeField := modelType.Field(i)

			var fieldRep interface{}
			fieldRep = recursivePrettyStringWithPadding(valField.Interface(), indentedPadding)

			prettyString += fmt.Sprintf("%s%-22s %v\n", indentedPadding, typeField.Name, fieldRep)
		}

		prettyString += fmt.Sprintf("%s}", padding)

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

// Println prints a log message and the pretty printed version of the model to the console.
func Println(messages ...interface{}) {

	prettyMessages := make([]interface{}, len(messages))
	for i, message := range messages {
		prettyMessages[i] = PrettyString(message)
	}

	fmt.Println(prettyMessages...)
}
