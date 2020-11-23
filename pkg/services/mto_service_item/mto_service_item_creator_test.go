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

type testMTOServiceItemQueryBuilder struct {
	fakeCreateOne   func(model interface{}) (*validate.Errors, error)
	fakeFetchOne    func(model interface{}, filters []services.QueryFilter) error
	fakeTransaction func(func(tx *pop.Connection) error) error
	fakeUpdateOne   func(models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testMTOServiceItemQueryBuilder) CreateOne(model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(model)
}

func (t *testMTOServiceItemQueryBuilder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	return t.fakeUpdateOne(model, eTag)
}

func (t *testMTOServiceItemQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	return t.fakeFetchOne(model, filters)
}

func (t *testMTOServiceItemQueryBuilder) Transaction(fn func(tx *pop.Connection) error) error {
	return t.fakeTransaction(fn)
}

func (suite *MTOServiceItemServiceSuite) TestCreateMTOServiceItem() {
	moveTaskOrder := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{Status: models.MoveStatusAPPROVED}})
	dimension := testdatagen.MakeDefaultMTOServiceItemDimension(suite.DB())
	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID: moveTaskOrder.ID,
		MoveTaskOrder:   moveTaskOrder,
		Dimensions: models.MTOServiceItemDimensions{
			dimension,
		},
	}

	// Happy path: If the service item is created successfully it should be returned
	suite.T().Run("success", func(t *testing.T) {
		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			return nil, nil
		}
		fakeFetchOne := func(model interface{}, filters []services.QueryFilter) error {
			return nil
		}
		fakeTx := func(fn func(tx *pop.Connection) error) error {
			return fn(&pop.Connection{})
		}
		fakeUpdateOne := func(model interface{}, etag *string) (*validate.Errors, error) {
			return nil, nil
		}

		builder := &testMTOServiceItemQueryBuilder{
			fakeCreateOne:   fakeCreateOne,
			fakeFetchOne:    fakeFetchOne,
			fakeTransaction: fakeTx,
			fakeUpdateOne:   fakeUpdateOne,
		}

		fakeCreateNewBuilder := func(db *pop.Connection) createMTOServiceItemQueryBuilder {
			return builder
		}

		creator := mtoServiceItemCreator{
			builder:          builder,
			createNewBuilder: fakeCreateNewBuilder,
		}
		createdServiceItem, verrs, err := creator.CreateMTOServiceItem(&serviceItem)

		suite.NoError(err)
		suite.Nil(verrs)
		suite.NotNil(createdServiceItem)

		createdServiceItemList := *createdServiceItem
		suite.NotEmpty(createdServiceItemList[0].Dimensions)
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

		builder := &testMTOServiceItemQueryBuilder{
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

		createdServiceItem, verrs, _ := creator.CreateMTOServiceItem(&serviceItem)
		suite.Error(verrs)
		suite.Nil(createdServiceItem)
	})

	// Should return a "NotFoundError" if the MTO ID is nil
	suite.T().Run("moveTaskOrderID not found", func(t *testing.T) {
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		notFoundID := uuid.Nil
		serviceItemNoMTO := models.MTOServiceItem{
			MoveTaskOrderID: notFoundID,
		}

		createdServiceItemNoMTO, _, err := creator.CreateMTOServiceItem(&serviceItemNoMTO)
		suite.Nil(createdServiceItemNoMTO)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundID.String())
	})

	// Should return a "NotFoundError" if the reServiceCode passed in isn't found on the table
	suite.T().Run("reServiceCode not found", func(t *testing.T) {
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		fakeCode := models.ReServiceCode("FAKE")
		serviceItemBadCode := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrder.ID,
			MoveTaskOrder:   moveTaskOrder,
			ReService: models.ReService{
				Code: fakeCode,
			},
		}

		createdServiceItemBadCode, _, err := creator.CreateMTOServiceItem(&serviceItemBadCode)
		suite.Nil(createdServiceItemBadCode)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), fakeCode)
	})

	// Should be able to create a service item with code ReServiceCodeMS or ReServiceCodeCS without a shipment,
	// and it should come back as "APPROVED"
	suite.T().Run("ReServiceCodeCS creation approved", func(t *testing.T) {
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		reServiceCS := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeCS,
			},
		})
		serviceItemCS := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrder.ID,
			MoveTaskOrder:   moveTaskOrder,
			ReService:       reServiceCS,
		}

		createdServiceItemCS, _, err := creator.CreateMTOServiceItem(&serviceItemCS)
		suite.NotNil(createdServiceItemCS)
		suite.NoError(err)

		createdServiceItemCSList := *createdServiceItemCS
		suite.Equal(createdServiceItemCSList[0].Status, models.MTOServiceItemStatus("APPROVED"))
	})

	// Should return a "NotFoundError" if the mtoShipmentID passed in isn't found
	// OR if it isn't linked to the mtoID passed in
	suite.T().Run("mtoShipmentID not found", func(t *testing.T) {
		builder := query.NewQueryBuilder(suite.DB())
		creator := NewMTOServiceItemCreator(builder)

		shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		reService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "ANY",
			},
		})
		serviceItemBadShip := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrder.ID,
			MoveTaskOrder:   moveTaskOrder,
			MTOShipmentID:   &shipment.ID,
			MTOShipment:     shipment,
			ReService:       reService,
		}

		createdServiceItemBadShip, _, err := creator.CreateMTOServiceItem(&serviceItemBadShip)
		suite.Nil(createdServiceItemBadShip)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), shipment.ID.String())
		suite.Contains(err.Error(), moveTaskOrder.ID.String())
	})

	// If the service item we're trying to create is shuttle service and there is no estimated weight, it fails.
	suite.T().Run("MTOServiceItemShuttle no prime weight", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultMTOShipment(suite.DB())

		serviceItemNoWeight := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrder.ID,
			MoveTaskOrder:   moveTaskOrder,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService: models.ReService{
				Code: models.ReServiceCodeDDSHUT,
			},
		}

		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			return nil, nil
		}
		fakeFetchOne := func(model interface{}, filters []services.QueryFilter) error {
			return nil
		}

		builder := &testMTOServiceItemQueryBuilder{
			fakeCreateOne: fakeCreateOne,
			fakeFetchOne:  fakeFetchOne,
		}

		creator := NewMTOServiceItemCreator(builder)
		createdServiceItem, _, err := creator.CreateMTOServiceItem(&serviceItemNoWeight)
		suite.Nil(createdServiceItem)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	// The timeMilitary fields need to be in the correct format.
	suite.T().Run("timeMilitary formatting for DDFSIT", func(t *testing.T) {
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{MTOShipment: models.MTOShipment{MoveTaskOrder: moveTaskOrder}})
		contact := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			FirstAvailableDeliveryDate: time.Now(),
		}
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrder.ID,
			MoveTaskOrder:   moveTaskOrder,
			MTOShipment:     shipment,
			MTOShipmentID:   &shipment.ID,
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		}

		fakeCreateOne := func(model interface{}) (*validate.Errors, error) {
			return nil, nil
		}
		fakeFetchOne := func(model interface{}, filters []services.QueryFilter) error {
			return nil
		}
		fakeTx := func(fn func(tx *pop.Connection) error) error {
			return fn(&pop.Connection{})
		}
		builder := &testMTOServiceItemQueryBuilder{
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
				MoveTaskOrderID: moveTaskOrder.ID,
				MoveTaskOrder:   moveTaskOrder,
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
	moveTaskOrder := testdatagen.MakeAvailableMove(suite.DB())
	moveTaskOrder.Status = models.MoveStatusAPPROVED
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
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

	serviceItemDOASIT := models.MTOServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		MTOShipment:     mtoShipment,
		MTOShipmentID:   &mtoShipment.ID,
		ReService:       reServiceDOASIT,
	}

	sitEntryDate := time.Date(2020, time.October, 24, 0, 0, 0, 0, time.UTC)
	sitPostalCode := "99999"
	reason := "lorem ipsum"

	suite.T().Run("Create DOFSIT service item", func(t *testing.T) {
		serviceItemDOFSIT := models.MTOServiceItem{
			MoveTaskOrder:   moveTaskOrder,
			MoveTaskOrderID: moveTaskOrder.ID,
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

		suite.NotNil(createdServiceItems)
		suite.NoError(err)
	})

	suite.T().Run("Create DOASIT item for shipment with DOFSIT", func(t *testing.T) {
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

	suite.T().Run("Do not create DOASIT if there is no DOFSIT on shipment", func(t *testing.T) {
		mtoShipmentNoServiceItems := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrder,
		})

		serviceItemDOASIT := models.MTOServiceItem{
			MoveTaskOrder:   moveTaskOrder,
			MoveTaskOrderID: moveTaskOrder.ID,
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
			MoveTaskOrder:   moveTaskOrder,
			MoveTaskOrderID: moveTaskOrder.ID,
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

	// Failed creation of DDFSIT because DDASIT/DDDSIT codes are not found in DB
	suite.T().Run("no DDASIT/DDDSIT codes", func(t *testing.T) {
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), "service code")
	})

	_ = testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
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
	suite.T().Run("bad CustomerContacts", func(t *testing.T) {
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
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.Contains(err.Error(), "timeMilitary")
	})

	// Successful creation of DDFSIT service item and the extra DDASIT/DDDSIT items
	suite.T().Run("DDFSIT creation approved", func(t *testing.T) {
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
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
				suite.Equal(len(item.CustomerContacts), len(serviceItemDDFSIT.CustomerContacts))
			}
		}
		suite.Equal(numDDASITFound, 1)
		suite.Equal(numDDDSITFound, 1)
		suite.Equal(numDDFSITFound, 1)
	})

	// Failed creation of DDFSIT because of duplicate service for shipment
	suite.T().Run("duplicate DDFSIT", func(t *testing.T) {
		serviceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID:  shipment.MoveTaskOrderID,
			MoveTaskOrder:    shipment.MoveTaskOrder,
			MTOShipmentID:    &shipment.ID,
			MTOShipment:      shipment,
			ReService:        reServiceDDFSIT,
			SITEntryDate:     &sitEntryDate,
			CustomerContacts: contacts,
		}

		createdServiceItems, _, err := creator.CreateMTOServiceItem(&serviceItemDDFSIT)
		suite.Nil(createdServiceItems)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
		suite.Contains(err.Error(), "A service item with reServiceCode DDFSIT already exists for this move and/or shipment.")
	})

	// Failed creation of DDDSIT service item
	suite.T().Run("cannot create DDDSIT", func(t *testing.T) {
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
