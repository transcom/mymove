package sequence

import "github.com/gobuffalo/pop"

// Package for using PostgreSQL Sequences
// https://www.postgresql.org/docs/10/sql-createsequence.html

// NextVal returns the next value of the given sequence
func NextVal(db *pop.Connection, sequence string) (int64, error) {
	var nextVal int64
	query := "SELECT nextval($1);"
	err := db.RawQuery(query, sequence).First(&nextVal)
	return nextVal, err
}
