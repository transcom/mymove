import PropTypes from 'prop-types';

import { ShipmentOptionsOneOf } from 'types/shipment';

export const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  shipmentType: ShipmentOptionsOneOf,
  serviceItemName: PropTypes.string,
  amount: PropTypes.number,
  createdAt: PropTypes.string,
});

// eslint-disable-next-line import/prefer-default-export
export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);
