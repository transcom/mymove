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
	Model    interface{}
	Type     *CustomType
	LinkOnly bool
}

// CustomType is a string that represents what kind of customization it is
type CustomType string

// Model customizations will each have a "type" here
// This does not have to match the model type but generally will
// You can have CustomType like ResidentialAddress to define specifically
// where this address will get created and nested
var control CustomType = "Control"
var Address CustomType = "Address"
var User CustomType = "User"
var ServiceMember CustomType = "ServiceMember"
var OfficeUser CustomType = "OfficeUser"

// defaultTypesMap allows us to assign CustomTypes for most default types
var defaultTypesMap = map[string]CustomType{
	"models.Address":       Address,
	"models.OfficeUser":    OfficeUser,
	"models.ServiceMember": ServiceMember,
	"models.User":          User,
}

// Instead of nesting structs, we create specific CustomTypes here to give devs
// a code-completion friendly way to select the right type

// addressGroup is a grouping of all address related fields
type addressGroup struct {
	PickupAddress            CustomType
	DeliveryAddress          CustomType
	SecondaryDeliveryAddress CustomType
	ResidentialAddress       CustomType
}

// Addresses is the struct to access the various fields externally
var Addresses = addressGroup{
	PickupAddress:            "PickupAddress",
	DeliveryAddress:          "DeliveryAddress",
	SecondaryDeliveryAddress: "SecondaryDeliveryAddress",
	ResidentialAddress:       "ResidentialAddress",
}

// dimensionGroup is a grouping of all the Dimension related fields
type dimensionGroup struct {
	CrateDimension CustomType
	ItemDimension  CustomType
}

// Dimensions is the struct to access the fields externally
var Dimensions = dimensionGroup{
	// MTOServiceItems may include:
	CrateDimension: "CrateDimension",
	ItemDimension:  "ItemDimension",
}

// dutyLocationsGroup is a grouping of all the duty location related fields
type dutyLocationsGroup struct {
	OriginDutyLocation CustomType
	NewDutyLocation    CustomType
}

// DutyLocations is the struct to access the fields externally
var DutyLocations = dutyLocationsGroup{
	// Orders may include:
	OriginDutyLocation: "OriginDutyLocation",
	NewDutyLocation:    "NewDutyLocation",
}

// controlObject is a struct used to control the global behavior of a
// set of customizations
type controlObject struct {
	isValid bool // has this set of customizations been validated
}

// Trait is a function that returns a set of customizations
// Every Trait should start with GetTrait for discoverability
type Trait func() []Customization

// assignType uses the model name to assign the CustomType
// if it's already assigned, do not reassign
func assignType(custom *Customization) error {
	if custom.Type != nil {
		return nil
	}
	// Get the model and check that it's a struct
	model := custom.Model
	mv := reflect.ValueOf(model)
	if mv.Kind() != reflect.Struct {
		return fmt.Errorf("Customization.Model field had type %v - should contain a struct", mv.Kind())
	}
	// Get the model type and find the default type
	typestring, ok := defaultTypesMap[mv.Type().String()]
	if ok {
		custom.Type = &typestring
	} else {
		return fmt.Errorf("Customization.Model field had type %v which is not supported in defaultTypesMap", mv.Type().String())
	}
	return nil
}

// setDefaultTypes assigns types to all customizations in the list provided
func setDefaultTypes(clist []Customization) {
	for idx := 0; idx < len(clist); idx++ {
		if err := assignType(&clist[idx]); err != nil {
			log.Panic(err.Error())
		}
	}
}

// setDefaultTypesTraits assigns types to all customizations in the traits
//func setDefaultTypesTraits()

// setupCustomizations prepares the customizations customs for the factory
// by applying and merging the traits.
// customs is a slice that will be modified by setupCustomizations.
//
// - Ensures a control object has been created
// - Assigns default types to all default customizations
// - Merges customizations and traits
// - Ensure there's only one customization per type
func setupCustomizations(customs []Customization, traits []Trait) []Customization {

	// If a valid control object does not exist, create
	_, controlCustom := findCustomWithIdx(customs, control)
	if controlCustom == nil {
		controlCustom = &Customization{
			Model: controlObject{
				isValid: false,
			},
			Type: &control,
		}
		customs = append(customs, *controlCustom)
	}
	// If it exists and is valid, return, this list has been setup and validated
	controller := controlCustom.Model.(controlObject)
	if controller.isValid {
		return customs
	}

	// If not valid:
	// Merge customizations with traits (also sets default types)
	customs = mergeCustomization(customs, traits)
	// Ensure unique customizations
	err := isUnique(customs)
	if err != nil {
		controller.isValid = false
		log.Panic(err)
	}
	// Store the validation result
	controller.isValid = true
	return customs

}

// isUnique ensures there's only one customization per type.
// Requirement: All customizations should already have a type assigned.
func isUnique(customs []Customization) error {
	// validate that there are no repeat CustomTypes
	m := make(map[CustomType]int)
	for i, custom := range customs {
		if custom.Type == nil {
			log.Panic("All customizations should have type.")
		}
		// if custom type already exists
		idx, exists := m[*custom.Type]
		if exists {
			return fmt.Errorf("Found more than one instance of %s Customization at index %d and %d",
				*custom.Type, idx, i)
		}
		// Add to hashmap
		m[*custom.Type] = i
	}
	return nil

}

// findCustomWithIdx is a helper function to find a customization of a specific type and its index
func findCustomWithIdx(customs []Customization, customType CustomType) (int, *Customization) {
	for i, custom := range customs {
		if custom.Type != nil && *custom.Type == customType {
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
// Required : All customizations MUST have types
//
// The order of application is
//   - Earlier traits override later traits in the trait list
//   - Customizations override the traits
//
// So if you have [trait1, trait2] customization
// and all three contain the same object:
//   - trait 1 will override trait 2 (so start with the highest priority)
//   - customization will override trait 2
func mergeCustomization(customs []Customization, traits []Trait) []Customization {
	// Make sure all customs have a proper type
	setDefaultTypes(customs)

	// Iterate list of traits, each could return a list of customizations
	for _, trait := range traits {

		// Get customizations and set default types
		traitCustomizations := trait()
		setDefaultTypes(traitCustomizations)

		// for each trait custom, merge or replace the one in user supplied customizations
		for _, traitCustom := range traitCustomizations {
			j, callerCustom := findCustomWithIdx(customs, *traitCustom.Type)
			if callerCustom != nil {
				// If a customization is marked as LinkOnly, it means we use that precreated object
				// Therefore we can't merge a trait with it, as we don't update fields on pre-created
				// objects. So we only merge if LinkOnly is false.
				if !callerCustom.LinkOnly {
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
// TBD should we validate again here?
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
		newCustoms[idx].Type = &to
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
		log.Panic(fmt.Errorf("Errors encountered in customization for %s: %w", *custom.Type, err))
	}
	return custom
}

// checkNestedModels ensures we have no nested models.
// - If a field is a struct that is a model, it should be empty
// - If a field is a pointer to struct, and that struct is a model it should be nil
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

func mustCreate(db *pop.Connection, model interface{}) {
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

// Source chars for random string
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// Returns a random alphanumeric string of specified length
func makeRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		randInt, err := random.GetRandomInt(len(letterBytes))
		if err != nil {
			log.Panicf("failed to create random string %v", err)
			return ""
		}
		b[i] = letterBytes[randInt]

	}
	return string(b)
}
