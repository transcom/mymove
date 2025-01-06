package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/entitlements"
)

func (suite *FactorySuite) TestBuildEntitlement() {
	fetcher := entitlements.NewWeightAllotmentFetcher()

	setupE1Allotment := func() {
		pg := BuildPayGrade(suite.DB(), []Customization{
			{
				Model: models.PayGrade{
					Grade: "E_1",
				},
			},
		}, nil)
		BuildHHGAllowance(suite.DB(), []Customization{
			{
				Model:    pg,
				LinkOnly: true,
			},
		}, nil)
	}
	setupO9Allotment := func() {
		pg := BuildPayGrade(suite.DB(), []Customization{
			{
				Model: models.PayGrade{
					Grade: "O_9",
				},
			},
		}, nil)
		BuildHHGAllowance(suite.DB(), []Customization{
			{
				Model:    pg,
				LinkOnly: true,
			},
		}, nil)
	}
	suite.Run("Successful creation of default entitlement", func() {
		// Under test:      BuildEntitlement
		// Mocked:          None
		// Set up:          Create an entitlement with no customizations or traits
		// Expected outcome:User should be created with default values

		// SETUP
		// Create a default entitlement to compare values
		setupE1Allotment()
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
		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), "E_1", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
		suite.NoError(err)
		defEnt.WeightAllotted = &allotment

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
		setupE1Allotment()
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
		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), "E_1", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
		suite.NoError(err)
		custEnt.WeightAllotted = &allotment
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
		setupO9Allotment()
		testEnt := BuildEntitlement(nil, nil, nil)
		// Set the weight allotment on the custom object to O_9
		testEnt.DBAuthorizedWeight = nil // clear original value
		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), "O_9", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
		suite.NoError(err)
		testEnt.WeightAllotted = &allotment
		testEnt.DBAuthorizedWeight = testEnt.AuthorizedWeight()

		// FUNCTION UNDER TEST
		grade := internalmessages.OrderPayGrade(models.ServiceMemberGradeO9)
		entitlement := BuildEntitlement(suite.DB(), []Customization{
			{Model: models.Order{
				Grade: &grade,
			}},
		}, nil)

		// VALIDATE RESULTS
		// Builder should have pulled the O_9 grade from the order to calculate
		// weight allotment
		suite.Equal(*testEnt.DBAuthorizedWeight, *entitlement.DBAuthorizedWeight)

	})

}

func (suite *FactorySuite) TestBuildPayGrade() {
	suite.Run("Successful creation of PayGrade with default values", func() {
		// Default grade should be "E_5"
		payGrade := BuildPayGrade(suite.DB(), nil, nil)

		suite.NotNil(payGrade.ID)
		suite.Equal("E_5", payGrade.Grade)
		suite.Equal("Enlisted Grade E-5", *payGrade.GradeDescription)

		pgCount, err := suite.DB().Count(models.PayGrade{})
		suite.NoError(err)
		suite.True(pgCount > 0)
	})

	suite.Run("BuildPayGrade with customization", func() {
		customGrade := "X-5"
		customDescription := "Custom Grade X-5"
		customPayGrade := models.PayGrade{
			Grade:            customGrade,
			GradeDescription: &customDescription,
		}

		payGrade := BuildPayGrade(
			suite.DB(),
			[]Customization{
				{Model: customPayGrade},
			},
			nil,
		)

		suite.Equal(customGrade, payGrade.Grade)
		suite.Equal(customDescription, *payGrade.GradeDescription)
	})

	suite.Run("Finds existing record", func() {

		persistedPayGrade := BuildPayGrade(suite.DB(), nil, nil)

		pg := BuildPayGrade(suite.DB(), []Customization{
			{
				Model:    persistedPayGrade,
				LinkOnly: true,
			},
		}, nil)

		suite.Equal(persistedPayGrade.ID, pg.ID)
		suite.Equal(persistedPayGrade.Grade, pg.Grade)

	})
}

func (suite *FactorySuite) TestBuildHHGAllowance() {
	suite.Run("Successful creation of HHGAllowance with default values", func() {
		// Default allowance and grade of E_5
		hhgAllowance := BuildHHGAllowance(suite.DB(), nil, nil)
		suite.NotNil(hhgAllowance.PayGradeID)
		suite.NotEmpty(hhgAllowance.PayGrade)
		suite.NotEmpty(hhgAllowance.ProGearWeight)
		suite.NotEmpty(hhgAllowance.ProGearWeightSpouse)
		suite.NotEmpty(hhgAllowance.TotalWeightSelf)
		suite.NotEmpty(hhgAllowance.TotalWeightSelfPlusDependents)
	})

	suite.Run("BuildHHGAllowance with customization", func() {
		hhgAllowance := BuildHHGAllowance(
			suite.DB(),
			[]Customization{
				{Model: models.HHGAllowance{
					TotalWeightSelf:               8000,
					TotalWeightSelfPlusDependents: 12000,
					ProGearWeight:                 3000,
					ProGearWeightSpouse:           600,
				}},
			},
			nil,
		)

		// E_5 default allowances
		suite.Equal(8000, hhgAllowance.TotalWeightSelf)
		suite.Equal(12000, hhgAllowance.TotalWeightSelfPlusDependents)
		suite.Equal(3000, hhgAllowance.ProGearWeight)
		suite.Equal(600, hhgAllowance.ProGearWeightSpouse)
	})

	suite.Run("Finds existing record", func() {
		pg := BuildPayGrade(suite.DB(), nil, nil)

		existingHhg := models.HHGAllowance{
			PayGradeID:                    pg.ID,
			TotalWeightSelf:               8000,
			TotalWeightSelfPlusDependents: 12000,
			ProGearWeight:                 3000,
			ProGearWeightSpouse:           600,
		}
		suite.MustCreate(&existingHhg)

		newHhg := BuildHHGAllowance(
			suite.DB(),
			[]Customization{
				{Model: models.HHGAllowance{PayGradeID: pg.ID}},
			},
			nil,
		)

		suite.Equal(existingHhg.ID, newHhg.ID)
		suite.Equal(3000, newHhg.ProGearWeight)
	})
}

func (suite *FactorySuite) TestSetupAllAllotments() {
	suite.Run("Successful creation of allotments for all known grades", func() {
		err := suite.DB().TruncateAll()
		suite.NoError(err)

		SetupDefaultAllotments(suite.DB())

		// Validate the allotments
		for grade, allowance := range knownAllowances {
			// Ensure pay grade exists
			pg := &models.PayGrade{}
			err := suite.DB().Where("grade = ?", grade).First(pg)
			suite.NoError(err, grade)
			suite.NotNil(pg.ID, grade)

			// Ensure HHGAllowance was created and matches the expected values
			hhgAllowance := &models.HHGAllowance{}
			err = suite.DB().Where("pay_grade_id = ?", pg.ID).First(hhgAllowance)
			suite.NoError(err, grade)
			suite.Equal(allowance.TotalWeightSelf, hhgAllowance.TotalWeightSelf, grade)
			suite.Equal(allowance.TotalWeightSelfPlusDependents, hhgAllowance.TotalWeightSelfPlusDependents, grade)
			suite.Equal(allowance.ProGearWeight, hhgAllowance.ProGearWeight, grade)
			suite.Equal(allowance.ProGearWeightSpouse, hhgAllowance.ProGearWeightSpouse, grade)
		}
	})
}
