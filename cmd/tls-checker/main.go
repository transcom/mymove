package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type stringSlice []string

func (s stringSlice) Contains(value string) bool {
	for _, x := range s {
		if value == x {
			return true
		}
	}
	return false
}

type intSlice []int

func (s intSlice) Contains(value int) bool {
	for _, x := range s {
		if value == x {
			return true
		}
	}
	return false
}

type errInvalidScheme struct {
	Scheme string
}

func (e *errInvalidScheme) Error() string {
	return "invalid scheme " + e.Scheme + ", only http and https are supported"
}

type errInvalidPath struct {
	Path string
}

func (e *errInvalidPath) Error() string {
	return "invalid path " + e.Path
}

type errTLSCheck struct {
	URL        string
	TLSVersion uint16
}

func (e *errTLSCheck) Error() string {
	tlsName := getTLSName(e.TLSVersion)
	return fmt.Sprintf("invalid request to url %s connected using %s", e.URL, tlsName)
}

func checkConfig(v *viper.Viper) error {
	schemesString := strings.TrimSpace(v.GetString("schemes"))

	if len(schemesString) == 0 {
		return errors.New("missing schemes")
	}

	schemes := stringSlice(strings.Split(schemesString, ","))

	for _, scheme := range schemes {
		if scheme != "http" && scheme != "https" {
			return &errInvalidScheme{Scheme: scheme}
		}
	}

	hosts := v.GetString("hosts")

	if len(hosts) == 0 {
		return errors.New("missing hosts")
	}

	pathsString := v.GetString("paths")

	if len(pathsString) == 0 {
		return errors.New("missing paths")
	}

	paths := stringSlice(strings.Split(pathsString, ","))

	for _, path := range paths {
		if !strings.HasPrefix(path, "/") {
			return &errInvalidPath{Path: path}
		}
	}

	clientKeyEncoded := v.GetString("key")
	clientCertEncoded := v.GetString("cert")
	clientKeyFile := v.GetString("key-file")
	clientCertFile := v.GetString("cert-file")

	if len(clientKeyEncoded) > 0 || len(clientCertEncoded) > 0 || len(clientKeyFile) > 0 || len(clientCertFile) > 0 {
		if schemes.Contains("http") {
			return errors.New("cannot use scheme http with client certificate, can only use https")
		}
	}

	return nil
}

func createTLSConfig(clientKey []byte, clientCert []byte, ca []byte, insecureSkipVerify bool, tlsVersion uint16) (*tls.Config, error) {

	keyPair, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	// #nosec b/c gosec triggers on InsecureSkipVerify
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{keyPair},
		InsecureSkipVerify: insecureSkipVerify,
		MinVersion:         tlsVersion,
		MaxVersion:         tlsVersion,
	}

	if len(ca) > 0 {
		rootCAs := x509.NewCertPool()
		rootCAs.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = rootCAs
	}

	return tlsConfig, nil
}

func createHTTPClient(v *viper.Viper, logger *zap.Logger, tlsVersion uint16) (*http.Client, error) {

	verbose := v.GetBool("verbose")

	clientKeyEncoded := v.GetString("key")
	clientCertEncoded := v.GetString("cert")
	skipVerify := v.GetBool("skip-verify")
	timeout := v.GetDuration("timeout")

	if verbose {
		if skipVerify {
			logger.Info("Skipping client-side certificate validation")
		}
	}

	// Supported TLS versions
	tlsConfig := &tls.Config{
		MinVersion: tlsVersion,
		MaxVersion: tlsVersion,
	}

	if len(clientKeyEncoded) > 0 && len(clientCertEncoded) > 0 {

		clientKey, clientKeyErr := base64.StdEncoding.DecodeString(clientKeyEncoded)
		if clientKeyErr != nil {
			return nil, errors.Wrap(clientKeyErr, "error decoding client key")
		}

		clientCert, clientCertErr := base64.StdEncoding.DecodeString(clientCertEncoded)
		if clientCertErr != nil {
			return nil, errors.Wrap(clientCertErr, "error decoding client cert")
		}

		caBytes := make([]byte, 0)
		if caEncoded := v.GetString("ca"); len(caEncoded) > 0 {
			caString, err := base64.StdEncoding.DecodeString(caEncoded)
			if err != nil {
				return nil, errors.Wrap(err, "error decoding certificate authority")
			}
			caBytes = []byte(caString)
		}

		var tlsConfigErr error
		tlsConfig, tlsConfigErr = createTLSConfig([]byte(clientKey), []byte(clientCert), caBytes, false, tlsVersion)
		if tlsConfigErr != nil {
			return nil, errors.Wrap(tlsConfigErr, "error creating TLS config")
		}

	} else {

		clientKeyFile := v.GetString("key-file")
		clientCertFile := v.GetString("cert-file")

		if len(clientKeyFile) > 0 && len(clientCertFile) > 0 {

			clientKey, clientKeyErr := ioutil.ReadFile(clientKeyFile) // #nosec b/c we need to read a file from a user-defined path
			if clientKeyErr != nil {
				return nil, errors.Wrap(clientKeyErr, "error reading client key file at "+clientKeyFile)
			}

			clientCert, clientCertErr := ioutil.ReadFile(clientCertFile) // #nosec b/c we need to read a file from a user-defined path
			if clientCertErr != nil {
				return nil, errors.Wrap(clientCertErr, "error reading client cert file at "+clientKeyFile)
			}

			caBytes := make([]byte, 0)
			if caFile := v.GetString("ca-file"); len(caFile) > 0 {
				content, err := ioutil.ReadFile(caFile) // #nosec b/c we need to read a file from a user-defined path
				if err != nil {
					return nil, errors.Wrap(err, "error reading ca file at "+caFile)
				}
				caBytes = content
			}
			var tlsConfigErr error
			tlsConfig, tlsConfigErr = createTLSConfig(clientKey, clientCert, caBytes, false, tlsVersion)
			if tlsConfigErr != nil {
				return nil, errors.Wrap(tlsConfigErr, "error creating TLS config")
			}
		}
	}

	httpTransport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: httpTransport,
	}
	return httpClient, nil

}

