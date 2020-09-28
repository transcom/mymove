package scripts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	"github.com/transcom/mymove/pkg/gen/primemessages"

	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

type MenuType string
func (m MenuType) String() string {
	return string(m)
}
const (
	MenuMain MenuType = "MAIN"
	SelectMTOMenu MenuType = "SELECT_MTO"
	UpdateMTOMenu MenuType = "UPDATE_MTO"
	UpdateShipmentMenu MenuType = "UPDATE_SHIPMENT"
	SelectServiceItems MenuType = "SELECT_SERVICE_ITEMS"
)

menus := []MenuType {
	MenuMain,
}

// ReferenceID string

type mtoDisplay struct {
	mtoID string
	description string
}

type serviceItemDisplay struct {
	serviceItemID string
	description string

}

type paymentRequestDisplay struct {
	paymentRequestID string
	description string
}

type mtoShipmentDisplay struct {
	mtoShipmentID string
	description string
}


type PaymentRequestsData struct {
	mtos primemessages.MoveTaskOrders
	mtoDisplayList []mtoDisplay
	currentMTO primemessages.MoveTaskOrder
	mtoShipmentDisplayList []mtoShipmentDisplay
	serviceItemDisplayList map[string][]serviceItemDisplay // map of MTO Shipment IDs to a list of service items
	paymentRequestDisplayList []paymentRequestDisplay
	v *viper.Viper
	logger *log.Logger
}


// InitPaymentRequestsFlags declares which flags are enabled
func InitPaymentRequestsFlags(flag *pflag.FlagSet) {
	flag.SortFlags = false
}

func checkPaymentRequestsConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	return nil
}

// PaymentRequests TBD
func PaymentRequests(cmd *cobra.Command, args []string) error {
	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	v := viper.New()

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkPaymentRequestsConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	displayMainMenu(logger)


	return nil
}

func displayMTOS() {

}

func selectMTOMenu(logger *log.Logger) {

}

func displayMainMenu(logger *log.Logger) (bool, MenuType, error) {
	const (
		FetchDisplay = iota
		Display
		SelectMTO
		ExitApp
	)
	exitApp := false
	var err error
	var selection int
	currentMenu := MenuMain

	display := []struct{
		option int
		description string
		nextMenu MenuType
	} {
		{
			option:      FetchDisplay,
			description: "Fetch and display MTOs",
			nextMenu: currentMenu,
		},
		{
			option:      Display,
			description: "Display MTOs",
			nextMenu: currentMenu,
		},
		{
			option:      SelectMTO,
			description: "Select MTO",
			nextMenu: SelectMTOMenu,
		},
		{
			option:      ExitApp,
			description: "Exit",
			nextMenu: currentMenu,
		},
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress <0> to display Payment Request: ")

	for _, option := range display {
		fmt.Printf("%d: %s\n", option.option, option.description)
	}
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	selection, err = strconv.Atoi(text)

	switch selection {
	case FetchDisplay:
		mtos, err := fetchMTOs(logger)
		displayMTOS()
		return exitApp, display[selection-1].nextMenu, nil
	case Display:
		displayMTOS()
		return exitApp, display[selection-1].nextMenu, nil
	}

	return exitApp, currentMenu, nil
}

func fetchMTOs(logger *log.Logger, v *viper.Viper) (primemessages.MoveTaskOrders, error) {
	var mtos primemessages.MoveTaskOrders



	primeGateway, cacStore, errCreateClient := utils.CreatePrimeClient(v)
	if errCreateClient != nil {
		return mtos, errCreateClient
	}

	var params mto.FetchMTOUpdatesParams
	params.SetTimeout(time.Second * 30)
	resp, err := primeGateway.MoveTaskOrder.FetchMTOUpdates(&params)
	if err != nil {
		return mtos, utils.HandleGatewayError(err, logger)
	}

	// primemessages.MoveTaskOrders
	payload := resp.GetPayload()
	if payload != nil {
		return payload, nil //payload = mtos
	} else {
		logger.Fatal(resp.Error())
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}
	return mtos, nil
}