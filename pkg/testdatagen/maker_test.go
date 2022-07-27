package testdatagen

import (
	"fmt"
	"testing"

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
	userFactory := NewUserMaker(models.User{}, nil)
	// Call create to create the object

	err := userFactory.Make(suite.DB(), Customization{
		Model: models.User{
			LoginGovEmail: "shimonatests@onetwothree.com",
		},
		Name:   "User",
		Create: false,
	})
	suite.NoError(err)
	fmt.Println(userFactory.Model.LoginGovEmail)

}

// func (suite *FactorySuite) TestUserMaker() {
// 	user := makeUserNew(suite.DB(), Variants{
// 		User: models.User{
// 			LoginGovEmail: "shimonatests@onetwothree.com",
// 		},
// 		MTOShipment: models.MTOShipment{
// 			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
// 			MoveTaskOrder: models.Move{
// 				Locator: "12024",
// 			},
// 		},
// 	})
// 	// showDetails(models.User{})
// 	showDetails(Varry{
// 		User: models.User{
// 			LoginGovEmail: "shimonatests@onetwothree.com",
// 		},
// 		MTOShipment: models.MTOShipment{
// 			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
// 			MoveTaskOrder: models.Move{
// 				Locator: "12024",
// 			},
// 		},
// 	})

// 	fmt.Println(user.LoginGovEmail)
// }
