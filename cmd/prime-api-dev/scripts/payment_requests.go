package scripts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"pault.ag/go/pksigner"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	primeClient "github.com/transcom/mymove/pkg/gen/primeclient"
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
	mtoShipment "github.com/transcom/mymove/pkg/gen/primeclient/mto_shipment"
	paymentrequestclient "github.com/transcom/mymove/pkg/gen/primeclient/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

type menuType string

func (m menuType) String() string {
	return string(m)
}

const (
	// MenuMain is the ID for the Main Menu display
	MenuMain menuType = "MAIN"
	// MTOMenu is the ID for the MTO Menu display
	MTOMenu menuType = "MTO_MENU"
	// UpdateMTOMenu is the ID for the Update MTO Menu display
	UpdateMTOMenu menuType = "UPDATE_MTO"
	// UpdateShipmentMenu is the ID for the Update Shipment Menu display
	UpdateShipmentMenu menuType = "UPDATE_SHIPMENT"
	// CreatePaymentRequestMenu is the ID for the Create Payment Request Menu display
	CreatePaymentRequestMenu menuType = "CREATE_PR"
)
const (
	layoutISO = "2006-01-02"

//	layoutUS  = "January 2, 2006"
)

/*
menus := []menuType {
	MenuMain,
}

*/

// ReferenceID string

type mtoDisplay struct {
	mtoID       string
	description string
}

type serviceItemDisplay struct {
	mtoShipmentID string
	serviceItemID string
	description   string
}

// TODO payment request display not implemented
/*
type paymentRequestDisplay struct {
	paymentRequestID string
	description      string
}
*/

type mtoShipmentDisplay struct {
	mtoShipmentID string
	description   string
	etag          string
}

type paymentRequestsData struct {
	mtos                   primemessages.MoveTaskOrders
	mtoDisplayList         []mtoDisplay
	currentMTO             primemessages.MoveTaskOrder
	mtoShipmentDisplayList []mtoShipmentDisplay
	serviceItemDisplayList []serviceItemDisplay
	// TODO payment request implemented
	// paymentRequestDisplayList []paymentRequestDisplay
	v      *viper.Viper
	logger *log.Logger
	store  *pksigner.Store
}

type updateInfo struct {
	value    string
	isString bool
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

	pr := paymentRequestsData{
		logger: logger,
		v:      v,
	}

	nextMenu := MenuMain
	exitApp := false

	for ok := true; ok; ok = !exitApp {
		switch nextMenu {
		case MenuMain:
			exitApp, nextMenu, err = pr.displayMainMenu()
			if err != nil {
				fmt.Printf("error returned from Main Menu display: %s\n", err.Error())
			}
		case MTOMenu:
			exitApp, nextMenu, err = pr.displayMTOMenu()
			if err != nil {
				fmt.Printf("error returned from MTO Menu display: %s\n", err.Error())
			}
		case UpdateMTOMenu:
			exitApp, nextMenu, err = pr.displayUpdateMTOMenu()
			if err != nil {
				fmt.Printf("error returned from Update MTO Menu display: %s\n", err.Error())
			}
		case UpdateShipmentMenu:
			exitApp, nextMenu, err = pr.displayUpdateShipmentMenu()
			if err != nil {
				fmt.Printf("error returned from Update Shipment Menu display: %s\n", err.Error())
			}
		case CreatePaymentRequestMenu:
			exitApp, nextMenu, err = pr.displayCreatePaymentRequestMenu()
			if err != nil {
				fmt.Printf("error returned from Create Payment Request Menu display: %s\n", err.Error())
			}
		}

	}

	pr.cleanup()
	cleanup()

	return nil
}

func cleanup() {

}

func (pr *paymentRequestsData) cleanup() {
	// Defer closing the store until after the API call has completed
	if pr.store != nil {
		defer func() {
			if closeErr := pr.store.Close(); closeErr != nil {
				pr.logger.Fatal(closeErr)
			}
		}()
	}
}

func (pr *paymentRequestsData) printMTO(mto *primemessages.MoveTaskOrder) bool {

	if mto == nil {
		return false
	}

	if mto.AvailableToPrimeAt != nil && mto.Order != nil && mto.Order.DestinationDutyStation != nil &&
		mto.Order.DestinationDutyStation.Address != nil && mto.Order.DestinationDutyStation.Address.City != nil &&
		mto.Order.Customer != nil {
		return true
	}

	return false
}

