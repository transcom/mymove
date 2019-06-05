package auth

import "errors"

func AuthorizeAdminUser(session *Session) error {
	if !session.IsSuperuser {
		return errors.New("USER_UNAUTHORIZED")
	}

	return nil
}
