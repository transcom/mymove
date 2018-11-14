package sequence

import "github.com/gobuffalo/pop"

// Package for using PostgreSQL Sequences
// https://www.postgresql.org/docs/10/sql-createsequence.html

// NextVal returns the next value of the given sequence
func NextVal(db *pop.Connection, sequence string) (int64, error) {
	var nextVal int64
	err := db.RawQuery("SELECT nextval($1);", sequence).First(&nextVal)
	return nextVal, err
}

// SetVal sets the current value of a sequence
func SetVal(db *pop.Connection, sequence string, val int64) error {
	err := db.RawQuery("SELECT setval($1, $2)", sequence, val).Exec()
	return err
}
