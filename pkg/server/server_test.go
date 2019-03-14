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

	"github.com/transcom/mymove/pkg/testingsuite"
)

type serverSuite struct {
	logger Logger
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

	httpsServer, err := CreateNamedServer(&CreateNamedServerInput{
		Host:         "127.0.0.1",
		Port:         8443,
		ClientAuth:   tls.NoClientCert,
		HTTPHandler:  suite.httpHandler,
		Logger:       suite.logger,
		Certificates: []tls.Certificate{keyPair},
	})
	suite.Nil(err)
	suite.Equal(len(httpsServer.TLSConfig.Certificates), 1)
	suite.Contains(httpsServer.TLSConfig.NameToCertificate, "localhost")
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

	httpsServer, err := CreateNamedServer(&CreateNamedServerInput{
		Host:        "127.0.0.1",
		Port:        8443,
		ClientAuth:  tls.NoClientCert,
		HTTPHandler: suite.httpHandler,
		Logger:      suite.logger,
		Certificates: []tls.Certificate{
			keyPairLocalhost,
			keyPairOffice,
		},
	})
	suite.Nil(err)
	suite.Equal(len(httpsServer.TLSConfig.Certificates), 2)
	suite.Contains(httpsServer.TLSConfig.NameToCertificate, "localhost")
	suite.Contains(httpsServer.TLSConfig.NameToCertificate, "officelocal")
}

func (suite *serverSuite) TestTLSConfigWithClientAuth() {

	keyPair, err := tls.X509KeyPair(
		suite.readFile("localhost.pem"),
		suite.readFile("localhost.key"))

	suite.Nil(err)

	caFile := suite.readFile("ca.pem")
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caFile)

	_, err = CreateNamedServer(&CreateNamedServerInput{
		Host:         "127.0.0.1",
		Port:         8443,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		HTTPHandler:  suite.httpHandler,
		Logger:       suite.logger,
		Certificates: []tls.Certificate{keyPair},
	})
	suite.Nil(err)
}

func (suite *serverSuite) TestTLSConfigWithMissingCA() {

	keyPair, err := tls.X509KeyPair(
		suite.readFile("localhost.pem"),
		suite.readFile("localhost.key"))

	suite.Nil(err)

	_, err = CreateNamedServer(&CreateNamedServerInput{
		Host:         "127.0.0.1",
		Port:         8443,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		HTTPHandler:  suite.httpHandler,
		Logger:       suite.logger,
		Certificates: []tls.Certificate{keyPair},
	})
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

	_, err = CreateNamedServer(&CreateNamedServerInput{
		Host:         "127.0.0.1",
		Port:         8443,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		HTTPHandler:  suite.httpHandler,
		Logger:       suite.logger,
		Certificates: []tls.Certificate{keyPair},
	})
	suite.Equal(ErrMissingCACert, err)
}

func (suite *serverSuite) TestHTTPServerConfig() {
	httpsServer, err := CreateNamedServer(&CreateNamedServerInput{
		Host:        "127.0.0.1",
		Port:        8080,
		HTTPHandler: suite.httpHandler,
		Logger:      suite.logger,
	})
	suite.Nil(err)
	suite.Equal(httpsServer.Addr, "127.0.0.1:8080")
	suite.Equal(suite.httpHandler, httpsServer.Handler)
}
