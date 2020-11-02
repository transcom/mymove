package paymentrequest

import "github.com/gobuffalo/pop/v5"

// RequestPaymentHelper is a helper to connect to the DB
type RequestPaymentHelper struct {
	DB *pop.Connection
}