func (pr *paymentRequestsData) displaySelectedMTO() {

	mto := pr.currentMTO
	pr.mtoShipmentDisplayList = []mtoShipmentDisplay{}
	pr.serviceItemDisplayList = []serviceItemDisplay{}

	fmt.Printf("\n\n============ DISPLAY SELECTED MTO ============ \n\n")

	if mto.AvailableToPrimeAt != nil {
		fmt.Printf("AvailableToPrime: %s\n", mto.AvailableToPrimeAt.String())
	}

	if mto.IsCanceled != nil {
		fmt.Printf("Is Canceled: %s\n", strconv.FormatBool(*mto.IsCanceled))
	}

	fmt.Printf("%s, %s\n", mto.Order.Customer.LastName, mto.Order.Customer.FirstName)

	fmt.Printf("Dest. Duty Station: %s, %s, %s\n", *mto.Order.DestinationDutyStation.Address.City,
		*mto.Order.DestinationDutyStation.Address.State, *mto.Order.DestinationDutyStation.Address.PostalCode)

	// Build shipment display descriptions
	for _, s := range mto.MtoShipments {
		var sstrs []string
		sstrs = append(sstrs, fmt.Sprintf("TOO approval date: %s\n", s.ApprovedDate.String()))
		if s.PickupAddress == nil {
			sstrs = append(sstrs, fmt.Sprint("Pickup address: <missing>\n"))
		} else {
			sstrs = append(sstrs, fmt.Sprintf("Pickup address: %s, %s, %s\n", *s.PickupAddress.City,
				*s.PickupAddress.State, *s.PickupAddress.PostalCode))
		}

		if s.DestinationAddress == nil {
			sstrs = append(sstrs, fmt.Sprint("Dest. address: <missing>\n"))
		} else {
			sstrs = append(sstrs, fmt.Sprintf("Dest. address: %s, %s, %s\n", *s.DestinationAddress.City,
				*s.DestinationAddress.State, *s.DestinationAddress.PostalCode))
		}

		sstrs = append(sstrs, fmt.Sprintf("Estimated weight: %d\n", s.PrimeEstimatedWeight))
		sstrs = append(sstrs, fmt.Sprintf("Actual weight: %d\n", s.PrimeActualWeight))

		sstrs = append(sstrs, fmt.Sprintf("Requested pickup date: %s\n", s.RequestedPickupDate.String()))
		sstrs = append(sstrs, fmt.Sprintf("Scheduled pickup date: %s\n", s.ScheduledPickupDate.String()))
		sstrs = append(sstrs, fmt.Sprintf("Actual pickup date: %s\n", s.ActualPickupDate.String()))

		shipmentDisplay := mtoShipmentDisplay{
			mtoShipmentID: s.ID.String(),
			description:   strings.Join(sstrs, ""),
			etag:          s.ETag,
		}
		pr.mtoShipmentDisplayList = append(pr.mtoShipmentDisplayList, shipmentDisplay)

	}

	for _, service := range mto.MtoServiceItems() {
		var strs []string

		strs = append(strs, fmt.Sprintf("%s status: %s\n", service.ReServiceName(), service.Status()))

		serviceItem := serviceItemDisplay{
			serviceItemID: service.ID().String(),
			description:   strings.Join(strs, ""),
			mtoShipmentID: service.MtoShipmentID().String(),
		}

		pr.serviceItemDisplayList = append(pr.serviceItemDisplayList, serviceItem)
	}

	// MTO level service items
	fmt.Print("\nService Items:\n")
	for ii, service := range pr.serviceItemDisplayList {
		if service.mtoShipmentID == "" {
			fmt.Printf("%d: %s", ii, service.description)
		}
	}

	for i, shipment := range pr.mtoShipmentDisplayList {
		fmt.Printf("\n%d: ============== Shipment %d =============\n", i, i)
		fmt.Printf("%s\n", shipment.description)

		// Shipment level service items
		fmt.Print("\nService Items:\n")
		for ii, service := range pr.serviceItemDisplayList {
			if service.mtoShipmentID == shipment.mtoShipmentID {
				fmt.Printf("%d: %s", ii, service.description)
			}
		}
	}

	/*
				AvailableToPrime:
				isCanceled:
				Branch:
			 	Lasttname, Firstname
				Dest Duty Station:

		        Service Items:
				#: code, description, status

				#: ============== Shipment 1 =============
				Shipment 1:
				TOO approval date:

				pickup address: city, state, zip
				Destination address: city, State, Zip

				estimated Weight <missing>
				actual weight <missing>
			    request pickup date <missing>
			    scheduled pickup date <missing>
			    actual pickup date <mission>
				first Availabe shipment date

			    Service Items:
				#: code, description, status

				#: ============== Payment Request Number =============
			    code, description, price, status
	*/

}

func (pr *paymentRequestsData) displayMTOS() {
	/*
		mtos primemessages.MoveTaskOrders
		mtoDisplayList []mtoDisplay
	*/
	// clear the list in case the MTOs fetch has been updates
	pr.mtoDisplayList = []mtoDisplay{}

	header := "#|\tAvail. to Prime|\tDest. City|\tCust. Name\n"

	// fill in new list of mtos
	for _, mto := range pr.mtos {

		if pr.printMTO(mto) == true {

			description := fmt.Sprintf("%s|\t%s|\t%s,%s\n", mto.AvailableToPrimeAt.String(),
				*mto.Order.DestinationDutyStation.Address.City,
				mto.Order.Customer.LastName, mto.Order.Customer.FirstName)
			info := mtoDisplay{
				mtoID:       mto.ID.String(),
				description: description,
			}
			pr.mtoDisplayList = append(pr.mtoDisplayList, info)
		}
	}

	// TODO
	// add a line for shipment: code, code, code, code under MTO to know which MTOs have service items already
	// maybe something line "existing payment reuest: #"

	// display to screen

	fmt.Printf("\n\n -------------- Returning MTOs -------------- \n\n")
	fmt.Print(header)
	for i, description := range pr.mtoDisplayList {
		fmt.Printf("%d: %s\n", i, description.description)
	}
}

