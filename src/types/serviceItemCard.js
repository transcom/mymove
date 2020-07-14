import PropTypes from 'prop-types';

import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';

export const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  shipmentId: PropTypes.string,
  shipmentType: ShipmentOptionsOneOf,
  serviceItemName: PropTypes.string,
  amount: PropTypes.number,
  status: PropTypes.oneOf(Object.values(SERVICE_ITEM_STATUS)),
  createdAt: PropTypes.string,
});

// eslint-disable-next-line import/prefer-default-export
export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);
