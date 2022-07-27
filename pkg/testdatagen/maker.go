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

// this is a method

func (uf UserMaker) Make(db *pop.Connection, custom Customization) error {
	user := uf.Model

	// Find assertions
	var customUser models.User
	if custom.Name == "User" {
		// If its a user, convert to model
		customUser = custom.Model.(models.User)
	}

	loginGovUUID := uuid.Must(uuid.NewV4())
	user.LoginGovUUID = &loginGovUUID
	user.LoginGovEmail = "first.last@login.gov.test"
	user.Active = false

	// Overwrite values with those from assertions
	mergeModels(user, customUser)

	mustCreate(db, user, false)

	return nil
}

type Customization struct {
	Model  interface{}
	Name   string
	Create bool
}
type Maker interface {
	Make(db *pop.Connection, custom Customization) error
}
