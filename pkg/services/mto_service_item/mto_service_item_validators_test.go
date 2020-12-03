package mtoserviceitem

import (
	"testing"
	"time"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOServiceItemServiceSuite) TestUpdateMTOServiceItemData() {
	// Set up the data needed for updateMTOServiceItemData obj
	checker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
	now := time.Now()

	// Set up service item models for successful and unsuccessful tests
	successServiceItem := oldServiceItem
	errorServiceItem := oldServiceItem

	// Test successful check for linked IDs
	suite.T().Run("checkLinkedIDs - success", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: successServiceItem, // as-is, should succeed
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkLinkedIDs()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for linked IDs
	suite.T().Run("checkLinkedIDs - failure", func(t *testing.T) {
		fakeUUID := uuid.FromStringOrNil("00010001-0001-0001-0001-000100010001")
		errorServiceItem.MoveTaskOrderID = fakeUUID
		errorServiceItem.MTOShipmentID = &fakeUUID
		errorServiceItem.ReServiceID = fakeUUID

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: errorServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkLinkedIDs()

		suite.NoError(err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Contains(serviceItemData.verrs.Keys(), "moveTaskOrderID")
		suite.Contains(serviceItemData.verrs.Keys(), "mtoShipmentID")
		suite.Contains(serviceItemData.verrs.Keys(), "reServiceID")
	})

	// Test successful check for Prime availability
	suite.T().Run("checkPrimeAvailability - success", func(t *testing.T) {
		oldServiceItemPrime := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
		})
		newServiceItemPrime := oldServiceItemPrime

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := serviceItemData.checkPrimeAvailability()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for Prime availability
	suite.T().Run("checkPrimeAvailability - failure", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  errorServiceItem, // the default errorServiceItem should not be Prime-available
			oldServiceItem:      oldServiceItem,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := serviceItemData.checkPrimeAvailability()

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.NoVerrs(serviceItemData.verrs) // this check doesn't add a validation error
	})

	// Test successful check for non-Prime fields
	suite.T().Run("checkNonPrimeFields - success", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: successServiceItem, // as-is, should succeed because all the values are the same
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkNonPrimeFields()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for non-Prime fields
	suite.T().Run("checkNonPrimeFields - failure", func(t *testing.T) {
		// Update the non-updateable fields:
		errorServiceItem.Status = models.MTOServiceItemStatusApproved
		errorServiceItem.RejectionReason = handlers.FmtString("reason")
		errorServiceItem.ApprovedAt = &now
		errorServiceItem.RejectedAt = &now

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: errorServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkNonPrimeFields()

		suite.NoError(err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Contains(serviceItemData.verrs.Keys(), "status")
		suite.Contains(serviceItemData.verrs.Keys(), "rejectionReason")
		suite.Contains(serviceItemData.verrs.Keys(), "approvedAt")
		suite.Contains(serviceItemData.verrs.Keys(), "rejectedAt")
	})

	// Test successful check for SIT departure service item - not updating SITDepartureDate
	suite.T().Run("checkSITDeparture w/ no SITDepartureDate update - success", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: successServiceItem, // default is not DDDSIT/DOPSIT
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkSITDeparture()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test successful check for SIT departure service item - DDDSIT
	suite.T().Run("checkSITDeparture w/ DDDSIT - success", func(t *testing.T) {
		oldDDDSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		})
		newDDDSIT := oldDDDSIT
		newDDDSIT.SITDepartureDate = &now

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newDDDSIT,
			oldServiceItem:     oldDDDSIT,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkSITDeparture()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for SIT departure service item - not a departure SIT item
	suite.T().Run("checkSITDeparture w/ non-departure SIT - failure", func(t *testing.T) {
		errorServiceItem.SITDepartureDate = &now
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: errorServiceItem, // default is not DDDSIT/DOPSIT
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkSITDeparture()

		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
		suite.NoVerrs(serviceItemData.verrs) // this check doesn't add a validation error
		suite.Contains(err.Error(), "SIT Departure Date may only be manually updated for DDDSIT and DOPSIT service items")
	})

	// Test successful check for service item w/out payment request
	suite.T().Run("checkPaymentRequests - success", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: successServiceItem, // as-is, should succeed
			oldServiceItem:     oldServiceItem,
			db:                 suite.DB(),
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkPaymentRequests()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check service item with an existing payment request
	suite.T().Run("checkPaymentRequests - failure", func(t *testing.T) {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentRequest: paymentRequest,
			MTOServiceItem: errorServiceItem,
		})

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: errorServiceItem,
			oldServiceItem:     oldServiceItem,
			db:                 suite.DB(),
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkPaymentRequests()

		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
		suite.NoVerrs(serviceItemData.verrs) // this check doesn't add a validation error
		suite.Contains(err.Error(), "this service item has an existing payment request and can no longer be updated")
	})

	// Test getVerrs for successful example
	suite.T().Run("getVerrs - success", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: successServiceItem, // as-is, should succeed
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		_ = serviceItemData.checkLinkedIDs() // this test should pass regardless of potential errors here
		_ = serviceItemData.checkNonPrimeFields()
		err := serviceItemData.getVerrs()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test getVerrs for unsuccessful example
	suite.T().Run("getVerrs - failure", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: errorServiceItem, // as-is, should fail
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		_ = serviceItemData.checkLinkedIDs() // this test should pass regardless of potential errors here
		_ = serviceItemData.checkNonPrimeFields()
		err := serviceItemData.getVerrs()

		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.True(serviceItemData.verrs.HasAny())
	})

	// Test setNewMTOServiceItem for successful example
	suite.T().Run("setNewMTOServiceItem - success", func(t *testing.T) {
		successServiceItem.Description = handlers.FmtString("testing update service item validators")
		successServiceItem.Reason = handlers.FmtString("")
		successServiceItem.SITEntryDate = &now
		successServiceItem.ApprovedAt = new(time.Time) // this is the zero time, what we need to nullify the field

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: successServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newServiceItem := serviceItemData.setNewMTOServiceItem()

		suite.NoVerrs(serviceItemData.verrs)
		suite.Nil(newServiceItem.Reason)
		suite.Nil(newServiceItem.ApprovedAt)
		suite.Equal(newServiceItem.SITEntryDate, successServiceItem.SITEntryDate)
		suite.Equal(newServiceItem.Description, successServiceItem.Description)
		suite.NotEqual(newServiceItem.Description, oldServiceItem.Description)
		suite.NotEqual(newServiceItem.Description, serviceItemData.oldServiceItem.Description)
	})
}
