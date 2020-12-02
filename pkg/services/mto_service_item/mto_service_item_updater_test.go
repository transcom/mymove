package mtoserviceitem

import (
	"testing"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/services"
)

type testUpdateMTOServiceItemQueryBuilder struct {
	fakeFetchOne  func(model interface{}, filters []services.QueryFilter) error
	fakeUpdateOne func(models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testUpdateMTOServiceItemQueryBuilder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	return t.fakeUpdateOne(model, eTag)
}

func (t *testUpdateMTOServiceItemQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	return t.fakeFetchOne(model, filters)
}

func (suite *MTOServiceItemServiceSuite) TestMTOServiceItemUpdater() {
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewMTOServiceItemUpdater(builder)

	serviceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	// Test not found error
	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundServiceItem := serviceItem
		notFoundServiceItem.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedServiceItem, err := updater.UpdateMTOServiceItemBase(&notFoundServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})

	// Test validation error
	suite.T().Run("Validation Error", func(t *testing.T) {
		invalidServiceItem := serviceItem
		invalidServiceItem.MoveTaskOrderID = serviceItem.ID // invalid Move ID

		updatedServiceItem, err := updater.UpdateMTOServiceItemBase(&invalidServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		invalidInputError := err.(services.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "moveTaskOrderID")
	})

	// Test precondition failed (stale eTag)
	suite.T().Run("Precondition Failed", func(t *testing.T) {
		newServiceItem := serviceItem
		updatedServiceItem, err := updater.UpdateMTOServiceItemBase(&newServiceItem, "bloop")

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	// Test successful update
	suite.T().Run("Success", func(t *testing.T) {
		reason := "because we did this service"
		sitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)

		newServiceItem := serviceItem
		newServiceItem.Reason = &reason
		newServiceItem.SITEntryDate = &sitEntryDate
		newServiceItem.Status = "" // should keep the status from the original service item

		updatedServiceItem, err := updater.UpdateMTOServiceItemBase(&newServiceItem, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(updatedServiceItem.ID, serviceItem.ID)
		suite.Equal(updatedServiceItem.MTOShipmentID, serviceItem.MTOShipmentID)
		suite.Equal(updatedServiceItem.MoveTaskOrderID, serviceItem.MoveTaskOrderID)
		suite.Equal(updatedServiceItem.Reason, newServiceItem.Reason)
		suite.Equal(updatedServiceItem.SITEntryDate.Local(), newServiceItem.SITEntryDate.Local())
		suite.Equal(updatedServiceItem.Status, serviceItem.Status) // should not have been updated
		suite.NotEqual(updatedServiceItem.Status, newServiceItem.Status)
	})
}

func (suite *MTOServiceItemServiceSuite) TestValidateUpdateMTOServiceItem() {
	// TODO
}
