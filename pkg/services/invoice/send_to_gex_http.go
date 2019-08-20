package invoice

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
)

// gexRequestTimeout is how long to wait on Gex request before timing out (30 seconds).
const gexRequestTimeout = time.Duration(30) * time.Second

// NewGexSenderHTTP creates a new GexSender service object
func NewGexSenderHTTP(url string, isTrueGexURL bool, tlsConfig *tls.Config, gexBasicAuthUsername string, gexBasicAuthPassword string) services.GexSender {
	return &gexSenderHTTP{
		url,
		isTrueGexURL,
		tlsConfig,
		gexBasicAuthUsername,
		gexBasicAuthPassword,
	}
}

// gexSenderHTTP represents a struct to contain an actual gex request function
type gexSenderHTTP struct {
	url                  string
	isTrueGexURL         bool
	tlsConfig            *tls.Config
	gexBasicAuthUsername string
	gexBasicAuthPassword string
}

// SendToGex sends an edi file string as a POST to the gex api
// To set local dev to send a real GEX request, replace your env.local:
// export GEX_URL=""  with "export GEX_URL=https://gexweba.daas.dla.mil/msg_data/submit/"
func (s *gexSenderHTTP) SendToGex(edi string, transactionName string) (resp *http.Response, err error) {
	// Ensure that the transaction body ends with a newline, otherwise the GEX EDI parser will fail silently
	edi = strings.TrimSpace(edi) + "\n"
	URL := s.url
	parsedURL, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}

	if s.isTrueGexURL {
		parsedURL.Path = parsedURL.Path + transactionName
		URL = parsedURL.String()
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
	request.SetBasicAuth(s.gexBasicAuthUsername, s.gexBasicAuthPassword)

	tr := &http.Transport{TLSClientConfig: s.tlsConfig}
	client := &http.Client{Transport: tr, Timeout: gexRequestTimeout}

	resp, err = client.Do(request)
	if err != nil {
		return resp, errors.Wrap(err, "Sending GEX POST request")
	}

	return resp, err
}
