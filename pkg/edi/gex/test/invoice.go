package test

import (
	"net/http"
	"strings"
)

// GexTestSend mocks a request to gex
type GexTestSend struct{}

// SendRequest returns an automatic response with status 200
func (s GexTestSend) SendRequest(edi string, transactionName string) (resp *http.Response, err error) {
	// EDI parser will fail silently
	edi = strings.TrimSpace(edi) + "\n"
	//resp = http.Response(http.ResponseWriter.WriteHeader(http.StatusOK, 200))
	return resp, err
}
