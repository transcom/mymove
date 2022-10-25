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

type MakerSuite struct {
	*testingsuite.PopTestSuite
}

func TestMakerSuite(t *testing.T) {

	ts := &MakerSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *MakerSuite) TestElevateCustomization() {

	suite.Run("Customization converted ", func() {

		customEmail := "leospaceman123@example.com"
		streetAddress := "235 Prospect Valley Road SE"

		customizationList :=
			[]Customization{
				{
					Model: models.User{LoginGovEmail: customEmail},
					Type:  User,
				},
				{
					Model: models.Address{StreetAddress1: streetAddress},
					Type:  Addresses.ResidentialAddress,
				},
			}
		tempCustoms := convertCustomizationInList(customizationList, Addresses.ResidentialAddress, Address)
		// Nothing wrong with customizations
		tempCustoms, err := validateCustomizations(tempCustoms)
		suite.NoError(err)
		// Customization has new type
		suite.Equal(Address, tempCustoms[1].Type)
		// Old customization list is unchanged
		suite.Equal(Addresses.ResidentialAddress, customizationList[1].Type)
	})
}

// getTraitActiveArmy is a custom Trait
func getTraitActiveArmy() []Customization {
	army := models.AffiliationARMY
	return []Customization{
		{
			Model: models.User{
				LoginGovEmail:         "trait@army.mil",
				CurrentAdminSessionID: "my-session-id",
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
}

func (suite *MakerSuite) TestMergeCustomization() {

	suite.Run("Customizations and traits merged into result", func() {
		// Under test:       mergeCustomization, which merges traits and customizations
		// Set up:           Create a customization without a matching trait, and traits with no customizations
		// Expected outcome: All should exist in the result list
		streetAddress := "235 Prospect Valley Road SE"
		result := mergeCustomization(
			// Address customization
			[]Customization{
				{
					Model: models.Address{
						StreetAddress1: streetAddress,
					},
					Type: Address,
				},
			},
			// User and ServiceMember customization
			[]Trait{
				getTraitActiveArmy,
			},
		)
		suite.Len(result, 3)
		_, custom := findCustomWithIdx(result, Address)
		suite.NotNil(custom)
		address := custom.Model.(models.Address)
		suite.Equal(streetAddress, address.StreetAddress1)

		// Check that User customization from trait is available
		_, custom = findCustomWithIdx(result, User)
		suite.NotNil(custom)
		user := custom.Model.(models.User)
		suite.Equal("trait@army.mil", user.LoginGovEmail)
		// Check that ServiceMember customization from trait is available
		_, custom = findCustomWithIdx(result, ServiceMember)
		suite.NotNil(custom)
	})

	suite.Run("Customization should override traits", func() {
		// Under test:       mergeCustomization, which merges traits and customizations
		// Set up:           Create a customization with a user email and a trait with a user email
		// Expected outcome: Customization should override the trait email
		//                   If an object exists and no customization, it should become a customization
		uuidval := uuid.Must(uuid.NewV4())
		result := mergeCustomization(
			[]Customization{
				{
					Model: models.User{
						LoginGovUUID:  &uuidval,
						LoginGovEmail: "custom@army.mil",
					},
					Type: User,
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

func (suite *MakerSuite) TestMergeInterfaces() {
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

func (suite *MakerSuite) TestHasID() {
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

func (suite *MakerSuite) TestNestedModelsCheck() {
	suite.Run("Must not call with pointer", func() {
		// Under test:       checkNestedModels, uses reflection to check if other models are nested
		// Set up:           Call with a pointer, instead of a struct
		// Expected outcome: Error

		c := Customization{
			Model: models.ServiceMember{},
			Type:  ServiceMember,
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
			Type: ServiceMember,
		}
		err := checkNestedModels(c)
		suite.Error(err)
		suite.Contains(err.Error(), "must contain a model")
	})

	suite.Run("Customization contains nested model", func() {
		// Under test:       checkNestedModels
		// Set up:           Call with a struct that contains a nested model
		// Expected outcome: Error
		user, err := BuildUser(suite.DB(), nil, nil)
		suite.NoError(err)
		c := Customization{
			Model: models.ServiceMember{
				User: user,
			},
			Type: ServiceMember,
		}
		err = checkNestedModels(c)
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
			Type: ServiceMember,
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
			Type: ServiceMember,
		}
		err := checkNestedModels(c)
		suite.NoError(err)

	})

}

func (suite *MakerSuite) TestValidateCustomizations() {
	suite.Run("Control obj added if missing", func() {
		customs := getTraitActiveArmy()
		suite.Len(customs, 2)

		customs, err := validateCustomizations(customs)
		suite.Nil(err)
		suite.Len(customs, 3)
		_, controlCustom := findCustomWithIdx(customs, control)
		suite.NotNil(controlCustom)
	})

	suite.Run("Control obj not added if not missing", func() {
		customs := getTraitActiveArmy()
		customs = append(customs, Customization{
			Model: controlObject{},
			Type:  control,
		})
		suite.Len(customs, 3)

		customs, err := validateCustomizations(customs)
		suite.Len(customs, 3)
		suite.NoError(err)
		_, controlCustom := findCustomWithIdx(customs, control)
		suite.NotNil(controlCustom)
	})

	suite.Run("Pass if customizations not repeated", func() {
		customs := getTraitActiveArmy()
		customs = append(customs, Customization{
			Model: models.Address{},
			Type:  Addresses.ResidentialAddress,
		},
			Customization{
				Model: models.Address{},
				Type:  Addresses.PickupAddress,
			},
		)
		suite.Len(customs, 4)

		customs, err := validateCustomizations(customs)
		suite.Len(customs, 5)
		suite.NoError(err)
		_, controlCustom := findCustomWithIdx(customs, control)
		suite.NotNil(controlCustom)
	})
	suite.Run("Error if duplicate customization is used", func() {
		customs := getTraitActiveArmy()
		customs = append(customs, Customization{
			Model: models.User{},
			Type:  User,
		})
		suite.Len(customs, 3)

		customs, err := validateCustomizations(customs)
		suite.Len(customs, 4) // control object should be added
		suite.ErrorContains(err, "Found more than one instance")

		// Check that control object was updated
		_, controlCustom := findCustomWithIdx(customs, control)
		controller := (*controlCustom).Model.(controlObject)
		suite.False(controller.isValid)
	})

}
