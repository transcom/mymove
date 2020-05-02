package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
)

// initReviewPaymentRequestFlags initializes flags.
func initReviewPaymentRequestFlags(flag *pflag.FlagSet) {
	flag.String(PaymentRequestID, "", "Payment Request ID to review")

	flag.SortFlags = false
}

// checkReviewPaymentRequestConfig checks the args.
func checkReviewPaymentRequestConfig(v *viper.Viper, args []string, logger logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		return err
	}

	return nil
}

func displayPaymentRequest(pr models.PaymentRequest, serviceItems models.PaymentServiceItems, tx *pop.Connection) (bool, int, error) {
	fmt.Print("\nDisplay Payment Request\n")

	var err error
	var input int

	s, _ := json.MarshalIndent(pr, "", "\t")
	fmt.Println(string(s))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress <0> to exit, <1> to display Service Item List: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	input, err = strconv.Atoi(text)
	fmt.Println(text)

	if err != nil {
		return true, 0, fmt.Errorf("error with user input %w", err)
	}

	if input == 0 {
		// exit program
		return false, input, nil
	} else if input == 1 {
		// display list of service items for this payment request
		input, err = displayPaymentRequestServiceItems(serviceItems)
		if err != nil {
			return false, 0, err
		}
		if input == 0 {
			// return and display payment request
			return true, input, nil
		} else if input > 0 {
			// display service item details (not params)
			serviceItem := serviceItems[input-1]
			input, err = displayPaymentRequestServiceItemDetails(serviceItem)
			if err != nil {
				return false, 0, err
			}
			if input == 0 {
				// return and display payment request
				return true, input, nil
			} else if input == 1 {
				// display service item params list (all of them) for current service item
				input, err = displayPaymentRequestServiceItemParams(serviceItem, tx)
				if err != nil {
					return false, 0, err
				}
				// return to display payment request
				return true, input, nil
			}
		}
	}
	return false, 0, nil
}

func displayPaymentRequestServiceItems(serviceItems models.PaymentServiceItems) (int, error) {
	fmt.Print("\nDisplay Service Item List\n")

	fmt.Printf("\nService Items\n")

	// index, ReService code, ReService name
	for i, item := range serviceItems {
		fmt.Printf("\t [%d] (%s) %s\n", i+1, item.MTOServiceItem.ReService.Code, item.MTOServiceItem.ReService.Name)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress <0> to display Payment Request, index number [#] to display Service Item: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	num, err := strconv.Atoi(text)
	fmt.Println(text)

	if err != nil {
		return 0, errors.New("error with user input")
	}

	return num, nil
}

func displayPaymentRequestServiceItemDetails(si models.PaymentServiceItem) (int, error) {
	fmt.Print("\nDisplay Service Item Details\n")

	s, _ := json.MarshalIndent(si, "", "\t")
	fmt.Println(string(s))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress <0> to display Payment Request, <1> to display Service Item Params: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	num, err := strconv.Atoi(text)
	fmt.Println(text)

	if err != nil {
		return 0, errors.New("error with user input")
	}

	return num, nil
}

func displayPaymentRequestServiceItemParams(si models.PaymentServiceItem, tx *pop.Connection) (int, error) {
	fmt.Print("\nDisplay Service Item Params\n")

	params := models.PaymentServiceItemParams{}
	err := tx.Q().Where("payment_service_item_id = $1", si.ID).Eager("ServiceItemParamKey").All(&params)
	if err != nil {
		return 0, fmt.Errorf("failure fetching service item params: %w", err)
	}

	for i, param := range params {
		fmt.Printf("(%d) Value: %s, IncomingKey: %s\n", i+1, param.Value, param.IncomingKey)
		fmt.Printf("ServiceItemParamKey:\n")
		s, _ := json.MarshalIndent(param.ServiceItemParamKey, "", "\t")
		fmt.Println(string(s))
	}

	var num int
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress <0> to display Payment Request: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	num, err = strconv.Atoi(text)
	fmt.Println(text)

	if err != nil {
		return 0, err
	}

	return num, nil
}


// reviewPaymentRequest reviews the payment request for an MTO
func reviewPaymentRequest(cmd *cobra.Command, args []string) error {
	v := viper.New()
	cli.InitDatabaseFlags(cmd.Flags())

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Create logger
	logger, err := logging.Config(v.GetString(cli.DbEnvFlag), v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// Check the config before talking to the CAC
	err = checkReviewPaymentRequestConfig(v, args, logger)
	if err != nil {
		logger.Fatal("check config", zap.Error(err))
	}

	// cac and api gateway
	_, cacStore, errCreateClient := CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Connect to the database
	// Set up DB flags
	err = cli.CheckDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}
	db, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Fatal("Could not close database", zap.Error(closeErr))
		}
	}()

	// Get params from CLI
	paymentRequestID := v.GetString(PaymentRequestID)
	if paymentRequestID == "" {
		return fmt.Errorf("paymentRequestID required: %s", paymentRequestID)
	}

	queryBuilder := query.NewQueryBuilder(db)
	paymentRequestFetcher := paymentrequest.NewPaymentRequestFetcher(queryBuilder)

	// Let's fetch the existing payment request using the PaymentRequestFetcher service object
	filter := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentRequestID)}
	existingPaymentRequest, err := paymentRequestFetcher.FetchPaymentRequest(filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Error finding Payment Request for status update with ID: %s", paymentRequestID), zap.Error(err))
		return fmt.Errorf("cannot find paymentRequestID: %s", paymentRequestID)
	}

	// Fetch all of the service items and keep in a slice for display
	prServiceItems := models.PaymentServiceItems{}
	err = db.Q().Where("payment_request_id = $1", paymentRequestID).Eager("MTOServiceItem.ReService").All(&prServiceItems)
	if err != nil {
		return fmt.Errorf("failure fetching service items: %w", err)
	}

	for cont := true; cont; cont, _, err = displayPaymentRequest(existingPaymentRequest, prServiceItems, db) {
		if err != nil {
			return fmt.Errorf("failed to display payment request with ID %s %w: ", paymentRequestID, err)
		}
	}

	return nil
}
