package server

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/transcom/mymove/pkg/storage"
	"golang.org/x/crypto/ocsp"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	idleTimeout       = 120 * time.Second // 2 minutes
	readHeaderTimeout = 60 * time.Second  // 1 minute
	maxHeaderSize     = 1 * 1000 * 1000   // 1 Megabyte
)

// the contextKey is typed so as not to conflict between similar keys from different pkgs
type contextKey string

var namedServerContextKey = contextKey("named_server")

// ErrMissingCACert represents an error caused by server config that requires
// certificate verification, but is missing a CA certificate
var ErrMissingCACert = errors.New("missing required CA certificate")

var cipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}

var curvePreferences = []tls.CurveID{
	tls.CurveP256,
	tls.X25519,
}

// CreateNamedServerInput contains the input for the CreateServer function.
type CreateNamedServerInput struct {
	Name                  string
	Host                  string
	Port                  int
	Logger                *zap.Logger
	HTTPHandler           http.Handler
	ClientAuth            tls.ClientAuthType
	Certificates          []tls.Certificate
	ClientCAs             *x509.CertPool // CaCertPool
	VerifyPeerCertificate func()
}

// NamedServer wraps *http.Server to override the definition of ListenAndServeTLS, but bypasses some restrictions.
type NamedServer struct {
	*http.Server
	Name               string
	IsServerReady      bool
	IsServerReadyMutex sync.Mutex
}

// Port returns the port the server binds to.  Returns -1 if any error.
func (s *NamedServer) Port() int {
	if !strings.Contains(s.Addr, ":") {
		return -1
	}
	port, err := strconv.Atoi(strings.SplitN(s.Addr, ":", 2)[1])
	if err != nil {
		return -1
	}
	return port
}

// ListenAndServeTLS is similar to (*http.Server).ListenAndServeTLS, but bypasses some restrictions.
func (s *NamedServer) ListenAndServeTLS() error {
	listener, err := tls.Listen("tcp", s.Addr, s.TLSConfig)
	if err != nil {
		return err
	}
	s.IsServerReadyMutex.Lock()
	s.IsServerReady = true
	s.IsServerReadyMutex.Unlock()
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			fmt.Println(fmt.Errorf("Failed to close listener due to %w", closeErr))
		}
	}()
	return s.Serve(listener)
}

// IsReady returns if a server is ready
func (s *NamedServer) IsReady() bool {
	s.IsServerReadyMutex.Lock()
	defer s.IsServerReadyMutex.Unlock()
	return s.IsServerReady
}

// WaitUntilReady waits until the server is ready
func (s *NamedServer) WaitUntilReady() {
	times := 0
	// Wait for server to be ready
	for !s.IsReady() && times < 4 {
		times++
		time.Sleep(500 * time.Millisecond)
	}
}

// Create client cert using TLS
func getClientCert(request *http.Request) *x509.Certificate {
	clientCert := request.TLS.PeerCertificates[0]
	return clientCert
}

type Fetcher struct {
	Fetch storage.FileStorer
}

func NewFetcher(fetch storage.FileStorer) (*Fetcher, error) {
	return &Fetcher{
		Fetch: fetch,
	}, nil
}

