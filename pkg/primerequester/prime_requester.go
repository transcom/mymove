package primerequester

import "github.com/transcom/mymove/pkg/models"

// SaveAccess saves prime requester field to the database
func SaveAccess(requester *string, clientCert *models.ClientCert) error {

	// TODO

	// 1. Query prime_requester for an existing user using name and client cert
	// 2. If the user exists update last_seen_at=now()
	// 3. If the user does not exist, insert new record and set allow_access=true and last_seen_at=now()

	return nil
}