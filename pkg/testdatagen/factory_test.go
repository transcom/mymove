package testdatagen

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

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

func (suite *FactorySuite) TestServiceMemberFactory() {
	sm := ServiceMemberFactory(suite.DB(), Variants{
		ServiceMemberCurrentAddress: models.Address{
			StreetAddress1: "This is my street",
		},

		User: models.User{
			LoginGovEmail: "shimonatests@onetwothree.com",
		},
	})
	fmt.Println(unsafe.Sizeof(sm))
	fmt.Println(unsafe.Sizeof(Variants{}))
	fmt.Println(*sm.FirstName)
	fmt.Println(sm.User.LoginGovEmail)
}

func (suite *FactorySuite) TestServiceMemberFactoryS() {
	sm := makeSMX(suite.DB(), Variants{
		ServiceMemberCurrentAddress: models.Address{
			StreetAddress1: "This is my street",
		},

		User: models.User{
			LoginGovEmail: "shimonatests@onetwothree.com",
		},
	})
	fmt.Println(*sm.FirstName)
	fmt.Println(sm.User.LoginGovEmail)
}

func (suite *FactorySuite) TestUserFactory() {
	userFactory := NewUserFac(models.User{}, nil)
	userFactory.Create(suite.DB(), Variants{
		User: models.User{
			LoginGovEmail: "shimonatests@onetwothree.com",
		},
	})
	fmt.Println(userFactory.Model.LoginGovEmail)

}

func checkNestedVariant(rV reflect.Value) {
	// this is a variant struct within the variants object.
	// We want to instrospect and check that it does not contain a second level

	if rV.Kind() == reflect.Struct {
		value := rV
		numberOfFields := value.NumField()
		fmt.Println("    >> struct with ", numberOfFields, "fields.")
		for i := 0; i < numberOfFields; i++ {
			field := value.Field(i)
			fmt.Printf("    %s || %s \n",
				field.Type(), field.Kind())
			if field.Kind() == reflect.Pointer && !field.IsNil() {
				fmt.Println("Second level nesting of", field.Type(), " no allowed")
			}
			if field.Kind() == reflect.Struct && !field.IsZero() {
				fmt.Println("Second level nesting of", field.Type(), " no allowed")
			}
		}
	}

}
func showDetails(i interface{}) {
	t1 := reflect.TypeOf(i)
	v1 := reflect.ValueOf(i)
	k1 := v1.Kind()
	fmt.Println("Type of interface:", t1)
	// fmt.Println("Value of interface", v1)
	fmt.Println("Kind of interface:", k1)
	if reflect.ValueOf(i).Kind() == reflect.Pointer {
		//pointer := reflect.ValueOf(i)
		//value := pointer.Elem().Field(0)
		fmt.Println("It is a pointer.")
		t := reflect.TypeOf(i).Elem()
		k := t.Kind()
		fmt.Println("Type of interface:", t)
		fmt.Println("Kind of interface:", k)
	} else {
		fmt.Println("not a pointer")
	}

	if reflect.ValueOf(i).Kind() == reflect.Struct {

		value := reflect.ValueOf(i)
		numberOfFields := value.NumField()
		if value.IsZero() {
			fmt.Println("value is zero")
		}
		fmt.Println("It is a struct with ", numberOfFields, "fields.")
		for i := 0; i < numberOfFields; i++ {
			field := value.Field(i)
			fmt.Printf("%d. %s || %s \n",
				(i + 1), field.Type(), field.Kind())
			if field.Kind() == reflect.Struct {
				checkNestedVariant(field)
			}
		}
	}
}

type SpecialString string
type Square struct{ dim int }
type Varry struct {
	Address     models.Address
	MTOShipment models.MTOShipment
	User        models.User
}

func (suite *FactorySuite) TestUserFactoryX() {
	user := makeUserX(suite.DB(), Variants{
		User: models.User{
			LoginGovEmail: "shimonatests@onetwothree.com",
		},
		MTOShipment: models.MTOShipment{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
			MoveTaskOrder: models.Move{
				Locator: "12024",
			},
		},
	})
	// showDetails(models.User{})
	showDetails(Varry{
		User: models.User{
			LoginGovEmail: "shimonatests@onetwothree.com",
		},
		MTOShipment: models.MTOShipment{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
			MoveTaskOrder: models.Move{
				Locator: "12024",
			},
		},
	})

	fmt.Println(user.LoginGovEmail)
}