func (pr *paymentRequestsData) displayShipment(index int) error {
	if index < len(pr.mtoShipmentDisplayList) {
		shipment := pr.mtoShipmentDisplayList[index]
		fmt.Println(shipment.description)

		fmt.Printf("\n\nService Items:\n")
		for ii, service := range pr.serviceItemDisplayList {
			if service.mtoShipmentID == shipment.mtoShipmentID {
				fmt.Printf("%d: %s", ii, service.description)
			}
		}
	} else {
		return fmt.Errorf("bad shipment index %d", index)
	}

	return nil
}

func getIntInput() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	selection, err := strconv.Atoi(text)
	return selection, err
}

func getStringInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text, nil
}

func getStrFmtDateInput() (strfmt.Date, error) {
	var dateString string
	var err error
	dateString, err = getStringInput()
	if err != nil {
		return strfmt.Date{}, err
	}
	var t time.Time
	t, err = time.Parse(layoutISO, dateString)
	return strfmt.Date(t), err
}

func (pr *paymentRequestsData) fetchMTOUpdates() error {

	fmt.Println("Fetching MTO updates... ")
	err := pr.fetchMTOs()
	if err != nil {
		fmt.Print("Error fetching updates, MTOs might be out of date, retry fetching or restart application.")
		return err
	}

	// if Current Selection is valid, update it with MTO updates
	for _, mto := range pr.mtos {
		if mto != nil {
			if mto.ID == pr.currentMTO.ID {
				pr.currentMTO = *mto
				fmt.Println("Current MTO was updated.")
			}
		}
	}

	return nil
}

func (pr *paymentRequestsData) updateShipmentsJSONToFile(f *os.File, shipmentUpdates map[string]updateInfo, shipmentIndex int) error {
	var strs []string

	/***************************************************
	{
		"mtoShipmentID": "5834016b-bc7b-421f-b87d-313acdddcc77",
		"ifMatch": "MjAyMC0wOS0wM1QyMTo1Nzo0NC43NjQ0ODFa",
		"body": {
		    "primeEstimatedWeight": 1000,
			"primeActualWeight": 3000
	    }
	}
	***************************************************/

	// {
	strs = append(strs, fmt.Sprint("{\n"))
	//		"mtoShipment": "ca9aeb58-e5a9-44b0-abe8-81d233dbdebf",
	strs = append(strs, fmt.Sprintf("\"mtoShipmentID\": \"%s\",\n", pr.mtoShipmentDisplayList[shipmentIndex].mtoShipmentID))
	//		"ifMatch": "MjAyMC0wOS0yOFQxNTo1OTozOC4zOTA0MjFa",
	strs = append(strs, fmt.Sprintf("\"ifMatch\": \"%s\",\n", pr.mtoShipmentDisplayList[shipmentIndex].etag))
	// 		"body": {
	strs = append(strs, fmt.Sprint("\"body\": {\n"))
	last := len(shipmentUpdates)
	counter := 0
	for key, value := range shipmentUpdates {
		counter++
		var fieldUpdate string
		//		"<field>": "<value>"
		if value.isString == true {
			fieldUpdate = fmt.Sprintf("\"%s\": \"%s\"", key, value.value)
		} else {
			fieldUpdate = fmt.Sprintf("\"%s\": %s", key, value.value)
		}

		if counter == last {
			// do not add ","
			fieldUpdate += "\n"
		} else {
			// need to add ","
			fieldUpdate += ",\n"
		}
		strs = append(strs, fieldUpdate)
	}

	// 		}  # close body{
	strs = append(strs, fmt.Sprint("}\n"))
	// }  # close json
	strs = append(strs, fmt.Sprint("}\n"))

	text := []byte(strings.Join(strs, ""))

	if _, err := f.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	return nil
}

