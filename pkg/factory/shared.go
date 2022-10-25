package factory

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/random"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Customization type is the building block for passing in customizations and traits
type Customization struct {
	Model     interface{}
	Type      CustomType
	ForceUUID bool
}

// CustomType is a string that represents what kind of customization it is
type CustomType string

// Model customizations will each have a "type" here
// This does not have to match the model type but generally will
// You can have CustomType like ResidentialAddress to define specifically
// where this address will get created and nested
const (
	control       CustomType = "Control"
	Address       CustomType = "Address"
	User          CustomType = "User"
	ServiceMember CustomType = "ServiceMember"
)

// Instead of nesting structs, we create specific CustomTypes here to give devs
// a code-completion friendly way to select the right type
type AddressesGroup struct {
	PickupAddress            CustomType
	DeliveryAddress          CustomType
	SecondaryDeliveryAddress CustomType
	ResidentialAddress       CustomType
}

var Addresses = AddressesGroup{
	PickupAddress:            "PickupAddress",
	DeliveryAddress:          "DeliveryAddress",
	SecondaryDeliveryAddress: "SecondaryDeliveryAddress",
	ResidentialAddress:       "ResidentialAddress",
}

type DimensionsGroup struct {
	CrateDimension CustomType
	ItemDimension  CustomType
}

var Dimensions = DimensionsGroup{
	// MTOServiceItems may include:
	CrateDimension: "CrateDimension",
	ItemDimension:  "ItemDimension",
}

type DutyLocationsGroup struct {
	OriginDutyLocation CustomType
	NewDutyLocation    CustomType
}

var DutyLocations = DutyLocationsGroup{
	// Orders may include:
	OriginDutyLocation: "OriginDutyLocation",
	NewDutyLocation:    "NewDutyLocation",
}

// Control is a struct used with CustomType Control to
// set flags on overall behaviour and status of the customizations
type controlObject struct {
	isValid bool // has this set of customizations been validated
	//stub    bool // if stub is false, only in-memory objects are created, not db
}

// Trait is a function that returns a set of customizations
// Every Trait should start with GetTrait for discoverability
type Trait func() []Customization

// validateCustomizations
func validateCustomizations(customs []Customization) ([]Customization, error) {
	_, controlCustom := findCustomWithIdx(customs, control)

	// if it does exist, and customization list has been validated already, return
	if controlCustom != nil {
		controls := controlCustom.Model.(controlObject)
		if controls.isValid {
			return customs, nil
		}
	} else {
		// If control object does not exist, create
		controlCustom = &Customization{
			Model: controlObject{},
			Type:  control,
		}
		customs = append(customs, *controlCustom)
	}
	controller := (*controlCustom).Model.(controlObject)
	// validate that there are no repeat model types
	m := make(map[CustomType]int)
	for i, custom := range customs {
		// if custom type already exists
		idx, exists := m[custom.Type]
		if exists {
			controller.isValid = false
			return customs, fmt.Errorf("Found more than one instance of %s Customization at index %d and %d",
				custom.Type, idx, i)
		}
		// Add to hashmap
		m[custom.Type] = i
	}
	// Store the validation result
	controller.isValid = true
	return customs, nil

}

// findCustomWithIdx is a helper function to find a customization of a specific type and its index
func findCustomWithIdx(customs []Customization, customType CustomType) (int, *Customization) {
	for i, custom := range customs {
		if custom.Type == customType {
			return i, &custom
		}
	}
	return -1, nil
}

// toStructPtr takes an interface wrapping a struct and
// returns an interface wrapping a pointer to the struct
// For e.g. interface{}(models.User) → interface{}(*models.User)
func toStructPtr(obj interface{}) interface{} {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp.Interface()
}

// toInterfacePtr takes an interface wrapping a pointer to a struct and
// returns an interface wrapping the struct
// For e.g. interface{}(*models.User) → interface{}(models.User)
func toInterfacePtr(obj interface{}) interface{} {
	rv := reflect.ValueOf(obj).Elem()
	return rv.Interface()
}

// mergeInterfaces transforms the interfaces to match what mergeModels expects
func mergeInterfaces(model1 interface{}, model2 interface{}) interface{} {
	// This wrapper function is needed because merge models expects the first
	// model to be an interface containing a pointer to struct.
	// This function converts
	//    an interface containing the struct
	//    → an interface containing a pointer to the struct
	modelPtr := toStructPtr(model1)
	testdatagen.MergeModels(modelPtr, model2)
	model := toInterfacePtr(modelPtr)
	return model
}

func hasID(model interface{}) bool {
	mv := reflect.ValueOf(model)

	// mv should be a model of type struct
	if mv.Kind() != reflect.Struct || !strings.HasPrefix(mv.Type().String(), "models.") {
		log.Panic("Expecting interface containing a model")
	}

	return !mv.FieldByName("ID").IsZero()
}

