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

// func (suite *FactorySuite) TestServiceMemberMaker() {
// 	sm := makeServiceMemberNew(suite.DB(), Variants{
// 		ServiceMemberCurrentAddress: models.Address{
// 			StreetAddress1: "This is my street",
// 		},

// 		User: models.User{
// 			LoginGovEmail: "shimonatests@onetwothree.com",
// 		},
// 	})
// 	fmt.Println(*sm.FirstName)
// 	fmt.Println(sm.User.LoginGovEmail)
// }

func (suite *MakerSuite) TestUserMaker() {
	// Create a factory with the end object

	user, err := userMaker(suite.DB(), []Customization{
		{
			Model: models.User{
				LoginGovEmail: "shimonatests@onetwothree.com",
			},
			Type:   CustomUser,
			Create: false,
		}}, []Trait{
		getTraitActiveUser,
	})
	suite.NoError(err)
	fmt.Println(user.LoginGovEmail, user.Active)

}

func (suite *MakerSuite) TestMergeCustomization() {
	uuidval := uuid.Must(uuid.NewV4())
	result := mergeCustomization(
		[]Trait{
			getTraitArmy,
		},
		[]Customization{
			{
				Model: models.User{
					LoginGovUUID:  &uuidval,
					LoginGovEmail: "custom@army.mil",
				},
				Type:   CustomUser,
				Create: false,
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