func checkURLWillNotConnect(httpClient *http.Client, url string, logger *zap.Logger) error {
	resp, err := httpClient.Get(url)
	if err == nil {
		return &errTLSCheck{URL: url, TLSVersion: resp.TLS.Version}
	}
	return nil
}

func createLogger(env string, level string) (*zap.Logger, error) {
	loglevel := zapcore.Level(uint8(0))
	err := (&loglevel).UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(loglevel)
	var loggerConfig zap.Config
	if env == "production" || env == "prod" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.Level = atomicLevel
	loggerConfig.DisableStacktrace = true
	return loggerConfig.Build(zap.AddStacktrace(zap.ErrorLevel))
}

func getTLSName(tlsVersion uint16) string {
	var tlsName string
	switch tlsVersion {
	case tls.VersionTLS10:
		tlsName = "TLS v1.0"
	case tls.VersionTLS11:
		tlsName = "TLS v1.1"
	case tls.VersionTLS12:
		tlsName = "TLS v1.2"
	case tls.VersionTLS13:
		tlsName = "TLS v1.3"
	}
	return tlsName
}

func main() {

	flag := pflag.CommandLine

	flag.StringP("schemes", "s", "https", "slice of schemes to check")
	flag.String("hosts", "", "comma-separated list of host names to check")
	flag.StringP("paths", "p", "/health", "slice of paths to check on each host")
	flag.String("key", "", "path to file of base64-encoded private key for client TLS")
	flag.String("key-file", "", "path to file of base64-encoded private key for client TLS")
	flag.String("cert", "", "base64-encoded public key for client TLS")
	flag.String("cert-file", "", "path to file of base64-encoded public key for client TLS")
	flag.String("ca", "", "base64-encoded certificate authority for mutual TLS")
	flag.String("ca-file", "", "path to file of base64-encoded certificate authority for mutual TLS")
	flag.Bool("skip-verify", false, "skip certifiate validation")
	flag.Duration("timeout", 5*time.Minute, "timeout duration")
	flag.Bool("exit-on-error", false, "exit on first tls check error")
	flag.String("log-env", "development", "logging config: development or production")
	flag.String("log-level", "error", "log level: debug, info, warn, error, dpanic, panic, or fatal")
	flag.Bool("verbose", false, "output extra information")

	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvPrefix("TLSCHECKER")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	tlsCheckErrors := make([]error, 0)
	defer func() {
		if len(tlsCheckErrors) > 0 {
			os.Exit(1)
		}
	}()

	logger, err := createLogger(v.GetString("log-env"), v.GetString("log-level"))
	if err != nil {
		log.Fatal(err.Error())
	}

	defer logger.Sync()

	err = checkConfig(v)
	if err != nil {
		switch e := err.(type) {
		case *errInvalidPath:
			logger.Fatal(e.Error(), zap.String("path", e.Path))
		case *errInvalidScheme:
			logger.Fatal(e.Error(), zap.String("scheme", e.Scheme))
		}
		logger.Fatal(err.Error())
	}

	verbose := v.GetBool("verbose")
	schemes := strings.Split(strings.TrimSpace(v.GetString("schemes")), ",")
	hosts := strings.Split(strings.TrimSpace(v.GetString("hosts")), ",")
	paths := strings.Split(strings.TrimSpace(v.GetString("paths")), ",")

	// TLS Versions that should not work
	var invalidTLSVersions = []uint16{
		tls.VersionTLS10,
		tls.VersionTLS11,
		// For Testing use these values
		// tls.VersionTLS12,
		// tls.VersionTLS13,
	}

	for _, tlsVersion := range invalidTLSVersions {

		tlsName := getTLSName(tlsVersion)

		httpClient, err := createHTTPClient(v, logger, tlsVersion)
		if err != nil {
			logger.Fatal(errors.Wrap(err, "error creating http client").Error())
		}

		exitOnError := v.GetBool("exit-on-error")

		for _, scheme := range schemes {
			for _, host := range hosts {
				for _, path := range paths {
					url := scheme + "://" + host + path
					if verbose {
						logger.Info("checking url will not connect with invalid TLS", zap.String("url", url), zap.String("tlsVersion", tlsName))
					}
					err := checkURLWillNotConnect(httpClient, url, logger)
					if err != nil {
						if exitOnError {
							logger.Fatal(err.Error())
						} else {
							logger.Warn(err.Error())
							tlsCheckErrors = append(tlsCheckErrors, err)
						}
					}
				}
			}
		}
	}
}
