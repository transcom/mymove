package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceMember creates a single ServiceMember and associated User.
func MakeServiceMember(db *pop.Connection) (models.ServiceMember, error) {
	user, err := MakeUser(db)
	if err != nil {
		return models.ServiceMember{}, err
	}

	serviceMember := models.ServiceMember{
		UserID:        user.ID,
		User:          user,
		FirstName:     models.StringPointer("Leo"),
		LastName:      models.StringPointer("Spacemen"),
		PersonalEmail: models.StringPointer("leo@example.com"),
	}

	verrs, err := db.ValidateAndSave(&serviceMember)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return serviceMember, err
}

// MakeServiceMemberData created 5 ServiceMembers (and in turn a User for each)
func MakeServiceMemberData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeServiceMember(db)
	}
}
