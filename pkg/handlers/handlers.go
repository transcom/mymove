package handlers

import (
	"github.com/markbates/pop"
)

// DB is the global *pop.Connection used by handlers.
var DB *pop.Connection

// Init the API package with its database connection
func Init(dbInitialConnection *pop.Connection) {
	DB = dbInitialConnection
}
