package factory

import (
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type FactorySuite struct {
	*testingsuite.PopTestSuite
}

func TestFactorySuite(t *testing.T) {

	ts := &FactorySuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

// getTraitActiveArmy is a custom Trait, used for testing
func getTraitActiveArmy() []Customization {
	army := models.AffiliationARMY
	return []Customization{
		{
			Model: models.User{
				LoginGovEmail:         "trait@army.mil",
				CurrentAdminSessionID: "my-session-id",
			},
			Type: &User,
		},
		{
			Model: models.ServiceMember{
				Affiliation: &army,
			},
			Type: &ServiceMember,
		},
	}
}

func (suite *FactorySuite) TestMergeCustomization() {

	suite.Run("Customizations and traits merged into result", func() {
		// Under test:       mergeCustomization, which merges traits and customizations
		// Set up:           Create a customization without a matching trait, and traits with no customizations
		// Expected outcome: All should exist in the result list
		streetAddress := "235 Prospect Valley Road SE"

		// CALL FUNCTION UNDER TEST
		result := mergeCustomization(
			// Address customization
			[]Customization{
				{
					Model: models.Address{
						StreetAddress1: streetAddress,
					},
					Type: &Address, // ← Address customization
				},
			},
			// User and ServiceMember customization
			[]Trait{
				getTraitActiveArmy, // ← User and ServiceMember customization
			},
		)
		suite.Len(result, 3)

		// VALIDATE RESULTS
		// Check that result included our Address customization
		_, custom := findCustomWithIdx(result, Address)
		suite.NotNil(custom)
		address := custom.Model.(models.Address)
		suite.Equal(streetAddress, address.StreetAddress1)

		// Check that result included our User customization
		_, custom = findCustomWithIdx(result, User)
		suite.NotNil(custom)
		user := custom.Model.(models.User)
		suite.Equal("trait@army.mil", user.LoginGovEmail)

		// Check that result included our ServiceMember customization
		_, custom = findCustomWithIdx(result, ServiceMember)
		sm := custom.Model.(models.ServiceMember)
		suite.Equal(models.AffiliationARMY, *sm.Affiliation)
		suite.NotNil(custom)
	})

	suite.Run("Customization should override traits", func() {
		// Under test:       mergeCustomization, which merges traits and customizations
		// Set up:           Create a customization with a user email and a trait with a user email
		// Expected outcome: Customization should override the trait email
		uuidval := uuid.Must(uuid.NewV4())
		// RUN FUNCTION UNDER TEST
		result := mergeCustomization(
			[]Customization{
				{
					Model: models.User{
						LoginGovUUID:  &uuidval,
						LoginGovEmail: "custom@army.mil",
					},
					Type: &User, // ← User customization
				},
			},
			[]Trait{
				getTraitActiveArmy, // ← Address and User customization
			},
		)

		// VALIDATE RESULTS
		userModel := result[0].Model.(models.User)
		// Customization email should be used
		suite.Equal("custom@army.mil", userModel.LoginGovEmail)
		// But other fields could come from trait
		suite.Equal("my-session-id", userModel.CurrentAdminSessionID)
	})

	suite.Run("Customization should override traits in priority order", func() {
		// Under test:       mergeCustomization, which merges traits and customizations
		// Set up:           Check the priority order, which is customization, then traits in order
		// Expected outcome: Customization should override getTrait1
		//                   getTrait1 should override getTrait1

		// customization is  custom   ______   ______
		// getTrait1 is      trait1   trait1   ______
		// getTrait1 is      trait2   trait2   trait2

		// Result should be  custom   trait1   trait2

		getTrait1 := func() []Customization {
			return []Customization{
				{
					Model: models.User{
						CurrentAdminSessionID:  "trait1",
						CurrentOfficeSessionID: "trait1",
						CurrentMilSessionID:    "",
					},
					Type: &User,
				},
			}
		}
		getTrait2 := func() []Customization {
			return []Customization{
				{
					Model: models.User{
						CurrentAdminSessionID:  "trait2",
						CurrentOfficeSessionID: "trait2",
						CurrentMilSessionID:    "trait2",
					},
					Type: &User,
				},
			}
		}
		// RUN FUNCTION UNDER TEST
		result := mergeCustomization(
			[]Customization{
				{
					Model: models.User{
						CurrentAdminSessionID: "custom",
					},
					Type: &User,
				},
			},
			[]Trait{
				getTrait1,
				getTrait2,
			},
		)

		// VALIDATE RESULTS
		user := result[0].Model.(models.User)
		suite.Equal("custom", user.CurrentAdminSessionID)
		suite.Equal("trait1", user.CurrentOfficeSessionID)
		suite.Equal("trait2", user.CurrentMilSessionID)
	})

}

func (suite *FactorySuite) TestMergeInterfaces() {
	// Under test:       mergeInterfaces, wrapper function for calling mergeModels
	// Set up:           Create two interface types and call mergeInterfaces
	// Expected outcome: Underlying model should contain fields from both models.
	//                   user1 fields should overwrite user2 fields
	user1 := models.User{
		LoginGovEmail: "user1@example.com",
		Active:        true,
	}
	uuidNew := uuid.Must(uuid.NewV4())
	user2 := models.User{
		LoginGovEmail: "user2@example.com",
		LoginGovUUID:  &uuidNew,
	}

	result := mergeInterfaces(user2, user1)
	user := result.(models.User)
	// user1 email should overwrite user2 email
	suite.Equal(user1.LoginGovEmail, user.LoginGovEmail)
	// All other fields set in interfaces should persist
	suite.Equal(user1.Active, user.Active)
	suite.Equal(user2.LoginGovUUID, user.LoginGovUUID)
}

func (suite *FactorySuite) TestHasID() {
	suite.Run("True if ID is provided", func() {
		// Under test:       hasIDs, uses reflection to check if a model has a populated ID
		// Set up:           Test a model with an id
		// Expected outcome: True
		testid := uuid.Must(uuid.NewV4())
		result := hasID(models.ServiceMember{
			ID: testid,
		})
		suite.True(result)
	})

	suite.Run("False if no id", func() {
		// Under test:       hasIDs
		// Set up:           Test a model with no id
		// Expected outcome: False
		result := hasID(models.ServiceMember{})
		suite.False(result)
	})

	suite.Run("False if nil id", func() {
		// Under test:       hasIDs
		// Set up:           Test a model with a nil id
		// Expected outcome: False
		result := hasID(models.ServiceMember{
			ID: uuid.Nil,
		})
		suite.False(result)
	})
}

func (suite *FactorySuite) TestNestedModelsCheck() {
	suite.Run("Must not call with pointer", func() {
		// Under test:       checkNestedModels, uses reflection to check if other models are nested
		// Set up:           Call with a pointer, instead of a struct
		// Expected outcome: Error

		c := Customization{
			Model: models.ServiceMember{},
			Type:  &ServiceMember,
		}
		err := checkNestedModels(&c)
		suite.Error(err)
		suite.Contains(err.Error(), "received a pointer")
	})

	suite.Run("Must not have missing model", func() {
		// Under test:       checkNestedModels
		// Set up:           Call with a struct that doesn't contain the main model
		// Expected outcome: Error
		c := Customization{
			Type: &ServiceMember,
		}
		err := checkNestedModels(c)
		suite.Error(err)
		suite.Contains(err.Error(), "must contain a model")
	})

	suite.Run("Customization contains nested model", func() {
		// Under test:       checkNestedModels
		// Set up:           Call with a struct that contains a nested model
		// Expected outcome: Error
		user := BuildUser(suite.DB(), nil, nil)
		c := Customization{
			Model: models.ServiceMember{
				User: user,
			},
			Type: &ServiceMember,
		}
		err := checkNestedModels(c)
		suite.Error(err)
		suite.Contains(err.Error(), "no nested models")

	})
	suite.Run("Customization contains ptr to nested model", func() {
		// Under test:       checkNestedModels
		// Set up:           Call with a struct with nested address
		// Expected outcome: Error
		resiAddress := models.Address{
			StreetAddress1: "142 E Barrel Hoop Circle #4A",
		}
		c := Customization{
			Model: models.ServiceMember{
				ResidentialAddress: &resiAddress,
			},
			Type: &ServiceMember,
		}
		err := checkNestedModels(c)
		suite.Error(err)
		suite.Contains(err.Error(), "no nested models")

	})
	suite.Run("Customization allows all other fields", func() {
		// Under test:       checkNestedModels
		// Set up:           Call with a struct with fields populated but no nested model
		// Expected outcome: No Error
		navy := models.AffiliationNAVY
		testid := uuid.Must(uuid.NewV4())
		edipi := RandomEdipi()
		timestamp := time.Now()
		rank := models.ServiceMemberRankE4
		name := "Riley Baker"
		phone := "555-777-9929"

		c := Customization{
			Model: models.ServiceMember{
				ID:                     uuid.Must(uuid.NewV4()),
				CreatedAt:              timestamp,
				UpdatedAt:              timestamp,
				UserID:                 testid,
				Edipi:                  &edipi,
				Affiliation:            &navy,
				Rank:                   &rank,
				FirstName:              &name,
				MiddleName:             &name,
				LastName:               &name,
				Suffix:                 &name,
				Telephone:              &phone,
				SecondaryTelephone:     &phone,
				PersonalEmail:          &name,
				PhoneIsPreferred:       swag.Bool(true),
				EmailIsPreferred:       swag.Bool(false),
				ResidentialAddressID:   &testid,
				BackupMailingAddressID: &testid,
				DutyLocationID:         &testid,
			},
			Type: &ServiceMember,
		}
		err := checkNestedModels(c)
		suite.NoError(err)

	})

}
func (suite *FactorySuite) TestDefaultTypes() {

	suite.Run("Default types added if missing", func() {
		// TESTCASE SCENARIO
		// Under test:       setDefaultTypes
		// Set up:           Pass customizations with known models
		// Expected outcome: Type is set on all models
		customs := []Customization{
			{
				Model: models.User{
					LoginGovEmail: "string",
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "string",
				},
			},
		}
		suite.Len(customs, 2)
		setDefaultTypes(customs)
		for _, c := range customs {
			suite.NotNil(c.Type)
		}
	})
	suite.Run("Error if type is unknown", func() {
		// TESTCASE SCENARIO
		// Under test:       assignType
		// Set up:           Create a customization with a type that isn't supported
		// Expected outcome: Error
		custom := Customization{
			Model: models.MoveHistory{
				Locator: "rock",
			},
		}
		err := assignType(&custom)
		suite.Error(err)
		suite.ErrorContains(err, "models.MoveHistory")
	})
}
func (suite *FactorySuite) TestSetupCustomizations() {

	suite.Run("Customizations and traits merged into result", func() {
		// Under test:       setupCustomizations which calls mergeCustomization,
		//                   which merges traits and customizations
		// Set up:           Create a customization without a matching trait, and traits with no customizations
		// Expected outcome: All should exist in the result list
		streetAddress := "235 Prospect Valley Road SE"
		// Call function under test
		result := setupCustomizations(
			// Address customization
			[]Customization{
				{
					Model: models.Address{
						StreetAddress1: streetAddress,
					},
				},
				{
					Model: controlObject{
						isValid: false,
					},
					Type: &control,
				},
			},
			// User and ServiceMember customization
			[]Trait{
				getTraitActiveArmy,
			},
		)

		// Expect to get 4 customizations, address, user, servicemember and control
		suite.Len(result, 4)

		// Find Address, check details
		_, custom := findCustomWithIdx(result, Address)
		suite.NotNil(custom)
		address := custom.Model.(models.Address)
		suite.Equal(streetAddress, address.StreetAddress1)

		// Find User, check details
		_, custom = findCustomWithIdx(result, User)
		suite.NotNil(custom)
		user := custom.Model.(models.User)
		suite.Equal("trait@army.mil", user.LoginGovEmail)

		// Find ServiceMember, check details
		_, custom = findCustomWithIdx(result, ServiceMember)
		suite.NotNil(custom)
	})

	suite.Run("Customization should override traits", func() {
		// Under test:       mergeCustomization, which merges traits and customizations
		// Set up:           Create a customization with a user email and a trait with a user email
		// Expected outcome: Customization should override the trait email
		//                   If an object exists and no customization, it should become a customization
		uuidval := uuid.Must(uuid.NewV4())
		result := setupCustomizations(
			[]Customization{
				{
					Model: models.User{
						LoginGovUUID:  &uuidval,
						LoginGovEmail: "custom@army.mil",
					},
				},
			},
			[]Trait{
				getTraitActiveArmy,
			},
		)
		userModel := result[0].Model.(models.User)
		// Customization email should be used
		suite.Equal("custom@army.mil", userModel.LoginGovEmail)
		// But other fields could come from trait
		suite.Equal("my-session-id", userModel.CurrentAdminSessionID)

	})

}
func (suite *FactorySuite) TestValidateCustomizations() {
	suite.Run("Control obj added if missing", func() {
		// Under test:       setupCustomizations
		// Set up:           Create some customizations without a control object
		// Expected outcome: A control object should be added to the list

		customs := getTraitActiveArmy()
		suite.Len(customs, 2)

		customs = setupCustomizations(customs, nil)
		suite.Len(customs, 3) // ← one was added
		_, controlCustom := findCustomWithIdx(customs, control)
		suite.NotNil(controlCustom)
	})

	suite.Run("Control obj not added if not missing", func() {
		// Under test:       setupCustomizations
		// Set up:           Create some customizations with a control object
		// Expected outcome: No objects should be added to the list
		customs := getTraitActiveArmy()
		customs = append(customs, Customization{
			Model: controlObject{},
			Type:  &control,
		})
		suite.Len(customs, 3)

		customs = setupCustomizations(customs, nil)
		suite.Len(customs, 3) // ← nothing added
		_, controlCustom := findCustomWithIdx(customs, control)
		suite.NotNil(controlCustom)
	})

	suite.Run("Pass if customizations not repeated", func() {
		// Under test:       uniqueCustomizations checks that there's only one
		//                   customization of each type
		// Set up:           Create some customizations of different types
		// Expected outcome: No error

		customs := getTraitActiveArmy()
		customs = append(customs, Customization{
			Model: models.Address{},
			Type:  &Addresses.ResidentialAddress,
		},
			Customization{
				Model: models.Address{},
				Type:  &Addresses.PickupAddress,
			},
		)
		suite.Len(customs, 4)

		customs, err := uniqueCustomizations(customs)
		suite.Len(customs, 4)
		suite.NoError(err)
	})
	suite.Run("Error if duplicate customization is used", func() {
		// Under test:       uniqueCustomizations checks that there's only one
		//                   customization of each type
		// Set up:           Create some customizations of the same type
		// Expected outcome: Error
		customs := getTraitActiveArmy() // contains a User type
		customs = append(customs, Customization{
			Model: models.User{},
			Type:  &User,
		})
		suite.Len(customs, 3)

		customs, err := uniqueCustomizations(customs)
		suite.Len(customs, 3)
		suite.ErrorContains(err, "Found more than one instance")

	})

}

func (suite *FactorySuite) TestElevateCustomization() {

	suite.Run("Customization converted ", func() {
		// Under test:       convertCustomizationInList converts the type of the customization
		//                   It's needed because dealing with the slice is finicky
		// Set up:           Create a ResidentialAddress customization, convert to Address
		// Expected outcome: No error
		customEmail := "leospaceman123@example.com"
		streetAddress := "235 Prospect Valley Road SE"

		customizationList :=
			[]Customization{
				{
					Model: models.User{LoginGovEmail: customEmail},
					Type:  &User,
				},
				{
					Model: models.Address{StreetAddress1: streetAddress},
					Type:  &Addresses.ResidentialAddress,
				},
			}

		// convert customization type from residentialAddress to Address
		tempCustoms := convertCustomizationInList(customizationList, Addresses.ResidentialAddress, Address)
		// Nothing wrong with customizations
		tempCustoms, err := uniqueCustomizations(tempCustoms)
		suite.NoError(err)
		// Customization has new type
		suite.Equal(Address, *tempCustoms[1].Type)
		// Old customization list is unchanged
		suite.Equal(Addresses.ResidentialAddress, *customizationList[1].Type)
	})
}
