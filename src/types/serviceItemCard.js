import PropTypes from 'prop-types';

import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { PaymentServiceItemParam } from 'types/order';

export const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  mtoShipmentID: PropTypes.string,
  mtoShipmentType: ShipmentOptionsOneOf,
  mtoServiceItemCode: PropTypes.string,
  mtoServiceItemName: PropTypes.string,
  amount: PropTypes.number,
  status: PropTypes.oneOf(Object.values(PAYMENT_SERVICE_ITEM_STATUS)),
  rejectionReason: PropTypes.string,
  createdAt: PropTypes.string,
  paymentServiceItemParams: PropTypes.arrayOf(PaymentServiceItemParam),
});

// eslint-disable-next-line import/prefer-default-export
export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);
