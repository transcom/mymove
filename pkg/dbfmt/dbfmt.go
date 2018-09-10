package dbfmt

import (
	"fmt"
	"reflect"

	"github.com/gobuffalo/uuid"
)

//PrettyModel

// TODO
// make sure pointers are not priting addresses

func recursivePrettyStringWithPadding(model interface{}, padding string) string {
	fmt.Println("RECURSE")

	prettyString := ""

	modelValue := reflect.Indirect(reflect.ValueOf(model))

	if !modelValue.IsValid() {
		return "<nil>"
	}

	if modelValue.Kind() == reflect.Slice {
		return "[REDACTED SLICE]"
	}

	modelType := modelValue.Type()
	fmt.Println("TPYE", modelType)

	// should check that we are a struct here.
	for i := 0; i < modelValue.NumField(); i++ {
		valField := modelValue.Field(i)
		typeField := modelType.Field(i)

		var fieldRep interface{}

		if typeField.Name == "ID" {
			id := valField.Interface().(uuid.UUID)
			fmt.Println("THATS A UUID", id)
			emptyUUID := uuid.UUID{}
			if id == emptyUUID {
				// fmt.Println("DONT LOAD")
				return "<<not loaded>>"
			}
		}

		if typeField.Tag.Get("belongs_to") != "" {
			// prettyString += fmt.Sprintf("%s %s = %v  (%v)\n", typeField.Name, valField.Type(), "[REDACTED]", typeField.Tag.Get("belongs_to"))
			fmt.Println("REUCRSING", typeField.Name)
			fieldRep = recursivePrettyStringWithPadding(valField.Interface(), padding+"    ")
		} else if typeField.Tag.Get("has_many") != "" {
			fieldRep = "[REDACTED MANY]"
		} else {
			fieldRep = valField.Interface()
		}

		prettyString += fmt.Sprintf("%s%s %s = %v\n", padding, typeField.Name, valField.Type(), fieldRep)
	}

	return prettyString
}

func PrettyString(model interface{}) string {

	fmt.Println("HIHIHIHIHIHIHIHIHIHIHIHI")

	return recursivePrettyStringWithPadding(model, "")
}
