package models_test

import (
	"log"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ModelSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ModelSuite) SetupTest() {
	errTruncateAll := //RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return value in the file is used for test database teardown
		//RA: Given the database is being reset for unit test use, there are no unexpected states and conditions to account for
		//RA Developer Status: Mitigated
		//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
		//RA Validator: jneuner@mitre.org
		//RA Modified Severity:
		suite.DB().TruncateAll() // nolint:errcheck
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
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

func TestModelSuite(t *testing.T) {
	hs := &ModelSuite{PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage())}
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
