package test

import (
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"strings"
)

// GexTestSend mocks a request to gex
type GexTestSend struct{}

// SendRequest returns an automatic response with status 200
func (s GexTestSend) SendRequest(edi string, transactionName string) (resp *http.Response, err error) {
	// EDI parser will fail silently
	edi = strings.TrimSpace(edi) + "\n"

	var statusOKApiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	request, err := http.NewRequest(
		"POST",
		statusOKApiStub.URL,
		strings.NewReader(edi),
	)
	if err != nil {
		return resp, errors.Wrap(err, "Creating GEX POST request")
	}
	resp, err = http.DefaultClient.Do(request)

	return resp, err
}
