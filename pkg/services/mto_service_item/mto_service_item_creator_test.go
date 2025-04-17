// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package mtoserviceitem

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
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
	move := factory.BuildMove(suite.DB(), nil, nil)
	reServiceDDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

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
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	dimension := models.MTOServiceItemDimension{
		Type:      models.DimensionTypeItem,
		Length:    12000,
		Height:    12000,
		Width:     12000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reServiceDDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	destAddress := factory.BuildDefaultAddress(suite.DB())

	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID:              move.ID,
		MoveTaskOrder:                move,
		ReService:                    reServiceDDFSIT,
		MTOShipmentID:                &shipment.ID,
		MTOShipment:                  shipment,
		Dimensions:                   models.MTOServiceItemDimensions{dimension},
		Status:                       models.MTOServiceItemStatusSubmitted,
		SITDestinationFinalAddressID: &destAddress.ID,
		SITDestinationFinalAddress:   &destAddress,
	}

	return serviceItem
}

func (suite *MTOServiceItemServiceSuite) buildValidIDFSITServiceItemWithValidMove() models.MTOServiceItem {
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	dimension := models.MTOServiceItemDimension{
		Type:      models.DimensionTypeItem,
		Length:    12000,
		Height:    12000,
		Width:     12000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reServiceIDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIDFSIT)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				MarketCode: models.MarketCodeInternational,
			},
		},
	}, nil)
	destAddress := factory.BuildDefaultAddress(suite.DB())

	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID:              move.ID,
		MoveTaskOrder:                move,
		ReService:                    reServiceIDFSIT,
		MTOShipmentID:                &shipment.ID,
		MTOShipment:                  shipment,
		Dimensions:                   models.MTOServiceItemDimensions{dimension},
		Status:                       models.MTOServiceItemStatusSubmitted,
		SITDestinationFinalAddressID: &destAddress.ID,
		SITDestinationFinalAddress:   &destAddress,
	}

	return serviceItem
}

