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

	suite.Run("Successful creation of entitlement with custom grade", func() {
		// Under test:      BuildEntitlement
		// Mocked:          None
		// Set up:          Create Entitlement with customization of orders for grade O_9
		// Expected outcome:Entitlement should contain weight allotments appropriate
		// .                to the grade included in the orders

		// SETUP
		// Create a default stubbed entitlement to compare values
		testEnt := BuildEntitlement(nil, nil, nil)
		// Set the weight allotment on the custom object to O_9
		testEnt.DBAuthorizedWeight = nil // clear original value
		testEnt.SetWeightAllotment("O_9")
		testEnt.DBAuthorizedWeight = testEnt.AuthorizedWeight()
		// Now DBAuthorizedWeight should be appropriate for O_9 grade

		// FUNCTION UNDER TEST
		entitlement := BuildEntitlement(suite.DB(), []Customization{
			{Model: models.Order{
				Grade: models.StringPointer("O_9"),
			}},
		}, nil)

		// VALIDATE RESULTS
		// Builder should have pulled the O_9 grade from the order to calculate
		// weight allotment
		suite.Equal(*testEnt.DBAuthorizedWeight, *entitlement.DBAuthorizedWeight)

	})

}
