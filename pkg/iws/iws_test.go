package iws

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"testing"

	// This flag package accepts ENV vars as well as cmd line flags
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type iwsSuite struct {
	suite.Suite
	logger  *zap.Logger
	client  http.Client
	host    string
	custNum string
}

var (
	transport *http.Transport
)

func (suite *iwsSuite) SetupSuite() {
	var err error
	suite.logger, err = zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	host := flag.String("iws_rbs_host", "", "hostname of the IWS RBS environment")
	custNum := flag.String("iws_rbs_cust_num", "", "customer number to present when connecting to IWS RBS")
	moveMilDODCACert := flag.String("move_mil_dod_ca_cert", "", "The DoD CA certificate used to sign the move.mil TLS certificates.")
	moveMilDODTLSCert := flag.String("move_mil_dod_tls_cert", "", "the DoD signed tls certificate for various move.mil services.")
	moveMilDODTLSKey := flag.String("move_mil_dod_tls_key", "", "the DoD signed tls key for various move.mil services.")

	flag.Parse()

	suite.host = *host
	suite.custNum = *custNum

	// Load client cert
	cert, err := tls.X509KeyPair([]byte(*moveMilDODTLSCert), []byte(*moveMilDODTLSKey))
	if err != nil {
		log.Fatal(err)
	}

	// Load CA certs
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(*moveMilDODCACert))

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport = &http.Transport{TLSClientConfig: tlsConfig}
}

func (suite *iwsSuite) SetupTest() {
	// Fresh HTTP client for each test
	suite.client = http.Client{Transport: transport}
}

func TestIwsSuite(t *testing.T) {
	suite.Run(t, new(iwsSuite))
}
