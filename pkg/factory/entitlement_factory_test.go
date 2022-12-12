package factory

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *FactorySuite) TestBuildEntitlement() {
	suite.Run("Successful creation of default entitlement", func() {
		// Under test:      BuildEntitlement
		// Mocked:          None
		// Set up:          Create a User with no customizations or traits
		// Expected outcome:User should be created with default values

		entitlement := BuildEntitlement(suite.DB(), nil, nil)
		suite.True(*entitlement.DependentsAuthorized)
		fmt.Printf("1 >>  %+v\n", entitlement)

		en2 := testdatagen.MakeEntitlement(suite.DB(), testdatagen.Assertions{})
		suite.True(*en2.DependentsAuthorized)
		fmt.Printf("2 >> %+v\n", en2)
	})

	// 	suite.Run("Successful creation of entitlement with matched email", func() {
	// 		// Under test:      BuildEntitlement
	// 		// Mocked:          None
	// 		// Set up:          Create a User but pass in a trait that sets
	// 		//                  both the entitlement and user email to a random
	// 		//                  value, as entitlement has uniqueness constraints
	// 		// Expected outcome:Entitlement should have the same random email as User

	// 		entitlement := BuildEntitlement(suite.DB(), nil, []Trait{
	// 			GetTraitEntitlementEmail,
	// 		})
	// 		suite.Equal(entitlement.Email, entitlement.User.LoginGovEmail)
	// 		suite.False(entitlement.User.Active)
	// 	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:      BuildEntitlement
		// Set up:          Create an entitlement and pass in specified emails
		// Expected outcome:entitlement and User should be created with specified emails
		entitlement := BuildEntitlement(suite.DB(), []Customization{
			{
				Model: models.Entitlement{
					DependentsAuthorized:  models.BoolPointer(true),
					PrivatelyOwnedVehicle: models.BoolPointer(false),
				},
			},
		}, nil)
		fmt.Printf("%+v", entitlement)
		suite.True(*entitlement.DependentsAuthorized)
		suite.False(*entitlement.PrivatelyOwnedVehicle)
		suite.True(*entitlement.NonTemporaryStorage)

		en2 := testdatagen.MakeEntitlement(suite.DB(), testdatagen.Assertions{
			Entitlement: models.Entitlement{
				DependentsAuthorized: models.BoolPointer(false),
			},
		})
		suite.False(*en2.DependentsAuthorized)
	})
	// }
	// func (suite *FactorySuite) TestBuildEntitlementExtra() {
	// 	// Under test:      BuildEntitlement
	// 	// Mocked:          None
	// 	// Set up:          Create a Entitlement but pass in a role
	// 	// Expected outcome:Created User should have the associated Role

	// 	suite.Run("Successful creation of TIO Admin User", func() {

	// 		// Create the TIO Role
	// 		tioRole := roles.Role{
	// 			ID:       uuid.Must(uuid.NewV4()),
	// 			RoleType: roles.RoleTypeTIO,
	// 			RoleName: "Transportation Invoicing Officer",
	// 		}
	// 		verrs, err := suite.DB().ValidateAndCreate(&tioRole)
	// 		suite.NoError(err)
	// 		suite.False(verrs.HasAny())

	// 		// FUNCTION UNDER TEST
	// 		entitlement := BuildEntitlement(suite.DB(), []Customization{
	// 			{
	// 				Model: models.User{
	// 					Roles: []roles.Role{tioRole},
	// 				},
	// 				Type: &User,
	// 			},
	// 		}, []Trait{
	// 			GetTraitEntitlementEmail,
	// 		})

	// 		// VALIDATE RESULT
	// 		// Check that the email trait worked
	// 		suite.Equal(entitlement.Email, entitlement.User.LoginGovEmail)
	// 		suite.False(entitlement.User.Active)
	// 		// Check that the user has the admin user role
	// 		_, hasRole := entitlement.User.Roles.GetRole(roles.RoleTypeTIO)
	// 		suite.True(hasRole)
	// 	})

	// 	suite.Run("Successful creation of Entitlement with linked User", func() {
	// 		// Under test:       BuildEntitlement
	// 		// Set up:           Create an entitlement and pass in a precreated user
	// 		// Expected outcome: entitlement should link in the precreated user

	// 		user := BuildUser(suite.DB(), []Customization{
	// 			{
	// 				Model: models.User{
	// 					CurrentAdminSessionID: "breathe",
	// 				},
	// 			},
	// 		}, nil)
	// 		entitlement := BuildEntitlement(suite.DB(), []Customization{
	// 			{
	// 				Model:    user,
	// 				LinkOnly: true,
	// 			},
	// 		}, []Trait{
	// 			GetTraitEntitlementEmail,
	// 		})
	// 		// Check that the linked user was used
	// 		suite.Equal(user.ID, *entitlement.UserID)
	// 		suite.Equal(user.ID, entitlement.User.ID)
	// 		suite.Equal("breathe", entitlement.User.CurrentAdminSessionID)
	// 		suite.False(entitlement.Active)

	// 	})
	// 	suite.Run("Successful creation of Entitlement with forced id User", func() {
	// 		// Under test:       BuildEntitlement
	// 		// Set up:           Create an entitlement and pass in an ID for User
	// 		// Expected outcome: entitlement and User should be created
	// 		//                   User should have specified ID

	// 		defaultLoginGovEmail := "first.last@login.gov.test"
	// 		uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
	// 		entitlement := BuildEntitlement(suite.DB(), []Customization{
	// 			{
	// 				Model: models.User{
	// 					ID: uuid,
	// 				},
	// 			},
	// 		}, []Trait{
	// 			GetTraitEntitlementEmail,
	// 		})
	// 		// Check that the forced ID was used
	// 		suite.Equal(uuid, *entitlement.UserID)
	// 		suite.Equal(uuid, entitlement.User.ID)

	// 		// Check that id can be found in DB
	// 		foundUser := models.User{}
	// 		err := suite.DB().Find(&foundUser, uuid)
	// 		suite.NoError(err)

	// 		// Check that email was applied to user
	// 		suite.NotContains(defaultLoginGovEmail, entitlement.User.LoginGovEmail)
	// 		suite.Equal(entitlement.Email, entitlement.User.LoginGovEmail)
	// 	})

	// 	suite.Run("Successful creation of stubbed Entitlement with forced id User", func() {
	// 		// Under test:       BuildEntitlement
	// 		// Set up:           Create an entitlement and pass in a precreated user
	// 		// Expected outcome: entitlement and User should be created with specified emails
	// 		uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
	// 		entitlement := BuildEntitlement(nil, []Customization{
	// 			{
	// 				Model: models.User{
	// 					ID: uuid,
	// 				},
	// 			},
	// 		}, []Trait{
	// 			GetTraitEntitlementEmail,
	// 		})
	// 		// Check that the forced ID was used
	// 		suite.Equal(uuid, *entitlement.UserID)
	// 		suite.Equal(uuid, entitlement.User.ID)

	// 		// Check that id cannot be found in DB
	// 		foundUser := models.User{}
	// 		err := suite.DB().Find(&foundUser, uuid)
	// 		suite.Error(err)

	//			// Check that email was applied to user
	//			suite.Equal(entitlement.Email, entitlement.User.LoginGovEmail)
	//		})
	//	}
}
