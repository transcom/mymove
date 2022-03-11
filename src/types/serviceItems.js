import PropTypes from 'prop-types';

import { LOA_TYPE, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { MTOServiceItemCustomerContactShape, MTOServiceItemDimensionShape, PaymentServiceItemParam } from 'types/order';

export const ServiceItemCardShape = PropTypes.shape({
  id: PropTypes.string, // service item id
  mtoShipmentID: PropTypes.string,
  mtoShipmentType: ShipmentOptionsOneOf,
  mtoShipmentTacType: PropTypes.oneOf(Object.values(LOA_TYPE)),
  mtoShipmentSacType: PropTypes.oneOf(Object.values(LOA_TYPE)),
  mtoServiceItemCode: PropTypes.string,
  mtoServiceItemName: PropTypes.string,
  amount: PropTypes.number,
  status: PropTypes.oneOf(Object.values(PAYMENT_SERVICE_ITEM_STATUS)),
  createdAt: PropTypes.string,
  paymentServiceItemParams: PropTypes.arrayOf(PaymentServiceItemParam),
});

export const ServiceItemCardsShape = PropTypes.arrayOf(ServiceItemCardShape);

export const ServiceItemDetailsShape = PropTypes.shape({
  id: PropTypes.string,
  mtoShipmentID: PropTypes.string,
  createdAt: PropTypes.string,
  submittedAt: PropTypes.string,
  approvedAt: PropTypes.string,
  rejectedAt: PropTypes.string,
  serviceItem: PropTypes.string,
  code: PropTypes.string,
  status: PropTypes.string,
  details: PropTypes.shape({
    reason: PropTypes.string,
    rejectionReason: PropTypes.string,
    description: PropTypes.string,
    pickupPostalCode: PropTypes.string,
    SITPostalCode: PropTypes.string,
    itemDimensions: MTOServiceItemDimensionShape,
    crateDimensions: MTOServiceItemDimensionShape,
    firstCustomerContact: MTOServiceItemCustomerContactShape,
    secondCustomerContact: MTOServiceItemCustomerContactShape,
    estimatedWeight: PropTypes.number,
  }),
});

export const ShipmentPaymentSITBalanceShape = PropTypes.shape({
  previouslyBilledDays: PropTypes.number,
  previouslyBilledEndDate: PropTypes.string,
  pendingSITDaysInvoiced: PropTypes.number.isRequired,
  pendingBilledEndDate: PropTypes.string.isRequired,
  totalSITDaysAuthorized: PropTypes.number.isRequired,
  totalSITDaysRemaining: PropTypes.number.isRequired,
  totalSITEndDate: PropTypes.string.isRequired,
});