func (pr *paymentRequestsData) displayUpdateShipmentMenu() (bool, menuType, error) {
	exitApp := false
	//var err error
	currentMenu := UpdateShipmentMenu

	// Select Shipment
	fmt.Println("Select shipment to update")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	shipmentIndex, err := strconv.Atoi(text)
	if err != nil {
		return exitApp, currentMenu, fmt.Errorf("error with selecting shipment: %w", err)
	}

	// Display Shipment
	fmt.Printf("\n\n================ Updating Shipment ===============\n")
	err = pr.displayShipment(shipmentIndex)
	if err != nil {
		fmt.Print("Could not display selected MTO shipment")
		fmt.Printf("error when selection shipment with index %d: %s", shipmentIndex, err.Error())
		return exitApp, MTOMenu, nil
	}

	// fields to update
	const (
		actualPickupDate = iota
		requestedPickupDate
		scheduledPickupDate
		primeEstimatedWeight
		primeActualWeight
		sendShipmentUpdate
		cancelUpdateShipment
	)

	fields := []struct {
		field       int
		description string
		json        string
	}{
		{
			field:       actualPickupDate,
			description: "Actual Pickup Date",
			json:        "actualPickupDate",
		},
		{
			field:       requestedPickupDate,
			description: "Requested Pickup Date",
			json:        "requestedPickupDate",
		},
		{
			field:       scheduledPickupDate,
			description: "Scheduled Pickup Date",
			json:        "scheduledPickupDate",
		},
		{
			field:       primeEstimatedWeight,
			description: "Prime Estimated Weight",
			json:        "primeEstimatedWeight",
		},
		{
			field:       primeActualWeight,
			description: "Prime Actual Weight",
			json:        "primeActualWeight",
		},
		{
			field:       sendShipmentUpdate,
			description: "Send shipment updates",
		},
		{
			field:       cancelUpdateShipment,
			description: "Cancel updating shipment",
		},
	}

	menuTitle("UPDATE SHIPMENT")
	shipment := primemessages.MTOShipment{}

	// print out update field options for selection
	for i, field := range fields {
		fmt.Printf("%d: %s\n", i, field.description)
	}

	shipmentUpdates := make(map[string]updateInfo)
	update := true
	for ok := true; ok; ok = update {
		fmt.Printf("\nSelect field to update: ")
		var selection int
		selection, err = getIntInput()
		if err != nil {
			log.Fatal("Cannot get int input", err)
		}
		selectedField := fields[selection]
		switch selectedField.field {
		case actualPickupDate:
			fmt.Printf("Updating %s\nEnter date as format YYYY-MM-DD: ", selectedField.description)
			var strFmtDate strfmt.Date
			strFmtDate, err = getStrFmtDateInput()
			if err != nil {
				log.Fatal("Cannot get date input", err)
			}
			shipment.ActualPickupDate = strFmtDate
			fieldValue := updateInfo{
				value:    strFmtDate.String(),
				isString: true,
			}
			shipmentUpdates[selectedField.json] = fieldValue
		case requestedPickupDate:
			fmt.Printf("Updating %s\nEnter date as format YYYY-MM-DD: ", selectedField.description)
			var strFmtDate strfmt.Date
			strFmtDate, err = getStrFmtDateInput()
			if err != nil {
				log.Fatal("Cannot get date input", err)
			}
			shipment.RequestedPickupDate = strFmtDate
			fieldValue := updateInfo{
				value:    strFmtDate.String(),
				isString: true,
			}
			shipmentUpdates[selectedField.json] = fieldValue
		case scheduledPickupDate:
			fmt.Printf("Updating %s\nEnter date as format YYYY-MM-DD: ", selectedField.description)
			var strFmtDate strfmt.Date
			strFmtDate, err = getStrFmtDateInput()
			if err != nil {
				log.Fatal("Cannot get date input", err)
			}
			shipment.ScheduledPickupDate = strFmtDate
			fieldValue := updateInfo{
				value:    strFmtDate.String(),
				isString: true,
			}
			shipmentUpdates[selectedField.json] = fieldValue
		case primeEstimatedWeight:
			fmt.Printf("Updating %s\nEnter weight: ", selectedField.description)
			var weight int
			weight, err = getIntInput()
			if err != nil {
				log.Fatal("Cannot get int input", err)
			}
			shipment.PrimeEstimatedWeight = int64(weight)
			fieldValue := updateInfo{
				value:    strconv.Itoa(weight),
				isString: false,
			}
			shipmentUpdates[selectedField.json] = fieldValue
		case primeActualWeight:
			fmt.Printf("Updating %s\nEnter weight: ", selectedField.description)
			var weight int
			weight, err = getIntInput()
			if err != nil {
				log.Fatal("Cannot get int input", err)
			}
			shipment.PrimeActualWeight = int64(weight)
			fieldValue := updateInfo{
				value:    strconv.Itoa(weight),
				isString: false,
			}
			shipmentUpdates[selectedField.json] = fieldValue
		case sendShipmentUpdate:
			payload := mtoShipment.UpdateMTOShipmentParams{}
			payload.Body = &shipment
			payload.IfMatch = pr.mtoShipmentDisplayList[shipmentIndex].etag
			payload.MtoShipmentID = strfmt.UUID(pr.mtoShipmentDisplayList[shipmentIndex].mtoShipmentID)
			// err = pr.updateMTOShipment(payload)

			tmpFile, err := ioutil.TempFile(os.TempDir(), "update-shipment-")
			if err != nil {
				log.Fatal("Cannot create temporary file", err)
			}

			/*
				// Remember to clean up the file afterwards
				defer os.Remove(tmpFile.Name())

			*/

			fmt.Println("Created File: " + tmpFile.Name())

			/*
				{
				  "mtoShipmentID": "5834016b-bc7b-421f-b87d-313acdddcc77",
				  "ifMatch": "MjAyMC0wOS0wM1QyMTo1Nzo0NC43NjQ0ODFa",
				  "body": {
				    "primeEstimatedWeight": 1000,
				    "primeActualWeight": 3000
				  }
				}
			*/

			err = pr.updateShipmentsJSONToFile(tmpFile, shipmentUpdates, shipmentIndex)

			if err != nil {
				fmt.Println("Create update shipments JSON file failed")
				fmt.Printf("error message: %s\n", err.Error())
			}

			err = pr.updateMTOShipment2(tmpFile.Name())

			/*
				// Close the file
				if err := tmpFile.Close(); err != nil {
					log.Fatal(err)
				}
			*/
			if err != nil {
				fmt.Println("Shipment update failed :( ")
				fmt.Printf("error message: %s\n", err.Error())
			} else {
				fmt.Printf("\nShipment update was successfully sent for processing (see reesponse for update success/fail)...\n")

				err = pr.fetchMTOUpdates()
				if err != nil {
					fmt.Print("Could not fetch MTO updates")
					return exitApp, MTOMenu, nil
				}

				// re-display update shipment menu and the current shipment that was updated

				// Display Shipment
				fmt.Printf("\n\n================ Updating Shipment ===============\n")
				err = pr.displayShipment(shipmentIndex)
				if err != nil {
					fmt.Print("Could not display selected MTO shipment")
					fmt.Printf("error when selection shipment with index %d: %s", shipmentIndex, err.Error())
					return exitApp, MTOMenu, nil
				}

				menuTitle("UPDATE SHIPMENT")

				// print out update field options for selection
				for i, field := range fields {
					fmt.Printf("%d: %s\n", i, field.description)
				}
			}
		case cancelUpdateShipment:
			update = false
		}
	}

	return exitApp, MTOMenu, nil
}

