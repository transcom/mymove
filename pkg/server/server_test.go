package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type serverSuite struct {
	logger *zap.Logger
	testingsuite.BaseTestSuite
	httpHandler http.Handler
}

func TestServerSuite(t *testing.T) {
	var httpHandler http.Handler
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	ss := &serverSuite{
		logger:      logger,
		httpHandler: httpHandler,
	}
	suite.Run(t, ss)
}

func (suite *serverSuite) readFile(filename string) []byte {
	testDataDir := "testdata"
	filePath := strings.Join([]string{testDataDir, filename}, "/")

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		suite.T().Fatalf("failed to read file %s: %s", filename, err)
	}
	return contents

}

func (suite *serverSuite) TestParseSingleTLSCert() {

	keyPair, err := tls.X509KeyPair(
		suite.readFile("localhost.pem"),
		suite.readFile("localhost.key"))

	suite.Nil(err)

	httpsServer := Server{
		ClientAuthType: tls.NoClientCert,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           8443,
		TLSCerts:       []tls.Certificate{keyPair},
	}

	tlsConfig, err := httpsServer.tlsConfig()
	suite.Nil(err)
	suite.Equal(len(tlsConfig.Certificates), 1)
	suite.Contains(tlsConfig.NameToCertificate, "localhost")
}

func (suite *serverSuite) TestParseBadTLSCert() {

	_, err := tls.X509KeyPair(
		suite.readFile("localhost-bad.pem"),
		suite.readFile("localhost.key"))

	suite.NotNil(err)
}

func (suite *serverSuite) TestParseMultipleTLSCerts() {

	keyPairLocalhost, err := tls.X509KeyPair(
		suite.readFile("localhost.pem"),
		suite.readFile("localhost.key"))

	suite.Nil(err)

	keyPairOffice, err := tls.X509KeyPair(
		suite.readFile("officelocal.pem"),
		suite.readFile("officelocal.key"))

	suite.Nil(err)

	httpsServer := Server{
		ClientAuthType: tls.NoClientCert,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           8443,
		TLSCerts: []tls.Certificate{
			keyPairLocalhost,
			keyPairOffice},
	}

	tlsConfig, err := httpsServer.tlsConfig()
	suite.Nil(err)
	suite.Equal(len(tlsConfig.Certificates), 2)
	suite.Contains(tlsConfig.NameToCertificate, "localhost")
	suite.Contains(tlsConfig.NameToCertificate, "officelocal")
}

func (suite *serverSuite) TestTLSConfigWithClientAuth() {

	keyPair, err := tls.X509KeyPair(
		suite.readFile("localhost.pem"),
		suite.readFile("localhost.key"))

	suite.Nil(err)

	caFile := suite.readFile("ca.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caFile)

	httpsServer := Server{
		ClientAuthType: tls.RequireAndVerifyClientCert,
		CaCertPool:     caCertPool,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           8443,
		TLSCerts:       []tls.Certificate{keyPair},
	}

	_, err = httpsServer.tlsConfig()
	suite.Nil(err)
}

func (suite *serverSuite) TestTLSConfigWithMissingCA() {

	keyPair, err := tls.X509KeyPair(
		suite.readFile("localhost.pem"),
		suite.readFile("localhost.key"))

	suite.Nil(err)

	httpsServer := Server{
		ClientAuthType: tls.RequireAndVerifyClientCert,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           8443,
		TLSCerts:       []tls.Certificate{keyPair},
	}

	_, err = httpsServer.tlsConfig()
	suite.Equal(ErrMissingCACert, err)
}

func (suite *serverSuite) TestTLSConfigWithMisconfiguredCA() {

	keyPair, err := tls.X509KeyPair(
		suite.readFile("localhost.pem"),
		suite.readFile("localhost.key"))

	suite.Nil(err)

	caFile := suite.readFile("localhost-bad.pem")
	caCertPool := x509.NewCertPool()
	certOk := caCertPool.AppendCertsFromPEM(caFile)
	suite.False(certOk)

	httpsServer := Server{
		ClientAuthType: tls.RequireAndVerifyClientCert,
		CaCertPool:     caCertPool,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           8443,
		TLSCerts:       []tls.Certificate{keyPair},
	}

	_, err = httpsServer.tlsConfig()
	suite.Equal(ErrMissingCACert, err)
}

func (suite *serverSuite) TestHTTPServerConfig() {
	var tlsConfig *tls.Config

	httpServer := Server{
		ListenAddress: "127.0.0.1",
		HTTPHandler:   suite.httpHandler,
		Logger:        suite.logger,
		Port:          8080,
	}

	config, err := httpServer.serverConfig(tlsConfig)

	suite.Nil(err)
	suite.Equal(config.Addr, "127.0.0.1:8080")
	suite.Equal(suite.httpHandler, config.Handler)
}
