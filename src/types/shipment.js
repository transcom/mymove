/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

import { AddressShape } from './address';
import { AgentShape } from './agent';

import { SHIPMENT_OPTIONS } from 'shared/constants';

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
  diversion: PropTypes.bool,
  counselorRemarks: PropTypes.string,
  customerRemarks: PropTypes.string,
  status: PropTypes.string,
});
