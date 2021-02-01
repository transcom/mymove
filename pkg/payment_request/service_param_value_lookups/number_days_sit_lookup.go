package serviceparamvaluelookups

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/models"
)

// TODO: Write NumberDaysSITLookup description
type NumberDaysSITLookup struct {
	MTOShipment models.MTOShipment
}

func (s NumberDaysSITLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	db := *keyData.db

	fmt.Printf("%v", keyData.MTOServiceItem.SITEntryDate)
	fmt.Println()
	fmt.Printf("%v", keyData.MTOServiceItem.SITDepartureDate)
	fmt.Println()

	var allSITPaymentServiceItems models.PaymentServiceItems
	err := db.Q().
		Join("mto_service_items msi", "msi.id = payment_service_items.mto_service_item_id").
		Join("re_services rs", "rs.id = msi.re_service_id").
		Eager("MTOServiceItem.ReService").
		Where("payment_service_items.status IN ($1, $2, $3, $4) AND msi.move_id = ($5) AND rs.code IN ($6, $7, $8, $9)", models.PaymentServiceItemStatusRequested, models.PaymentServiceItemStatusApproved, models.PaymentServiceItemStatusSentToGex, models.PaymentServiceItemStatusPaid, s.MTOShipment.MoveTaskOrderID, models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT).
		All(&allSITPaymentServiceItems)
	if err != nil {
		return "", err
	}

	remainingSITDays := 90
	originSITDays, destinationSITDays := 0, 0

	var originSITEntryDate time.Time
	var originSITDepartureDate time.Time
	var destinationSITEntryDate time.Time
	var destinationSITDepartureDate time.Time

	for _, sitPaymentServiceItem := range allSITPaymentServiceItems {
		if (keyData.MTOServiceItem.ReService.Code == models.ReServiceCodeDOFSIT || keyData.MTOServiceItem.ReService.Code == models.ReServiceCodeDOASIT) && (sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOFSIT || sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOASIT) && (sitPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil) {
			return "", fmt.Errorf("a previous Origin SIT MTO Service Item already has a departure date of %v", sitPaymentServiceItem.MTOServiceItem.SITDepartureDate)
		} else if (keyData.MTOServiceItem.ReService.Code == models.ReServiceCodeDDFSIT || keyData.MTOServiceItem.ReService.Code == models.ReServiceCodeDDASIT) && (sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDFSIT || sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDASIT) && sitPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil {
			return "", fmt.Errorf("a previous Destination SIT MTO Service Item already has a departure date of %v", sitPaymentServiceItem.MTOServiceItem.SITDepartureDate)
		}

		if (sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOFSIT || sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOASIT) && sitPaymentServiceItem.MTOServiceItem.SITEntryDate != nil {
			sitEntryDate := sitPaymentServiceItem.MTOServiceItem.SITEntryDate
			if originSITEntryDate.IsZero() || originSITEntryDate == *sitEntryDate {
				originSITEntryDate = *sitEntryDate
			} else {
				return "", fmt.Errorf("a different SIT Entry Date for SIT Origin MTO Service Items already exists: %v", originSITEntryDate)
			}
		} else if (sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDFSIT || sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDASIT) && sitPaymentServiceItem.MTOServiceItem.SITEntryDate != nil {
			sitEntryDate := sitPaymentServiceItem.MTOServiceItem.SITEntryDate
			if destinationSITEntryDate.IsZero() || destinationSITEntryDate == *sitEntryDate {
				destinationSITEntryDate = *sitEntryDate
			} else {
				return "", fmt.Errorf("a different SIT Entry Date for SIT Destination MTO Service Items already exists: %v", destinationSITEntryDate)
			}
		}

		if (sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOFSIT || sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOASIT) && sitPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil {
			sitDepartureDate := sitPaymentServiceItem.MTOServiceItem.SITDepartureDate
			if originSITDepartureDate.IsZero() || originSITDepartureDate == *sitDepartureDate {
				originSITDepartureDate = *sitDepartureDate
			} else {
				return "", fmt.Errorf("a different SIT Departure Date for SIT Origin MTO Service Items already exists: %v", originSITDepartureDate)
			}
		} else if (sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDFSIT || sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDASIT) && sitPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil {
			sitDepartureDate := sitPaymentServiceItem.MTOServiceItem.SITDepartureDate
			if destinationSITDepartureDate.IsZero() || destinationSITDepartureDate == *sitDepartureDate {
				destinationSITDepartureDate = *sitDepartureDate
			} else {
				return "", fmt.Errorf("a different SIT Departure Date for SIT Destination MTO Service Items already exists: %v", destinationSITDepartureDate)
			}
		}

		if sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOFSIT {
			// fmt.Println(string(sitPaymentServiceItem.MTOServiceItem.ReService.Code))
			originSITDays++
			// fmt.Printf("Origin SIT Days requested: %s", strconv.Itoa(originSITDays))
			// fmt.Println()
		} else if sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDOASIT {
			// fmt.Println(string(sitPaymentServiceItem.MTOServiceItem.ReService.Code))
			originSITDays += 29
			// fmt.Printf("Origin SIT Days requested: %s", strconv.Itoa(originSITDays))
			// fmt.Println()
		} else if sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDFSIT {
			// fmt.Println(string(sitPaymentServiceItem.MTOServiceItem.ReService.Code))
			destinationSITDays++
			// fmt.Printf("Destination SIT Days requested: %s", strconv.Itoa(destinationSITDays))
			// fmt.Println()
		} else if sitPaymentServiceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDDASIT {
			// fmt.Println(string(sitPaymentServiceItem.MTOServiceItem.ReService.Code))
			destinationSITDays += 29
			// fmt.Printf("Destination SIT Days requested: %s", strconv.Itoa(destinationSITDays))
			// fmt.Println()
		}
	}

	remainingSITDays = remainingSITDays - originSITDays - destinationSITDays

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Printf("ORIGIN SIT ENTRY DATE: %v", originSITEntryDate)
	fmt.Println()
	fmt.Printf("ORIGIN SIT DEPARTURE DATE: %v", originSITDepartureDate)
	fmt.Println()
	fmt.Printf("DESTINATION SIT ENTRY DATE: %v", destinationSITEntryDate)
	fmt.Println()
	fmt.Printf("DESTINATION SIT DEPARTURE DATE: %v", destinationSITDepartureDate)
	fmt.Println()
	fmt.Println()
	fmt.Println()

	fmt.Printf("Total Origin SIT Days requested: %s", strconv.Itoa(originSITDays))
	fmt.Println()
	fmt.Printf("Total Destination SIT Days requested: %s", strconv.Itoa(destinationSITDays))
	fmt.Println()
	fmt.Printf("Total SIT Days requested: %s", strconv.Itoa(originSITDays+destinationSITDays))
	fmt.Println()
	fmt.Printf("Number of SIT days remaining: %s", strconv.Itoa(remainingSITDays))
	fmt.Println()
	fmt.Println()

	return strconv.Itoa(len(allSITPaymentServiceItems)), nil
}

// Happy path OriginSIT no departure date
// Happy path OriginSIT with departure date
// Happy path destinationSIT no departure date
// Happy path destinationSIT with departure date
// Departure date already exists for OriginSIT
// Departure date already exists for DestinationSIT
// Less than 29 remainingSITDays
// Less than 29 days since sitEntryDate

// Relationship between PaymentServiceItems and MTOServiceItems
// How the lookup is called on a specific MTOServiceItem before that MTOServiceItem has a PaymentServiceItem'
// Getting a basic test stubbed so I can start itteratively testing my code
// The business logic for remaining SIT days - multiple discussions with Jacquie and the documentation that came out of that
// GO syntax - Date
// How to write query to retreive appropriate records
