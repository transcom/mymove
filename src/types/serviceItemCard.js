import PropTypes from 'prop-types';

import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';

export const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  mtoShipmentID: PropTypes.string,
  mtoShipmentType: ShipmentOptionsOneOf,
  mtoServiceItemName: PropTypes.string,
  amount: PropTypes.number,
  status: PropTypes.oneOf(Object.values(PAYMENT_SERVICE_ITEM_STATUS)),
  rejectionReason: PropTypes.string,
  createdAt: PropTypes.string,
});

// eslint-disable-next-line import/prefer-default-export
export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);
