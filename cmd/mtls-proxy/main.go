package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pault.ag/go/pksigner"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	proxyAddrFlag   = "addr"
	tlsInsecureFlag = "insecure"
	// CertPathFlag is the path to the certificate to use for the CA
	caCertPathFlag string = "cacertpath"
	// KeyPathFlag is the path to the key to use for the CA
	caKeyPathFlag string = "cakeypath"
	// CertPathFlag is the path to the certificate to use for TLS
	certPathFlag string = "certpath"
	// KeyPathFlag is the path to the key to use for TLS
	keyPathFlag   string = "keypath"
	primeHostFlag string = "primehost"
	primePortFlag string = "primeport"
)

func initFlags(flags *pflag.FlagSet) {
	cli.InitCACFlags(flags)
	flags.String(proxyAddrFlag, ":8080", "Proxy address")
	flags.String(caCertPathFlag, "./config/tls/devlocal-ca.pem", "Path to the CA cert")
	flags.String(caKeyPathFlag, "./config/tls/devlocal-ca.key", "Path to the CA key")
	flags.String(certPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flags.String(keyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flags.Bool(tlsInsecureFlag, false, "Insecure")
	flags.String(primeHostFlag, "primelocal", "Hostname of prime server")
	flags.Int(primePortFlag, 9443, "Port of prime server")
}

// Stolen from go's reverseproxy.go
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func serveFunction(cmd *cobra.Command, args []string) error {
	v := viper.New()
	err := cmd.ParseFlags(args)
	if err != nil {
		return err
	}
	flags := cmd.Flags()
	errBindPFlags := v.BindPFlags(flags)
	if errBindPFlags != nil {
		return fmt.Errorf("Could not bind flags: %w", errBindPFlags)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err = cli.CheckCAC(v)
	if err != nil {
		return err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil
	}
	logger.Info("Starting to parse proxy options")

	insecure := v.GetBool(tlsInsecureFlag)

	var httpTransport *http.Transport
	if v.GetBool(cli.CACFlag) {
		var errGetCACStore error
		var store *pksigner.Store
		store, errGetCACStore = cli.GetCACStore(v)
		if errGetCACStore != nil {
			return errGetCACStore
		}
		cert, err := store.TLSCertificate()
		if err != nil {
			return err
		}

		// must explicitly state what signature algorithms we allow as of Go 1.14 to disable RSA-PSS signatures
		cert.SupportedSignatureAlgorithms = []tls.SignatureScheme{tls.PKCS1WithSHA256}

		// #nosec b/c gosec triggers on InsecureSkipVerify
		tlsConfig := tls.Config{
			Certificates:       []tls.Certificate{*cert},
			InsecureSkipVerify: insecure,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
		}
		tlsConfig.BuildNameToCertificate()
		httpTransport = &http.Transport{
			TLSClientConfig: &tlsConfig,
		}

	} else {
		certPath := v.GetString(certPathFlag)
		keyPath := v.GetString(keyPathFlag)
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		caCertPath := v.GetString(caCertPathFlag)
		caKeyPath := v.GetString(caKeyPathFlag)
		cacert, err := tls.LoadX509KeyPair(caCertPath, caKeyPath)
		if err != nil {
			return err
		}
		cacert.Leaf, err = x509.ParseCertificate(cacert.Certificate[0])
		if err != nil {
			return err
		}
		certpool := x509.NewCertPool()
		tlsConfig := tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: insecure,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
			RootCAs:            certpool,
		}
		tlsConfig.BuildNameToCertificate()
		httpTransport = &http.Transport{
			TLSClientConfig: &tlsConfig,
		}
	}

	primeHost := v.GetString(primeHostFlag)
	primePort := v.GetInt(primePortFlag)
	primeHostPort := fmt.Sprintf("https://%s:%d", primeHost, primePort)
	primeUrl, err := url.Parse(primeHostPort)
	if err != nil {
		return err
	}

	targetQuery := primeUrl.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = primeUrl.Scheme
		req.URL.Host = primeUrl.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(primeUrl, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	mtlsReverseProxy := httputil.ReverseProxy{
		Director:  director,
		Transport: httpTransport,
	}
	addr := v.GetString(proxyAddrFlag)
	logger.Info("Starting proxy", zap.Any("addr", addr))
	err = http.ListenAndServe(addr, &mtlsReverseProxy)
	logger.Info("Finished proxy", zap.Error(err))
	return nil
}

func main() {
	root := cobra.Command{
		Use:   "mtls-proxy [flags]",
		Short: "MTLS Proxy for MilMove",
		Long:  "MTLS Proxy for MilMove",
	}

	serveCommand := &cobra.Command{
		Use:          "serve",
		Short:        "Runs MTLS proxy",
		Long:         "Runs MTLS proxy",
		RunE:         serveFunction,
		SilenceUsage: true,
	}

	flags := serveCommand.Flags()
	initFlags(flags)
	root.AddCommand(serveCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
