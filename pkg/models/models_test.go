package models_test

import (
	"reflect"
	"sort"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ModelSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *ModelSuite) SetupTest() {

}

func (suite *ModelSuite) verifyValidationErrors(model ValidateableModel, exp map[string][]string) {
	t := suite.T()
	t.Helper()

	verrs, err := model.Validate(nil)
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

func TestModelSuite(t *testing.T) {
	hs := &ModelSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
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
