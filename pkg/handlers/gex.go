package handlers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	gexop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/gex"
)

// SendGexRequestHandler sends a request to GEX
type SendGexRequestHandler HandlerContext

// Handle sends a request to GEX
func (h SendGexRequestHandler) Handle(params gexop.SendGexRequestParams) middleware.Responder {
	transactionName := *params.SendGexRequestPayload.TransactionName
	transactionBody := *params.SendGexRequestPayload.TransactionBody

	request, err := http.NewRequest(
		"POST",
		"https://gexweba.daas.dla.mil/msg_data/submit/"+transactionName,
		strings.NewReader(transactionBody),
	)
	if err != nil {
		h.logger.Error("Creating GEX POST request", zap.Error(err))
		return gexop.NewSendGexRequestInternalServerError()
	}

	cert := os.Getenv("CLIENT_TLS_CERT")
	key := os.Getenv("CLIENT_TLS_KEY")
	certificate, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		h.logger.Error("Creating client certificate", zap.Error(err))
		return gexop.NewSendGexRequestInternalServerError()
	}
	config := tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAnyClientCert,
		RootCAs:      getDoDRootCAs(),
	}
	tr := &http.Transport{TLSClientConfig: &config}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		h.logger.Error("Sending GEX POST request", zap.Error(err))
		return gexop.NewSendGexRequestInternalServerError()
	}

	fmt.Println("Server response: " + resp.Status)
	return gexop.NewSendGexRequestOK()
}

func getDoDRootCAs() *x509.CertPool {
	pool := x509.NewCertPool()
	cas := os.Getenv("GEX_DOD_CA")
	pool.AppendCertsFromPEM([]byte(cas))
	return pool
}
