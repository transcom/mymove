package models_test

import (
	"context"
	"log"
	"reflect"
	"sort"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
)

type ModelSuite struct {
	suite.Suite
	db *pop.Connection
}

func (suite *ModelSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *ModelSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		t.Fatalf("errors encountered saving %v: %v", model, verrs)
	}
}

func (suite *ModelSuite) verifyValidationErrors(model models.ValidateableModel, exp map[string][]string) {
	t := suite.T()
	t.Helper()

	verrs, err := model.Validate(suite.db)
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

	return !suite.NoError(err) && noVerr
}

func TestModelSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	hs := &ModelSuite{db: db}
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
