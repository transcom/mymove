package invoice

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

// SendToGex sends a body string as a POST to the gex api
// To set local dev to send a real GEX request, replace your env.local:
// export GEX_URL=""  with "export GEX_URL=https://gexweba.daas.dla.mil/msg_data/submit/"
func (s *gexSenderHTTP) SendToGex(channel services.GEXChannel, body string, filename string) (resp *http.Response, err error) {
	// Ensure that the transaction body ends with a newline, otherwise the GEX EDI parser will fail silently
	body = strings.TrimSpace(body) + "\n"
	URL := s.url
	parsedURL, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}

	if s.isTrueGexURL {
		URL = parsedURL.String()
	}

	// ensure channel is one of the valid expected ones
	if !isChannelValid(channel) {
		return nil, fmt.Errorf("Invalid channel type, expected %q", validGEXChannels())
	}

	request, err := http.NewRequest(
		"POST",
		URL,
		strings.NewReader(body),
	)
	if err != nil {
		return resp, fmt.Errorf("Creating GEX POST request: %w", err)
	}

	q := request.URL.Query()
	q.Add("fname", filename)
	q.Add("channel", string(channel))
	request.URL.RawQuery = q.Encode()

	// We need to provide basic auth credentials for the GEX server, as well as
	// our client certificate for the proxy in front of the GEX server.
	request.SetBasicAuth(s.gexBasicAuthUsername, s.gexBasicAuthPassword)

	tr := &http.Transport{TLSClientConfig: s.tlsConfig}
	client := &http.Client{Transport: tr, Timeout: gexRequestTimeout}

	resp, err = client.Do(request)
	if err != nil {
		return resp, fmt.Errorf("Sending GEX POST request: %w", err)
	}

	return resp, err
}

func validGEXChannels() []services.GEXChannel {
	return []services.GEXChannel{
		services.GEXChannelInvoice,
		services.GEXChannelDataWarehouse,
	}
}

func isChannelValid(channel services.GEXChannel) bool {
	for _, c := range validGEXChannels() {
		if channel == c {
			return true
		}
	}
	return false
}
