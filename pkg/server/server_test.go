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
	"go.uber.org/zap"
)

type serverSuite struct {
	logger *zap.Logger
	suite.Suite
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
	tlsCert := TLSCert{
		CertPEMBlock: suite.readFile("localhost.pem"),
		KeyPEMBlock:  suite.readFile("localhost.key"),
	}
	httpsServer := Server{
		ClientAuthType: tls.NoClientCert,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           "8443",
		TLSCerts:       []TLSCert{tlsCert},
	}

	tlsConfig, err := httpsServer.tlsConfig()
	suite.Nil(err)
	suite.Equal(len(tlsConfig.Certificates), 1)
	suite.Contains(tlsConfig.NameToCertificate, "localhost")
}

func (suite *serverSuite) TestParseBadTLSCert() {
	tlsCert1 := TLSCert{
		CertPEMBlock: suite.readFile("localhost-bad.pem"),
		KeyPEMBlock:  suite.readFile("localhost.key"),
	}

	httpsServer := Server{
		ClientAuthType: tls.NoClientCert,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           "8443",
		TLSCerts:       []TLSCert{tlsCert1},
	}

	tlsConfig, err := httpsServer.tlsConfig()
	suite.NotNil(err)
	suite.Nil(tlsConfig)
}

func (suite *serverSuite) TestParseMultipleTLSCerts() {
	tlsLocalhost := TLSCert{
		CertPEMBlock: suite.readFile("localhost.pem"),
		KeyPEMBlock:  suite.readFile("localhost.key"),
	}

	tlsOfficelocal := TLSCert{
		CertPEMBlock: suite.readFile("officelocal.pem"),
		KeyPEMBlock:  suite.readFile("officelocal.key"),
	}

	httpsServer := Server{
		ClientAuthType: tls.NoClientCert,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           "8443",
		TLSCerts: []TLSCert{
			tlsLocalhost,
			tlsOfficelocal},
	}

	tlsConfig, err := httpsServer.tlsConfig()
	suite.Nil(err)
	suite.Equal(len(tlsConfig.Certificates), 2)
	suite.Contains(tlsConfig.NameToCertificate, "localhost")
	suite.Contains(tlsConfig.NameToCertificate, "officelocal")
}

func (suite *serverSuite) TestTLSConfigWithClientAuth() {
	tlsCert := TLSCert{
		CertPEMBlock: suite.readFile("localhost.pem"),
		KeyPEMBlock:  suite.readFile("localhost.key"),
	}
	caFile := suite.readFile("ca.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caFile)

	httpsServer := Server{
		ClientAuthType: tls.RequireAndVerifyClientCert,
		CaCertPool:     caCertPool,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           "8443",
		TLSCerts:       []TLSCert{tlsCert},
	}

	_, err := httpsServer.tlsConfig()
	suite.Nil(err)
}

func (suite *serverSuite) TestTLSConfigWithMissingCA() {
	tlsCert := TLSCert{
		CertPEMBlock: suite.readFile("localhost.pem"),
		KeyPEMBlock:  suite.readFile("localhost.key"),
	}
	httpsServer := Server{
		ClientAuthType: tls.RequireAndVerifyClientCert,
		ListenAddress:  "127.0.0.1",
		HTTPHandler:    suite.httpHandler,
		Logger:         suite.logger,
		Port:           "8443",
		TLSCerts:       []TLSCert{tlsCert},
	}

	_, err := httpsServer.tlsConfig()
	suite.Equal(ErrMissingCACert, err)
}

func (suite *serverSuite) TestTLSConfigWithMisconfiguredCA() {
	tlsCert := TLSCert{
		CertPEMBlock: suite.readFile("localhost.pem"),
		KeyPEMBlock:  suite.readFile("localhost.key"),
	}
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
		Port:           "8443",
		TLSCerts:       []TLSCert{tlsCert},
	}

	_, err := httpsServer.tlsConfig()
	suite.Equal(ErrMissingCACert, err)
}

func (suite *serverSuite) TestHTTPServerConfig() {
	var tlsConfig *tls.Config

	httpServer := Server{
		ListenAddress: "127.0.0.1",
		HTTPHandler:   suite.httpHandler,
		Logger:        suite.logger,
		Port:          "8080",
	}

	config, err := httpServer.serverConfig(tlsConfig)

	suite.Nil(err)
	suite.Equal(config.Addr, "127.0.0.1:8080")
	suite.Equal(suite.httpHandler, config.Handler)
}
