package handlers

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	gexop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/gex"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// SendGexRequestHandler sends a request to GEX
type SendGexRequestHandler HandlerContext

// Handle sends a request to GEX
func (h SendGexRequestHandler) Handle(params gexop.SendGexRequestParams) middleware.Responder {
	transactionName := *params.SendGexRequestPayload.TransactionName
	transactionBody := *params.SendGexRequestPayload.TransactionBody

	// Ensure that the transaction body ends with a newline, otherwise the GEX
	// EDI parser will fail silently
	transactionBody = strings.TrimSpace(transactionBody) + "\n"

	request, err := http.NewRequest(
		"POST",
		"https://gexweba.daas.dla.mil/msg_data/submit/"+transactionName,
		strings.NewReader(transactionBody),
	)
	if err != nil {
		h.logger.Error("Creating GEX POST request", zap.Error(err))
		return gexop.NewSendGexRequestInternalServerError()
	}

	// We need to provide basic auth credentials for the GEX server, as well as
	// our client certificate for the proxy in front of the GEX server.
	request.SetBasicAuth(os.Getenv("GEX_BASIC_AUTH_USERNAME"), os.Getenv("GEX_BASIC_AUTH_PASSWORD"))

	config, err := getTLSConfig()
	if err != nil {
		h.logger.Error("Creating TLS config", zap.Error(err))
		return gexop.NewSendGexRequestInternalServerError()
	}

	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		h.logger.Error("Sending GEX POST request", zap.Error(err))
		return gexop.NewSendGexRequestInternalServerError()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	responseBody := buf.String()

	responsePayload := internalmessages.GexResponsePayload{
		GexResponse: resp.Status + "; " + responseBody,
	}
	return gexop.NewSendGexRequestOK().WithPayload(&responsePayload)
}

func getTLSConfig() (*tls.Config, error) {
	clientCert := os.Getenv("CLIENT_TLS_CERT")
	clientKey := os.Getenv("CLIENT_TLS_KEY")
	certificate, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM([]byte(os.Getenv("GEX_DOD_CA")))

	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      rootCAs,
	}, nil
}
