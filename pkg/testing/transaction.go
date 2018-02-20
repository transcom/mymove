package testing

import (
	"testing"

	"github.com/markbates/pop"
)

// StartTransaction starts a database transaction and returns that transaction for use in place of a *pop.Connection and a rollback func, which must be called to rollback the transaction at the end of the test.
//
// Note that Pop simply returns the current transaction if one has already been started. Using this helper within tests that exercise code that works with transactions itself may have unexpected results.
//
// Examples
//
//     func ExampleTest(t *testing.T) {
//         tx, rollback := StartTransaction(t, dbConnection)
//         defer rollback()
//
//         // ...rest of test, using tx in place of dbConnection
//     }
//
func StartTransaction(t *testing.T, db *pop.Connection) (*pop.Connection, func()) {
	t.Helper()

	tx, err := db.NewTransaction()
	if err != nil {
		t.Fatalf("transaction could not be started: %s", err)
	}
	rollback := func() {
		t.Helper()

		err := tx.TX.Rollback()
		if err != nil {
			t.Fatal(err)
		}
	}
	return tx, rollback
}

// MustSave handles saving and error checking for Pop model creation inside tests.
func MustSave(t *testing.T, db *pop.Connection, s interface{}) {
	t.Helper()

	verrs, err := db.ValidateAndSave(s)
	if err != nil {
		t.Fatal(err)
	}
	if verrs.Count() > 0 {
		t.Fatalf("errors encountered saving %v: %v", s, verrs)
	}
}
