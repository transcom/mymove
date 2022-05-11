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
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testCreateMTOServiceItemQueryBuilder struct {
	fakeCreateOne   func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	fakeFetchOne    func(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	fakeTransaction func(appCtx appcontext.AppContext, fn func(txnAppCtx appcontext.AppContext) error) error
	fakeUpdateOne   func(appCtx appcontext.AppContext, models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testCreateMTOServiceItemQueryBuilder) CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(appCtx, model)
}

func (t *testCreateMTOServiceItemQueryBuilder) UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error) {
	return t.fakeUpdateOne(appCtx, model, eTag)
}

func (t *testCreateMTOServiceItemQueryBuilder) FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
	return t.fakeFetchOne(appCtx, model, filters)
}

func (t *testCreateMTOServiceItemQueryBuilder) Transaction(appCtx appcontext.AppContext, fn func(txnAppCtx appcontext.AppContext) error) error {
	return t.fakeTransaction(appCtx, fn)
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

func (suite *MTOServiceItemServiceSuite) buildValidDDFSITServiceItemWithValidMove() models.MTOServiceItem {
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

func (suite *MTOServiceItemServiceSuite) buildValidDOSHUTServiceItemWithValidMove() models.MTOServiceItem {
	move := testdatagen.MakeAvailableMove(suite.DB())
	reServiceDOSHUT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOSHUT,
		},
	})

	estimatedPrimeWeight := unit.Pound(6000)
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedPrimeWeight,
		},
	})

	estimatedWeight := unit.Pound(4200)
	actualWeight := unit.Pound(4000)

	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID: move.ID,
		MoveTaskOrder:   move,
		ReService:       reServiceDOSHUT,
		MTOShipmentID:   &shipment.ID,
		MTOShipment:     shipment,
		EstimatedWeight: &estimatedWeight,
		ActualWeight:    &actualWeight,
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

	// TESTCASE SCENARIO
	// Under test: CreateMTOServiceItem function
	// Set up:     We create an unapproved move and attempt to create service items on it.
	// Expected outcome:
	//             Error because we cannot create service items before move is approved.

	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	creator := NewMTOServiceItemCreator(builder, moveRouter)
	serviceItemForUnapprovedMove := suite.buildValidServiceItemWithInvalidMove()

	createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemForUnapprovedMove)

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

	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	creator := NewMTOServiceItemCreator(builder, moveRouter)

	// Happy path: If the service item is created successfully it should be returned
	suite.Run("200 Success - SIT Service Item Creation", func() {

		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     We create an approved move and attempt to create DDFSIT service item on it.
		// Expected outcome:
		//             3 SIT items are created, status of move is APPROVALS_REQUESTED

		sitServiceItem := suite.buildValidDDFSITServiceItemWithValidMove()
		sitMove := sitServiceItem.MoveTaskOrder

		createdServiceItems, verrs, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &sitServiceItem)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(createdServiceItems)

		var foundMove models.Move
		err = suite.DB().Find(&foundMove, sitMove.ID)
		suite.NoError(err)

		createdServiceItemList := *createdServiceItems
		suite.Equal(len(createdServiceItemList), 3)
		suite.NotEmpty(createdServiceItemList[2].Dimensions)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, foundMove.Status)
	})

	// Happy path: If the service item is created successfully it should be returned
	suite.Run("200 Success - SHUT Service Item Creation", func() {

		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     We create an approved move and attempt to create DOSHUT service item on it.
		// Expected outcome:
		//             DOSHUT service item is successfully created and returned

		shutServiceItem := suite.buildValidDOSHUTServiceItemWithValidMove()
		shutMove := shutServiceItem.MoveTaskOrder

		createdServiceItem, verrs, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &shutServiceItem)

		var foundMove models.Move
		suite.DB().Find(&foundMove, shutMove.ID)

		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(createdServiceItem)

		createdServiceItemList := *createdServiceItem
		suite.Require().Equal(len(createdServiceItemList), 1)
		suite.Equal(unit.Pound(4200), *createdServiceItemList[0].EstimatedWeight)
		suite.Equal(unit.Pound(4000), *createdServiceItemList[0].ActualWeight)
	})

	// Status default value: If we try to create an mto service item and haven't set the status, we default to SUBMITTED
	suite.Run("success using default status value", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     We create an approved move and attempt to create a service item without a status
		// Expected outcome:
		//             Service item is created and has a status of Submitted

		serviceItemNoStatus := suite.buildValidServiceItemWithNoStatusAndValidMove()
		createdServiceItems, verrs, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemNoStatus)
		suite.NoError(err)
		suite.NoVerrs(verrs)
		suite.NoError(err)
		serviceItemsToCheck := *createdServiceItems
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItemsToCheck[0].Status)
	})

	// If error when trying to create, the create should fail.
	// Bad data which could be IDs that doesn't exist (MoveTaskOrderID or REServiceID)
	suite.Run("creation error", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Mocked:     QueryBuilder
		// Set up:     We create an approved move and mock the query builder to return an error
		// Expected outcome:
		//             Handler returns an error

		sitServiceItem := suite.buildValidDDFSITServiceItemWithValidMove()

		expectedError := "Can't create service item for some reason"
		verrs := validate.NewErrors()
		verrs.Add("test", expectedError)

		fakeCreateOne := func(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
			return verrs, errors.New(expectedError)
		}
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error {
			return nil
		}
		fakeTx := func(appCtx appcontext.AppContext, fn func(txnAppCtx appcontext.AppContext) error) error {
			return fn(appCtx)
		}

		builder := &testCreateMTOServiceItemQueryBuilder{
			fakeCreateOne:   fakeCreateOne,
			fakeFetchOne:    fakeFetchOne,
			fakeTransaction: fakeTx,
		}

		fakeCreateNewBuilder := func() createMTOServiceItemQueryBuilder {
			return builder
		}

		creator := mtoServiceItemCreator{
			builder:          builder,
			createNewBuilder: fakeCreateNewBuilder,
		}

		createdServiceItems, verrs, _ := creator.CreateMTOServiceItem(suite.AppContextForTest(), &sitServiceItem)
		suite.Error(verrs)
		suite.Nil(createdServiceItems)
	})

	// Should return a "NotFoundError" if the MTO ID is nil
	suite.Run("moveID not found", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create service item on a non-existent move ID
		// Expected outcome:
		//             Not found error returned, no new service items created
		notFoundID := uuid.Must(uuid.NewV4())
		serviceItemNoMTO := models.MTOServiceItem{
			MoveTaskOrderID: notFoundID,
		}

		createdServiceItemsNoMTO, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemNoMTO)
		suite.Nil(createdServiceItemsNoMTO)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundID.String())
	})

	// Should return a "NotFoundError" if the reServiceCode passed in isn't found on the table
	suite.Run("reServiceCode not found", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create service item with a nonexistent service code
		// Expected outcome:
		//             Not found error returned, no new service items created

		sitServiceItem := suite.buildValidDDFSITServiceItemWithValidMove()
		sitMove := sitServiceItem.MoveTaskOrder

		fakeCode := models.ReServiceCode("FAKE")
		serviceItemBadCode := models.MTOServiceItem{
			MoveTaskOrderID: sitMove.ID,
			MoveTaskOrder:   sitMove,
			ReService: models.ReService{
				Code: fakeCode,
			},
		}

		createdServiceItemsBadCode, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemBadCode)
		suite.Nil(createdServiceItemsBadCode)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), fakeCode)
	})

	// Should be able to create a service item with code ReServiceCodeMS or ReServiceCodeCS without a shipment,
	// and it should come back as "APPROVED"
	suite.Run("ReServiceCodeCS creation approved", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create an approved move with no shipments. Then create service items for CS or MS.
		// Expected outcome:
		//             Success, CS and MS can be created on moves without shipments.

		move := testdatagen.MakeAvailableMove(suite.DB())
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

		createdServiceItemsCS, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemCS)
		suite.NotNil(createdServiceItemsCS)
		suite.NoError(err)

		createdServiceItemCSList := *createdServiceItemsCS
		suite.Equal(createdServiceItemCSList[0].Status, models.MTOServiceItemStatus("APPROVED"))

		reServiceMS := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		})
		serviceItemMS := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService:       reServiceMS,
		}

		createdServiceItemsMS, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemMS)
		suite.NotNil(createdServiceItemsMS)
		suite.NoError(err)

		createdServiceItemMSList := *createdServiceItemsMS
		suite.Equal(createdServiceItemMSList[0].Status, models.MTOServiceItemStatus("APPROVED"))
	})

	// Should return a "NotFoundError" if the mtoShipmentID isn't linked to the mtoID passed in
	suite.Run("mtoShipmentID not found", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create service item on a shipment that is not related to the move
		// Expected outcome:
		//             Not found error returned, no new service items created

		move := testdatagen.MakeAvailableMove(suite.DB())
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

		createdServiceItemsBadShip, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemBadShip)
		suite.Nil(createdServiceItemsBadShip)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), shipment.ID.String())
		suite.Contains(err.Error(), move.ID.String())
	})

	// If the service item we're trying to create is shuttle service and there is no estimated weight, it fails.
	suite.Run("MTOServiceItemShuttle no prime weight", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDSHUT service item on a shipment without estimated weight
		// Expected outcome:
		//             Conflict error returned, no new service items created

		move := testdatagen.MakeAvailableMove(suite.DB())
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

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemNoWeight)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	setupDDFSITData := func() (models.MTOServiceItemCustomerContact, models.MTOServiceItem) {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		reServiceDDFSIT := testdatagen.MakeDDFSITReService(suite.DB())

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
				Code: reServiceDDFSIT.Code,
			},
		}
		return contact, serviceItemDDFSIT
	}
	// The timeMilitary fields need to be in the correct format.
	suite.Run("Check DDFSIT timeMilitary=HH:MMZ", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDFSIT service item with a bad time "10:30Z"
		// Expected outcome: InvalidInput error returned, no new service items created
		contact, serviceItemDDFSIT := setupDDFSITData()
		contact.TimeMilitary = "10:30Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "timeMilitary")
	})

	suite.Run("timeMilitary=XXMMZ bad hours", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDFSIT service item with a bad time "2645Z"
		// Expected outcome: InvalidInput error returned, no new service items created
		contact, serviceItemDDFSIT := setupDDFSITData()
		contact.TimeMilitary = "2645Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "timeMilitary")
		suite.Contains(err.Error(), "hours must be between 00 and 23")
	})

	suite.Run("timeMilitary=HHXXZ bad minutes", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDFSIT service item with a bad time "2167Z"
		// Expected outcome: InvalidInput error returned, no new service items created
		contact, serviceItemDDFSIT := setupDDFSITData()
		contact.TimeMilitary = "2167Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "timeMilitary")
		suite.Contains(err.Error(), "minutes must be between 00 and 59")
	})

	suite.Run("timeMilitary=HHMMX bad suffix", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDFSIT service item with a bad time "2050M"
		// Expected outcome: InvalidInput error returned, no new service items created
		contact, serviceItemDDFSIT := setupDDFSITData()
		contact.TimeMilitary = "2050M"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "timeMilitary")
		suite.Contains(err.Error(), "must end with 'Z'")
	})

	suite.Run("timeMilitary=HHMMZ success", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDFSIT service item with a correctly formatted time"
		// Expected outcome: Success, service items created.
		contact, serviceItemDDFSIT := setupDDFSITData()
		contact.TimeMilitary = "1405Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contact}
		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)

		suite.NotNil(createdServiceItems)
		suite.NoError(err)
	})
}

