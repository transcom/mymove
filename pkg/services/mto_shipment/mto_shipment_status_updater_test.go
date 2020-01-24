package mtoshipment

import (
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testPaymentRequestStatusQueryBuilder struct {
	fakeFetchOne func(model interface{}) error
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
	params := mtoshipmentops.PatchMTOShipmentStatusParams{
		IfUnmodifiedSince: strfmt.DateTime(shipment.UpdatedAt),
		Body:              &ghcmessages.MTOShipment{Status: "APPROVED"},
	}

	suite.T().Run("If we get a mto shipment pointer with a status it should update and return no error", func(t *testing.T) {

		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			reflect.ValueOf(model).Elem().FieldByName("MoveTaskOrderID").Set(reflect.ValueOf(shipment.MoveTaskOrderID))
			reflect.ValueOf(model).Elem().FieldByName("PickupAddressID").Set(reflect.ValueOf(shipment.PickupAddressID))
			reflect.ValueOf(model).Elem().FieldByName("DestinationAddressID").Set(reflect.ValueOf(shipment.DestinationAddressID))
			reflect.ValueOf(model).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(shipment.UpdatedAt))
			return nil
		}

		builder := &testPaymentRequestStatusQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		updater := NewMTOShipmentStatusUpdater(suite.DB(), builder)

		_, err := updater.UpdateMTOShipmentStatus(params)
		suite.NoError(err)
	})

	suite.T().Run("If there is an error updating the shipment status we should get one returned", func(t *testing.T) {
		fakeFetchOne := func(model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		builder := &testPaymentRequestStatusQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}

		updater := NewMTOShipmentStatusUpdater(suite.DB(), builder)

		_, err := updater.UpdateMTOShipmentStatus(params)
		suite.Error(err)
		suite.IsType(&ValidationError{}, err)
	})
}