// mergeCustomization takes the original set of customizations
// and merges with the traits.
// The order of application is
//   - Earlier traits override later traits in the trait list
//   - Customizations override the traits
//
// So if you have [trait1, trait2] customization
// and all three contain the same object:
//   - trait 1 will override trait 2 (so start with the highest priority)
//   - customization will override trait 2
func mergeCustomization(customs []Customization, traits []Trait) []Customization {
	// Get a list of traits, each could return a list of customizations
	for _, trait := range traits {
		traitCustomizations := trait()

		// for each trait custom, merge or replace the one in user supplied customizations
		for _, traitCustom := range traitCustomizations {
			j, callerCustom := findCustomWithIdx(customs, traitCustom.Type)
			if callerCustom != nil {
				// If a customization has an ID, it means we use that precreated object
				// Therefore we can't merge a trait with it, as those fields will not get
				// updated.
				if !hasID(callerCustom.Model) {
					result := mergeInterfaces(traitCustom.Model, callerCustom.Model)
					callerCustom.Model = result
					customs[j] = *callerCustom
				}
			} else {
				customs = append(customs, traitCustom)
			}
		}
	}
	return customs
}

// Helper function for when you need to elevate the type of customization from say
// ResidentialAddress to Address before you call makeAddress
// This is a little finicky because we want to be careful not to harm the existing list
func convertCustomizationInList(customs []Customization, from CustomType, to CustomType) []Customization {
	if _, custom := findCustomWithIdx(customs, to); custom != nil {
		log.Panic(fmt.Errorf("A customization of type %s already exists", to))
	}
	if idx, custom := findCustomWithIdx(customs, from); custom != nil {
		// Create a slice in new memory
		var newCustoms []Customization
		// Populate with copies of objects
		newCustoms = append(newCustoms, customs...)
		// Update the type
		newCustoms[idx].Type = to
		return newCustoms
	}
	log.Panic(fmt.Errorf("No customization of type %s found", from))
	return nil
}

func findValidCustomization(customs []Customization, customType CustomType) *Customization {
	_, custom := findCustomWithIdx(customs, customType)
	if custom == nil {
		return nil
	}

	// Else check that the customization is valid
	if err := checkNestedModels(*custom); err != nil {
		log.Panic(fmt.Errorf("Errors encountered in customization for %s: %w", custom.Type, err))
	}
	return custom
}

func checkNestedModels(c interface{}) error {

	// c IS THE CUSTOMIZATION, SHOULD NOT BE A POINTER
	// c SHOULD NOT BE A POINTER
	if reflect.ValueOf(c).Kind() == reflect.Pointer {
		return fmt.Errorf("function expects a struct, received a pointer")
	}

	// mv IS THE MODEL VALUE, SHOULD NOT BE EMPTY
	mv := reflect.ValueOf(c).FieldByName("Model") // get the interface
	if mv.IsNil() {
		return fmt.Errorf("customization must contain a model")
	}
	mv = mv.Elem() // get the model from the interface

	// mv SHOULD BE A STRUCT
	if mv.Kind() == reflect.Struct {
		numberOfFields := mv.NumField()
		mt := mv.Type() // get the model type

		// CHECK ALL FIELDS IN THE STRUCT
		for i := 0; i < numberOfFields; i++ {
			fieldName := mt.Field(i).Name
			field := mv.Field(i)

			// There are a couple conditions we want to check for
			// - If a field is a struct that is a model, it should be empty
			// - If a field is a pointer to struct, and that struct is a model it should be nil

			// IF A FIELD IS A MODELS STRUCT - SHOULD BE EMPTY
			ft := field.Type()
			if field.Kind() == reflect.Struct && strings.HasPrefix(ft.String(), "models.") {
				if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
					return fmt.Errorf("%s cannot be populated, no nested models allowed", fieldName)
				}
			}

			// IF A FIELD IS A POINTER TO A MODELS STRUCT - SHOULD ALSO BE EMPTY
			if field.Kind() == reflect.Pointer {
				nf := field.Elem()
				if !field.IsNil() && nf.Kind() == reflect.Struct && strings.HasPrefix(nf.Type().String(), "models.") {
					return fmt.Errorf("%s cannot be populated, no nested models allowed", fieldName)
				}
			}
		}
	}
	return nil
}

func mustCreate(db *pop.Connection, model interface{}, stub bool) {
	if stub {
		return
	}

	verrs, err := db.ValidateAndCreate(model)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("Validation errors encountered saving %#v: %v", model, verrs))
	}
}

// RandomEdipi creates a random Edipi for a service member
func RandomEdipi() string {
	low := 1000000000
	high := 9999999999
	randInt, err := random.GetRandomIntAddend(low, high)
	if err != nil {
		log.Panicf("Failure to generate random Edipi %v", err)
	}
	return strconv.Itoa(low + int(randInt))
}
