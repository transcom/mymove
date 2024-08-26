package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/transcom/mymove/cmd/pptas-api-client/pptas"
	"github.com/transcom/mymove/cmd/pptas-api-client/utils"
	"github.com/transcom/mymove/pkg/cli"
)

// initRootFlags initializes flags relating to the prime api
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitLoggingFlags(flag)

	flag.String(utils.CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(utils.KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(utils.HostnameFlag, cli.HTTPPrimeServerNameLocal, "The hostname to connect to")
	flag.Int(utils.PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(utils.InsecureFlag, false, "Skip TLS verification and validation")
	flag.String(utils.FilenameFlag, "", "The name of the file being passed in")
	flag.String(utils.IDFlag, "", "The UUID of the object being retrieved or updated")
	flag.Duration(utils.WaitFlag, time.Second*80, "duration to wait for server to respond")
}

func main() {
	root := cobra.Command{
		Use:   "prime-api-client [flags]",
		Short: "Prime API client",
		Long:  "Prime API client",
	}
	initRootFlags(root.PersistentFlags())

	PPTASReportsCommand := &cobra.Command{
		Use:          "list-moves",
		Short:        "An optimized fetch for all moves available to Prime",
		Long:         "Fetches moves that are available to Prime quickly, without all the data for nested objects.",
		RunE:         pptas.PPTASReports,
		SilenceUsage: true,
	}
	pptas.InitPPTASReportsFlags(PPTASReportsCommand.Flags())
	root.AddCommand(PPTASReportsCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
