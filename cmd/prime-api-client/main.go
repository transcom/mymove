package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"

	"github.com/spf13/cobra"
)

const (
	// CertPathFlag is the path to the certificate to use for TLS
	CertPathFlag string = "certpath"
	// KeyPathFlag is the path to the key to use for TLS
	KeyPathFlag string = "keypath"
	// HostnameFlag is the hostname to connect to
	HostnameFlag string = "hostname"
	// PortFlag is the port to connect to
	PortFlag string = "port"
	// InsecureFlag indicates that TLS verification and validation can be skipped
	InsecureFlag string = "insecure"
)

// initRootFlags initializes flags relating to the prime api
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitVerboseFlags(flag)

	flag.String(CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(HostnameFlag, cli.HTTPPrimeServerNameLocal, "The hostname to connect to")
	flag.Int(PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(InsecureFlag, false, "Skip TLS verification and validation")
}

// CheckRootConfig checks the validity of the prime api flags
func CheckRootConfig(v *viper.Viper) error {
	err := cli.CheckCAC(v)
	if err != nil {
		return err
	}

	err = cli.CheckVerbose(v)
	if err != nil {
		return err
	}

	if (v.GetString(CertPathFlag) != "" && v.GetString(KeyPathFlag) == "") || (v.GetString(CertPathFlag) == "" && v.GetString(KeyPathFlag) != "") {
		return fmt.Errorf("Both TLS certificate and key paths must be provided")
	}

	return nil
}

func main() {
	root := cobra.Command{
		Use:   "prime-api-client [flags]",
		Short: "Prime API client",
		Long:  "Prime API client",
	}
	initRootFlags(root.PersistentFlags())

	fetchMTOsCommand := &cobra.Command{
		Use:          "fetch-mtos",
		Short:        "fetch mtos",
		Long:         "fetch move task orders",
		RunE:         fetchMTOs,
		SilenceUsage: true,
	}
	root.AddCommand(fetchMTOsCommand)

	updateMTOShipmentCommand := &cobra.Command{
		Use:          "update-mto-shipment",
		Short:        "update mto shipment",
		Long:         "update move task order shipment",
		RunE:         updateMTOShipment,
		SilenceUsage: true,
	}
	initUpdateMTOShipmentFlags(updateMTOShipmentCommand.Flags())
	root.AddCommand(updateMTOShipmentCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\nprime-api-client completion > /usr/local/etc/bash_completion.d/prime-api-client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
