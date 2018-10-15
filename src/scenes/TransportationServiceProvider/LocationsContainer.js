import { connect } from 'react-redux';
import { get } from 'lodash';
import Locations from './Locations';
import { patchShipment } from './ducks';
import { getFormValues } from 'redux-form';

const mapStateToProps = (state, ownProps) => {
  const shipment = get(state, 'tsp.shipment', {});
  const formName = 'shipment_locations';
  const newDutyStation = get(shipment, 'move.new_duty_station.address', {});
  // if they do not have a delivery address, default to the station's address info
  const deliveryAddress = shipment.has_delivery_address
    ? shipment.delivery_address
    : {
        city: newDutyStation.city,
        state: newDutyStation.state,
        postal_code: newDutyStation.postal_code,
      };

  return {
    addressSchema: get(state, 'swaggerPublic.spec.definitions.Address'),
    deliveryAddress,
    initialValues: {
      pickupAddress: shipment.pickup_address,
      deliveryAddress: deliveryAddress,
      secondaryPickupAddress: shipment.secondary_pickup_address,
    },
    shipment,
    title: 'Locations',
    update: ownProps.update,

    // TO-DO: don't set has_properties automatically to true
    getUpdateArgs: () => {
      const values = getFormValues(formName)(state);
      return [
        shipment.id,
        {
          delivery_address: values.deliveryAddress,
          pickup_address: values.pickupAddress,
          secondary_pickup_address: values.secondaryPickupAddress,
          has_secondary_pickup_address: true,
          has_delivery_address: true,
        },
      ];
    },
  };
};

export default connect(
  mapStateToProps,
  // mapDispatchToProps,
)(Locations);
