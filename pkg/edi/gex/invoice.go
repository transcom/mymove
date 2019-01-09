package gex

import (
	"crypto/tls"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// gexRequestTimeout is how long to wait on Gex request before timing out (30 seconds).
const gexRequestTimeout = time.Duration(30) * time.Second

// SendToGex is an interface for sending and receiving a request
type SendToGex interface {
	Call(edi string, transactionName string) (resp *http.Response, err error)
}

// SendToGexHTTP represents a struct to contain an actual gex request function
type SendToGexHTTP struct {
	URL                  string
	IsTrueGexURL         bool
	TLSConfig            *tls.Config
	GEXBasicAuthUsername string
	GEXBasicAuthPassword string
}

// Call sends an edi file string as a POST to the gex api
func (s SendToGexHTTP) Call(edi string, transactionName string) (resp *http.Response, err error) {
	// Ensure that the transaction body ends with a newline, otherwise the GEX EDI parser will fail silently
	edi = strings.TrimSpace(edi) + "\n"
	URL := s.URL
	if s.IsTrueGexURL {
		URL = filepath.Join(s.URL, transactionName)
	}
	request, err := http.NewRequest(
		"POST",
		URL,
		strings.NewReader(edi),
	)
	if err != nil {
		return resp, errors.Wrap(err, "Creating GEX POST request")
	}

	// We need to provide basic auth credentials for the GEX server, as well as
	// our client certificate for the proxy in front of the GEX server.
	request.SetBasicAuth(s.GEXBasicAuthUsername, s.GEXBasicAuthPassword)

	tr := &http.Transport{TLSClientConfig: s.TLSConfig}
	client := &http.Client{Transport: tr, Timeout: gexRequestTimeout}

	resp, err = client.Do(request)
	if err != nil {
		return resp, errors.Wrap(err, "Sending GEX POST request")
	}

	return resp, err
}
