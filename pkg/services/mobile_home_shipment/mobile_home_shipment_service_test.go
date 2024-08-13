package mobilehomeshipment

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MobileHomeShipmentSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *MobileHomeShipmentSuite) SetupTest() {

}

func TestMobileHomeShipmentServiceSuite(t *testing.T) {
	ts := &MobileHomeShipmentSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}



// func (suite *ModelSuite) verifyValidationErrors(model m.ValidateableModel, exp map[string][]string) {
// 	t := suite.T()
// 	t.Helper()

// 	verrs, err := model.Validate(nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if verrs.Count() != len(exp) {
// 		t.Errorf("expected %d errors, got %d", len(exp), verrs.Count())
// 	}

// 	var expKeys []string
// 	for key, errors := range exp {
// 		e := verrs.Get(key)
// 		expKeys = append(expKeys, key)
// 		if !sameStrings(e, errors) {
// 			t.Errorf("expected errors on %s to be %v, got %v", key, errors, e)
// 		}
// 	}

// 	for _, key := range verrs.Keys() {
// 		if !sliceContains(key, expKeys) {
// 			errors := verrs.Get(key)
// 			t.Errorf("unexpected validation errors on %s: %v", key, errors)
// 		}
// 	}
// }

// func TestModelSuite(t *testing.T) {
// 	hs := &ModelSuite{
// 		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
// 	}
// 	suite.Run(t, hs)
// 	hs.PopTestSuite.TearDown()
// }

// func sliceContains(needle string, haystack []string) bool {
// 	for _, s := range haystack {
// 		if s == needle {
// 			return true
// 		}
// 	}
// 	return false
// }

// func sameStrings(a []string, b []string) bool {
// 	sort.Strings(a)
// 	sort.Strings(b)
// 	return reflect.DeepEqual(a, b)
// }
