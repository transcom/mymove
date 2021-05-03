//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package mtoserviceitem

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testCreateMTOServiceItemQueryBuilder struct {
	fakeCreateOne   func(model interface{}) (*validate.Errors, error)
	fakeFetchOne    func(model interface{}, filters []services.QueryFilter) error
	fakeTransaction func(func(tx *pop.Connection) error) error
	fakeUpdateOne   func(models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testCreateMTOServiceItemQueryBuilder) CreateOne(model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(model)
}

func (t *testCreateMTOServiceItemQueryBuilder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	return t.fakeUpdateOne(model, eTag)
}

func (t *testCreateMTOServiceItemQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	return t.fakeFetchOne(model, filters)
}

func (t *testCreateMTOServiceItemQueryBuilder) Transaction(fn func(tx *pop.Connection) error) error {
	return t.fakeTransaction(fn)
}

func (suite *MTOServiceItemServiceSuite) buildValidServiceItemWithInvalidMove() models.MTOServiceItem {
	// Default move has status DRAFT, which is invalid for this test because
	// service items can only be created if a Move's status is Approved or
	// Approvals Requested
	move := testdatagen.MakeDefaultMove(suite.DB())
	reServiceDDFSIT := testdatagen.MakeDDFSITReService(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	serviceItemForUnapprovedMove := models.MTOServiceItem{
		MoveTaskOrderID: move.ID,
		MoveTaskOrder:   move,
		ReService:       reServiceDDFSIT,
		MTOShipmentID:   &shipment.ID,
		MTOShipment:     shipment,
	}

	return serviceItemForUnapprovedMove
}

func (suite *MTOServiceItemServiceSuite) buildValidServiceItemWithValidMove() models.MTOServiceItem {
	move := testdatagen.MakeAvailableMove(suite.DB())
	dimension := models.MTOServiceItemDimension{
		Type:      models.DimensionTypeItem,
		Length:    12000,
		Height:    12000,
		Width:     12000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reServiceDDFSIT := testdatagen.MakeDDFSITReService(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID: move.ID,
		MoveTaskOrder:   move,
		ReService:       reServiceDDFSIT,
		MTOShipmentID:   &shipment.ID,
		MTOShipment:     shipment,
		Dimensions:      models.MTOServiceItemDimensions{dimension},
		Status:          models.MTOServiceItemStatusSubmitted,
	}

	return serviceItem
}

func (suite *MTOServiceItemServiceSuite) buildValidServiceItemWithNoStatusAndValidMove() models.MTOServiceItem {
	move := testdatagen.MakeAvailableMove(suite.DB())
	dimension := models.MTOServiceItemDimension{
		Type:      models.DimensionTypeItem,
		Length:    12000,
		Height:    12000,
		Width:     12000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{})
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID: move.ID,
		MoveTaskOrder:   move,
		ReService:       reService,
		MTOShipmentID:   &shipment.ID,
		MTOShipment:     shipment,
		Dimensions:      models.MTOServiceItemDimensions{dimension},
	}

	return serviceItem
}

// Should return a message stating that service items can't be created if
// the move is not in approved status.
func (suite *MTOServiceItemServiceSuite) TestCreateMTOServiceItemWithInvalidMove() {
	builder := query.NewQueryBuilder(suite.DB())
	creator := NewMTOServiceItemCreator(builder)
	serviceItemForUnapprovedMove := suite.buildValidServiceItemWithInvalidMove()

	createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemForUnapprovedMove)

	move := serviceItemForUnapprovedMove.MoveTaskOrder
	suite.DB().Find(&move, move.ID)

	var serviceItem models.MTOServiceItem
	suite.DB().Where("move_id = ?", move.ID).First(&serviceItem)

	suite.Nil(createdServiceItems)
	suite.Zero(serviceItem.ID)
	suite.Error(err)
	suite.Contains(err.Error(), "Cannot create service items before a move has been approved")
	suite.Equal(models.MoveStatusDRAFT, move.Status)
}

func (suite *MTOServiceItemServiceSuite) TestCreateMTOServiceItem() {
	serviceItem := suite.buildValidServiceItemWithValidMove()
	move := serviceItem.MoveTaskOrder
	builder := query.NewQueryBuilder(suite.DB())
	creator := NewMTOServiceItemCreator(builder)

	// Happy path: If the service item is created successfully it should be returned
	suite.T().Run("success", func(t *testing.T) {
		createdServiceItems, verrs, err := creator.CreateMTOServiceItem(&serviceItem)

		var foundMove models.Move
		suite.DB().Find(&foundMove, move.ID)

		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(createdServiceItems)

		createdServiceItemList := *createdServiceItems
		suite.Equal(len(createdServiceItemList), 3)
		suite.NotEmpty(createdServiceItemList[2].Dimensions)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, foundMove.Status)
	})

	// Status default value: If we try to create an mto service item and haven't set the status, we default to SUBMITTED
	suite.T().Run("success using default status value", func(t *testing.T) {
		serviceItemNoStatus := suite.buildValidServiceItemWithNoStatusAndValidMove()
		createdServiceItems, verrs, err := creator.CreateMTOServiceItem(&serviceItemNoStatus)
		suite.NoError(err)
		suite.NoVerrs(verrs)
		suite.NoError(err)
		serviceItemsToCheck := *createdServiceItems
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItemsToCheck[0].Status)
	})

	// If error when trying to create, the create should fail.
	// Bad data which could be IDs that doesn't exist (MoveTaskOrderID or REServiceID)
	suite.T().Run("creation error", func(t *testing.T) {
		expectedError := "Can't create service item for some reason"
		verrs := validate.NewErrors()
		verrs.Add("test", expectedError)

		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			return verrs, errors.New(expectedError)
		}
		fakeFetchOne := func(model interface{}, filters []services.QueryFilter) error {
			return nil
		}
		fakeTx := func(fn func(tx *pop.Connection) error) error {
			return fn(&pop.Connection{})
		}

		builder := &testCreateMTOServiceItemQueryBuilder{
			fakeCreateOne:   fakeCreateOne,
			fakeFetchOne:    fakeFetchOne,
			fakeTransaction: fakeTx,
		}

		fakeCreateNewBuilder := func(db *pop.Connection) createMTOServiceItemQueryBuilder {
			return builder
		}

		creator := mtoServiceItemCreator{
			builder:          builder,
			createNewBuilder: fakeCreateNewBuilder,
		}

		createdServiceItems, verrs, _ := creator.CreateMTOServiceItem(&serviceItem)
		suite.Error(verrs)
		suite.Nil(createdServiceItems)
	})

	// Should return a "NotFoundError" if the MTO ID is nil
	suite.T().Run("moveID not found", func(t *testing.T) {
		notFoundID := uuid.Must(uuid.NewV4())
		serviceItemNoMTO := models.MTOServiceItem{
			MoveTaskOrderID: notFoundID,
		}

		createdServiceItemsNoMTO, _, err := creator.CreateMTOServiceItem(&serviceItemNoMTO)
		suite.Nil(createdServiceItemsNoMTO)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundID.String())
	})

	// Should return a "NotFoundError" if the reServiceCode passed in isn't found on the table
	suite.T().Run("reServiceCode not found", func(t *testing.T) {
		fakeCode := models.ReServiceCode("FAKE")
		serviceItemBadCode := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService: models.ReService{
				Code: fakeCode,
			},
		}

		createdServiceItemsBadCode, _, err := creator.CreateMTOServiceItem(&serviceItemBadCode)
		suite.Nil(createdServiceItemsBadCode)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), fakeCode)
	})

	// Should be able to create a service item with code ReServiceCodeMS or ReServiceCodeCS without a shipment,
	// and it should come back as "APPROVED"
	suite.T().Run("ReServiceCodeCS creation approved", func(t *testing.T) {
		reServiceCS := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeCS,
			},
		})
		serviceItemCS := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService:       reServiceCS,
		}

		createdServiceItemsCS, _, err := creator.CreateMTOServiceItem(&serviceItemCS)
		suite.NotNil(createdServiceItemsCS)
		suite.NoError(err)

		createdServiceItemCSList := *createdServiceItemsCS
		suite.Equal(createdServiceItemCSList[0].Status, models.MTOServiceItemStatus("APPROVED"))
	})

	// Should return a "NotFoundError" if the mtoShipmentID passed in isn't found
	// OR if it isn't linked to the mtoID passed in
	suite.T().Run("mtoShipmentID not found", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		reService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "ANY",
			},
		})
		serviceItemBadShip := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			MTOShipmentID:   &shipment.ID,
			MTOShipment:     shipment,
			ReService:       reService,
		}

		createdServiceItemsBadShip, _, err := creator.CreateMTOServiceItem(&serviceItemBadShip)
		suite.Nil(createdServiceItemsBadShip)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), shipment.ID.String())
		suite.Contains(err.Error(), move.ID.String())
	})

	// If the service item we're trying to create is shuttle service and there is no estimated weight, it fails.
	suite.T().Run("MTOServiceItemShuttle no prime weight", func(t *testing.T) {
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		reService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDSHUT,
			},
		})

		serviceItemNoWeight := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reService,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemNoWeight)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	// The timeMilitary fields need to be in the correct format.
	suite.T().Run("timeMilitary formatting for DDFSIT", func(t *testing.T) {
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		contact := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			FirstAvailableDeliveryDate: time.Now(),
		}
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			Status:          models.MTOServiceItemStatusSubmitted,
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		}

		suite.T().Run("timeMilitary=HH:MMZ", func(t *testing.T) {
			contact.TimeMilitary = "10:30Z"
			serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
			createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)

			suite.Nil(createdServiceItems)
			suite.Error(err)
			suite.IsType(services.InvalidInputError{}, err)
			suite.Contains(err.Error(), "timeMilitary")
		})

		suite.T().Run("timeMilitary=XXMMZ bad hours", func(t *testing.T) {
			contact.TimeMilitary = "2645Z"
			serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
			createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)

			suite.Nil(createdServiceItems)
			suite.Error(err)
			suite.IsType(services.InvalidInputError{}, err)
			suite.Contains(err.Error(), "timeMilitary")
			suite.Contains(err.Error(), "hours must be between 00 and 23")
		})

		suite.T().Run("timeMilitary=HHXXZ bad minutes", func(t *testing.T) {
			contact.TimeMilitary = "2167Z"
			serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
			createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)

			suite.Nil(createdServiceItems)
			suite.Error(err)
			suite.IsType(services.InvalidInputError{}, err)
			suite.Contains(err.Error(), "timeMilitary")
			suite.Contains(err.Error(), "minutes must be between 00 and 59")
		})

		suite.T().Run("timeMilitary=HHXXZ bad minutes", func(t *testing.T) {
			contact.TimeMilitary = "2167Z"
			serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
			createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)

			suite.Nil(createdServiceItems)
			suite.Error(err)
			suite.IsType(services.InvalidInputError{}, err)
			suite.Contains(err.Error(), "timeMilitary")
			suite.Contains(err.Error(), "minutes must be between 00 and 59")
		})

		suite.T().Run("timeMilitary=HHMMX bad suffix", func(t *testing.T) {
			contact.TimeMilitary = "2050M"
			serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
			createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)

			suite.Nil(createdServiceItems)
			suite.Error(err)
			suite.IsType(services.InvalidInputError{}, err)
			suite.Contains(err.Error(), "timeMilitary")
			suite.Contains(err.Error(), "must end with 'Z'")
		})

		suite.T().Run("timeMilitary=HHMMZ success", func(t *testing.T) {
			contact := models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				FirstAvailableDeliveryDate: time.Now(),
			}
			serviceItemDDFSIT := models.MTOServiceItem{
				MoveTaskOrderID: move.ID,
				MoveTaskOrder:   move,
				Status:          models.MTOServiceItemStatusSubmitted,
				ReService: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			}

			contact.TimeMilitary = "1405Z"
			serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
			createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)

			suite.NotNil(createdServiceItems)
			suite.NoError(err)
		})
	})
}

