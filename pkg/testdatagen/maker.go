package testdatagen

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

type CustomType string

type Customization struct {
	Model     interface{}
	Type      CustomType
	ForceUUID bool
}
type AddressesCustomType struct {
	PickupAddress            CustomType
	DeliveryAddress          CustomType
	SecondaryDeliveryAddress CustomType
	ResidentialAddress       CustomType
}

const (
	Address       CustomType = "Address"
	User          CustomType = "User"
	ServiceMember CustomType = "ServiceMember"
)

type Trait func() []Customization

var Addresses = AddressesCustomType{
	PickupAddress:            "PickupAddress",
	DeliveryAddress:          "DeliveryAddress",
	SecondaryDeliveryAddress: "SecondaryDeliveryAddress",
	ResidentialAddress:       "ResidentialAddress",
}

func getTraitActiveUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
			},
			Type: User,
		},
	}
}

func getTraitArmy() []Customization {
	army := models.AffiliationARMY
	var VariantUserArmy = []Customization{
		{
			Model: models.User{
				LoginGovEmail: "trait@army.mil",
			},
			Type: User,
		},
		{
			Model: models.ServiceMember{
				Affiliation: &army,
			},
			Type: ServiceMember,
		},
	}
	return VariantUserArmy
}
func getTraitNavy() []Customization {
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

func findCustomWithIdx(customs []Customization, customType CustomType) (int, *Customization) {
	for i, custom := range customs {
		if custom.Type == customType {
			return i, &custom
		}
	}
	return -1, nil
}
func findCustom(customs []Customization, customType CustomType) *Customization {
	_, custom := findCustomWithIdx(customs, customType)
	return custom
}

// This function takes an interface wrapping a struct and
// returns an interface wrapping a pointer to the struct
// For e.g. interface{}(models.User) → interface{}(*models.User)
func toStructPtr(obj interface{}) interface{} {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp.Interface()
}

// This function takes an interface wrapping a pointer to a struct and
// returns an interface wrapping the struct
// For e.g. interface{}(*models.User) → interface{}(models.User)
func toInterfacePtr(obj interface{}) interface{} {
	rv := reflect.ValueOf(obj).Elem()
	return rv.Interface()
}

// This function transforms the interfaces to match what mergeModels expects
func mergeInterfaces(model1 interface{}, model2 interface{}) interface{} {
	modelPtr := toStructPtr(model1)
	mergeModels(modelPtr, model2)
	model := toInterfacePtr(modelPtr)
	return model
}

// This function take the original set of customizations
// and merges with the traits.
// The order of application is
// - Earlier traits override later traits in the trait list
// - Customizations override the traits
// So if you have [trait1, trait2] customization
// and all three contain the same object:
// - trait 1 will override trait 2 (so start with the highest priority)
// - customization will override trait 2
// MYTODO if a customization has an id, it should not be merged, just hooked up
func mergeCustomization(traits []Trait, customs []Customization) []Customization {
	// Get a list of traits, each could return a list of customizations
	fmt.Println("Found ", len(traits), "traits")
	for i, trait := range traits {
		traitCustomizations := trait()
		fmt.Println(i, ": Trait with ", len(traitCustomizations), "customizations")
		// for each customization, merge of replace the one in user supplied customizations
		for _, traitCustom := range traitCustomizations {
			j, callerCustom := findCustomWithIdx(customs, traitCustom.Type)
			if callerCustom != nil {
				fmt.Println("   ", traitCustom.Type, ": Found matching customization")
				result := mergeInterfaces(callerCustom.Model, traitCustom.Model)
				callerCustom.Model = result
				customs[j] = *callerCustom
			} else {
				fmt.Println("   ", traitCustom.Type, ": No matching customization")
				customs = append(customs, traitCustom)
			}
		}
	}
	return customs
}

func UserMaker(db *pop.Connection, customs []Customization, traits []Trait) (models.User, error) {

	// Combine all traits into the customization list,
	// do not pass on the traits in downstream function
	// so this merge is not repeated downstream
	// MYTODO validate the customizations for nested objects
	if len(traits) != 0 {
		// The order of application is that the customizations override the traits
		customs = mergeCustomization(traits, customs)
	}

	// Find user assertion and convert to models user
	var cUser models.User
	if result := findCustom(customs, User); result != nil {
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

	// MYTODO: Add stub functionality
	mustCreate(db, &user, false)

	return user, nil
}

// MakeServiceMember creates a single ServiceMember
// If not provided, it will also create an associated
// - User
// - ResidentialAddress
func ServiceMemberMaker(db *pop.Connection, customs []Customization, traits []Trait) (models.ServiceMember, error) {
	// Apply traits
	if len(traits) != 0 {
		customs = mergeCustomization(traits, customs)
	}

	// Find the customization for service member
	var cServiceMember models.ServiceMember
	if result := findCustom(customs, ServiceMember); result != nil {
		cServiceMember = result.Model.(models.ServiceMember)
	}

	// Find the customization for user
	var user models.User
	if result := findCustom(customs, User); result != nil {
		user = result.Model.(models.User)
	}
	if isZeroUUID(user.ID) {
		user, _ = UserMaker(db, customs, nil)
	}
	// User is now either provided or created user

	// Find the customization for residential address
	var resiAddress models.Address
	result := findCustom(customs, Addresses.ResidentialAddress)
	if result == nil {
		// No customization
		resiAddress, _ = AddressMaker(db, nil, nil)
	} else if isZeroUUID(resiAddress.ID) {
		// Customization exists but had no ID
		result.Type = Address
		resiAddress, _ = AddressMaker(db,
			[]Customization{*result}, nil)
	} else {
		// Customization exists and had an ID
		// This means we just need to use this object as-is
		resiAddress = result.Model.(models.Address)
	}
	// resiAddress is now either provided or created residential address

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

func AddressMaker(db *pop.Connection, customs []Customization, traits []Trait) (models.Address, error) {
	// Apply traits
	if len(traits) != 0 {
		customs = mergeCustomization(traits, customs)
	}

	// Find the customization for service member
	var cAddress models.Address
	if result := findCustom(customs, Address); result != nil {
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
