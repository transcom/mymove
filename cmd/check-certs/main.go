package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// Certs
	cli.InitCertFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}
	certificates, rootCAs, err := cli.InitDoDCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	logger.Debug("Server DOD Key Pair Loaded")
	logger.Debug("Trusted Certificate Authorities", zap.Any("subjects", rootCAs.Subjects()))

	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs}
	conn, err := tls.Dial("tcp", "gexweba.daas.dla.mil:443", tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	for _, chain := range conn.ConnectionState().VerifiedChains {
		for certNum, cert := range chain {
			fmt.Println(certNum, cert.Subject, cert.NotAfter)
		}
	}
}
