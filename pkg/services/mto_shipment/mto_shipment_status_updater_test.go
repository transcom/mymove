package mtoshipment

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gobuffalo/validate"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testPaymentRequestStatusQueryBuilder struct {
	fakeUpdateOne func(model interface{}) (*validate.Errors, error)
	fakeFetchOne  func(model interface{}) error
}

func (t *testPaymentRequestStatusQueryBuilder) UpdateOne(model interface{}) (*validate.Errors, error) {
	v, m := t.fakeUpdateOne(model)
	return v, m
}

func (t *testPaymentRequestStatusQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	m := t.fakeFetchOne(model)
	return m
}

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentStatus() {
	id, err := uuid.NewV4()
	suite.NoError(err)

	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	shipment.Status = models.MTOShipmentStatusSubmitted

	suite.T().Run("If we get a mto shipment pointer with a status it should update and return no error", func(t *testing.T) {

		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		fakeUpdateOne := func(model interface{}) (*validate.Errors, error) {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return &validate.Errors{}, nil
		}

		builder := &testPaymentRequestStatusQueryBuilder{
			fakeUpdateOne: fakeUpdateOne,
			fakeFetchOne:  fakeFetchOne,
		}

		updater := NewMTOShipmentStatusUpdater(builder)

		verrs, err := updater.UpdateMTOShipmentStatus(&shipment, "APPROVED")
		suite.NoError(err)
		suite.NoVerrs(verrs)
	})

	suite.T().Run("If there is an error updating the payment request status we should get one returned", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		fakeUpdateOne := func(model interface{}) (*validate.Errors, error) {
			return nil, errors.New("Update error")
		}

		builder := &testPaymentRequestStatusQueryBuilder{
			fakeUpdateOne: fakeUpdateOne,
			fakeFetchOne:  fakeFetchOne,
		}

		updater := NewMTOShipmentStatusUpdater(builder)

		_, err := updater.UpdateMTOShipmentStatus(&shipment, "Invalid status")
		suite.Error(err)
		suite.Equal(err.Error(), "Update error")

	})
}
