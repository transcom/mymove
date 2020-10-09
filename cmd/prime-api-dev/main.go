package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/transcom/mymove/cmd/prime-api-dev/scripts"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	"github.com/transcom/mymove/pkg/cli"
)

// initRootFlags initializes flags relating to the prime api
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitVerboseFlags(flag)

	flag.String(utils.CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(utils.KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(utils.HostnameFlag, cli.HTTPPrimeServerNameLocal, "The hostname to connect to")
	flag.Int(utils.PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(utils.InsecureFlag, false, "Skip TLS verification and validation")
}

func main() {
	root := cobra.Command{
		Use:   "prime-api-dev [flags]",
		Short: "Prime API Dev Scripts",
		Long:  "Prime API Dev Scripts",
	}
	initRootFlags(root.PersistentFlags())

	paymentRequestsCommand := &cobra.Command{
		Use:          "pr",
		Short:        "Use Prime API Payment Request functions",
		Long:         "Use Prime API Payment Request functions",
		RunE:         scripts.PaymentRequests,
		SilenceUsage: true,
	}
	scripts.InitPaymentRequestsFlags(paymentRequestsCommand.Flags())
	root.AddCommand(paymentRequestsCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
