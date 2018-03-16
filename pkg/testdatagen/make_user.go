package testdatagen

import (
	"log"

	"github.com/markbates/pop"
	"github.com/satori/go.uuid"

	"github.com/transcom/mymove/pkg/models"
)

func MakeUser(db *pop.Connection) (models.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Panic(err)
	}

	user := models.User{
		LoginGovUUID:  id,
		LoginGovEmail: "first.last@login.gov.test",
	}

	verrs, err := db.ValidateAndSave(&user)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return user, err
}
