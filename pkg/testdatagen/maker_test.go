package testdatagen

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

func (suite *MakerSuite) TestUserMaker() {
	defaultEmail := "first.last@login.gov.test"
	customEmail := "leospaceman123@example.com"
	suite.Run("Successful creation of default user", func() {
		// Under test:      UserMaker
		// Mocked:          None
		// Set up:          Create a User with no customizations or traits
		// Expected outcome:User should be created with default values

		user, err := UserMaker(suite.DB(), nil, nil)
		suite.NoError(err)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.False(user.Active)
	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:      UserMaker
		// Set up:          Create a User with a customized email and no trait
		// Expected outcome:User should be created with email and inactive status
		user, err := UserMaker(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
				Type: User,
			},
		}, nil)
		suite.NoError(err)
		suite.Equal(customEmail, user.LoginGovEmail)
		suite.False(user.Active)

	})

	suite.Run("Successful creation of user with trait", func() {
		// Under test:      UserMaker
		// Set up:          Create a User with a trait
		// Expected outcome:User should be created with default email and active status

		user, err := UserMaker(suite.DB(), nil,
			[]GetTraitFunc{
				GetTraitActiveUser,
			})
		suite.NoError(err)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

	suite.Run("Successful creation of user with both", func() {
		// Under test:      UserMaker
		// Set up:          Create a User with a customized email and active trait
		// Expected outcome:User should be created with email and active status

		user, err := UserMaker(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
				Type: User,
			}}, []GetTraitFunc{
			GetTraitActiveUser,
		})
		suite.NoError(err)
		suite.Equal(customEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

}

func (suite *MakerSuite) TestServiceMemberMaker() {
	defaultEmail := "first.last@login.gov.test"
	customEmail := "leospaceman47@example.com"

	suite.Run("Successful creation of servicemember", func() {
		// Under test:      ServiceMemberMaker
		// Mocked:          None
		// Set up:          Create a service member with no customizations or traits
		// Expected outcome:ServiceMember should be created with default values
		serviceMember, err := ServiceMemberMaker(suite.DB(), nil, nil)
		suite.NoError(err)
		suite.NotEqual(uuid.Nil, serviceMember.UserID)
		suite.Equal(defaultEmail, serviceMember.User.LoginGovEmail)
		suite.False(serviceMember.User.Active)

	})

	suite.Run("Successful creation of servicemember with customization only", func() {
		// Under test:       ServiceMemberMaker
		// Set up:           Create a service member with a customized User and ResidentialAddress
		// Expected outcome: ServiceMember should be created with the right customizations
		streetAddress := "448 Washington Blvd NE"
		sm, err := ServiceMemberMaker(suite.DB(),
			[]Customization{
				{
					Model: models.User{LoginGovEmail: customEmail},
					Type:  User,
				},
				{
					Model: models.Address{StreetAddress1: streetAddress},
					Type:  Addresses.ResidentialAddress,
				},
			}, nil)

		suite.NoError(err)
		suite.Equal(customEmail, sm.User.LoginGovEmail)
		suite.Equal(streetAddress, sm.ResidentialAddress.StreetAddress1)
		suite.NotEqual(streetAddress, sm.BackupMailingAddress)

	})

	suite.Run("Successful creation of servicemember with customization and trait", func() {
		// Under test:       ServiceMemberMaker
		// Set up:           Create a service member with a customized User and ResidentialAddress
		//                   as well as the active user trait
		// Expected outcome: ServiceMember should be created with the right customizations, and active
		streetAddress := "448 Washington Blvd NE"
		sm, err := ServiceMemberMaker(suite.DB(),
			[]Customization{
				{
					Model: models.User{LoginGovEmail: customEmail},
					Type:  User,
				},
				{
					Model: models.Address{StreetAddress1: streetAddress},
					Type:  Addresses.ResidentialAddress,
				},
			},
			[]GetTraitFunc{
				GetTraitActiveUser,
			})

		suite.NoError(err)
		suite.True(sm.User.Active)
		suite.Equal(customEmail, sm.User.LoginGovEmail)
		suite.Equal(streetAddress, sm.ResidentialAddress.StreetAddress1)
		suite.NotEqual(streetAddress, sm.BackupMailingAddress)

	})

	suite.Run("Successful creation of servicemember with premade user", func() {
		// Under test:       ServiceMemberMaker
		// Set up:           Create a service member with a pre-created user, add a trait that affects user
		// Expected outcome: ServiceMember should not create a user, but link to the provided user
		//                   Trait should not be applied because user has already been created
		user, err := UserMaker(suite.DB(), nil, nil)
		suite.NoError(err)
		sm, err := ServiceMemberMaker(suite.DB(),
			[]Customization{
				{
					Model: user,
					Type:  User,
				},
			},
			[]GetTraitFunc{
				GetTraitActiveUser,
			})

		suite.NoError(err)
		suite.Equal(user.ID, sm.UserID)
		suite.False(sm.User.Active)

	})

	suite.Run("Successful creation of servicemember with only traits", func() {
		// Under test:       ServiceMemberMaker
		// Set up:           Create a service member with only traits
		// Expected outcome: ServiceMember should create an active User, Navymember

		// Create a service member with getTraitArmy
		// Both getTraitArmy and getTraitActiveUser change fields in User
		// They should merge successfully and show up in the end objects
		serviceMember, err := ServiceMemberMaker(suite.DB(),
			nil,
			[]GetTraitFunc{
				GetTraitActiveUser,
				getTraitArmy,
			})
		suite.NoError(err)
		suite.Equal("trait@army.mil", serviceMember.User.LoginGovEmail)
		suite.True(serviceMember.User.Active)
		suite.Equal(models.AffiliationARMY, *serviceMember.Affiliation)
	})

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
			[]GetTraitFunc{
				getTraitArmy,
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
			[]GetTraitFunc{
				getTraitArmy,
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
		user, err := UserMaker(suite.DB(), nil, nil)
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

// GetTraitArmy is a custom GetTraitFunc
func getTraitArmy() []Customization {
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
