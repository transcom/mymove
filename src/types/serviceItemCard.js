import PropTypes from 'prop-types';

import { ShipmentTypeOneOf } from 'types/shipment';

const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  shipmentType: ShipmentTypeOneOf,
  serviceItemName: PropTypes.string,
  amount: PropTypes.number,
});

// eslint-disable-next-line import/prefer-default-export
export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);
