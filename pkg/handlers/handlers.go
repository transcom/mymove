package handlers

import (
	"github.com/markbates/pop"
)

// this file defines globals used by all handlers.

var dbConnection *pop.Connection

// Init the API package with its database connection
func Init(dbInitialConnection *pop.Connection) {
	dbConnection = dbInitialConnection
}