func (pr *paymentRequestsData) paymentRequestJSONToFile(f *os.File, serviceItems map[string]updateInfo) error {
	var strs []string

	/***************************************************
	{
	  "body": {
	    "isFinal": false,
	    "moveTaskOrderID": "49abcdbf-d4ed-4c9c-9ce1-677ee7653f77",
	    "serviceItems": [
	      {
	        "id": "22e1bf91-47db-45df-bade-f52d7252ed3e"
	      },
	      {
	        "id": "2c4fb2a4-a5c3-4438-92ba-64e8ba1502e4"
	      }
	    ]
	  }
	}
	***************************************************/

	// {
	strs = append(strs, fmt.Sprint("{\n"))
	// 		"body": {
	strs = append(strs, fmt.Sprint("\"body\": {\n"))
	//		    "isFinal": false,
	strs = append(strs, fmt.Sprintf("\"isFinal\": %s,\n", serviceItems["isFinal"].value))
	//		    "moveTaskOrderID": "49abcdbf-d4ed-4c9c-9ce1-677ee7653f77",
	strs = append(strs, fmt.Sprintf("\"moveTaskOrderID\": \"%s\",\n", pr.currentMTO.ID.String()))

	//"serviceItems": [
	strs = append(strs, fmt.Sprint("\"serviceItems\": [\n"))
	last := len(serviceItems)
	last-- // need to account for skipping over 'isFinal' key
	counter := 0
	for key, value := range serviceItems {

		if key == "isFinal" {
			continue
		} else {
			counter++
		}
		var serviceItemStr string
		//{
		//		    "id": "ca9aeb58-e5a9-44b0-abe8-81d233dbdebf"
		//},
		serviceItemStr = fmt.Sprintf("{\n\"id\": \"%s\"\n}", value.value)

		if counter == last {
			// do not add ","
			serviceItemStr += "\n"
		} else {
			// need to add ","
			serviceItemStr += ",\n"
		}
		strs = append(strs, serviceItemStr)
	}

	//          ]  # close "serviceItems": [
	strs = append(strs, fmt.Sprint("]\n"))
	// 		}  # close body{
	strs = append(strs, fmt.Sprint("}\n"))
	// }  # close json
	strs = append(strs, fmt.Sprint("}\n"))

	text := []byte(strings.Join(strs, ""))

	if _, err := f.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	return nil
}

