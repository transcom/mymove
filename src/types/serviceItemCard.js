import PropTypes from 'prop-types';

import { ShipmentTypeOneOf } from 'types/shipment';
import { SERVICE_ITEM_STATUS } from 'shared/constants';

const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  shipmentType: ShipmentTypeOneOf,
  shipmentId: PropTypes.string,
  serviceItemName: PropTypes.string,
  amount: PropTypes.number,
  status: PropTypes.oneOf(Object.values(SERVICE_ITEM_STATUS)),
  createdAt: PropTypes.string,
});

// eslint-disable-next-line import/prefer-default-export
export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);
