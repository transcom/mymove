package gex

import (
	"crypto/tls"
	"crypto/x509"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

// SendInvoiceToGex sends an edi file string as a POST to the gex api
func SendInvoiceToGex(logger *zap.Logger, edi string, transactionName string) (status int, err error) {
	// Ensure that the transaction body ends with a newline, otherwise the GEX
	// EDI parser will fail silently
	edi = strings.TrimSpace(edi) + "\n"
	request, err := http.NewRequest(
		"POST",
		"https://gexweba.daas.dla.mil/msg_data/submit/"+transactionName,
		strings.NewReader(edi),
	)
	if err != nil {
		logger.Error("Creating GEX POST request", zap.Error(err))
		return 0, err
	}

	// We need to provide basic auth credentials for the GEX server, as well as
	// our client certificate for the proxy in front of the GEX server.
	request.SetBasicAuth(os.Getenv("GEX_BASIC_AUTH_USERNAME"), os.Getenv("GEX_BASIC_AUTH_PASSWORD"))

	config, err := GetTLSConfig()
	if err != nil {
		logger.Error("Creating TLS config", zap.Error(err))
		return 0, err
	}

	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		logger.Error("Sending GEX POST request", zap.Error(err))
		return 0, err
	}

	return resp.StatusCode, err
}

// GetTLSConfig gets the configuration certs for the GEX connection
func GetTLSConfig() (*tls.Config, error) {
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
