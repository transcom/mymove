package handlers

import (
	"github.com/markbates/pop"
)

func stubDB(newDB *pop.Connection, body func()) {
	originalDB := DB
	DB = newDB
	body()
	DB = originalDB
}
