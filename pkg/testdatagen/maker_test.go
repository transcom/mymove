package testdatagen

import (
	"fmt"
	"testing"

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
		[]GetTraitFunc{
			GetTraitArmy,
		},
		[]Customization{
			{
				Model: models.User{
					LoginGovUUID:  &uuidval,
					LoginGovEmail: "custom@army.mil",
				},
				Type: User,
			},
		})
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
