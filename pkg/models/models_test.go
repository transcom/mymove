package models_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gobuffalo/validate"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ModelSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ModelSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func (suite *ModelSuite) verifyValidationErrors(model models.ValidateableModel, exp map[string][]string) {
	t := suite.T()
	t.Helper()

	verrs, err := model.Validate(suite.DB())
	if err != nil {
		t.Fatal(err)
	}

	if verrs.Count() != len(exp) {
		t.Errorf("expected %d errors, got %d", len(exp), verrs.Count())
	}

	var expKeys []string
	for key, errors := range exp {
		e := verrs.Get(key)
		expKeys = append(expKeys, key)
		if !sameStrings(e, errors) {
			t.Errorf("expected errors on %s to be %v, got %v", key, errors, e)
		}
	}

	for _, key := range verrs.Keys() {
		if !sliceContains(key, expKeys) {
			errors := verrs.Get(key)
			t.Errorf("unexpected validation errors on %s: %v", key, errors)
		}
	}
}

func (suite *ModelSuite) noValidationErrors(verrs *validate.Errors, err error) bool {
	noVerr := true
	if !suite.False(verrs.HasAny()) {
		noVerr = false
		for _, k := range verrs.Keys() {
			suite.Empty(verrs.Get(k))
		}
	}

	return suite.NoError(err) && noVerr
}

func TestModelSuite(t *testing.T) {
	hs := &ModelSuite{PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage())}
	suite.Run(t, hs)
}

func sliceContains(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func sameStrings(a []string, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	return reflect.DeepEqual(a, b)
}
