package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services/invoice"
)

func checkPostFileToGEXConfig(v *viper.Viper, logger logger) error {
	logger.Debug("checking config")

	err := cli.CheckGEX(v)
	if err != nil {
		return err
	}

	err = cli.CheckCert(v)
	if err != nil {
		return err
	}

	ediFile := v.GetString("edi")
	if len(ediFile) == 0 {
		return errors.Errorf("%s is missing", "edi")
	}

	trasaction := v.GetString("transaction-name")
	if len(trasaction) == 0 {
		return errors.Errorf("%s is missing", "transaction-name")
	}

	return nil
}

func initPostFileToGEXFlags(flag *pflag.FlagSet) {
	// Verbose
	cli.InitVerboseFlags(flag)

	// GEX
	cli.InitGEXFlags(flag)

	// Certificate
	cli.InitCertFlags(flag)

	flag.String("edi", "", "The filepath to an edi file to send to GEX")
	flag.String("transaction-name", "test", "The required name sent in the url of the gex api request")
	flag.Parse(os.Args[1:])

	// Don't sort flags
	flag.SortFlags = false
}

// go run ./cmd/milmove-tasks post-file-to-gex --edi filepath --transaction-name transactionName --gex-url 'url'
func postFileToGEX(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse args")
	}
	flags := cmd.Flags()
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkPostFileToGEXConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	ediFile := v.GetString("edi")

	file, err := os.Open(ediFile) // #nosec
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	edi, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	ediString := string(edi[:])
	// make sure edi ends in new line
	ediString = strings.TrimSpace(ediString) + "\n"

	fmt.Println(ediString)

	certificates, rootCAs, err := initDODCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		log.Fatal("Error in getting tls certs", err)
	}
	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs}
	url := v.GetString("gex-url")
	if len(url) == 0 {
		log.Fatal("Not sending to GEX because no URL set.")
	}

	resp, err := invoice.NewGexSenderHTTP(
		url,
		true,
		tlsConfig,
		v.GetString("gex-basic-auth-username"),
		v.GetString("gex-basic-auth-password"),
	).SendToGex(ediString, "")
	if resp == nil || err != nil {
		log.Fatal("Gex Sender had no response", err)
	}

	fmt.Println("Sending to GEX. . .")
	fmt.Printf("status code: %v, error: %v \n", resp.StatusCode, err)

	return nil
}

//TODO: Infra will work to refactor and reduce duplication (also found in cmd/milmove/main.go)
func initDODCertificates(v *viper.Viper, logger *zap.Logger) ([]tls.Certificate, *x509.CertPool, error) {

	tlsCert := v.GetString("move-mil-dod-tls-cert")
	if len(tlsCert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-tls-cert")
	}

	caCert := v.GetString("move-mil-dod-ca-cert")
	if len(caCert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-ca-cert")
	}

	//Append move.mil cert with CA certificate chain
	cert := bytes.Join(
		[][]byte{
			[]byte(tlsCert),
			[]byte(caCert),
		},
		[]byte("\n"),
	)

	key := []byte(v.GetString("move-mil-dod-tls-key"))
	if len(key) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-tls-key")
	}

	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, "failed to parse DOD keypair for server")
	}

	pathToPackage := v.GetString("dod-ca-package")
	if len(pathToPackage) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is missing", "dod-ca-package"))
	}

	pkcs7Package, err := ioutil.ReadFile(pathToPackage) // #nosec
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, fmt.Sprintf("%s is invalid", "dod-ca-package"))
	}

	if len(pkcs7Package) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is an empty file", "dod-ca-package"))
	}

	dodCACertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		return make([]tls.Certificate, 0), dodCACertPool, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	return []tls.Certificate{keyPair}, dodCACertPool, nil

}

//TODO: Infra will refactor to reduce duplication
type errInvalidPKCS7 struct {
	Path string
}

//TODO: Infra will refactor to reduce duplication
func (e *errInvalidPKCS7) Error() string {
	return fmt.Sprintf("invalid DER encoded PKCS7 package: %s", e.Path)
}