func (pr *paymentRequestsData) displayCreatePaymentRequestMenu() (bool, menuType, error) {

	exitApp := false
	//var err error
	//currentMenu := CreatePaymentRequestMenu

	// Display MTO currently selected
	pr.displaySelectedMTO()

	// fields to update
	const (
		addServiceItem = iota
		updateIsFinal
		sendPaymentRequest
		cancelPaymentRequest
	)

	options := []struct {
		option      int
		description string
	}{
		{
			option:      addServiceItem,
			description: "Select service item index to add",
		},
		{
			option:      updateIsFinal,
			description: "Update isFinal flag (default is false)",
		},
		{
			option:      sendPaymentRequest,
			description: "Send Payment Request",
		},
		{
			option:      cancelPaymentRequest,
			description: "Cancel and return to previous menu",
		},
	}

	menuTitle("CREATE PAYMENT REQUEST")

	// print out update field options for selection
	for i, option := range options {
		fmt.Printf("%d: %s\n", i, option.description)
	}

	serviceItems := make(map[string]updateInfo)
	serviceItems["isFinal"] = updateInfo{
		value:    "false",
		isString: false,
	}
	var err error
	update := true
	for ok := true; ok; ok = update {
		fmt.Printf("\nSelect menu option: ")
		var selection int
		selection, err = getIntInput()
		if err != nil {
			fmt.Printf("\nBad option input with error: %s\n", err.Error())
		}
		selectedOption := options[selection]
		switch selectedOption.option {
		case addServiceItem:
			fmt.Print("Select service index to add to payment request: ")
			var index int
			index, err = getIntInput()
			if err != nil {
				fmt.Printf("\nBad index input with error: %s\n", err.Error())
			} else {
				if index < len(pr.serviceItemDisplayList) {
					serviceItem := pr.serviceItemDisplayList[index]
					fieldValue := updateInfo{
						value:    serviceItem.serviceItemID,
						isString: true,
					}
					serviceItems[serviceItem.serviceItemID] = fieldValue
				} else {
					fmt.Printf("\nEntered bad service item index <%d> try again.\n", index)
				}
			}
		case updateIsFinal:
			fmt.Print("Updating isFinal flag\nEnter \"0\" for false and \"1\" for true: ")
			var isFinal int
			isFinal, err = getIntInput()
			if err != nil {
				fmt.Printf("\nBad 0 or 1 input with error: %s\n", err.Error())
			} else {
				if isFinal == 0 {
					fieldValue := updateInfo{
						value:    "false",
						isString: true,
					}
					serviceItems["isFinal"] = fieldValue
				} else if isFinal == 1 {
					fieldValue := updateInfo{
						value:    "true",
						isString: true,
					}
					serviceItems["isFinal"] = fieldValue
				} else {
					fmt.Printf("\nDid not enter valid \"0\" or \"1\" for isFinal flag: %d.\n", isFinal)
				}
			}
		case sendPaymentRequest:

			tmpFile, err := ioutil.TempFile(os.TempDir(), "payment-request-")
			if err != nil {
				log.Fatal("Cannot create temporary file", err)
			}

			/*
				// Remember to clean up the file afterwards
				defer os.Remove(tmpFile.Name())

			*/

			fmt.Println("Created File: " + tmpFile.Name())

			/*
				{
				  "mtoShipmentID": "5834016b-bc7b-421f-b87d-313acdddcc77",
				  "ifMatch": "MjAyMC0wOS0wM1QyMTo1Nzo0NC43NjQ0ODFa",
				  "body": {
				    "primeEstimatedWeight": 1000,
				    "primeActualWeight": 3000
				  }
				}
			*/

			err = pr.paymentRequestJSONToFile(tmpFile, serviceItems)

			if err != nil {
				fmt.Println("Create payment request JSON file failed")
				fmt.Printf("error message: %s\n", err.Error())
			}

			err = pr.creatPaymentRequest(tmpFile.Name())

			if err != nil {
				fmt.Println("Create payment request failed :( ")
				fmt.Printf("error message: %s\n", err.Error())
			} else {
				fmt.Printf("\nCreate payment request was successfully sent for processing (see reesponse for update success/fail)...\n")

				err = pr.fetchMTOUpdates()
				if err != nil {
					fmt.Print("Could not fetch MTO updates")
					return exitApp, MTOMenu, nil
				}

				// re-display updated MTO

				// Display MTO currently selected
				pr.displaySelectedMTO()

				menuTitle("CREATE PAYMENT REQUEST")

				// print out update field options for selection
				for i, option := range options {
					fmt.Printf("%d: %s\n", i, option.description)
				}
			}
		case cancelPaymentRequest:
			update = false
		}
	}

	return exitApp, MTOMenu, nil
}

func (pr *paymentRequestsData) displayUpdateMTOMenu() (bool, menuType, error) {

	fmt.Printf("\n UPDATE MTO not implemented \n")

	return false, MTOMenu, nil
}

