package testdatagen

import (
	"fmt"
	"reflect"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// func (sf ServiceMemberMaker) Create(db *pop.Connection, custom Customization) error {
// 	sm := sf.Model

// 	// Get servicemember customization
// 	var customSM models.ServiceMember
// 	if custom.Name == "ServiceMember" {
// 		// If its a user, convert to model
// 		customSM = custom.Model.(models.ServiceMember)
// 	}

// 	// You need a user, check if one was provided
// 	var customUser models.User
// 	if custom.Name == "User" {
// 		customUser = custom.Model.(models.User)
// 	}
// 	if custom.Name == "User" && custom.Create == false {
// 		userMaker := NewUserMaker(models.User{}, nil)
// 		userMaker.Make(db, variants)
// 		user = *userMaker.Model
// 	}

// 	army := models.AffiliationARMY
// 	randomEdipi := RandomEdipi()
// 	rank := models.ServiceMemberRankE1
// 	email := "leo_spaceman_sm@example.com"

// 	sm.UserID = user.ID
// 	sm.User = user
// 	sm.Edipi = swag.String(randomEdipi)
// 	sm.Affiliation = &army
// 	sm.FirstName = swag.String("Leo")
// 	sm.LastName = swag.String("Spacemen")
// 	sm.Telephone = swag.String("212-123-4567")
// 	sm.PersonalEmail = &email
// 	sm.Rank = &rank

// 	// Overwrite values with those from assertions
// 	mergeModels(sm, variants.ServiceMember)

// 	mustCreate(db, sm, variants.Stub)

// 	return nil
// }

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
func mergeCustomization(traits []Trait, customs []Customization) []Customization {
	// Get a list of traits, each could return a list of customizations
	fmt.Println("We have", len(traits), "traits")
	for _, trait := range traits {
		traitCustomizations := trait()
		fmt.Println("this trait has", len(traitCustomizations))
		// for each customization, merge of replace the one in user supplied customizations
		for _, traitCustom := range traitCustomizations {
			fmt.Println("Found trait", traitCustom.Type)
			j, callerCustom := findCustomWithIdx(customs, traitCustom.Type)
			if callerCustom != nil {
				result := mergeInterfaces(callerCustom.Model, traitCustom.Model)
				callerCustom.Model = result
				customs[j] = *callerCustom
			} else {
				fmt.Println("No custom", traitCustom.Type)
				customs = append(customs, traitCustom)
			}
		}
	}
	return customs
}

func userMaker(db *pop.Connection, customs []Customization, traits []Trait) (models.User, error) {

	// Combine all traits into the customization list,
	// then clear traits so this is not repeated downstream
	if len(traits) != 0 {
		// This function take the original set of customizations
		// and merges with the traits.
		// The order of application is that the customizations override the trai
		customs = mergeCustomization(traits, customs)
		// traits = nil
	}

	// Find user assertion and convert to models user
	fmt.Print(traits)
	customUser := findCustom(customs, CustomUser).Model.(models.User)

	// create user
	loginGovUUID := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "first.last@login.gov.test",
		Active:        false,
	}

	// Overwrite values with those from assertions
	mergeModels(&user, customUser)

	mustCreate(db, &user, false)

	return user, nil
}

type CustomType string

const (
	CustomUser          CustomType = "User"
	CustomServiceMember CustomType = "ServiceMember"
)

type AddressesCustomType struct {
	PickupAddress            CustomType
	DeliveryAddress          CustomType
	SecondaryDeliveryAddress CustomType
}

var Addresses = AddressesCustomType{
	PickupAddress:            "PickupAddress",
	DeliveryAddress:          "DeliveryAddress",
	SecondaryDeliveryAddress: "SecondaryDeliveryAddress",
}

type Trait func() []Customization

func getTraitActiveUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
			},
			Type: CustomUser,
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
			Type: CustomUser,
		},
		{
			Model: models.ServiceMember{
				Affiliation: &army,
			},
			Type: CustomServiceMember,
		},
	}
	return VariantUserArmy
}

type Customization struct {
	Model       interface{}
	Type        CustomType
	Create      bool
	ReflectType reflect.Type
}
type Maker interface {
	Make(db *pop.Connection, customs []Customization, traits []Trait) error
}
