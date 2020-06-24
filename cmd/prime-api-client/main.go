package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/transcom/mymove/cmd/prime-api-client/prime"
	"github.com/transcom/mymove/cmd/prime-api-client/support"
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
		Use:   "prime-api-client [flags]",
		Short: "Prime API client",
		Long:  "Prime API client",
	}
	initRootFlags(root.PersistentFlags())

	fetchMTOsCommand := &cobra.Command{
		Use:          "fetch-mto-updates",
		Short:        "Fetch all MTOs available to prime",
		Long:         "fetch move task orders",
		RunE:         prime.FetchMTOUpdates,
		SilenceUsage: true,
	}
	prime.InitFetchMTOUpdatesFlags(fetchMTOsCommand.Flags())
	root.AddCommand(fetchMTOsCommand)

	listMTOsCommand := &cobra.Command{
		Use:          "support-list-mtos",
		Short:        "Fetch all MTOs",
		Long:         "fetch all move task orders",
		RunE:         prime.FetchMTOUpdates,
		SilenceUsage: true,
	}
	support.InitListMTOsFlags(listMTOsCommand.Flags())
	root.AddCommand(listMTOsCommand)

	createMTOCommand := &cobra.Command{
		Use:   "support-create-move-task-order",
		Short: "Create a MoveTaskOrder",
		Long: `
  This command creates a MoveTaskOrder object.
  It requires the caller to pass in a file using the --filename param.

  Endpoint path: /move-task-orders
  The file should contain json as follows:
    {
      "body": <MoveTaskOrder>
    }
  Please see API documentation for full details on the MoveTaskOrder definition.`,
		RunE:         support.CreateMTO,
		SilenceUsage: true,
	}
	support.InitCreateMTOFlags(createMTOCommand.Flags())
	root.AddCommand(createMTOCommand)

	createMTOShipmentCommand := &cobra.Command{
		Use:   "create-mto-shipment",
		Short: "Create MTO shipment",
		Long: `
	This command creates a MTO shipment.
	It requires the caller to pass in a file using the --filename arg.
	The file should contain a body defining the MTOShipment object.
	Endpoint path: /mto-shipments
	The file should contain json as follows:
		{
			"body": <MTOShipment>,
		}
	Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.CreateMTOShipment,
		SilenceUsage: true,
	}
	prime.InitCreateMTOShipmentFlags(createMTOShipmentCommand.Flags())
	root.AddCommand(createMTOShipmentCommand)

	updateMTOShipmentCommand := &cobra.Command{
		Use:   "update-mto-shipment",
		Short: "Update MTO shipment",
		Long: `
  This command updates an MTO shipment.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters, headers and a body for the payload.

  Endpoint path: /mto-shipments/{mtoShipmentID}
  The file should contain json as follows:
  	{
      "mtoShipmentID": <uuid string>,
      "ifMatch": <eTag>,
      "body": <MTOShipment>
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.UpdateMTOShipment,
		SilenceUsage: true,
	}
	prime.InitUpdateMTOShipmentFlags(updateMTOShipmentCommand.Flags())
	root.AddCommand(updateMTOShipmentCommand)

	updatePostCounselingInfo := &cobra.Command{
		Use:          "update-mto-post-counseling-information",
		Short:        "update post counseling info",
		Long:         "Update post counseling info such as discovering that customer has a PPM",
		RunE:         prime.UpdatePostCounselingInfo,
		SilenceUsage: true,
	}
	prime.InitUpdatePostCounselingInfoFlags(updatePostCounselingInfo.Flags())
	root.AddCommand(updatePostCounselingInfo)

	createMTOServiceItemCommand := &cobra.Command{
		Use:   "create-mto-service-item",
		Short: "Create mto service item",
		Long: `
  This command creates an MTO service item on an MTO shipment.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters and headers and a body for the payload.

  Endpoint path: /mto-service-items
  The file should contain json as follows:
  	{
  	"body": <MTOServiceItem>
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.CreateMTOServiceItem,
		SilenceUsage: true,
	}
	prime.InitCreateMTOServiceItemFlags(createMTOServiceItemCommand.Flags())
	root.AddCommand(createMTOServiceItemCommand)

	makeAvailableToPrimeCommand := &cobra.Command{
		Use:   "support-make-move-task-order-available",
		Short: "Make MTO available to prime",
		Long: `
  This command makes an MTO available for prime consumption.
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters and headers.

  Endpoint path: /move-task-orders/{moveTaskOrderID}/available-to-prime
  The file should contain json as follows:
  	{
  	"moveTaskOrderID": <uuid string>,
  	"ifMatch": <eTag>
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.MakeMTOAvailable,
		SilenceUsage: true,
	}
	support.InitMakeMTOAvailableFlags(makeAvailableToPrimeCommand.Flags())
	root.AddCommand(makeAvailableToPrimeCommand)

	updatePaymentRequestStatusCommand := &cobra.Command{
		Use:   "support-update-payment-request-status",
		Short: "Update payment request status for prime",
		Long: `
  This command allows prime to update payment request status.
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters and headers.

  Endpoint path: /payment-requests/{paymentRequestID}/status
  The file should contain json as follows:
    {
      "paymentRequestID": <uuid string>,
      "ifMatch": <etag>,
      "body" : <paymentRequestStatus>
    }
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.UpdatePaymentRequestStatus,
		SilenceUsage: true,
	}
	support.InitUpdatePaymentRequestStatusFlags(updatePaymentRequestStatusCommand.Flags())
	root.AddCommand(updatePaymentRequestStatusCommand)

	getMoveTaskOrder := &cobra.Command{
		Use:   "support-get-move-task-order",
		Short: "Get an individual mto",
		Long: `
  This command gets a single move task order by ID
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters and headers.

  Endpoint path: /move-task-orders/{moveTaskOrderID}
  The file should contain json as follows:
  	{
  	"moveTaskOrderID": <uuid string>,
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.GetMTO,
		SilenceUsage: true,
	}
	support.InitGetMTOFlags(getMoveTaskOrder.Flags())
	root.AddCommand(getMoveTaskOrder)

	updateMTOServiceItemStatus := &cobra.Command{
		Use:   "support-update-mto-service-item-status",
		Short: "Update service item status",
		Long: `
  This command allows prime to update the MTO service item status.
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain a body defining the request body.

  Endpoint path: service-items/{mtoServiceItemID}/status
    {
      "mtoServiceItemID": <uuid string>,
      "ifMatch": <etag>,
      "body": {
        "status": "APPROVED"
    }
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.UpdateMTOServiceItemStatus,
		SilenceUsage: true,
	}
	support.InitUpdateMTOServiceItemStatusFlags(updateMTOServiceItemStatus.Flags())
	root.AddCommand(updateMTOServiceItemStatus)

	createPaymentRequestCommand := &cobra.Command{
		Use:   "create-payment-request",
		Short: "Create payment request",
		Long: `
  This command gets a single move task order by ID
  It requires the caller to pass in a file using the --filename arg.
  The file should contain a body defining the PaymentRequest object.
  Endpoint path: /payment-requests
  The file should contain json as follows:
  	{
  	"body": <PaymentRequest>,
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.CreatePaymentRequest,
		SilenceUsage: true,
	}
	prime.InitCreatePaymentRequestFlags(createPaymentRequestCommand.Flags())
	root.AddCommand(createPaymentRequestCommand)

	listMTOPaymentRequestsCommand := &cobra.Command{
		Use:   "support-list-mto-payment-requests",
		Short: "Get all payment requests for a given MTO",
		Long: `
  This command allows the user to get all payment requests associated with an MTO.
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters.

  Endpoint path: /move-task-orders/{moveTaskOrderID}/payment-requests
  The file should contain json as follows:
    {
      "moveTaskOrderID": <uuid string>,
    }
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.ListMTOPaymentRequests,
		SilenceUsage: true,
	}
	support.InitListMTOPaymentRequestsFlags(listMTOPaymentRequestsCommand.Flags())
	root.AddCommand(listMTOPaymentRequestsCommand)

	createPaymentRequestUploadCommand := &cobra.Command{
		Use:          "create-upload",
		Short:        "Create payment request upload",
		Long:         "Create payment request upload for a payment request",
		RunE:         prime.CreatePaymentRequestUpload,
		SilenceUsage: true,
	}
	prime.InitCreatePaymentRequestUploadFlags(createPaymentRequestUploadCommand.Flags())
	root.AddCommand(createPaymentRequestUploadCommand)

	updateMTOShipmentStatusCommand := &cobra.Command{
		Use:   "support-update-mto-shipment-status",
		Short: "Update MTO shipment status for prime",
		Long: `
  This command allows prime to update the MTO shipment status.
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain a body defining the request body.

  Endpoint path: /mto-shipments/{mtoShipmentID}/status
  The file should contain json as follows:
    {
      "mtoShipmentID": <uuid string>,
      "ifMatch": <etag>,
      "body": <MtoShipmentRequestStatus>,
    }
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.UpdateMTOShipmentStatus,
		SilenceUsage: true,
	}
	support.InitUpdateMTOShipmentStatusFlags(updateMTOShipmentStatusCommand.Flags())
	root.AddCommand(updateMTOShipmentStatusCommand)

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
