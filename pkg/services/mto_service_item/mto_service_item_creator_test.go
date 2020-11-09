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
