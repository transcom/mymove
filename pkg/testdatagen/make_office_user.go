package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeOfficeUser creates a single office user and associated TransportOffice
func MakeOfficeUser(db *pop.Connection) (models.OfficeUser, error) {
	user, err := MakeUser(db)
	if err != nil {
		return models.OfficeUser{}, err
	}

	office, _ := MakeTransportationOffice(db)

	officeUser := models.OfficeUser{
		UserID:                 &user.ID,
		User:                   &user,
		TransportationOffice:   office,
		TransportationOfficeID: office.ID,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  "leo_spaceman@example.com",
		Telephone:              "415-555-1212",
	}

	verrs, err := db.ValidateAndSave(&officeUser)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return officeUser, nil
}
