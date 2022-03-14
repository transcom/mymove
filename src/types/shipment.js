/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

import { AddressShape, ResidentialAddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { shipmentStatuses } from 'constants/shipments';

export const ShipmentOptionsOneOf = PropTypes.oneOf([
  SHIPMENT_OPTIONS.HHG,
  SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
  SHIPMENT_OPTIONS.NTS,
  SHIPMENT_OPTIONS.NTSR,
  SHIPMENT_OPTIONS.PPM,
]);

export const ShipmentShape = PropTypes.shape({
  shipmentType: ShipmentOptionsOneOf,
  requestedPickupDate: PropTypes.string,
  scheduledPickupDate: PropTypes.string,
  actualPickupDate: PropTypes.string,
  pickupAddress: AddressShape,
  secondaryPickupAddress: AddressShape,
  destinationAddress: AddressShape,
  secondaryDeliveryAddress: AddressShape,
  agents: PropTypes.arrayOf(AgentShape),
  primeEstimatedWeight: PropTypes.number,
  primeActualWeight: PropTypes.number,
  ntsRecordedWeight: PropTypes.number,
  diversion: PropTypes.bool,
  counselorRemarks: PropTypes.string,
  customerRemarks: PropTypes.string,
  status: PropTypes.string,
  reweigh: PropTypes.shape({
    id: PropTypes.string,
  }),
  storageFacility: PropTypes.shape({
    address: AddressShape.isRequired,
    facilityName: PropTypes.string.isRequired,
    lotNumber: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
  }),
});

export const ShipmentStatusesOneOf = PropTypes.oneOf([
  shipmentStatuses.DRAFT,
  shipmentStatuses.SUBMITTED,
  shipmentStatuses.APPROVED,
  shipmentStatuses.CANCELLATION_REQUESTED,
  shipmentStatuses.DIVERSION_REQUESTED,
  shipmentStatuses.CANCELED,
  shipmentStatuses.REJECTED,
]);

export const StorageFacilityShape = PropTypes.shape({
  facilityName: PropTypes.string,
  phone: PropTypes.string,
  email: PropTypes.string,
  address: ResidentialAddressShape,
  lotNumber: PropTypes.string,
});