func (suite *MTOServiceItemServiceSuite) TestCreateOriginSITServiceItem() {
	// Set up data to use for all Origin SIT Service Item tests
	move := testdatagen.MakeAvailableMove(suite.DB())
	move.Status = models.MoveStatusAPPROVED
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	reServiceDOASIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOASIT,
		},
	})

	reServiceDOFSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOFSIT,
		},
	})

	reServiceDOPSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		},
	})

	serviceItemDOASIT := models.MTOServiceItem{
		MoveTaskOrder:   move,
		MoveTaskOrderID: move.ID,
		MTOShipment:     mtoShipment,
		MTOShipmentID:   &mtoShipment.ID,
		ReService:       reServiceDOASIT,
		Status:          models.MTOServiceItemStatusSubmitted,
	}

	sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
	sitPostalCode := "99999"
	reason := "lorem ipsum"

	suite.T().Run("Failure - 422 Cannot create DOFSIT service item with non-null address.ID", func(t *testing.T) {
		testMove := testdatagen.MakeAvailableMove(suite.DB())
		testMove.Status = models.MoveStatusAPPROVED
		testMtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testMove,
		})

		// Create and address where ID != uuid.Nil
		actualPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})

		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:             testMove,
			MoveTaskOrderID:           testMove.ID,
			MTOShipment:               testMtoShipment,
			MTOShipmentID:             &testMtoShipment.ID,
			ReService:                 reServiceDOFSIT,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, verr, err := creator.CreateMTOServiceItem(&serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(verr)
		suite.IsType(services.InvalidInputError{}, err)

	})

	suite.T().Run("Create DOFSIT service item and auto-create DOASIT, DOPSIT", func(t *testing.T) {
		// Customer gets new pickup address for SIT Origin Pickup (DOPSIT) which gets added when
		// creating DOFSIT (SIT origin first day).

		// Do not create Address in the database (Assertions.Stub = true) because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		actualPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true})
		actualPickupAddress.ID = uuid.Nil

		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:             move,
			MoveTaskOrderID:           move.ID,
			MTOShipment:               mtoShipment,
			MTOShipmentID:             &mtoShipment.ID,
			ReService:                 reServiceDOFSIT,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDOFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		createdServiceItemsList := *createdServiceItems
		suite.Equal(3, len(createdServiceItemsList))

		numDOFSITFound := 0
		numDOASITFound := 0
		numDOPSITFound := 0

		for _, item := range createdServiceItemsList {
			suite.Equal(serviceItemDOFSIT.MoveTaskOrderID, item.MoveTaskOrderID)
			suite.Equal(serviceItemDOFSIT.MTOShipmentID, item.MTOShipmentID)
			suite.Equal(serviceItemDOFSIT.SITEntryDate, item.SITEntryDate)
			suite.Equal(serviceItemDOFSIT.Reason, item.Reason)
			suite.Equal(serviceItemDOFSIT.SITPostalCode, item.SITPostalCode)

			switch item.ReService.Code {
			case models.ReServiceCodeDOFSIT:
				numDOFSITFound++
			case models.ReServiceCodeDOASIT:
				numDOASITFound++
			case models.ReServiceCodeDOPSIT:
				numDOPSITFound++
			}
		}

		suite.Equal(1, numDOFSITFound)
		suite.Equal(1, numDOASITFound)
		suite.Equal(1, numDOPSITFound)
	})

	suite.T().Run("Create standalone DOASIT item for shipment if existing DOFSIT", func(t *testing.T) {
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDOASIT)

		createdDOASITItem := (*createdServiceItems)[0]
		originalDate, _ := sitEntryDate.MarshalText()
		returnedDate, _ := createdDOASITItem.SITEntryDate.MarshalText()

		// Item is created successfully
		suite.NotNil(createdServiceItems)
		suite.NoError(err)
		// Item contains fields copied over from DOFSIT parent
		suite.EqualValues(originalDate, returnedDate)
		suite.EqualValues(*createdDOASITItem.Reason, reason)
		suite.EqualValues(*createdDOASITItem.SITPostalCode, sitPostalCode)
	})

	suite.T().Run("Failure - 422 Create standalone DOASIT item for shipment does not match existing DOFSIT addresses", func(t *testing.T) {
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		// Change pickup address
		actualPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true})
		existingServiceItem := &serviceItemDOASIT
		existingServiceItem.SITOriginHHGActualAddress = &actualPickupAddress

		createdServiceItems, verr, err := creator.CreateMTOServiceItem(existingServiceItem)
		suite.Nil(createdServiceItems)
		suite.Error(verr)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("Do not create DOFSIT if one already exists for the shipment", func(t *testing.T) {
		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
			MTOShipment:     mtoShipment,
			MTOShipmentID:   &mtoShipment.ID,
			ReService:       reServiceDOFSIT,
			SITEntryDate:    &sitEntryDate,
			SITPostalCode:   &sitPostalCode,
			Reason:          &reason,
		}
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	suite.T().Run("Do not create standalone DOPSIT service item", func(t *testing.T) {
		mtoShipmentNoServiceItems := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		serviceItemDOPSIT := models.MTOServiceItem{
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
			MTOShipment:     mtoShipmentNoServiceItems,
			MTOShipmentID:   &mtoShipmentNoServiceItems.ID,
			ReService:       reServiceDOPSIT,
		}

		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDOPSIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

	})

	suite.T().Run("Do not create standalone DOASIT if there is no DOFSIT on shipment", func(t *testing.T) {
		mtoShipmentNoServiceItems := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		serviceItemDOASIT := models.MTOServiceItem{
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
			MTOShipment:     mtoShipmentNoServiceItems,
			MTOShipmentID:   &mtoShipmentNoServiceItems.ID,
			ReService:       reServiceDOASIT,
		}

		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDOASIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Do not create DOASIT if the DOFSIT ReService Code is bad", func(t *testing.T) {
		badReService := models.ReService{
			Code: "bad code",
		}

		serviceItemDOASIT := models.MTOServiceItem{
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
			MTOShipment:     mtoShipment,
			MTOShipmentID:   &mtoShipment.ID,
			ReService:       badReService,
		}

		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDOASIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

}

func (suite *MTOServiceItemServiceSuite) TestCreateOriginSITServiceItemFailToCreateDOFSIT() {
	// Set up data to use for all Origin SIT Service Item tests
	move := testdatagen.MakeAvailableMove(suite.DB())
	move.Status = models.MoveStatusAPPROVED
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	reServiceDOFSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOFSIT,
		},
	})

	sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
	sitPostalCode := "99999"
	reason := "lorem ipsum"

	suite.T().Run("Fail to create DOFSIT service item due to missing SITOriginHHGActualAddress", func(t *testing.T) {

		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
			MTOShipment:     mtoShipment,
			MTOShipmentID:   &mtoShipment.ID,
			ReService:       reServiceDOFSIT,
			SITEntryDate:    &sitEntryDate,
			SITPostalCode:   &sitPostalCode,
			Reason:          &reason,
		}
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})
}

