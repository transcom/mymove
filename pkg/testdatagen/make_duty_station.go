package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeDutyStation creates a single DutyStation
func MakeDutyStation(db *pop.Connection, name string, branch internalmessages.MilitaryBranch, address models.Address) (models.DutyStation, error) {
	var station models.DutyStation

	verrs, err := db.ValidateAndSave(&address)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	station = models.DutyStation{
		Name:      name,
		Branch:    branch,
		AddressID: address.ID,
	}

	verrs, err = db.ValidateAndSave(&station)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return station, err
}