func (suite *MTOServiceItemServiceSuite) buildValidDOSHUTServiceItemWithValidMove() models.MTOServiceItem {
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	reServiceDOSHUT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOSHUT)

	estimatedPrimeWeight := unit.Pound(6000)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedPrimeWeight,
			},
		},
	}, nil)

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
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	dimension := models.MTOServiceItemDimension{
		Type:      models.DimensionTypeItem,
		Length:    12000,
		Height:    12000,
		Width:     12000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reService := factory.FetchReService(suite.DB(), nil, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

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
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
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
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	// Happy path: If the service item is created successfully it should be returned
	suite.Run("200 Success - Destination SIT Service Item Creation", func() {

		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     We create an approved move and attempt to create DDFSIT service item on it. Includes Dimensions
		//             and a SITDestinationFinalAddress
		// Expected outcome:
		//             4 SIT items are created, status of move is APPROVALS_REQUESTED

		sitServiceItem := suite.buildValidDDFSITServiceItemWithValidMove()
		sitMove := sitServiceItem.MoveTaskOrder
		sitShipment := sitServiceItem.MTOShipment

		createdServiceItems, verrs, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &sitServiceItem)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(createdServiceItems)

		var foundMove models.Move
		err = suite.DB().Find(&foundMove, sitMove.ID)
		suite.NoError(err)

		createdServiceItemList := *createdServiceItems
		suite.Equal(len(createdServiceItemList), 4)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, foundMove.Status)

		numDDFSITFound := 0
		numDDASITFound := 0
		numDDDSITFound := 0
		numDDSFSCFound := 0

		for _, createdServiceItem := range createdServiceItemList {
			// checking that the service item final destination address equals the shipment's final destination address
			suite.Equal(sitShipment.DestinationAddress.StreetAddress1, createdServiceItem.SITDestinationFinalAddress.StreetAddress1)
			suite.Equal(sitShipment.DestinationAddressID, createdServiceItem.SITDestinationFinalAddressID)

			switch createdServiceItem.ReService.Code {
			case models.ReServiceCodeDDFSIT:
				suite.NotEmpty(createdServiceItem.Dimensions)
				numDDFSITFound++
			case models.ReServiceCodeDDASIT:
				numDDASITFound++
			case models.ReServiceCodeDDDSIT:
				numDDDSITFound++
			case models.ReServiceCodeDDSFSC:
				numDDSFSCFound++
			}
		}
		suite.Equal(numDDASITFound, 1)
		suite.Equal(numDDDSITFound, 1)
		suite.Equal(numDDFSITFound, 1)
		suite.Equal(numDDSFSCFound, 1)
	})

	suite.Run("200 Success - International Destination SIT Service Item Creation", func() {

		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     We create an approved move and attempt to create IDFSIT service item on it. Includes Dimensions
		//             and a SITDestinationFinalAddress
		// Expected outcome:
		//             4 SIT items are created, status of move is APPROVALS_REQUESTED

		sitServiceItem := suite.buildValidIDFSITServiceItemWithValidMove()
		sitMove := sitServiceItem.MoveTaskOrder
		sitShipment := sitServiceItem.MTOShipment

		createdServiceItems, verrs, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &sitServiceItem)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(createdServiceItems)

		var foundMove models.Move
		err = suite.DB().Find(&foundMove, sitMove.ID)
		suite.NoError(err)

		createdServiceItemList := *createdServiceItems
		suite.Equal(len(createdServiceItemList), 4)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, foundMove.Status)

		numIDFSITFound := 0
		numIDASITFound := 0
		numIDDSITFound := 0
		numIDSFSCFound := 0

		for _, createdServiceItem := range createdServiceItemList {
			// checking that the service item final destination address equals the shipment's final destination address
			suite.Equal(sitShipment.DestinationAddress.StreetAddress1, createdServiceItem.SITDestinationFinalAddress.StreetAddress1)
			suite.Equal(sitShipment.DestinationAddressID, createdServiceItem.SITDestinationFinalAddressID)

			switch createdServiceItem.ReService.Code {
			case models.ReServiceCodeIDFSIT:
				suite.NotEmpty(createdServiceItem.Dimensions)
				numIDFSITFound++
			case models.ReServiceCodeIDASIT:
				numIDASITFound++
			case models.ReServiceCodeIDDSIT:
				numIDDSITFound++
			case models.ReServiceCodeIDSFSC:
				numIDSFSCFound++
			}
		}
		suite.Equal(numIDASITFound, 1)
		suite.Equal(numIDDSITFound, 1)
		suite.Equal(numIDFSITFound, 1)
		suite.Equal(numIDSFSCFound, 1)
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

		fakeCreateOne := func(_ appcontext.AppContext, _ interface{}) (*validate.Errors, error) {
			return verrs, errors.New(expectedError)
		}
		fakeFetchOne := func(_ appcontext.AppContext, _ interface{}, _ []services.QueryFilter) error {
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

	// Should be able to create a service item with code ReServiceCodeMS or ReServiceCodeCS that uses a shipment's requested pickup date,
	// and it should come back as "APPROVED"
	suite.Run("ReServiceCodeCS & ReServiceCodeMS creation approved", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create an approved move with a shipment. Then create service items for CS & MS.
		// Expected outcome:
		//             Success, CS & MS can be created as long as requested pickup date exists on a shipment

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		startDate := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, time.December, 31, 12, 0, 0, 0, time.UTC)
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		reServiceCS := factory.FetchReServiceByCode(suite.DB(), "CS")
		csTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceCS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&csTaskOrderFee)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		pickupDate := time.Date(2024, time.July, 31, 12, 0, 0, 0, time.UTC)
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &pickupDate,
				},
			},
		}, nil)
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

		reServiceMS := factory.FetchReServiceByCode(suite.DB(), "MS")
		msTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceMS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&msTaskOrderFee)

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

	// Should not be able to create a service item with code ReServiceCodeMS if there is one already created for the move.
	suite.Run("ReServiceCodeMS multiple creation error", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Then create service items for CS or MS. Then try to create again.
		// Expected outcome:
		//             Return empty MTOServiceItems and continue, MS cannot be created if there is one already created for the move.

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		startDate := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, time.December, 31, 12, 0, 0, 0, time.UTC)
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		reServiceMS := factory.FetchReServiceByCode(suite.DB(), "MS")
		msTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceMS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&msTaskOrderFee)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		pickupDate := time.Date(2024, time.July, 31, 12, 0, 0, 0, time.UTC)
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &pickupDate,
				},
			},
		}, nil)

		serviceItemMS := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService:       reServiceMS,
		}

		createdServiceItemsMS, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemMS)
		suite.NotNil(createdServiceItemsMS)
		suite.NoError(err)

		createdServiceItemsMSDupe, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemMS)

		suite.Nil(err)
		suite.NotNil(createdServiceItemsMSDupe)
		suite.Equal(*createdServiceItemsMSDupe, models.MTOServiceItems(nil))
	})

	// Should not be able to create CS or MS service items unless a shipment within the move has a requested pickup date
	suite.Run("ReServiceCodeCS & ReServiceCodeMS creation error due to lack of shipment requested pickup date", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create an approved move with a shipment that does not have a requested pickup date. Then attempt to create service items for CS & MS.
		// Expected outcome:
		//             Error, CS & MS cannot be created unless a shipment within the move has a requested pickup date

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		startDate := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate := time.Date(2020, time.December, 31, 12, 0, 0, 0, time.UTC)
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		reServiceCS := factory.FetchReServiceByCode(suite.DB(), "CS")
		csTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceCS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&csTaskOrderFee)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: nil,
				},
			},
		}, nil)
		shipment.RequestedPickupDate = nil
		suite.MustSave(&shipment)
		serviceItemCS := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService:       reServiceCS,
		}

		createdServiceItemsCS, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemCS)
		suite.Nil(createdServiceItemsCS)
		suite.Error(err)
		suite.Contains(err.Error(), "cannot create fee for service item CS: missing requested pickup date (non-PPMs) or expected departure date (PPMs) for shipment")

		reServiceMS := factory.FetchReServiceByCode(suite.DB(), "MS")
		msTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceMS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&msTaskOrderFee)

		serviceItemMS := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService:       reServiceMS,
		}

		createdServiceItemsMS, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemMS)
		suite.Nil(createdServiceItemsMS)
		suite.Error(err)
		suite.Contains(err.Error(), "cannot create fee for service item MS: missing requested pickup date (non-PPMs) or expected departure date (PPMs) for shipment")
	})

	// Should be able to create CS service item for full PPM that has expected departure date
	suite.Run("ReServiceCodeCS creation for Full PPM", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create an approved move with a PPM shipment that has an expected departure date
		//             Success, CS can be created

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		startDate := time.Date(testdatagen.GHCTestYear, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate := time.Date(testdatagen.GHCTestYear, time.December, 31, 12, 0, 0, 0, time.UTC)
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		reServiceCS := factory.FetchReServiceByCode(suite.DB(), "CS")
		csTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceCS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&csTaskOrderFee)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItemCS := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			ReService:       reServiceCS,
		}

		createdServiceItemsCS, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemCS)
		suite.NotNil(createdServiceItemsCS)
		suite.NoError(err)
	})

	suite.Run("ReServiceCodeCS & ReServiceCodeMS use the correct contract year based on a shipment's requested pickup date", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create an approved move with a shipment that has a requested pickup date. Then create service items for CS & MS.
		// Expected outcome:
		//             Success and the service items should have the correct price based off of the contract year/requested pickup date

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		startDate := time.Date(testdatagen.GHCTestYear, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate := time.Date(testdatagen.GHCTestYear, time.December, 31, 12, 0, 0, 0, time.UTC)
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		contract2 := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		startDate2 := time.Date(2021, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate2 := time.Date(2021, time.December, 31, 12, 0, 0, 0, time.UTC)
		contractYear2 := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract2,
				ContractID:           contract2.ID,
				StartDate:            startDate2,
				EndDate:              endDate2,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		reServiceCS := factory.FetchReServiceByCode(suite.DB(), "CS")
		csTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceCS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&csTaskOrderFee)

		// creating second fee that we will test against
		csTaskOrderFee2 := models.ReTaskOrderFee{
			ContractYearID: contractYear2.ID,
			ServiceID:      reServiceCS.ID,
			PriceCents:     100000,
		}
		suite.MustSave(&csTaskOrderFee2)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		// going to link a shipment that has a requested pickup date falling under the second contract period
		pickupDate := time.Date(2021, time.July, 1, 12, 0, 0, 0, time.UTC)
		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &pickupDate,
				},
			},
		}, nil)
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
		suite.Equal(*createdServiceItemCSList[0].LockedPriceCents, csTaskOrderFee2.PriceCents)

		reServiceMS := factory.FetchReServiceByCode(suite.DB(), "MS")
		msTaskOrderFee := models.ReTaskOrderFee{
			ContractYearID: contractYear.ID,
			ServiceID:      reServiceMS.ID,
			PriceCents:     90000,
		}
		suite.MustSave(&msTaskOrderFee)
		msTaskOrderFee2 := models.ReTaskOrderFee{
			ContractYearID: contractYear2.ID,
			ServiceID:      reServiceMS.ID,
			PriceCents:     100000,
		}
		suite.MustSave(&msTaskOrderFee2)

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
		suite.Equal(*createdServiceItemMSList[0].LockedPriceCents, csTaskOrderFee2.PriceCents)
	})

	// Should return a "NotFoundError" if the mtoShipmentID isn't linked to the mtoID passed in
	suite.Run("mtoShipmentID not found", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create service item on a shipment that is not related to the move
		// Expected outcome:
		//             Not found error returned, no new service items created

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		reService := factory.FetchReServiceByCode(suite.DB(), "ANY")
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
	suite.Run("MTOServiceItemDomesticShuttle no prime weight is okay", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDSHUT service item on a shipment without estimated weight
		// Expected outcome:
		//             Conflict error returned, no new service items created

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDSHUT)

		serviceItemNoWeight := models.MTOServiceItem{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reService,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemNoWeight)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)
	})

	setupDDFSITData := func() (models.MTOServiceItemCustomerContact, models.MTOServiceItemCustomerContact, models.MTOServiceItem) {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		reServiceDDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)

		contactOne := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              time.Now(),
			FirstAvailableDeliveryDate: time.Now(),
		}

		contactTwo := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeSecond,
			DateOfContact:              time.Now(),
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
		return contactOne, contactTwo, serviceItemDDFSIT
	}
	// The timeMilitary fields need to be in the correct format.
	suite.Run("Check DDFSIT timeMilitary=HH:MMZ", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DDFSIT service item with a bad time "10:30Z"
		// Expected outcome: InvalidInput error returned, no new service items created
		contactOne, contactTwo, serviceItemDDFSIT := setupDDFSITData()
		contactOne.TimeMilitary = "10:30Z"
		contactTwo.TimeMilitary = "14:00Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contactOne, contactTwo}
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
		contactOne, contactTwo, serviceItemDDFSIT := setupDDFSITData()
		contactOne.TimeMilitary = "2645Z"
		contactTwo.TimeMilitary = "3625Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contactOne, contactTwo}
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
		contactOne, contactTwo, serviceItemDDFSIT := setupDDFSITData()
		contactOne.TimeMilitary = "2167Z"
		contactTwo.TimeMilitary = "1253Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contactOne, contactTwo}
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
		contactOne, contactTwo, serviceItemDDFSIT := setupDDFSITData()
		contactOne.TimeMilitary = "2050M"
		contactTwo.TimeMilitary = "1224M"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contactOne, contactTwo}
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
		contactOne, contactTwo, serviceItemDDFSIT := setupDDFSITData()
		contactOne.TimeMilitary = "1405Z"
		contactTwo.TimeMilitary = "2013Z"
		serviceItemDDFSIT.CustomerContacts = models.MTOServiceItemCustomerContacts{contactOne, contactTwo}
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
	var reServiceDOSFSC models.ReService

	var reServiceIOFSIT models.ReService

	setupTestData := func() models.MTOShipment {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		reServiceDOASIT = factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOASIT)
		reServiceDOFSIT = factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)
		reServiceDOPSIT = factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOPSIT)
		reServiceDOSFSC = factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOSFSC)

		return mtoShipment
	}

	setupTestInternationalData := func(isOconusPickupAddress bool, isOconusDestinationAddress bool) models.MTOShipment {
		oconusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		conusAddress := factory.BuildAddress(suite.DB(), nil, nil)

		var pickupAddress models.Address
		var destinationAddress models.Address

		if isOconusPickupAddress {
			pickupAddress = oconusAddress
		} else {
			pickupAddress = conusAddress
		}

		if isOconusDestinationAddress {
			destinationAddress = oconusAddress
		} else {
			destinationAddress = conusAddress
		}

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					MarketCode:           models.MarketCodeInternational,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
				},
			},
		}, nil)

		reServiceIOFSIT = factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIOFSIT)

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
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})

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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		createdServiceItems, verr, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(verr)
		suite.IsType(apperror.InvalidInputError{}, err)

	})

	suite.Run("Create DOFSIT service item and auto-create DOASIT, DOPSIT, DOSFSC", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create DOFSIT service item with a new address
		// Expected outcome: Success, 4 service items created

		// Customer gets new pickup address for SIT Origin Pickup (DOPSIT) which gets added when
		// creating DOFSIT (SIT origin first day).
		shipment := setupTestData()

		// Do not create Address in the database (Assertions.Stub = true) because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
		actualPickupAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
		actualPickupAddress.ID = uuid.Nil
		actualPickupAddress.CountryId = &country.ID
		actualPickupAddress.Country = &country

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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		createdServiceItemsList := *createdServiceItems
		suite.Equal(4, len(createdServiceItemsList))

		numDOFSITFound := 0
		numDOASITFound := 0
		numDOPSITFound := 0
		numDOSFSCFound := 0

		for _, item := range createdServiceItemsList {
			suite.Equal(serviceItemDOFSIT.MoveTaskOrderID, item.MoveTaskOrderID)
			suite.Equal(serviceItemDOFSIT.MTOShipmentID, item.MTOShipmentID)
			suite.Equal(serviceItemDOFSIT.SITEntryDate, item.SITEntryDate)
			suite.Equal(serviceItemDOFSIT.Reason, item.Reason)
			suite.Equal(serviceItemDOFSIT.SITPostalCode, item.SITPostalCode)
			suite.Equal(actualPickupAddress.StreetAddress1, item.SITOriginHHGActualAddress.StreetAddress1)
			suite.Equal(actualPickupAddress.ID, *item.SITOriginHHGActualAddressID)

			if item.ReService.Code == models.ReServiceCodeDOPSIT || item.ReService.Code == models.ReServiceCodeDOSFSC {
				suite.Equal(*item.SITDeliveryMiles, 400)
			}

			switch item.ReService.Code {
			case models.ReServiceCodeDOFSIT:
				numDOFSITFound++
			case models.ReServiceCodeDOASIT:
				numDOASITFound++
			case models.ReServiceCodeDOPSIT:
				numDOPSITFound++
			case models.ReServiceCodeDOSFSC:
				numDOSFSCFound++
			}
		}

		suite.Equal(1, numDOFSITFound)
		suite.Equal(1, numDOASITFound)
		suite.Equal(1, numDOPSITFound)
		suite.Equal(1, numDOSFSCFound)
	})

	suite.Run("Create IOFSIT service item and auto-create IOASIT, IOPSIT, IOSFSC", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create IOFSIT service item with a new address
		// Expected outcome: Success, 4 service items created

		// Customer gets new pickup address for SIT Origin Pickup (IOPSIT) which gets added when
		// creating IOFSIT (SIT origin first day).
		shipment := setupTestInternationalData(false, true)

		// Do not create Address in the database (Assertions.Stub = true) because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
		actualPickupAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
		actualPickupAddress.ID = uuid.Nil
		actualPickupAddress.CountryId = &country.ID
		actualPickupAddress.Country = &country

		serviceItemIOFSIT := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceIOFSIT,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(50, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemIOFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		createdServiceItemsList := *createdServiceItems
		suite.Equal(4, len(createdServiceItemsList))

		numIOFSITFound := 0
		numIOASITFound := 0
		numIOPSITFound := 0
		numIOSFSCFound := 0

		for _, item := range createdServiceItemsList {
			suite.Equal(serviceItemIOFSIT.MoveTaskOrderID, item.MoveTaskOrderID)
			suite.Equal(serviceItemIOFSIT.MTOShipmentID, item.MTOShipmentID)
			suite.Equal(serviceItemIOFSIT.SITEntryDate, item.SITEntryDate)
			suite.Equal(serviceItemIOFSIT.Reason, item.Reason)
			suite.Equal(serviceItemIOFSIT.SITPostalCode, item.SITPostalCode)
			suite.Equal(actualPickupAddress.StreetAddress1, item.SITOriginHHGActualAddress.StreetAddress1)
			suite.Equal(actualPickupAddress.ID, *item.SITOriginHHGActualAddressID)

			if item.ReService.Code == models.ReServiceCodeIOPSIT || item.ReService.Code == models.ReServiceCodeIOSFSC {
				suite.Equal(*item.SITDeliveryMiles, 50)
			}

			switch item.ReService.Code {
			case models.ReServiceCodeIOFSIT:
				numIOFSITFound++
			case models.ReServiceCodeIOASIT:
				numIOASITFound++
			case models.ReServiceCodeIOPSIT:
				numIOPSITFound++
			case models.ReServiceCodeIOSFSC:
				numIOSFSCFound++
			}
		}

		suite.Equal(1, numIOFSITFound)
		suite.Equal(1, numIOASITFound)
		suite.Equal(1, numIOPSITFound)
		suite.Equal(1, numIOSFSCFound)
	})

	setupDOFSIT := func(shipment models.MTOShipment) services.MTOServiceItemCreator {
		// Create DOFSIT
		country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
		actualPickupAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
		actualPickupAddress.ID = uuid.Nil
		actualPickupAddress.CountryId = &country.ID
		actualPickupAddress.Country = &country

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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

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

		actualPickupAddress2 := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
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

	suite.Run("Do not create DOFSIT if departure date is after entry date", func() {
		shipment := setupTestData()
		originAddress := factory.BuildAddress(suite.DB(), nil, nil)
		reServiceDOFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)
		serviceItemDOFSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     models.TimePointer(time.Now().AddDate(0, 0, 1)),
					SITDepartureDate: models.TimePointer(time.Now()),
				},
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    originAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
			},
		}, nil)
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		_, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Departure Date (%s) must be after the SIT Entry Date (%s)",
			serviceItemDOFSIT.SITDepartureDate.Format("2006-01-02"),
			serviceItemDOFSIT.SITEntryDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("Do not create DOFSIT if departure date is the same as entry date", func() {
		today := models.TimePointer(time.Now())
		shipment := setupTestData()
		originAddress := factory.BuildAddress(suite.DB(), nil, nil)
		reServiceDOFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)
		serviceItemDOFSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     today,
					SITDepartureDate: today,
				},
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    originAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
			},
		}, nil)
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		_, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Departure Date (%s) must be after the SIT Entry Date (%s)",
			serviceItemDOFSIT.SITDepartureDate.Format("2006-01-02"),
			serviceItemDOFSIT.SITEntryDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOPSIT)

		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

	})

	suite.Run("Do not create standalone DOSFSC service item", func() {
		// TESTCASE SCENARIO
		// Under test: CreateMTOServiceItem function
		// Set up:     Create a shipment, then create a DOSFSC item on it
		// Expected outcome: Invalid input error, can't create standalone DOSFSC, no DOSFSC item created

		shipment := setupTestData()

		serviceItemDOPSIT := models.MTOServiceItem{
			MoveTaskOrder:   shipment.MoveTaskOrder,
			MoveTaskOrderID: shipment.MoveTaskOrderID,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService:       reServiceDOSFSC,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

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

	sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
	sitPostalCode := "99999"
	reason := "lorem ipsum"

	suite.Run("Fail to create DOFSIT service item due to missing SITOriginHHGActualAddress", func() {
		// Set up data to use for all Origin SIT Service Item tests
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		move.Status = models.MoveStatusAPPROVED
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		reServiceDOFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDOFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})
}

// TestCreateDestSITServiceItem tests the creation of destination SIT service items
func (suite *MTOServiceItemServiceSuite) TestCreateDestSITServiceItem() {

	setupTestData := func() (models.MTOShipment, services.MTOServiceItemCreator, models.ReService) {
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		reServiceDDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
		return shipment, creator, reServiceDDFSIT

	}

	setupTestInternationalData := func() (models.MTOShipment, services.MTOServiceItemCreator, models.ReService) {
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					MarketCode:           models.MarketCodeInternational,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
				},
			},
		}, nil)
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(125, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		reServiceIDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIDFSIT)
		return shipment, creator, reServiceIDFSIT
	}

	setupAdditionalSIT := func() (models.ReService, models.ReService, models.ReService) {
		// These codes will be needed for the following tests:
		reServiceDDASIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDASIT)
		reServiceDDDSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDDSIT)
		reServiceDDSFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDSFSC)
		return reServiceDDASIT, reServiceDDDSIT, reServiceDDSFSC
	}

	setupAdditionalInternationalSIT := func() (models.ReService, models.ReService, models.ReService) {
		// These codes will be needed for the following tests:
		reServiceIDASIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIDASIT)
		reServiceIDDSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIDDSIT)
		reServiceIDSFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIDSFSC)
		return reServiceIDASIT, reServiceIDDSIT, reServiceIDSFSC
	}

	getCustomerContacts := func() models.MTOServiceItemCustomerContacts {
		deliveryDate := time.Now()
		attemptedContact := time.Now()
		contact1 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              attemptedContact,
			FirstAvailableDeliveryDate: deliveryDate,
			TimeMilitary:               "0815Z",
		}
		contact2 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeSecond,
			DateOfContact:              attemptedContact,
			FirstAvailableDeliveryDate: deliveryDate,
			TimeMilitary:               "1430Z",
		}
		var contacts models.MTOServiceItemCustomerContacts
		contacts = append(contacts, contact1, contact2)
		return contacts
	}

	convertCustomerIDsToFindTestMap := func(contacts models.MTOServiceItemCustomerContacts) map[uuid.UUID]bool {
		customerContactIDMap := make(map[uuid.UUID]bool, len(contacts))
		// load all known customer IDs into map
		for _, contact := range contacts {
			customerContactIDMap[contact.ID] = true
		}
		return customerContactIDMap
	}

	sitEntryDate := time.Now().AddDate(0, 0, 1)
	sitDepartureDate := sitEntryDate.AddDate(0, 0, 7)
	attemptedContact := time.Now()

	// Successful creation of DDFSIT MTO service item.
	suite.Run("Success - Creation of DDFSIT MTO Service Item", func() {

		shipment, creator, reServiceDDFSIT := setupTestData()
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: getCustomerContacts(),
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		_, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.NoError(err)
	})

	// Failed creation of DDFSIT because CustomerContacts has invalid data
	suite.Run("Failure - bad CustomerContacts", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		setupAdditionalSIT()

		badContact1 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              attemptedContact,
			FirstAvailableDeliveryDate: sitEntryDate,
			TimeMilitary:               "2611B",
		}
		badContact2 := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeSecond,
			DateOfContact:              attemptedContact,
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
	suite.Run("Success - DDFSIT creation approved - no SITDestinationFinalAddress", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		setupAdditionalSIT()

		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			SITDepartureDate: &sitDepartureDate,
			CustomerContacts: getCustomerContacts(),
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		createdServiceItemList := *createdServiceItems
		suite.Equal(len(createdServiceItemList), 4)

		// check the returned items for the correct data
		numDDASITFound := 0
		numDDDSITFound := 0
		numDDFSITFound := 0
		numDDSFSCFound := 0
		for _, item := range createdServiceItemList {
			suite.Equal(item.MoveTaskOrderID, serviceItemDDFSIT.MoveTaskOrderID)
			suite.Equal(item.MTOShipmentID, serviceItemDDFSIT.MTOShipmentID)
			suite.Equal(item.SITEntryDate, serviceItemDDFSIT.SITEntryDate)
			suite.Equal(item.SITDepartureDate, serviceItemDDFSIT.SITDepartureDate)

			if item.ReService.Code == models.ReServiceCodeDDDSIT || item.ReService.Code == models.ReServiceCodeDDSFSC {
				suite.Equal(*item.SITDeliveryMiles, 400)
			}

			if item.ReService.Code == models.ReServiceCodeDDASIT {
				numDDASITFound++
			}
			if item.ReService.Code == models.ReServiceCodeDDDSIT {
				numDDDSITFound++
			}
			if item.ReService.Code == models.ReServiceCodeDDFSIT {
				numDDFSITFound++
				suite.Equal(len(item.CustomerContacts), len(serviceItemDDFSIT.CustomerContacts))
			}
			if item.ReService.Code == models.ReServiceCodeDDDSIT {
				numDDSFSCFound++
			}
		}
		suite.Equal(numDDASITFound, 1)
		suite.Equal(numDDDSITFound, 1)
		suite.Equal(numDDFSITFound, 1)
		suite.Equal(numDDSFSCFound, 1)

		// We create one set of customer contacts and attach them to each destination service item.
		// This portion verifies that.
		customerContactIDMap := convertCustomerIDsToFindTestMap(serviceItemDDFSIT.CustomerContacts)
		// Verify there are only 2 customers created
		suite.Equal(len(customerContactIDMap), 2)
		for _, createdServiceItem := range createdServiceItemList {
			for _, item := range createdServiceItem.CustomerContacts {
				// remove ID from map to denote it was found
				delete(customerContactIDMap, item.ID)
			}
		}
		// found all expected IDs. expect empty map
		suite.Equal(len(customerContactIDMap), 0)
	})

	// Successful creation of IDFSIT service item and the extra IDASIT/IDDSIT items
	suite.Run("Success - IDFSIT creation approved - no SITDestinationFinalAddress", func() {
		shipment, creator, reServiceIDFSIT := setupTestInternationalData()
		setupAdditionalInternationalSIT()

		serviceItemIDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceIDFSIT,
			SITEntryDate:     &sitEntryDate,
			SITDepartureDate: &sitDepartureDate,
			CustomerContacts: getCustomerContacts(),
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemIDFSIT)
		suite.NotNil(createdServiceItems)
		suite.NoError(err)

		createdServiceItemList := *createdServiceItems
		suite.Equal(len(createdServiceItemList), 4)

		// check the returned items for the correct data
		numIDASITFound := 0
		numIDDSITFound := 0
		numIDFSITFound := 0
		numIDSFSCFound := 0
		for _, item := range createdServiceItemList {
			suite.Equal(item.MoveTaskOrderID, serviceItemIDFSIT.MoveTaskOrderID)
			suite.Equal(item.MTOShipmentID, serviceItemIDFSIT.MTOShipmentID)
			suite.Equal(item.SITEntryDate, serviceItemIDFSIT.SITEntryDate)
			suite.Equal(item.SITDepartureDate, serviceItemIDFSIT.SITDepartureDate)

			suite.Equal(item.SITDestinationOriginalAddressID, serviceItemIDFSIT.SITDestinationOriginalAddressID)
			suite.Equal(item.SITDestinationFinalAddressID, serviceItemIDFSIT.SITDestinationFinalAddressID)

			if item.ReService.Code == models.ReServiceCodeIDDSIT || item.ReService.Code == models.ReServiceCodeIDSFSC {
				// if this fails check the mock in the setupdata func and/or if destination address is OCONUS
				suite.Equal(*item.SITDeliveryMiles, 125)
			}

			if item.ReService.Code == models.ReServiceCodeIDASIT {
				numIDASITFound++
			}
			if item.ReService.Code == models.ReServiceCodeIDDSIT {
				numIDDSITFound++
			}
			if item.ReService.Code == models.ReServiceCodeIDFSIT {
				numIDFSITFound++
				suite.Equal(len(item.CustomerContacts), len(serviceItemIDFSIT.CustomerContacts))
			}
			if item.ReService.Code == models.ReServiceCodeIDDSIT {
				numIDSFSCFound++
			}
		}
		suite.Equal(numIDASITFound, 1)
		suite.Equal(numIDDSITFound, 1)
		suite.Equal(numIDFSITFound, 1)
		suite.Equal(numIDSFSCFound, 1)

		// We create one set of customer contacts and attach them to each destination service item.
		// This portion verifies that.
		customerContactIDMap := convertCustomerIDsToFindTestMap(serviceItemIDFSIT.CustomerContacts)
		// Verify there are only 2 customers created
		suite.Equal(len(customerContactIDMap), 2)
		for _, createdServiceItem := range createdServiceItemList {
			for _, item := range createdServiceItem.CustomerContacts {
				// remove ID from map to denote it was found
				delete(customerContactIDMap, item.ID)
			}
		}

		// found all expected IDs. expect empty map
		suite.Equal(len(customerContactIDMap), 0)
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
			CustomerContacts: getCustomerContacts(),
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

	suite.Run("Failure - SIT entry date is before FADD for DDFSIT creation", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		setupAdditionalSIT()

		sitEntryDateBeforeToday := time.Now().AddDate(0, 0, -1)

		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDateBeforeToday,
			CustomerContacts: getCustomerContacts(),
			Status:           models.MTOServiceItemStatusSubmitted,
		}

		// Make a second attempt to add a DDFSIT
		serviceItem, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.Nil(serviceItem)
		suite.Error(err)
		suite.IsType(apperror.UnprocessableEntityError{}, err)
		expectedError := fmt.Sprintf(
			"the SIT Entry Date (%s) cannot be before the First Available Delivery Date (%s)",
			serviceItemDDFSIT.SITEntryDate.Format("2006-01-02"),
			serviceItemDDFSIT.CustomerContacts[0].FirstAvailableDeliveryDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("Do not create DDFSIT if departure date is after entry date", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		serviceItemDDFSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     models.TimePointer(time.Now().AddDate(0, 0, 1)),
					SITDepartureDate: models.TimePointer(time.Now()),
				},
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)
		_, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Departure Date (%s) must be after the SIT Entry Date (%s)",
			serviceItemDDFSIT.SITDepartureDate.Format("2006-01-02"),
			serviceItemDDFSIT.SITEntryDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("Do not create DDFSIT if departure date is the same as entry date", func() {
		today := models.TimePointer(time.Now())
		shipment, creator, reServiceDDFSIT := setupTestData()
		serviceItemDDFSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     today,
					SITDepartureDate: today,
				},
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)
		_, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDFSIT)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Departure Date (%s) must be after the SIT Entry Date (%s)",
			serviceItemDDFSIT.SITDepartureDate.Format("2006-01-02"),
			serviceItemDDFSIT.SITEntryDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})

	// Successful creation of DDASIT service item
	suite.Run("Success - DDASIT creation approved", func() {
		shipment, creator, reServiceDDFSIT := setupTestData()
		reServiceDDASIT, _, _ := setupAdditionalSIT()

		// First create a DDFSIT because it's required to request a DDASIT
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: getCustomerContacts(),
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
		suite.Equal(createdServiceItemsList[0].SITEntryDate.Day(), sitEntryDate.Day())
		suite.Equal(createdServiceItemsList[0].SITEntryDate.Month(), sitEntryDate.Month())
		suite.Equal(createdServiceItemsList[0].SITEntryDate.Year(), sitEntryDate.Year())
	})

	// Failed creation of DDASIT service item due to no DDFSIT on shipment
	suite.Run("Failure - DDASIT creation needs DDFSIT", func() {

		// Make the necessary SIT code objects
		reServiceDDASIT, _, _ := setupAdditionalSIT()
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)

		// Make a shipment with no DDFSIT
		now := time.Now()
		shipmentNoDDFSIT := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
					Status:             models.MoveStatusAPPROVED,
				},
			},
		}, nil)
		serviceItemDDASIT := models.MTOServiceItem{
			MoveTaskOrderID: shipmentNoDDFSIT.MoveTaskOrderID,
			MoveTaskOrder:   shipmentNoDDFSIT.MoveTaskOrder,
			MTOShipmentID:   &shipmentNoDDFSIT.ID,
			MTOShipment:     shipmentNoDDFSIT,
			ReService:       reServiceDDASIT,
			Status:          models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
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
		_, reServiceDDDSIT, _ := setupAdditionalSIT()

		serviceItemDDDSIT := models.MTOServiceItem{
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MTOShipment:      shipment,
			MTOShipmentID:    &shipment.ID,
			ReService:        reServiceDDDSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: getCustomerContacts(),
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

	// Failed creation of DDSFSC service item
	suite.Run("Failure - cannot create DDSFSC", func() {
		shipment, creator, _ := setupTestData()
		_, _, reServiceDDSFSC := setupAdditionalSIT()

		serviceItemDDSFSC := models.MTOServiceItem{
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MTOShipment:      shipment,
			MTOShipmentID:    &shipment.ID,
			ReService:        reServiceDDSFSC,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: getCustomerContacts(),
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &serviceItemDDSFSC)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), models.ReServiceCodeDDSFSC)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.NotEmpty(invalidInputError.ValidationErrors)
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "reServiceCode")
	})

	suite.Run("Failure - cannot create domestic service item international domestic shipment", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		dimension := models.MTOServiceItemDimension{
			Type:      models.DimensionTypeItem,
			Length:    12000,
			Height:    12000,
			Width:     12000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// setup domestic shipment
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					MarketCode: models.MarketCodeInternational,
				},
			},
		}, nil)
		destAddress := factory.BuildDefaultAddress(suite.DB())

		// setup international service item. must fail validation for a domestic shipment
		reServiceDDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)
		internationalServiceItem := models.MTOServiceItem{
			MoveTaskOrderID:              move.ID,
			MoveTaskOrder:                move,
			ReService:                    reServiceDDFSIT,
			MTOShipmentID:                &shipment.ID,
			MTOShipment:                  shipment,
			Dimensions:                   models.MTOServiceItemDimensions{dimension},
			Status:                       models.MTOServiceItemStatusSubmitted,
			SITDestinationFinalAddressID: &destAddress.ID,
			SITDestinationFinalAddress:   &destAddress,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &internationalServiceItem)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		suite.Contains(err.Error(), "cannot create domestic service items for international shipment")
	})

	suite.Run("Failure - cannot create international service item for domestic shipment", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		dimension := models.MTOServiceItemDimension{
			Type:      models.DimensionTypeItem,
			Length:    12000,
			Height:    12000,
			Width:     12000,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// setup domestic shipment
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					MarketCode: models.MarketCodeDomestic,
				},
			},
		}, nil)
		destAddress := factory.BuildDefaultAddress(suite.DB())

		// setup international service item. must fail validation for a domestic shipment
		reServiceIDFSIT := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIDFSIT)
		internationalServiceItem := models.MTOServiceItem{
			MoveTaskOrderID:              move.ID,
			MoveTaskOrder:                move,
			ReService:                    reServiceIDFSIT,
			MTOShipmentID:                &shipment.ID,
			MTOShipment:                  shipment,
			Dimensions:                   models.MTOServiceItemDimensions{dimension},
			Status:                       models.MTOServiceItemStatusSubmitted,
			SITDestinationFinalAddressID: &destAddress.ID,
			SITDestinationFinalAddress:   &destAddress,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		createdServiceItems, _, err := creator.CreateMTOServiceItem(suite.AppContextForTest(), &internationalServiceItem)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		suite.Contains(err.Error(), "cannot create international service items for domestic shipment")
	})
}

