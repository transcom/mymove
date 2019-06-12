package accesscontrols

import (
	"errors"

	"github.com/transcom/mymove/pkg/auth"
)

func AuthorizeAdminUser(session *auth.Session) error {
	if !session.IsSuperuser {
		return errors.New("USER_UNAUTHORIZED")
	}

	return nil
}