func transformCommonName(input string) string {
	//Remove spaces from common name
	noSpaces := strings.ReplaceAll(input, " ", "")

	//Convert dashes to underscores
	noDashes := strings.ReplaceAll(noSpaces, "-", "_")

	//Capitalize all letters
	capitalize := strings.ToUpper(noDashes)

	return capitalize
}
func fetchCRL(url string) (*x509.RevocationList, error) {
	httpResponse, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= 300 {
		return nil, errors.New("failed to retrieve CRL")
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	return x509.ParseRevocationList(body)
}

// Request CRL response from server.
// Returns error if the server can't get the certificate
// func getCRLResponse(fetch storage.FileStorer, v *viper.Viper, crlFile string, clientCert *x509.Certificate, issuerCert *x509.Certificate) error {
func getCRLResponse(clientCert *x509.Certificate, issuerCert *x509.Certificate) error {
	//Get the name of the Issuer (common name) from the client cert.
	//getIsserName := clientCert.Issuer.CommonName
	//
	//transformedIssuerName := transformCommonName(getIsserName)

	// NEW WORK:
	// Grab the common name, take out spaces, convert dashes to underscores, and capitalize
	// Once you get the filename you can pass in the path to the CRL and open that file

	//bucketName := "bucket_name" // This is the bucket name I am getting from Infra
	//folderPath := "path/to/folder"
	//fileStorer := storage.InitStorage(v, awsSession, appCtx.Logger())
	//actualNameOfS3Bucket := v.GetString(bucketName) // TODO: Is this something I actually need?
	//
	//// Build_URL
	//crlFilePath := path.Join(folderPath, transformedIssuerName) + ".crl"

	for _, url := range clientCert.CRLDistributionPoints {
		// TODO: Skip LDAP

		//x509.ParseRevocationList is not a direct ASN.1 representation, so leaves the option to add more detailed information
		parseCRL, err := fetchCRL(url)
		if err != nil {
			return err
		}

		//Parsed CRL against the issuer certificate
		err = parseCRL.CheckSignatureFrom(issuerCert)
		if err != nil {
			return err
		}

		// Check that the revocation list can be trusted
		if parseCRL.NextUpdate.Before(time.Now()) {
			return fmt.Errorf("CRL expired")
		}

		// Check id cert shows up in Revoked List
		for _, revokedCertificate := range parseCRL.RevokedCertificates {
			fmt.Printf("Revoked certificate serial number: %s\n", revokedCertificate.SerialNumber.String())
			if revokedCertificate.SerialNumber.Cmp(clientCert.SerialNumber) == 0 {
				return fmt.Errorf("The certificate is revoked!")
			}
		}
	}

	return nil
}

// Request OCSP response from server.
// Returns error if the server can't get the certificate
func getOCSPResponse(ocspServer string, request *http.Request, issuerCert *x509.Certificate) (*ocsp.Response, error) {
	var ocspRead = io.ReadAll

	var httpClient = &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	clientCert := getClientCert(request)

	// ocspRequestOpts is the hash function used in the request. We are using SHA256 instead of the default.
	ocspRequestOpts := &ocsp.RequestOptions{Hash: crypto.SHA256}

	// buffer contains the serialized request that will be sent to the server.
	buffer, err := ocsp.CreateRequest(clientCert, issuerCert, ocspRequestOpts)
	if err != nil {
		return nil, err
	}
	// HTTP requests must be made with TLS, and since client certs uses TLS this satisfies that requirement
	httpRequest, err := http.NewRequest(http.MethodPost, ocspServer, bytes.NewBuffer(buffer))
	if err != nil {
		return nil, err
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	body, err := ocspRead(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	ocspResponse, err := ocsp.ParseResponseForCert(body, clientCert, issuerCert)
	//return ocsp.ParseResponseForCert(body, leafCertificate, issuerCert)
	return ocspResponse, err
}

// If this callback returns nil, then the handshake continues and will not be aborted
// rawCerts contain chains of certificates in raw ASN.1 format
// each raw certs starts with the leafCert and ends with a root self-signed CA certificate
// verifiedChains have a certificate chain that verifies the signature validity and ends with a trusted certificate in the chain
func certRevokedCheck(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	var req *http.Request
	//var v *viper.Viper
	cert := verifiedChains[0][0]       // first argument verifies the client cert, second index 0 is the client cert
	issuerCert := verifiedChains[0][1] // second index of 1 is the issuer of the cert
	ocspResponse, err := getOCSPResponse(cert.OCSPServer[0], req, issuerCert)

	if err != nil {
		//return err // the revocation list was not checked and an error was encountered.
		return getCRLResponse(cert, issuerCert)
	}
	switch ocspResponse.Status {
	case ocsp.Good:
		fmt.Printf("[+] Certificate status: Good. It is still valid\n")
	case ocsp.Revoked:
		fmt.Printf("[!] Certificate status: Revoked.\n")
		return fmt.Errorf("The certificate was revoked!  The application can not trust the certificate.")
	case ocsp.Unknown:
		fmt.Printf("[?] Certificate status: Unknown\n")
		return fmt.Errorf("The certificate is unknown to OCSP server! The server does not know about the existence of the certificate serial number.")
	}

	fmt.Printf("Server certificate was allowed\n")
	return nil
}

// CreateNamedServer returns a no-tls, tls, or mutual-tls Server based on the input given and an error, if any.
func CreateNamedServer(input *CreateNamedServerInput) (*NamedServer, error) {

	address := fmt.Sprintf("%s:%d", input.Host, input.Port)

	var tlsConfig *tls.Config
	if len(input.Certificates) > 0 {

		if input.ClientAuth == tls.VerifyClientCertIfGiven || input.ClientAuth == tls.RequireAndVerifyClientCert {
			// RA Summary: staticcheck - SA1019 - Using a deprecated function, variable, constant or field
			// RA: Linter is flagging: input.ClientCAs.Subjects is deprecated: if s was returned by SystemCertPool, Subjects will not include the system roots.
			// RA: Why code valuable: It allows us to ensure we error if missing expected client certs.
			// RA: Mitigation: The deprecation notes this is a problem when reading SystemCertPool, but we do not use this here and are building up our own cert pool instead.
			// RA Developer Status: Mitigated
			// RA Validator Status: Mitigated
			// RA Validator: leodis.f.scott.civ@mail.mil
			// RA Modified Severity: CAT III
			// nolint:staticcheck
			if input.ClientCAs == nil || len(input.ClientCAs.Subjects()) == 0 {
				return nil, ErrMissingCACert
			}
		}

		// Follow Mozilla's "modern" server side TLS recommendations
		// https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
		// https://statics.tls.security.mozilla.org/server-side-tls-conf-4.0.json
		// This configuration is compatible with Firefox 27, Chrome 30, IE 11 on
		// Windows 7, Edge, Opera 17, Safari 9, Android 5.0, and Java 8
		tlsConfig = &tls.Config{
			CipherSuites:             cipherSuites,
			Certificates:             input.Certificates,
			ClientAuth:               input.ClientAuth,
			ClientCAs:                input.ClientCAs,
			CurvePreferences:         curvePreferences,
			MinVersion:               tls.VersionTLS12,
			NextProtos:               []string{"h2"},
			PreferServerCipherSuites: true,
			VerifyPeerCertificate:    certRevokedCheck,
		}
	}

	// wrappedHandler includes the name of the server in the context
	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), namedServerContextKey, input.Name)
		input.HTTPHandler.ServeHTTP(w, r.WithContext(ctx))
	})

	srv := &NamedServer{
		Name: input.Name,
		Server: &http.Server{
			Addr:              address,
			ErrorLog:          newStandardLogger(input.Logger),
			Handler:           wrappedHandler,
			IdleTimeout:       idleTimeout,
			MaxHeaderBytes:    maxHeaderSize,
			ReadHeaderTimeout: readHeaderTimeout,
			TLSConfig:         tlsConfig,
		},
	}
	return srv, nil

}

// NamedServerFromContext returns name name of the server that was previously added into the context, if any.
func NamedServerFromContext(ctx context.Context) string {
	name, ok := ctx.Value(namedServerContextKey).(string)
	if !ok {
		return ""
	}
	return name
}