// TestCreateDestSITServiceItem tests the creation of destination SIT service items
func (suite *MTOServiceItemServiceSuite) TestCreateDestSITServiceItem() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
	})
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
	})
	builder := query.NewQueryBuilder(suite.DB())
	creator := NewMTOServiceItemCreator(builder)

	reServiceDDFSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDDFSIT,
		},
	})
	sitEntryDate := time.Now()
	contact1 := models.MTOServiceItemCustomerContact{
		Type:                       models.CustomerContactTypeFirst,
		FirstAvailableDeliveryDate: sitEntryDate,
		TimeMilitary:               "0815Z",
	}
	contact2 := models.MTOServiceItemCustomerContact{
		Type:                       models.CustomerContactTypeSecond,
		FirstAvailableDeliveryDate: sitEntryDate,
		TimeMilitary:               "1430Z",
	}
	var contacts models.MTOServiceItemCustomerContacts
	contacts = append(contacts, contact1, contact2)

	var successfulDDFSIT models.MTOServiceItem // set in the success test for DDFSIT and used in other tests

	// Failed creation of DDFSIT because DDASIT/DDDSIT codes are not found in DB
	suite.T().Run("Failure - no DDASIT/DDDSIT codes", func(t *testing.T) {
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), "service code")
	})

	// These codes will be needed for the following tests:
	reServiceDDASIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDDASIT,
		},
	})
	reServiceDDDSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDDDSIT,
		},
	})

	// Failed creation of DDFSIT because CustomerContacts has invalid data
	suite.T().Run("Failure - bad CustomerContacts", func(t *testing.T) {
		badContact1 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			FirstAvailableDeliveryDate: sitEntryDate,
			TimeMilitary:               "2611B",
		}
		badContact2 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeSecond,
			FirstAvailableDeliveryDate: sitEntryDate,
			TimeMilitary:               "aaaaaaah",
		}
		var badContacts models.MTOServiceItemCustomerContacts
		badContacts = append(badContacts, badContact1, badContact2)

		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: badContacts,
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.Contains(err.Error(), "timeMilitary")
	})

	// Successful creation of DDFSIT service item and the extra DDASIT/DDDSIT items
	suite.T().Run("Success - DDFSIT creation approved", func(t *testing.T) {
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		createdServiceItemList := *createdServiceItems
		suite.Equal(len(createdServiceItemList), 3)

		// check the returned items for the correct data
		numDDASITFound := 0
		numDDDSITFound := 0
		numDDFSITFound := 0
		for _, item := range createdServiceItemList {
			suite.Equal(item.MoveTaskOrderID, serviceItemDDFSIT.MoveTaskOrderID)
			suite.Equal(item.MTOShipmentID, serviceItemDDFSIT.MTOShipmentID)
			suite.Equal(item.SITEntryDate, serviceItemDDFSIT.SITEntryDate)

			if item.ReService.Code == models.ReServiceCodeDDASIT {
				numDDASITFound++
			}
			if item.ReService.Code == models.ReServiceCodeDDDSIT {
				numDDDSITFound++
			}
			if item.ReService.Code == models.ReServiceCodeDDFSIT {
				numDDFSITFound++
				successfulDDFSIT = item
				suite.Equal(len(item.CustomerContacts), len(serviceItemDDFSIT.CustomerContacts))
			}
		}
		suite.Equal(numDDASITFound, 1)
		suite.Equal(numDDDSITFound, 1)
		suite.Equal(numDDFSITFound, 1)
	})

	// Failed creation of DDFSIT because of duplicate service for shipment
	suite.T().Run("Failure - duplicate DDFSIT", func(t *testing.T) {
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
		suite.Contains(err.Error(), "A service item with reServiceCode DDFSIT already exists for this move and/or shipment.")
	})

	// Successful creation of DDASIT service item
	suite.T().Run("Success - DDASIT creation approved", func(t *testing.T) {
		serviceItemDDASIT := models.MTOServiceItem{
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MTOShipmentID:   &shipment.ID,
			MTOShipment:     shipment,
			ReService:       reServiceDDASIT,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDASIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)
		suite.Equal(len(*createdServiceItems), 1)

		createdServiceItemsList := *createdServiceItems
		suite.Equal(createdServiceItemsList[0].ReService.Code, models.ReServiceCodeDDASIT)
		// The time on the date doesn't matter, so let's just check the date:
		suite.Equal(createdServiceItemsList[0].SITEntryDate.Day(), successfulDDFSIT.SITEntryDate.Day())
		suite.Equal(createdServiceItemsList[0].SITEntryDate.Month(), successfulDDFSIT.SITEntryDate.Month())
		suite.Equal(createdServiceItemsList[0].SITEntryDate.Year(), successfulDDFSIT.SITEntryDate.Year())
	})

	// Failed creation of DDASIT service item due to no DDFSIT on shipment
	suite.T().Run("Failure - DDASIT creation needs DDFSIT", func(t *testing.T) {
		shipmentNoDDFSIT := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		serviceItemDDASIT := models.MTOServiceItem{
			MoveTaskOrderID: shipmentNoDDFSIT.MoveTaskOrderID,
			MoveTaskOrder:   shipmentNoDDFSIT.MoveTaskOrder,
			MTOShipmentID:   &shipmentNoDDFSIT.ID,
			MTOShipment:     shipmentNoDDFSIT,
			ReService:       reServiceDDASIT,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDASIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), "No matching first-day SIT service item found")
		suite.Contains(err.Error(), shipmentNoDDFSIT.ID.String())
	})

	// Failed creation of DDDSIT service item
	suite.T().Run("Failure - cannot create DDDSIT", func(t *testing.T) {
		serviceItemDDDSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDDSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDDSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.Contains(err.Error(), "DDDSIT")

		invalidInputError := err.(services.InvalidInputError)
		suite.NotEmpty(invalidInputError.ValidationErrors)
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "reServiceCode")
	})
}
