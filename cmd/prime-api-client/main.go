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
	cli.InitLoggingFlags(flag)

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

	createWebhookCommand := &cobra.Command{
		Use:   "support-create-webhook-notification",
		Short: "Create a WebhookNotification",
		Long: `
  This command creates a WebhookNotification object.
  Passing in a file is optional, but when passed in a file the --filename param must be used.

  Endpoint path: /webhook-notifications
  The file should contain json as follows:
    {
      "body": <WebhookNotification>
    }
  Please see API documentation for full details on the WebhookNotification definition.`,
		RunE:         support.CreateWebhookNotification,
		SilenceUsage: true,
	}
	support.InitCreateWebhookNotificationFlags(createWebhookCommand.Flags())
	root.AddCommand(createWebhookCommand)

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

	updateMTOShipmentAddressCommand := &cobra.Command{
		Use:   "update-mto-shipment-address",
		Short: "Update MTO shipment address",
		Long: `
  This command updates an address associated with an MTO shipment.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters, headers and a body for the payload.

  Endpoint path: /mto-shipments/{mtoShipmentID}/addresses/{addressID}
  The file should contain json as follows:
  	{
	  "mtoShipmentID": <uuid string>,
	  "addressID": <uuid string>,
      "ifMatch": <eTag>,
      "body": <MTOShipmentAddress>
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.UpdateMTOShipmentAddress,
		SilenceUsage: true,
	}
	prime.InitUpdateMTOShipmentAddressFlags(updateMTOShipmentAddressCommand.Flags())
	root.AddCommand(updateMTOShipmentAddressCommand)

	updateMTOAgentCommand := &cobra.Command{
		Use:   "update-mto-agent",
		Short: "Update MTO agent",
		Long: `
  This command updates an agent associated with an MTO shipment.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters, headers and a body for the payload.

  Endpoint path: /mto-shipments/{mtoShipmentID}/agents/{agentID}
  The file should contain json as follows:
  	{
	  "mtoShipmentID": <uuid string>,
	  "agentID": <uuid string>,
      "ifMatch": <eTag>,
      "body": <MTOAgent>
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.UpdateMTOAgent,
		SilenceUsage: true,
	}
	prime.InitUpdateMTOAgentFlags(updateMTOAgentCommand.Flags())
	root.AddCommand(updateMTOAgentCommand)

	createMTOAgentCommand := &cobra.Command{
		Use:   "create-mto-agent",
		Short: "Create MTO agent",
		Long: `
  This command creates an agent associated with an MTO shipment.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters and a body for the payload.

  Endpoint path: /mto-shipments/{mtoShipmentID}/agents
  The file should contain json as follows:
  	{
	  "mtoShipmentID": <uuid string>,
      "body": <MTOAgent>
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.CreateMTOAgent,
		SilenceUsage: true,
	}
	prime.InitCreateMTOAgentFlags(createMTOAgentCommand.Flags())
	root.AddCommand(createMTOAgentCommand)

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

	updateMTOServiceItem := &cobra.Command{
		Use:   "update-mto-service-item",
		Short: "Update service item",
		Long: `
  	This command updates an MTO service item. It requires the caller to pass in
	a file using the --filename arg. The file should contain path parameters,
	headers and a body for the payload.

	Endpoint path: /mto-service-items/{mtoServiceItemID}
  	The file should contain json as follows:
 	  {
        "mtoServiceItemID": <uuid string>,
        "ifMatch": <etag>,
        "body" : <UpdateMTOServiceItem>
      }
  	Please see API documentation for full details on the endpoint definition.`,
		RunE:         prime.UpdateMTOServiceItem,
		SilenceUsage: false,
	}
	prime.InitUpdateMTOServiceItemFlags(updateMTOServiceItem.Flags())
	root.AddCommand(updateMTOServiceItem)

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

	getPaymentRequestEDI := &cobra.Command{
		Use:   "support-get-payment-request-edi",
		Short: "Get the EDI for a payment request",
		Long: `
  This command generates and returns the EDI for a given payment request.
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters.

  Endpoint path: /payment-requests/{paymentRequestID}/edi
  The file should contain json as follows:
  	{
  	"paymentRequestID": <uuid string>
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.GetPaymentRequestEDI,
		SilenceUsage: true,
	}
	support.InitGetPaymentRequestEDIFlags(getPaymentRequestEDI.Flags())
	root.AddCommand(getPaymentRequestEDI)

	processReviewedPaymentRequests := &cobra.Command{
		Use:   "support-reviewed-payment-requests",
		Short: "Use to test sending a payment request to syncada",
		Long: `
  This command gives the option to update the status of payment request to a given status.
  It also has the option to send the reviewed payment request to syncada.
  This is a support endpoint and is not available in production.
  It requires the caller to pass in a file using the --filename arg.
  The file should contain path parameters.

  Endpoint path: /payment-requests/process-reviewed
  The file should contain json as follows (only sendToSyncada is required):
  	{
	  body: {
		"paymentRequestID": <uuid string>,
		"sendToSyncada": <boolean>,
		"status": <string>
	  }
  	}
  Please see API documentation for full details on the endpoint definition.`,
		RunE:         support.ProcessReviewedPaymentRequests,
		SilenceUsage: true,
	}
	support.InitGetPaymentRequestEDIFlags(processReviewedPaymentRequests.Flags())
	root.AddCommand(processReviewedPaymentRequests)

	hideNonFakeMoveTaskOrdersCommand := &cobra.Command{
		Use:   "support-hide-non-fake-mtos",
		Short: "Hide moves not in the fake data spreadsheet",
		Long: `This command will trigger finding all of the moves in stg and env environments
		that do not use fake data from the fake names and addresses spreadsheet.
		To do this, all of the moves that do not match the data in the fake data spreadsheet
		will set the moves.show field to false.
		This will cause the move to not appear in the office applications.`,
		RunE:         support.HideNonFakeMoveTaskOrders,
		SilenceUsage: true,
	}
	support.InitHideNonFakeMoveTaskOrdersFlags(hideNonFakeMoveTaskOrdersCommand.Flags())
	root.AddCommand(hideNonFakeMoveTaskOrdersCommand)

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
