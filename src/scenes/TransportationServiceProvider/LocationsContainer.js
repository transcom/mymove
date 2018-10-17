import { connect } from 'react-redux';
import { get } from 'lodash';
import Locations from './Locations';
import { getFormValues } from 'redux-form';

import { getPublicSwaggerDefinition } from 'shared/Swagger/selectors';

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
  const schema = getPublicSwaggerDefinition(state, 'Shipment');
  const formValues = getFormValues(formName)(state);

  return {
    addressSchema: get(state, 'swaggerPublic.spec.definitions.Address'),
    schema,
    formValues,

    deliveryAddress,
    initialValues: {
      pickupAddress: shipment.pickup_address,
      deliveryAddress: deliveryAddress,
      secondaryPickupAddress: shipment.secondary_pickup_address,
      hasDeliveryAddress: shipment.has_delivery_address,
      hasSecondaryPickupAddress: shipment.has_secondary_pickup_address,
    },
    shipment,
    hasDeliveryAddress: shipment.has_delivery_address,
    hasSecondaryPickupAddress: shipment.has_secondary_pickup_address,
    title: 'Locations',
    update: ownProps.update,

    // TO-DO: don't set has_properties automatically to true
    getUpdateArgs: () => {
      // const values = getFormValues(formName)(state);
      return [
        shipment.id,
        {
          delivery_address: formValues.deliveryAddress,
          pickup_address: formValues.pickupAddress,
          secondary_pickup_address: formValues.secondaryPickupAddress,
          has_secondary_pickup_address: formValues.has_secondary_pickup_address,
          has_delivery_address: formValues.has_delivery_address,
        },
      ];
    },
  };
};

export default connect(
  mapStateToProps,
  // mapDispatchToProps,
)(Locations);
