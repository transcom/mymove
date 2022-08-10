package testdatagen

import (
	"fmt"
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
		user, err := UserMaker(suite.DB(), nil, nil)
		suite.NoError(err)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.Equal(false, user.Active)
	})

	suite.Run("Successful creation of user with customization", func() {
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
		suite.Equal(true, user.Active)

	})

	suite.Run("Successful creation of user with trait", func() {
		user, err := UserMaker(suite.DB(), nil,
			[]GetTraitFunc{
				GetTraitActiveUser,
			})
		suite.NoError(err)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.Equal(true, user.Active)
	})
	suite.Run("Successful creation of user with both", func() {
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
		suite.Equal(true, user.Active)
	})

}

func (suite *MakerSuite) TestServiceMemberMaker() {
	defaultEmail := "first.last@login.gov.test"

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

	suite.Run("Successful creation of servicemember with premade user", func() {
		// Under test:      ServiceMemberMaker
		// Set up:          Create a service member with a passed in user
		// Expected outcome:ServiceMember should not create a user, but link to the provided user
		user, err := UserMaker(suite.DB(), nil, nil)
		suite.NoError(err)
		ServiceMember, err := ServiceMemberMaker(suite.DB(),
			[]Customization{
				{
					Model: user,
					Type:  User,
				},
			}, nil)

		suite.NoError(err)
		suite.Equal(user.ID, ServiceMember.UserID)

	})

	serviceMember, err := ServiceMemberMaker(suite.DB(),
		nil,
		[]GetTraitFunc{
			GetTraitActiveUser,
			GetTraitNavy,
		})
	suite.NoError(err)
	fmt.Println(serviceMember.Affiliation)
	fmt.Println(serviceMember.User.LoginGovEmail, serviceMember.User.Active)
	fmt.Println(serviceMember.User)

}

func (suite *MakerSuite) TestMergeCustomization() {
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
			GetTraitArmy,
		},
	)
	fmt.Println("User")
	userModel := result[0].Model.(models.User)
	fmt.Println(userModel)
}

func (suite *MakerSuite) TestMergeInterfaces() {
	user1 := models.User{
		LoginGovEmail: "user1@email.com",
	}
	uuidNew := uuid.Must(uuid.NewV4())
	user2 := models.User{
		LoginGovUUID: &uuidNew,
	}
	result := mergeInterfaces(user1, user2)
	user := result.(models.User)
	fmt.Println(user.LoginGovEmail, user.LoginGovUUID)
}

func (suite *MakerSuite) TestHasID() {
	suite.Run("True if ID is provided", func() {
		testid := uuid.Must(uuid.NewV4())
		result := hasID(models.ServiceMember{
			ID: testid,
		})
		suite.True(result)
	})

	suite.Run("False if no id", func() {
		result := hasID(models.ServiceMember{})
		suite.False(result)
	})

	suite.Run("False if nil id", func() {
		result := hasID(models.ServiceMember{
			ID: uuid.Nil,
		})
		suite.False(result)
	})
}

func (suite *MakerSuite) TestNestedModelsCheck() {
	suite.Run("Must not call with pointer", func() {
		c := Customization{
			Model: models.ServiceMember{},
			Type:  ServiceMember,
		}
		err := checkNestedModels(&c)
		suite.Error(err)
		suite.Contains(err.Error(), "received a pointer")
	})

	suite.Run("Must not have missing model", func() {
		c := Customization{
			Type: ServiceMember,
		}
		err := checkNestedModels(c)
		suite.Error(err)
		suite.Contains(err.Error(), "must contain a model")
	})

	suite.Run("Customization contains nested model", func() {

		// Expect nested user to cause error
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

		// Expect nested user to cause error
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
