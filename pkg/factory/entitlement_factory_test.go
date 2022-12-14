package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildEntitlement() {
	suite.Run("Successful creation of default entitlement", func() {
		// Under test:      BuildEntitlement
		// Mocked:          None
		// Set up:          Create an entitlement with no customizations or traits
		// Expected outcome:User should be created with default values

		// SETUP
		// Create a default entitlement to compare values
		defEnt := models.Entitlement{
			DependentsAuthorized:                         models.BoolPointer(true),
			TotalDependents:                              models.IntPointer(0),
			NonTemporaryStorage:                          models.BoolPointer(true),
			PrivatelyOwnedVehicle:                        models.BoolPointer(true),
			StorageInTransit:                             models.IntPointer(90),
			ProGearWeight:                                2000,
			ProGearWeightSpouse:                          500,
			RequiredMedicalEquipmentWeight:               1000,
			OrganizationalClothingAndIndividualEquipment: true,
		}
		defEnt.SetWeightAllotment("E_1")
		defEnt.DBAuthorizedWeight = defEnt.AuthorizedWeight()

		// FUNCTION UNDER TEST
		entitlement := BuildEntitlement(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(*defEnt.DependentsAuthorized, *entitlement.DependentsAuthorized)
		suite.Equal(*defEnt.TotalDependents, *entitlement.TotalDependents)
		suite.Equal(*defEnt.NonTemporaryStorage, *entitlement.NonTemporaryStorage)
		suite.Equal(*defEnt.PrivatelyOwnedVehicle, *entitlement.PrivatelyOwnedVehicle)
		suite.Equal(*defEnt.StorageInTransit, *entitlement.StorageInTransit)
		suite.Equal(defEnt.DBAuthorizedWeight, entitlement.DBAuthorizedWeight)
		suite.Equal(defEnt.ProGearWeight, entitlement.ProGearWeight)
		suite.Equal(defEnt.ProGearWeightSpouse, entitlement.ProGearWeightSpouse)
		suite.Equal(defEnt.RequiredMedicalEquipmentWeight, entitlement.RequiredMedicalEquipmentWeight)
		suite.Equal(defEnt.OrganizationalClothingAndIndividualEquipment, entitlement.OrganizationalClothingAndIndividualEquipment)

	})

	suite.Run("Successful creation of customized entitlement", func() {
		// Under test:      BuildEntitlement
		// Mocked:          None
		// Set up:          Create Entitlement with customization
		// Expected outcome:Entitlement should customized values

		// SETUP
		// Create a default entitlement to compare values
		custEnt := models.Entitlement{
			DependentsAuthorized:                         models.BoolPointer(false),
			TotalDependents:                              models.IntPointer(0),
			NonTemporaryStorage:                          models.BoolPointer(true),
			PrivatelyOwnedVehicle:                        models.BoolPointer(true),
			StorageInTransit:                             models.IntPointer(90),
			ProGearWeight:                                2000,
			ProGearWeightSpouse:                          500,
			RequiredMedicalEquipmentWeight:               1000,
			OrganizationalClothingAndIndividualEquipment: true,
		}

		// FUNCTION UNDER TEST
		entitlement := BuildEntitlement(suite.DB(), []Customization{
			{Model: custEnt},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(*custEnt.DependentsAuthorized, *entitlement.DependentsAuthorized)
		suite.Equal(*custEnt.TotalDependents, *entitlement.TotalDependents)
		suite.Equal(*custEnt.NonTemporaryStorage, *entitlement.NonTemporaryStorage)
		suite.Equal(*custEnt.PrivatelyOwnedVehicle, *entitlement.PrivatelyOwnedVehicle)
		suite.Equal(*custEnt.StorageInTransit, *entitlement.StorageInTransit)
		suite.Equal(custEnt.ProGearWeight, entitlement.ProGearWeight)
		suite.Equal(custEnt.ProGearWeightSpouse, entitlement.ProGearWeightSpouse)
		suite.Equal(custEnt.RequiredMedicalEquipmentWeight, entitlement.RequiredMedicalEquipmentWeight)
		suite.Equal(custEnt.OrganizationalClothingAndIndividualEquipment, entitlement.OrganizationalClothingAndIndividualEquipment)

		// Set the weight allotment on the custom object so as to compare
		custEnt.SetWeightAllotment("E_1")
		custEnt.DBAuthorizedWeight = custEnt.AuthorizedWeight()

		// Check that the created object had the correct allotments set
		suite.Equal(custEnt.DBAuthorizedWeight, entitlement.DBAuthorizedWeight)

	})

	suite.Run("Successful return of linkOnly entitlement", func() {
		// Under test:       BuildEntitlement
		// Set up:           Create an entitlement and pass in a linkOnly entitlement
		// Expected outcome: No new entitlement should be created.

		// Check num entitlements
		precount, err := suite.DB().Count(&models.Entitlement{})
		suite.NoError(err)

		entitlement := BuildEntitlement(suite.DB(), []Customization{
			{
				Model: models.Entitlement{
					ID:            uuid.Must(uuid.NewV4()),
					ProGearWeight: 2765,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.Entitlement{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(2765, entitlement.ProGearWeight)

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
