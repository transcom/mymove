package paymentrequest

import "github.com/gobuffalo/pop"

// RequestPaymentHelper is a helper to connect to the DB
type RequestPaymentHelper struct {
	DB     *pop.Connection
	Logger Logger
}
