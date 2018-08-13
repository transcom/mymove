package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTspUser creates a single TSP user and associated Transportation Service Provider
func MakeTspUser(db *pop.Connection, assertions Assertions) models.TspUser {
	user := assertions.TspUser.User
	if assertions.TspUser.UserID == nil || isZeroUUID(*assertions.TspUser.UserID) {
		user = MakeUser(db, assertions)
	}

	tsp := MakeDefaultTSP(db)

	tspUser := models.TspUser{
		UserID: &user.ID,
		User:   user,
		TransportationServiceProvider:   tsp,
		TransportationServiceProviderID: tsp.ID,
		FirstName:                       "Leo",
		LastName:                        "Spaceman",
		Email:                           "leo_spaceman@example.com",
		Telephone:                       "415-555-1212",
	}

	mergeModels(&tspUser, assertions.TspUser)

	mustCreate(db, &tspUser)

	return tspUser
}

// MakeDefaultTspUser makes an TspUser with default values
func MakeDefaultTspUser(db *pop.Connection) models.TspUser {
	return MakeTspUser(db, Assertions{})
}
