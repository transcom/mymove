package models

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
)

var dbConnection *pop.Connection

func setupDBConnection() {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	dbConnection = conn
}

type validatableModel interface {
	Validate(*pop.Connection) (*validate.Errors, error)
}

func verifyValidationErrors(model validatableModel, exp map[string][]string, t *testing.T) {
	t.Helper()

	verrs, err := model.Validate(dbConnection)
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
			t.Errorf("expected errors on %s to be %v, got %v", key, e, errors)
		}
	}

	for _, key := range verrs.Keys() {
		if !sliceContains(key, expKeys) {
			errors := verrs.Get(key)
			t.Errorf("unexpected validation errors on %s: %v", key, errors)
		}
	}
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

func TestMain(m *testing.M) {
	setupDBConnection()

	dbConnection.TruncateAll()

	os.Exit(m.Run())
}
