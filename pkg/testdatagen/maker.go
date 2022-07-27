package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// this is the class - note type struct and member vairbal
type ServiceMemberMaker struct {
	Model     *models.ServiceMember
	ForceUUID *uuid.UUID
}

// this is the constructor
func NewServiceMemberMaker(serviceMember models.ServiceMember, forceUUID *uuid.UUID) ServiceMemberMaker {
	return ServiceMemberMaker{&serviceMember, forceUUID}
}

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

// this is the class - note type struct and member variable
type UserMaker struct {
	Model     *models.User
	ForceUUID *uuid.UUID
}

// this is the constructor
func NewUserMaker(user models.User, forceUUID *uuid.UUID) UserMaker {
	return UserMaker{&user, forceUUID}
}

func findCustom(customs []Customization, customType CustomType) *Customization {
	for _, custom := range customs {
		if custom.Type == customType {
			return &custom
		}
	}
	return nil
}

// this is a method

func (uf UserMaker) Make(db *pop.Connection, customs []Customization, traits []Trait) error {
	user := uf.Model

	// Find user assertion and convert to models user
	customUser := findCustom(customs, CustomUser).Model.(models.User)

	loginGovUUID := uuid.Must(uuid.NewV4())
	user.LoginGovUUID = &loginGovUUID
	user.LoginGovEmail = "first.last@login.gov.test"
	user.Active = false

	// Overwrite values with those from assertions
	mergeModels(user, customUser)

	mustCreate(db, user, false)

	return nil
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

func getTraitArmy() []Customization {
	var army = models.AffiliationARMY
	var VariantUserArmy = []Customization{
		{
			Model: models.User{
				LoginGovEmail: "testing@army.mil",
			},
			Type: CustomUser,
		},
		{
			Model: models.ServiceMember{
				Affiliation: &army,
			},
			Type: CustomUser,
		},
	}
	return VariantUserArmy
}

type Customization struct {
	Model  interface{}
	Type   CustomType
	Create bool
}
type Maker interface {
	Make(db *pop.Connection, customs []Customization, traits []Trait) error
}
