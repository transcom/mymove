package testharness

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func MakeNeedsOrdersUser(db *pop.Connection) models.User {
	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))
	user := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				OktaEmail: email,
				Active:    true,
			},
		},
	}, nil)

	suffix := strings.Split(user.ID.String(), "-")[0]

	factory.BuildExtendedServiceMember(db, []factory.Customization{
		{
			Model: models.ServiceMember{
				PersonalEmail: models.StringPointer(email),
				FirstName:     models.StringPointer("NEEDS" + suffix),
				LastName:      models.StringPointer("ORDERS" + suffix),
				CacValidated:  true,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	return user
}
