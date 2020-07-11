import PropTypes from 'prop-types';

import { ShipmentTypeOneOf } from 'types/shipment';

export const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  shipmentType: ShipmentTypeOneOf,
  serviceItemName: PropTypes.string,
  amount: PropTypes.number,
  createdAt: PropTypes.string,
});

// eslint-disable-next-line import/prefer-default-export
export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);