func (pr *paymentRequestsData) displayMTOMenu() (bool, menuType, error) {

	// See MTO details, see shipment fields and see service items
	//		fields we care about for pricing: available to prime, request pickup date, scheduled pickup date,
	//      estimated weight, actual weight
	//
	// Update MTO
	// 		take you to the Update MTO menu
	//
	// Update MTO Shipment
	//      take you to the Update MTO shipment menu
	//
	// Send payment request
	//      select which service items to use
	//
	// View payment request
	//

	const (
		UpdateMTO = iota
		UpdateShipment
		CreatePaymentRequest
		ViewPaymentRequest
		FetchMTO
		PreviousMenu
		ExitApp
	)
	exitApp := false
	var err error
	var selection int
	currentMenu := MTOMenu

	display := map[int]struct {
		option      int
		description string
		nextMenu    menuType
	}{
		UpdateMTO: {
			option:      UpdateMTO,
			description: "Update Current MTO",
			nextMenu:    UpdateMTOMenu,
		},
		UpdateShipment: {
			option:      UpdateShipment,
			description: "Update Shipment",
			nextMenu:    UpdateShipmentMenu,
		},
		CreatePaymentRequest: {
			option:      CreatePaymentRequest,
			description: "Create Payment Request",
			nextMenu:    CreatePaymentRequestMenu,
		},
		PreviousMenu: {
			option:      PreviousMenu,
			description: "Go to previous menu",
			nextMenu:    MenuMain,
		},
		ExitApp: {
			option:      ExitApp,
			description: "Exit",
			nextMenu:    currentMenu,
		},
	}

	// Display MTO Details for selected MTO
	pr.displaySelectedMTO()

	menuTitle("MTO MENU")
	// Display menu options
	for _, option := range display {
		fmt.Printf("%d: %s\n", option.option, option.description)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nSelect menu option: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	selection, err = strconv.Atoi(text)
	if err != nil {
		fmt.Printf("Error with selection <%s>, try again\n", err.Error())
		return exitApp, currentMenu, nil
	}

	switch selection {
	case UpdateMTO:
		return exitApp, display[selection].nextMenu, nil
	case UpdateShipment:
		return exitApp, display[selection].nextMenu, nil
	case CreatePaymentRequest:
		_, _, err := pr.displayCreatePaymentRequestMenu()
		if err != nil {
			fmt.Printf("Error with creating payment <%s>", err.Error())
		}
		return exitApp, display[selection].nextMenu, nil
	case PreviousMenu:
		return exitApp, display[selection].nextMenu, nil
	case ExitApp:
		exitApp = true
		return exitApp, display[selection].nextMenu, nil
	}

	return exitApp, currentMenu, nil
}

func menuTitle(name string) {
	fmt.Printf("\n\n")
	fmt.Println("################################################")
	fmt.Printf("                %s\n", name)
	fmt.Print("################################################\n\n")
}

func (pr *paymentRequestsData) displayMainMenu() (bool, menuType, error) {
	const (
		FetchDisplay = iota
		Display
		SelectMTO
		PreviousMenu
		ExitApp
	)
	exitApp := false
	var err error
	var selection int
	currentMenu := MenuMain

	display := map[int]struct {
		option      int
		description string
		nextMenu    menuType
	}{
		FetchDisplay: {
			option:      FetchDisplay,
			description: "Fetch and display MTOs",
			nextMenu:    currentMenu,
		},
		Display: {
			option:      Display,
			description: "Display MTOs",
			nextMenu:    currentMenu,
		},
		SelectMTO: {
			option:      SelectMTO,
			description: "Select MTO",
			nextMenu:    MTOMenu,
		},

		ExitApp: {
			option:      ExitApp,
			description: "Exit",
			nextMenu:    currentMenu,
		},
	}
	reader := bufio.NewReader(os.Stdin)
	//fmt.Print("\nPress <0> to display Payment Request: ")

	menuTitle("MAIN MENU")

	for _, option := range display {
		fmt.Printf("%d: %s\n", option.option, option.description)
	}
	fmt.Print("Enter selection: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	selection, err = strconv.Atoi(text)
	if err != nil {
		fmt.Printf("Error with selection <%s>, try again\n", err.Error())
		return exitApp, currentMenu, nil
	}

	switch selection {
	case FetchDisplay:
		err = pr.fetchMTOs()
		if err != nil {
			fmt.Printf("Error fetching MTO <%s>, try again", err.Error())
		}
		pr.displayMTOS()
		return exitApp, display[selection].nextMenu, nil
	case Display:
		pr.displayMTOS()
		return exitApp, display[selection].nextMenu, nil
	case SelectMTO:
		err = pr.selectMTO()
		if err != nil {
			fmt.Printf("Error selecting MTO index<%d>, try again %s\n", selection, err.Error())
			return exitApp, currentMenu, nil
		}
		return exitApp, display[selection].nextMenu, nil

	case ExitApp:
		exitApp = true
		return exitApp, display[selection].nextMenu, nil

	}

	return exitApp, currentMenu, nil
}
func (pr *paymentRequestsData) selectMTO() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter MTO index to select and use: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	selection, err := strconv.Atoi(text)

	if err != nil {
		return err
	}

	mtoID := pr.mtoDisplayList[selection].mtoID
	for _, mto := range pr.mtos {
		if mto != nil {
			if mtoID == mto.ID.String() {
				pr.currentMTO = *mto
				return nil
			}
		}
	}
	return nil
}

// updateMTOShipment creates a gateway and sends the request to the endpoint
func (pr *paymentRequestsData) updateMTOShipment2(filename string) error {

	/*
		// Show what we are sending
		showJsonPayload, errJSONMarshall := json.Marshal(shipmentPayload)
		if errJSONMarshall != nil {
			pr.logger.Fatal(errJSONMarshall)
		}
		fmt.Printf("Sending payload for shipment updates...\n")
		fmt.Println(string(showJsonPayload))

	*/

	// Decode json from file that was passed into MTOShipment
	var shipmentPayload mtoShipment.UpdateMTOShipmentParams
	err := utils.DecodeJSONFileToPayload(filename, false, &shipmentPayload)
	if err != nil {
		pr.logger.Fatal(err)
	}

	shipmentPayload.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	primeGateway, _, errCreateClient := pr.getPrimeClient()
	if errCreateClient != nil {
		return errCreateClient
	}

	// Make the API Call
	resp, err := primeGateway.MtoShipment.UpdateMTOShipment(&shipmentPayload)
	if err != nil {
		fmt.Printf("\n\nprimeGateway.MtoShipment.UpdateMTOShipment() failed with: [%s]\n\n", err.Error())
		return utils.HandleGatewayError(err, pr.logger)
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.MarshalIndent(payload, "", "    ")
		if errJSONMarshall != nil {
			pr.logger.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		pr.logger.Fatal(resp.Error())
	}

	// Defer closing the store until after the API call has completed
	/*
		if cacStore != nil {
			defer cacStore.Close()
		}

	*/

	return nil
}

// updateMTOShipment creates a gateway and sends the request to the endpoint
func (pr *paymentRequestsData) creatPaymentRequest(filename string) error {

	/*
		// Show what we are sending
		showJsonPayload, errJSONMarshall := json.Marshal(shipmentPayload)
		if errJSONMarshall != nil {
			pr.logger.Fatal(errJSONMarshall)
		}
		fmt.Printf("Sending payload for shipment updates...\n")
		fmt.Println(string(showJsonPayload))

	*/

	// Decode json from file that was passed into MTOShipment
	paymentRequestParams := paymentrequestclient.CreatePaymentRequestParams{}
	err := utils.DecodeJSONFileToPayload(filename, false, &paymentRequestParams)
	if err != nil {
		pr.logger.Fatal(err)
	}

	paymentRequestParams.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	primeGateway, _, errCreateClient := pr.getPrimeClient()
	if errCreateClient != nil {
		return errCreateClient
	}

	// Make the API Call
	resp, err := primeGateway.PaymentRequest.CreatePaymentRequest(&paymentRequestParams)
	if err != nil {
		fmt.Printf("\n\nprimeGateway.PaymentRequest.CreatePaymentRequest() failed with: [%s]\n\n", err.Error())
		return utils.HandleGatewayError(err, pr.logger)
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.MarshalIndent(payload, "", "    ")
		if errJSONMarshall != nil {
			pr.logger.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		pr.logger.Fatal(resp.Error())
	}

	// Defer closing the store until after the API call has completed
	/*
		if cacStore != nil {
			defer cacStore.Close()
		}

	*/

	return nil
}

func (pr *paymentRequestsData) getPrimeClient() (*primeClient.Mymove, *pksigner.Store, error) {
	// TODO When making a call and you already have the CAC store, if you get back CKR_USER_NOT_LOGGED_IN then
	// log in again
	// TODO added additional error logging to see where error is happening
	// TODO another error if already initialize is 2020/09/29 23:00:56 pkcs11: 0x191: CKR_CRYPTOKI_ALREADY_INITIALIZED
	// created new function CreatePrimeClientWithCACStoreParam to avoid this.
	if pr.store == nil {
		primeGateway, cacStore, errCreateClient := utils.CreatePrimeClient(pr.v)
		pr.store = cacStore
		return primeGateway, cacStore, errCreateClient
	}

	// else use existing CAC store
	primeGateway, cacStore, errCreateClient := utils.CreatePrimeClientWithCACStoreParam(pr.v, pr.store)
	pr.store = cacStore
	return primeGateway, cacStore, errCreateClient
}

func (pr *paymentRequestsData) fetchMTOs() error {
	pr.mtos = primemessages.MoveTaskOrders{}

	primeGateway, _, errCreateClient := pr.getPrimeClient()
	if errCreateClient != nil {
		return errCreateClient
	}

	var params mto.FetchMTOUpdatesParams
	params.SetTimeout(time.Second * 30)
	resp, err := primeGateway.MoveTaskOrder.FetchMTOUpdates(&params)
	if err != nil {
		fmt.Printf("\n\nprimeGateway.MoveTaskOrder.FetchMTOUpdates() failed with: [%s]\n\n", err.Error())
		return utils.HandleGatewayError(err, pr.logger)
	}

	// primemessages.MoveTaskOrders
	payload := resp.GetPayload()
	if payload != nil {
		pr.mtos = payload
		return nil
	}

	pr.logger.Fatal(resp.Error())

	// TODO you'll never hit this on success...
	// Defer closing the store until after the API call has completed
	/*
		if cacStore != nil {
			defer cacStore.Close()
		}
	*/
	return nil
}
