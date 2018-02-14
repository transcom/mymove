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

func verifyValidationErrors(v *validate.Errors, exp map[string][]string, t *testing.T) {
	if v.Count() != len(exp) {
		t.Errorf("expected %d errors, got %d", len(exp), v.Count())
	}

	for key, errors := range exp {
		e := v.Get(key)
		if !equalSlice(e, errors) {
			t.Errorf("expected errors on %s to be %v, got %v", key, e, errors)
		}
	}
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

	os.Exit(m.Run())
}
