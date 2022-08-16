package testdatagen

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
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

// GetTraitFunc is a function that returns a set of customizations
// Every GetTraitFunc should start with GetTrait for discoverability
type GetTraitFunc func() []Customization

// GetTraitActiveUser returns a customization to enable active on a user
func GetTraitActiveUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
			},
			Type: User,
		},
	}
}

// GetTraitNavy is a sample GetTraitFunc
func GetTraitNavy() []Customization {
	navy := models.AffiliationNAVY
	return []Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &navy,
			},
			Type: ServiceMember,
		},
	}
}

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
	mergeModels(modelPtr, model2)
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
//     - Earlier traits override later traits in the trait list
//     - Customizations override the traits
//
// So if you have [trait1, trait2] customization
// and all three contain the same object:
//     - trait 1 will override trait 2 (so start with the highest priority)
//     - customization will override trait 2
func mergeCustomization(customs []Customization, traits []GetTraitFunc) []Customization {
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

// UserMaker is the base maker function to create a user
// MYTODO Instead of error (not useful) can we return a list of the created objects?
func UserMaker(db *pop.Connection, customs []Customization, traits []GetTraitFunc) (models.User, error) {
	customs = mergeCustomization(customs, traits)
	customs, err := validateCustomizations(customs)
	if err != nil {
		log.Panic(err)
	}

	// Find user assertion and convert to models user
	var cUser models.User
	if result := findValidCustomization(customs, User); result != nil {
		cUser = result.Model.(models.User)
	}

	// create user
	// MYTODO: Add forceUUID functionality
	loginGovUUID := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "first.last@login.gov.test",
		Active:        false,
	}

	// Overwrite values with those from assertions
	mergeModels(&user, cUser)

	// MYTODO: Add back stub functionality
	mustCreate(db, &user, false)

	return user, nil
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

// MakeServiceMember creates a single ServiceMember
// If not provided, it will also create an associated
// - User
// - ResidentialAddress
func ServiceMemberMaker(db *pop.Connection, customs []Customization, traits []GetTraitFunc) (models.ServiceMember, error) {
	customs = mergeCustomization(customs, traits)
	customs, err := validateCustomizations(customs)
	if err != nil {
		log.Panic(err)
	}

	// Find/create the required user model
	var user models.User
	if result := findValidCustomization(customs, User); result != nil {
		user = result.Model.(models.User)
	}
	if isZeroUUID(user.ID) {
		user, _ = UserMaker(db, customs, nil)
	}
	// At this point, user exists. It's either the provided or created user

	// Find/create a residential address
	var resiAddress models.Address
	result := findValidCustomization(customs, Addresses.ResidentialAddress)
	if result == nil {
		// No customization
		resiAddress, _ = AddressMaker(db, customs, nil)
	} else {
		// Customization exists
		resiAddress = result.Model.(models.Address)
		if isZeroUUID(resiAddress.ID) {
			// Convert ResidentialAddress type to Address type before passing on to Address maker
			tempCustoms := convertCustomizationInList(customs, Addresses.ResidentialAddress, Address)
			resiAddress, _ = AddressMaker(db, tempCustoms, nil)
		}
	}
	// At this point, resiAddress exists. It's either the provided or created residential address

	// Find the customization for service member
	var cServiceMember models.ServiceMember
	if result := findValidCustomization(customs, ServiceMember); result != nil {
		cServiceMember = result.Model.(models.ServiceMember)
	}

	// MYTODO We can add randomization and control with a flag
	randomEdipi := RandomEdipi()
	rank := models.ServiceMemberRankE1
	army := models.AffiliationARMY
	email := "leospaceman@gmail.com"

	// Set default values, include any nested IDs
	serviceMember := models.ServiceMember{
		User:                 user,
		UserID:               user.ID,
		Edipi:                swag.String(randomEdipi),
		Affiliation:          &army,
		FirstName:            swag.String("Leo"),
		LastName:             swag.String("Spacemen"),
		Telephone:            swag.String("212-123-4567"),
		ResidentialAddressID: &resiAddress.ID,
		ResidentialAddress:   &resiAddress,
		PersonalEmail:        &email,
		Rank:                 &rank,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceMember, cServiceMember)

	mustCreate(db, &serviceMember, false)

	return serviceMember, nil
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

func AddressMaker(db *pop.Connection, customs []Customization, traits []GetTraitFunc) (models.Address, error) {
	customs = mergeCustomization(customs, traits)
	customs, err := validateCustomizations(customs)
	if err != nil {
		log.Panic(err)
	}

	// Find the customization for address
	var cAddress models.Address
	if result := findValidCustomization(customs, Address); result != nil {
		cAddress = result.Model.(models.Address)
	}

	address := models.Address{
		StreetAddress1: "123 Any Street",
		StreetAddress2: swag.String("P.O. Box 12345"),
		StreetAddress3: swag.String("c/o Some Person"),
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "90210",
		Country:        swag.String("US"),
	}

	mergeModels(&address, cAddress)

	mustCreate(db, &address, false)

	return address, nil
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
