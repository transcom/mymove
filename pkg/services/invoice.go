package services

import "net/http"

// GexSender is an interface for sending and receiving a request
type GexSender interface {
	SendToGex(edi string, transactionName string) (resp *http.Response, err error)
}
