package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
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
	initFetchMTOsFlags(fetchMTOsCommand.Flags())
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

	createMTOServiceItemCommand := &cobra.Command{
		Use:          "create-mto-service-item",
		Short:        "Create mto service item",
		Long:         "Create move task order service item for move task order and/or shipment",
		RunE:         createMTOServiceItem,
		SilenceUsage: true,
	}
	initCreateMTOServiceItemFlags(createMTOServiceItemCommand.Flags())
	root.AddCommand(createMTOServiceItemCommand)

	makeAvailableToPrimeCommand := &cobra.Command{
		Use:          "support-make-mto-available-to-prime",
		Short:        "Make mto available to prime",
		Long:         "Makes an mto available to the prime for prime-api consumption",
		RunE:         updateMTOStatus,
		SilenceUsage: true,
	}
	initUpdateMTOStatusFlags(makeAvailableToPrimeCommand.Flags())
	root.AddCommand(makeAvailableToPrimeCommand)

	updatePaymentRequestStatusCommand := &cobra.Command{
		Use:          "support-update-payment-request-status",
		Short:        "Update payment request status for prime",
		Long:         "Allows prime to update payment request status in non-prod envs",
		RunE:         updatePaymentRequestStatus,
		SilenceUsage: true,
	}
	initUpdatePaymentRequestStatusFlags(updatePaymentRequestStatusCommand.Flags())
	root.AddCommand(updatePaymentRequestStatusCommand)

	getMoveTaskOrder := &cobra.Command{
		Use:          "support-get-mto",
		Short:        "Get an individual mto",
		Long:         "Get an individual mto's information",
		RunE:         getMTO,
		SilenceUsage: true,
	}
	initGetMTOFlags(getMoveTaskOrder.Flags())
	root.AddCommand(getMoveTaskOrder)

	updateMTOServiceItemStatus := &cobra.Command{
		Use:          "support-update-mto-service-item-status",
		Short:        "Update service item status",
		Long:         "Approve or reject a service item",
		RunE:         updateMTOServiceItemStatus,
		SilenceUsage: true,
	}
	initUpdateMTOServiceItemStatusFlags(updateMTOServiceItemStatus.Flags())
	root.AddCommand(updateMTOServiceItemStatus)

	createPaymentRequestCommand := &cobra.Command{
		Use:          "create-payment-request",
		Short:        "Create payment request",
		Long:         "Create payment request for a move task order",
		RunE:         createPaymentRequest,
		SilenceUsage: true,
	}
	initCreatePaymentRequestFlags(createPaymentRequestCommand.Flags())
	root.AddCommand(createPaymentRequestCommand)

	patchMTOShipmentStatusCommand := &cobra.Command{
		Use:          "support-patch-mto-shipment-status",
		Short:        "Update MTO shipment status for prime",
		Long:         "Allows prime to update MTO shipment status in non-prod envs",
		RunE:         patchMTOShipmentStatus,
		SilenceUsage: true,
	}
	initPatchMTOShipmentStatusFlags(patchMTOShipmentStatusCommand.Flags())
	root.AddCommand(patchMTOShipmentStatusCommand)

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
