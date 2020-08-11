import PropTypes from 'prop-types';

import { AddressShape } from './address';

export const DestinationDutyStationShape = PropTypes.shape({
  name: PropTypes.string,
});

export const OriginDutyStationShape = PropTypes.shape({
  name: PropTypes.string,
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

export const MoveOrderShape = PropTypes.shape({
  date_issued: PropTypes.string,
  report_by_date: PropTypes.string,
  department_indicator: PropTypes.string, // TODO - is this in the API response?
  order_number: PropTypes.string,
  order_type: PropTypes.string,
  order_type_detail: PropTypes.string,
  tac: PropTypes.string,
  sacSDN: PropTypes.string,
  destinationDutyStation: DestinationDutyStationShape,
  originDutyStation: OriginDutyStationShape,
  entitlement: EntitlementShape,
});

export const CustomerShape = PropTypes.shape({
  agency: PropTypes.string,
  first_name: PropTypes.string,
  last_name: PropTypes.string,
  dodID: PropTypes.string,
  phone: PropTypes.string,
  email: PropTypes.string,
  current_address: AddressShape,
});

export const MTOShipmentShape = PropTypes.shape({
  id: PropTypes.string,
  shipmentType: PropTypes.string, // TODO - is this in API response?
  scheduledPickupDate: PropTypes.string,
  requestedPickupDate: PropTypes.string,
  pickupAddress: AddressShape,
  destinationAddress: AddressShape,
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
  moveOrderId: PropTypes.string,
  referenceId: PropTypes.string,
  requestedPickupDate: PropTypes.string,
  updatedAt: PropTypes.string,
});

export const MTOServiceItemShape = PropTypes.shape({
  approvedAt: PropTypes.string,
  createdAt: PropTypes.string,
  deletedAt: PropTypes.string,
  id: PropTypes.string,
  moveTaskOrderID: PropTypes.string,
  mtoShipmentID: PropTypes.string,
  pickupPostalCode: PropTypes.string,
  reServiceCode: PropTypes.string,
  reServiceID: PropTypes.string,
  reServiceName: PropTypes.string,
  reason: PropTypes.string,
  rejectedAt: PropTypes.string,
  submittedAt: PropTypes.string,
  status: PropTypes.string,
});

export const PaymentServiceItemShape = PropTypes.shape({
  id: PropTypes.string,
  createdAt: PropTypes.string,
  mtoServiceItemID: PropTypes.string,
  priceCents: PropTypes.number,
  status: PropTypes.string,
  rejectionReason: PropTypes.string,
});

export const PaymentRequestShape = PropTypes.shape({
  id: PropTypes.string,
  moveTaskOrderID: PropTypes.string,
  paymentRequestNumber: PropTypes.string,
  status: PropTypes.string,
  eTag: PropTypes.string,
  serviceItems: PropTypes.arrayOf(PropTypes.string),
  reviewedAt: PropTypes.string,
});