func (suite *MTOServiceItemServiceSuite) TestCreateOriginSITServiceItem() {

	// Set up data to use for all Origin SIT Service Item tests
	var reServiceDOASIT models.ReService
	var reServiceDOFSIT models.ReService
	var reServiceDOPSIT models.ReService

	setupTestData := func() models.MTOShipment {
		move := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		reServiceDOASIT = testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
		})

		reServiceDOFSIT = testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		})

		reServiceDOPSIT = testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})
		return mtoShipment
	}

	sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
	sitPostalCode := "99999"
	reason := "lorem ipsum"

	suite.Run("Failure - 422 Cannot create DOFSIT service item with non-null address.ID", func() {

		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DOFSIT service item with a non-null address ID
		// Expected outcome: InvalidInput error returned, no new service items created
		shipment := setupTestData()

		// Create and address where ID != uuid.Nil
		actualPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})

		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceDOFSIT,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)

		createdServiceItems, verr, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(verr)
		suite.IsType(apperror.InvalidInputError{}, err)

	})

	suite.Run("Create DOFSIT service item and auto-create DOASIT, DOPSIT", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DOFSIT service item with a new address
		// Expected outcome: Success, 3 service items created

		// Customer gets new pickup address for SIT Origin Pickup (DOPSIT) which gets added when
		// creating DOFSIT (SIT origin first day).
		shipment := setupTestData()

		// Do not create Address in the database (Assertions.Stub = true) because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		actualPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true})
		actualPickupAddress.ID = uuid.Nil

		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceDOFSIT,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
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

	setupDOFSIT := func(shipment models.MTOShipment) services.MTOServiceItemCreator {
		// Create DOFSIT
		actualPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true})
		actualPickupAddress.ID = uuid.Nil

		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceDOFSIT,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)

		// Successful creation of DOFSIT
		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		return creator
	}

	suite.Run("Create standalone DOASIT item for shipment if existing DOFSIT", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DOFSIT service item successfully
		//             Create DOASIT item on existing DOFSIT
		// Expected outcome: Success, DOASIT item created

		shipment := setupTestData()
		creator := setupDOFSIT(shipment)

		// Create DOASIT
		serviceItemDOASIT := models.MTOServiceItem{
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reServiceDOASIT,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOASIT)

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

	suite.Run("Failure - 422 Create standalone DOASIT item for shipment does not match existing DOFSIT addresses", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DOFSIT service item successfully
		//             Create DOASIT item on existing DOFSIT but with non-matching address
		// Expected outcome: Invalid input error, no service items created

		shipment := setupTestData()
		creator := setupDOFSIT(shipment)

		// Change pickup address
		serviceItemDOASIT := models.MTOServiceItem{
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reServiceDOASIT,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		actualPickupAddress2 := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true})
		existingServiceItem := &serviceItemDOASIT
		existingServiceItem.SITOriginHHGActualAddress = &actualPickupAddress2

		createdServiceItems, verr, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), existingServiceItem)
		suite.Nil(createdServiceItems)
		suite.Error(verr)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Do not create DOFSIT if one already exists for the shipment", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DOFSIT service item successfully
		//             Create another DOFSIT item on the same shipment
		// Expected outcome: Conflict error, no new DOFSIT item created

		shipment := setupTestData()
		creator := setupDOFSIT(shipment)

		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reServiceDOFSIT,
			SITEntryDate:    &sitEntryDate,
			SITPostalCode:   &sitPostalCode,
			Reason:          &reason,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Do not create standalone DOPSIT service item", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create a shipment, then create a DOPSIT item on it
		// Expected outcome: Invalid input error, can't create standalone DOPSIT, no DOPSIT item created

		shipment := setupTestData()

		serviceItemDOPSIT := models.MTOServiceItem{
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reServiceDOPSIT,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOPSIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

	})

	suite.Run("Do not create standalone DOASIT if there is no DOFSIT on shipment", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create a shipment, then create a DOASIT item on it
		// Expected outcome: Invalid input error, can't create standalone DOASIT, no DOASIT item created
		shipment := setupTestData()

		serviceItemDOASIT := models.MTOServiceItem{
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reServiceDOASIT,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOASIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Do not create DOASIT if the DOFSIT ReService Code is bad", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create a shipment, then create a DOFSIT item on it
		//             Create a serviceItem with type DOASIT but a bad reServiceCode
		// Expected outcome: Not found error, can't create DOASIT
		shipment := setupTestData()
		creator := setupDOFSIT(shipment)
		badReService := models.ReService{
			Code: "bad code",
		}

		serviceItemDOASIT := models.MTOServiceItem{
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       badReService,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOASIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
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

	suite.Run("Fail to create DOFSIT service item due to missing SITOriginHHGActualAddress", func() {

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
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}

// TestCreateDestSITServiceItem tests the creation of destination SIT service items
func (suite *MTOServiceItemServiceSuite) TestCreateDestSITServiceItem() {

	setupTestData := func() (models.MTOShipment, services.MTOServiceItemCreator, models.ReService) {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)

		reServiceDDFSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		})
		return shipment, creator, reServiceDDFSIT

	}

	setupAdditionalSIT := func() (models.ReService, models.ReService) {
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
		return reServiceDDASIT, reServiceDDDSIT
	}

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
	suite.Run("Failure - no DDASIT/DDDSIT codes", func() {

		shipment, creator, reServiceDDFSIT := setupTestData()
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

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "service code")
	})

	// Failed creation of DDFSIT because CustomerContacts has invalid data
	suite.Run("Failure - bad CustomerContacts", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		setupAdditionalSIT()

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

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "timeMilitary")
	})

	// Successful creation of DDFSIT service item and the extra DDASIT/DDDSIT items
	suite.Run("Success - DDFSIT creation approved", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		setupAdditionalSIT()

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

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
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
	suite.Run("Failure - duplicate DDFSIT", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		setupAdditionalSIT()

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

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		// Make a second attempt to add a DDFSIT
		createdServiceItems, _, err = creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), fmt.Sprintf("A service item with reServiceCode %s already exists for this move and/or shipment.", models.ReServiceCodeDDFSIT))
	})

	// Successful creation of DDASIT service item
	suite.Run("Success - DDASIT creation approved", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		reServiceDDASIT, _ := setupAdditionalSIT()

		// First create a DDFSIT because it's required to request a DDASIT
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

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		// Then attempt to create a DDASIT
		serviceItemDDASIT := models.MTOServiceItem{
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MTOShipmentID:   &shipment.ID,
			MTOShipment:     shipment,
			ReService:       reServiceDDASIT,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err = creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDASIT)
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
	suite.Run("Failure - DDASIT creation needs DDFSIT", func() {

		// Make the necessary SIT code objects
		reServiceDDASIT, _ := setupAdditionalSIT()
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		})

		// Make a shipment with no DDFSIT
		now := time.Now()
		shipmentNoDDFSIT := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
				Status:             models.MoveStatusAPPROVED,
			},
		})
		serviceItemDDASIT := models.MTOServiceItem{
			MoveTaskOrderID: shipmentNoDDFSIT.MoveTaskOrderID,
			MoveTaskOrder:   shipmentNoDDFSIT.MoveTaskOrder,
			MTOShipmentID:   &shipmentNoDDFSIT.ID,
			MTOShipment:     shipmentNoDDFSIT,
			ReService:       reServiceDDASIT,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		creator := NewMTOServiceItemCreator(builder, moveRouter)
		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDASIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "No matching first-day SIT service item found")
		suite.Contains(err.Error(), shipmentNoDDFSIT.ID.String())
	})

	// Failed creation of DDDSIT service item
	suite.Run("Failure - cannot create DDDSIT", func() {
		shipment, creator, _ := setupTestData()
		_, reServiceDDDSIT := setupAdditionalSIT()

		serviceItemDDDSIT := models.MTOServiceItem{
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MTOShipment:      shipment,
			MTOShipmentID:    &shipment.ID,
			ReService:        reServiceDDDSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDDSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), models.ReServiceCodeDDDSIT)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.NotEmpty(invalidInputError.ValidationErrors)
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "reServiceCode")
	})
}
