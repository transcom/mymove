package services

import "net/http"

// SendToGex is an interface for sending and receiving a request
type SendToGex interface {
	Call(edi string, transactionName string) (resp *http.Response, err error)
}
