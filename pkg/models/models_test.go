package models_test

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/stretchr/testify/suite"
)

type ModelSuite struct {
	suite.Suite
	db *pop.Connection
}

func (suite *ModelSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *ModelSuite) mustSave(model interface{}) {
	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		suite.T().Fatalf("errors encountered saving %v: %v", model, verrs)
	}
}

func (suite *ModelSuite) verifyValidationErrors(model validatableModel, exp map[string][]string) {
	t := suite.T()
	t.Helper()

	verrs, err := model.Validate(suite.db)
	if err != nil {
		t.Fatal(err)
	}

	if verrs.Count() != len(exp) {
		t.Errorf("expected %d errors, got %d", len(exp), verrs.Count())
	}

	expKeys := []string{}
	for key, errors := range exp {
		e := verrs.Get(key)
		expKeys = append(expKeys, key)
		if !equalSlice(e, errors) {
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
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	hs := &ModelSuite{db: db}
	suite.Run(t, hs)
}

type validatableModel interface {
	Validate(*pop.Connection) (*validate.Errors, error)
}

func sliceContains(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func equalSlice(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
