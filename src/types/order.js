import PropTypes from 'prop-types';

import { AddressShape } from './address';
import { BackupContactShape } from './backupContact';

import dimensionTypes from 'constants/dimensionTypes';
import customerContactTypes from 'constants/customerContactTypes';
import { ShipmentOptionsOneOf } from 'types/shipment';

export const DestinationDutyLocationShape = PropTypes.shape({
  name: PropTypes.string,
  address: AddressShape,
});

export const OriginDutyLocationShape = PropTypes.shape({
  id: PropTypes.string,
  name: PropTypes.string,
  address_id: PropTypes.string,
  address: AddressShape,
});

export const EntitlementShape = PropTypes.shape({
  authorizedWeight: PropTypes.number,
  dependentsAuthorized: PropTypes.bool,
  nonTemporaryStorage: PropTypes.bool,
  privatelyOwnedVehicle: PropTypes.bool,
  proGearWeight: PropTypes.number,
  proGearWeightSpouse: PropTypes.number,
  storageInTransit: PropTypes.number,
  totalWeight: PropTypes.number,
  totalDependents: PropTypes.number,
});

export const OrderShape = PropTypes.shape({
  date_issued: PropTypes.string,
  report_by_date: PropTypes.string,
  department_indicator: PropTypes.string,
  order_number: PropTypes.string,
  order_type: PropTypes.string,
  order_type_detail: PropTypes.string,
  tac: PropTypes.string,
  sac: PropTypes.string,
  destinationDutyLocation: DestinationDutyLocationShape,
  originDutyLocation: OriginDutyLocationShape,
  entitlement: EntitlementShape,
});

export const OrdersInfoShape = PropTypes.shape({
  id: PropTypes.string,
  currentDutyLocation: OriginDutyLocationShape,
  newDutyLocation: DestinationDutyLocationShape,
  issuedDate: PropTypes.string,
  reportByDate: PropTypes.string,
  departmentIndicator: PropTypes.string,
  ordersNumber: PropTypes.string,
  ordersType: PropTypes.string,
  ordersTypeDetail: PropTypes.string,
  tacMDC: PropTypes.string,
  sacSDN: PropTypes.string,
});

export const CustomerShape = PropTypes.shape({
  agency: PropTypes.string,
  first_name: PropTypes.string,
  last_name: PropTypes.string,
  dodID: PropTypes.string,
  phone: PropTypes.string,
  email: PropTypes.string,
  current_address: AddressShape,
  backup_contact: BackupContactShape,
});

export const MTOShipmentShape = PropTypes.shape({
  id: PropTypes.string,
  shipmentType: PropTypes.string, // TODO - is this in API response?
  scheduledPickupDate: PropTypes.string,
  requestedPickupDate: PropTypes.string,
  pickupAddress: AddressShape,
  destinationAddress: AddressShape,
  diversion: PropTypes.bool,
});

export const MTOAgentShape = PropTypes.shape({
  id: PropTypes.string,
  firstName: PropTypes.string,
  lastName: PropTypes.string,
  agentType: PropTypes.string,
  email: PropTypes.string,
  phone: PropTypes.string,
});

export const MoveTaskOrderShape = PropTypes.shape({
  id: PropTypes.string,
  availableToPrimeAt: PropTypes.string,
  createdAt: PropTypes.string,
  eTag: PropTypes.string,
  isCanceled: PropTypes.bool,
  orderId: PropTypes.string,
  referenceId: PropTypes.string,
  requestedPickupDate: PropTypes.string,
  updatedAt: PropTypes.string,
});

export const MTOServiceItemDimensionShape = PropTypes.shape({
  type: PropTypes.oneOf(Object.values(dimensionTypes)),
  length: PropTypes.number,
  height: PropTypes.number,
  width: PropTypes.number,
});

export const MTOServiceItemCustomerContactShape = PropTypes.shape({
  type: PropTypes.oneOf(Object.values(customerContactTypes)),
  timeMilitary: PropTypes.string,
  firstAvailableDeliveryDate: PropTypes.string,
});

export const MTOServiceItemShape = PropTypes.shape({
  approvedAt: PropTypes.string,
  createdAt: PropTypes.string,
  customerContacts: PropTypes.arrayOf(MTOServiceItemCustomerContactShape),
  deletedAt: PropTypes.string,
  dimensions: PropTypes.arrayOf(MTOServiceItemDimensionShape),
  id: PropTypes.string,
  moveTaskOrderID: PropTypes.string,
  mtoShipmentID: PropTypes.string,
  pickupPostalCode: PropTypes.string,
  SITPostalCode: PropTypes.string,
  reServiceCode: PropTypes.string,
  reServiceID: PropTypes.string,
  reServiceName: PropTypes.string,
  reason: PropTypes.string,
  rejectedAt: PropTypes.string,
  submittedAt: PropTypes.string,
  status: PropTypes.string,
});

export const PaymentServiceItemParam = PropTypes.shape({
  key: PropTypes.string,
  value: PropTypes.string,
});

export const PaymentServiceItemShape = PropTypes.shape({
  id: PropTypes.string,
  createdAt: PropTypes.string,
  mtoServiceItemID: PropTypes.string,
  mtoServiceItemCode: PropTypes.string,
  mtoServiceItemName: PropTypes.string,
  mtoShipmentType: ShipmentOptionsOneOf,
  priceCents: PropTypes.number,
  status: PropTypes.string,
  rejectionReason: PropTypes.string,
  paymentServiceItemParams: PropTypes.arrayOf(PaymentServiceItemParam),
});

export const PaymentRequestShape = PropTypes.shape({
  id: PropTypes.string,
  createdAt: PropTypes.string,
  moveTaskOrderID: PropTypes.string,
  paymentRequestNumber: PropTypes.string,
  status: PropTypes.string,
  eTag: PropTypes.string,
  serviceItems: PropTypes.arrayOf(PropTypes.oneOfType([PropTypes.string, PaymentServiceItemShape])),
  reviewedAt: PropTypes.string,
});

export const ExistingUploadsShape = PropTypes.arrayOf(
  PropTypes.shape({
    id: PropTypes.string.isRequired,
    created_at: PropTypes.string.isRequired,
    bytes: PropTypes.number.isRequired,
    url: PropTypes.string.isRequired,
    filename: PropTypes.string.isRequired,
  }),
);

export const OrdersLOAShape = PropTypes.shape({
  tac: PropTypes.string,
  sac: PropTypes.string,
  ntsTac: PropTypes.string,
  ntsSac: PropTypes.string,
});
