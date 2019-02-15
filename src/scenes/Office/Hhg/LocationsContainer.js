import { get } from 'lodash';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import { getPublicSwaggerDefinition } from 'shared/Swagger/selectors';
import { selectOrdersForMove } from 'shared/Entities/modules/orders';
import { selectShipment } from 'shared/Entities/modules/shipments';
import LocationsPanel from 'shared/LocationsPanel/LocationsPanel';

const mapStateToProps = (state, ownProps) => {
  const shipment = selectShipment(state, ownProps.shipmentId);
  const formName = 'shipment_locations';
  const orders = selectOrdersForMove(state, shipment.move_id);
  const newDutyStation = get(orders, 'new_duty_station.address', {});
  const schema = getPublicSwaggerDefinition(state, 'Shipment');
  const formValues = getFormValues(formName)(state);

  return {
    addressSchema: get(state, 'swaggerPublic.spec.definitions.Address'),
    schema,
    formValues,
    newDutyStation,

    initialValues: {
      pickup_address: get(shipment, 'pickup_address', {}),
      delivery_address: get(shipment, 'delivery_address', {}),
      secondary_pickup_address: get(shipment, 'secondary_pickup_address', {}),
      has_delivery_address: shipment.has_delivery_address,
      has_secondary_pickup_address: shipment.has_secondary_pickup_address,
    },
    shipment,
    title: 'Locations',

    getUpdateArgs: () => {
      const params = {
        pickup_address: formValues.pickup_address,
        has_secondary_pickup_address: formValues.has_secondary_pickup_address,
        has_delivery_address: formValues.has_delivery_address,
      };
      // Avoid sending empty objects as addresses
      if (formValues.has_secondary_pickup_address) {
        params.secondary_pickup_address = formValues.secondary_pickup_address;
      }
      if (formValues.has_delivery_address) {
        params.delivery_address = formValues.delivery_address;
      }
      return [shipment.id, params];
    },
  };
};

export default connect(mapStateToProps)(LocationsPanel);
