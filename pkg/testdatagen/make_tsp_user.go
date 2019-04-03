package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTspUser creates a single TSP user and associated Transportation Service Provider
func MakeTspUser(db *pop.Connection, assertions Assertions) models.TspUser {

	user := assertions.TspUser.User
	if assertions.TspUser.UserID == nil || isZeroUUID(*assertions.TspUser.UserID) {
		if assertions.User.LoginGovEmail == "" {
			// There's a uniqueness constraint on office user emails so add some randomness
			email := fmt.Sprintf("leo_spaceman_tsp_%s@example.com", makeRandomString(5))
			assertions.User.LoginGovEmail = email
		}
		user = MakeUser(db, assertions)
	}

	var tsp models.TransportationServiceProvider
	var tspAssertions = assertions.TransportationServiceProvider
	if &tspAssertions == nil {
		tsp = MakeDefaultTSP(db)
	} else {
		tsp = MakeTSP(db, Assertions{
			TransportationServiceProvider: tspAssertions,
		})
	}

	tspUser := models.TspUser{
		UserID:                          &user.ID,
		User:                            user,
		TransportationServiceProvider:   tsp,
		TransportationServiceProviderID: tsp.ID,
		FirstName:                       "Leo",
		LastName:                        "Spaceman",
		Email:                           user.LoginGovEmail,
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