func (suite *MTOServiceItemServiceSuite) TestPriceEstimator() {
	suite.Run("Calcuating price estimated on creation for HHG ", func() {
		setupTestData := func() models.MTOShipment {
			// Set up data to use for all Origin SIT Service Item tests

			move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
			estimatedPrimeWeight := unit.Pound(6000)
			pickupDate := time.Date(2024, time.July, 31, 12, 0, 0, 0, time.UTC)
			pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
			deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

			mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    pickupAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.PickupAddress,
				},
				{
					Model:    deliveryAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.MTOShipment{
						PrimeEstimatedWeight: &estimatedPrimeWeight,
						RequestedPickupDate:  &pickupDate,
					},
				},
			}, nil)

			return mtoShipment
		}

		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reServiceCodeDPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)
		reServiceCodeDDP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)
		reServiceCodeDUPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)
		reServiceCodeDLH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
		reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)
		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)

		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		reason := "lorem ipsum"

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					Name:                 "Test Contract Year",
					EscalationCompounded: 1.125,
					StartDate:            testdatagen.ContractStartDate,
					EndDate:              testdatagen.ContractEndDate,
				},
			})

		serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "945",
					ServicesSchedule: 1,
				},
			})

		serviceAreaDest := testdatagen.MakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "503",
					ServicesSchedule: 1,
				},
			})

		serviceAreaPriceDOP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDOP.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(1234),
		}

		serviceAreaPriceDPK := factory.FetchOrMakeDomesticOtherPrice(suite.DB(), []factory.Customization{
			{
				Model: models.ReDomesticOtherPrice{
					ContractID:   contractYear.Contract.ID,
					ServiceID:    reServiceCodeDPK.ID,
					IsPeakPeriod: true,
					Schedule:     1,
					PriceCents:   unit.Cents(121),
				},
			},
		}, nil)

		serviceAreaPriceDDP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDDP.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceAreaDest.ID,
			PriceCents:            unit.Cents(482),
		}

		serviceAreaPriceDUPK := factory.FetchOrMakeDomesticOtherPrice(suite.DB(), []factory.Customization{
			{
				Model: models.ReDomesticOtherPrice{
					ContractID:   contractYear.Contract.ID,
					ServiceID:    reServiceCodeDUPK.ID,
					IsPeakPeriod: true,
					Schedule:     1,
					PriceCents:   unit.Cents(945),
				},
			},
		}, nil)

		serviceAreaPriceDLH := models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			WeightLower:           500,
			WeightUpper:           10000,
			MilesLower:            1,
			MilesUpper:            10000,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceMillicents:       unit.Millicents(482),
		}

		serviceAreaPriceDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDSH.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(999),
		}

		testdatagen.FetchOrMakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
			GHCDieselFuelPrice: models.GHCDieselFuelPrice{
				FuelPriceInMillicents: unit.Millicents(281400),
				PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
				EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:               time.Date(2025, time.March, 17, 0, 0, 0, 0, time.UTC),
			},
		})

		suite.MustSave(&serviceAreaPriceDOP)
		suite.MustSave(&serviceAreaPriceDPK)
		suite.MustSave(&serviceAreaPriceDDP)
		suite.MustSave(&serviceAreaPriceDUPK)
		suite.MustSave(&serviceAreaPriceDLH)
		suite.MustSave(&serviceAreaPriceDSH)

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceArea,
				Zip3:                "945",
			},
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceAreaDest,
				Zip3:                "503",
			},
		})

		shipment := setupTestData()
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		serviceItemDOP := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDOP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDPK := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDDP := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDDP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDUPK := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDUPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDLH := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDLH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDSH := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDSH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemFSC := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeFSC,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		dopEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDOP, shipment)
		suite.Equal(unit.Cents(66330), dopEstimatedPriceInCents)

		dpkEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDPK, shipment)
		suite.Equal(unit.Cents(586080), dpkEstimatedPriceInCents)

		ddpEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDDP, shipment)
		suite.Equal(unit.Cents(45870), ddpEstimatedPriceInCents)

		dupkEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDUPK, shipment)
		suite.Equal(unit.Cents(47652), dupkEstimatedPriceInCents)

		dlhEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDLH, shipment)
		suite.Equal(unit.Cents(13437600), dlhEstimatedPriceInCents)

		dshEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDSH, shipment)
		suite.Equal(unit.Cents(10929600), dshEstimatedPriceInCents)

		fscEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemFSC, shipment)
		// negative because we are using 2020 fuel rates
		suite.Equal(unit.Cents(-168), fscEstimatedPriceInCents)
	})

	suite.Run("Calcuating price estimated on creation for NTS shipment ", func() {
		setupTestData := func() models.MTOShipment {
			// Set up data to use for all Origin SIT Service Item tests

			move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
			estimatedPrimeWeight := unit.Pound(6000)
			pickupDate := time.Date(2024, time.July, 31, 12, 0, 0, 0, time.UTC)
			pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
			deliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})

			mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    pickupAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.PickupAddress,
				},
				{
					Model:    deliveryAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model: models.MTOShipment{
						PrimeEstimatedWeight: &estimatedPrimeWeight,
						RequestedPickupDate:  &pickupDate,
						ShipmentType:         models.MTOShipmentTypeHHGOutOfNTS,
					},
				},
			}, nil)

			return mtoShipment
		}

		reServiceCodeDOP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
		reServiceCodeDPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDPK)
		reServiceCodeDDP := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDP)
		reServiceCodeDUPK := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUPK)
		reServiceCodeDLH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
		reServiceCodeDSH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDSH)
		reServiceCodeFSC := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)

		startDate := time.Now().AddDate(-1, 0, 0)
		endDate := startDate.AddDate(1, 1, 1)
		sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "99999"
		reason := "lorem ipsum"

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					Name:                 "Test Contract Year",
					EscalationCompounded: 1.125,
					StartDate:            startDate,
					EndDate:              endDate,
				},
			})

		serviceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "945",
					ServicesSchedule: 1,
				},
			})

		serviceAreaDest := testdatagen.MakeReDomesticServiceArea(suite.DB(),
			testdatagen.Assertions{
				ReDomesticServiceArea: models.ReDomesticServiceArea{
					Contract:         contractYear.Contract,
					ServiceArea:      "503",
					ServicesSchedule: 1,
				},
			})

		serviceAreaPriceDOP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDOP.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(1234),
		}

		serviceAreaPriceDPK := factory.FetchOrMakeDomesticOtherPrice(suite.DB(), []factory.Customization{
			{
				Model: models.ReDomesticOtherPrice{
					ContractID:   contractYear.Contract.ID,
					ServiceID:    reServiceCodeDPK.ID,
					IsPeakPeriod: true,
					Schedule:     1,
					PriceCents:   unit.Cents(121),
				},
			},
		}, nil)

		serviceAreaPriceDDP := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDDP.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceAreaDest.ID,
			PriceCents:            unit.Cents(482),
		}

		serviceAreaPriceDUPK := factory.FetchOrMakeDomesticOtherPrice(suite.DB(), []factory.Customization{
			{
				Model: models.ReDomesticOtherPrice{
					ContractID:   contractYear.Contract.ID,
					ServiceID:    reServiceCodeDUPK.ID,
					IsPeakPeriod: true,
					Schedule:     1,
					PriceCents:   unit.Cents(945),
				},
			},
		}, nil)

		serviceAreaPriceDLH := models.ReDomesticLinehaulPrice{
			ContractID:            contractYear.Contract.ID,
			WeightLower:           500,
			WeightUpper:           10000,
			MilesLower:            1,
			MilesUpper:            10000,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceMillicents:       unit.Millicents(482),
		}

		serviceAreaPriceDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            contractYear.Contract.ID,
			ServiceID:             reServiceCodeDSH.ID,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: serviceArea.ID,
			PriceCents:            unit.Cents(999),
		}

		testdatagen.FetchOrMakeGHCDieselFuelPrice(suite.DB(), testdatagen.Assertions{
			GHCDieselFuelPrice: models.GHCDieselFuelPrice{
				FuelPriceInMillicents: unit.Millicents(281400),
				PublicationDate:       time.Date(2020, time.March, 9, 0, 0, 0, 0, time.UTC),
				EffectiveDate:         time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC),
				EndDate:               time.Date(2025, time.March, 17, 0, 0, 0, 0, time.UTC),
			},
		})

		suite.MustSave(&serviceAreaPriceDOP)
		suite.MustSave(&serviceAreaPriceDPK)
		suite.MustSave(&serviceAreaPriceDDP)
		suite.MustSave(&serviceAreaPriceDUPK)
		suite.MustSave(&serviceAreaPriceDLH)
		suite.MustSave(&serviceAreaPriceDSH)

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceArea,
				Zip3:                "945",
			},
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            contract,
				ContractID:          contract.ID,
				DomesticServiceArea: serviceAreaDest,
				Zip3:                "503",
			},
		})

		shipment := setupTestData()
		actualPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		serviceItemDOP := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDOP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDPK := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDDP := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDDP,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDUPK := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDUPK,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDLH := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDLH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemDSH := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeDSH,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		serviceItemFSC := models.MTOServiceItem{
			MoveTaskOrder:             shipment.MoveTaskOrder,
			MoveTaskOrderID:           shipment.MoveTaskOrderID,
			MTOShipment:               shipment,
			MTOShipmentID:             &shipment.ID,
			ReService:                 reServiceCodeFSC,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			Reason:                    &reason,
			SITOriginHHGActualAddress: &actualPickupAddress,
			Status:                    models.MTOServiceItemStatusSubmitted,
		}

		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(800, nil)
		creator := NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

		dopEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDOP, shipment)
		suite.Equal(unit.Cents(66330), dopEstimatedPriceInCents)

		dpkEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDPK, shipment)
		suite.Equal(unit.Cents(586080), dpkEstimatedPriceInCents)

		ddpEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDDP, shipment)
		suite.Equal(unit.Cents(45870), ddpEstimatedPriceInCents)

		dupkEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDUPK, shipment)
		suite.Equal(unit.Cents(47652), dupkEstimatedPriceInCents)

		dlhEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDLH, shipment)
		suite.Equal(unit.Cents(29589120), dlhEstimatedPriceInCents)

		dshEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemDSH, shipment)
		suite.Equal(unit.Cents(21859200), dshEstimatedPriceInCents)

		fscEstimatedPriceInCents, _ := creator.FindEstimatedPrice(suite.AppContextForTest(), &serviceItemFSC, shipment)
		// negative because we are using 2020 fuel rate
		suite.Equal(unit.Cents(-335), fscEstimatedPriceInCents)
	})

}
func (suite *MTOServiceItemServiceSuite) TestGetAdjustedWeight() {
	suite.Run("If no weight is provided", func() {
		var incomingWeight unit.Pound
		adjustedWeight := GetAdjustedWeight(incomingWeight, false)
		suite.Equal(unit.Pound(0), *adjustedWeight)
	})
	suite.Run("If a weight of 0 is provided", func() {
		incomingWeight := unit.Pound(0)
		adjustedWeight := GetAdjustedWeight(incomingWeight, false)
		suite.Equal(unit.Pound(0), *adjustedWeight)
	})
	suite.Run("If weight of 100 is provided", func() {
		incomingWeight := unit.Pound(100)
		adjustedWeight := GetAdjustedWeight(incomingWeight, false)
		suite.Equal(unit.Pound(500), *adjustedWeight)
	})
	suite.Run("If weight of 454 is provided", func() {
		incomingWeight := unit.Pound(454)
		adjustedWeight := GetAdjustedWeight(incomingWeight, false)
		suite.Equal(unit.Pound(500), *adjustedWeight)
	})
	suite.Run("If weight of 456 is provided", func() {
		incomingWeight := unit.Pound(456)
		adjustedWeight := GetAdjustedWeight(incomingWeight, false)
		suite.Equal(unit.Pound(501), *adjustedWeight)
	})
	suite.Run("If weight of 1000 is provided", func() {
		incomingWeight := unit.Pound(1000)
		adjustedWeight := GetAdjustedWeight(incomingWeight, false)
		suite.Equal(unit.Pound(1100), *adjustedWeight)
	})

	suite.Run("If no weight is provided UB", func() {
		var incomingWeight unit.Pound
		adjustedWeight := GetAdjustedWeight(incomingWeight, true)
		suite.Equal(unit.Pound(0), *adjustedWeight)
	})
	suite.Run("If a weight of 0 is provided UB", func() {
		incomingWeight := unit.Pound(0)
		adjustedWeight := GetAdjustedWeight(incomingWeight, true)
		suite.Equal(unit.Pound(0), *adjustedWeight)
	})
	suite.Run("If weight of 100 is provided UB", func() {
		incomingWeight := unit.Pound(100)
		adjustedWeight := GetAdjustedWeight(incomingWeight, true)
		suite.Equal(unit.Pound(300), *adjustedWeight)
	})
	suite.Run("If weight of 272 is provided UB", func() {
		incomingWeight := unit.Pound(272)
		adjustedWeight := GetAdjustedWeight(incomingWeight, true)
		suite.Equal(unit.Pound(300), *adjustedWeight)
	})
	suite.Run("If weight of 274 is provided UB", func() {
		incomingWeight := unit.Pound(274)
		adjustedWeight := GetAdjustedWeight(incomingWeight, true)
		suite.Equal(unit.Pound(301), *adjustedWeight)
	})
	suite.Run("If weight of 1000 is provided UB", func() {
		incomingWeight := unit.Pound(1000)
		adjustedWeight := GetAdjustedWeight(incomingWeight, true)
		suite.Equal(unit.Pound(1100), *adjustedWeight)
	})
}
