/* eslint-disable import/prefer-default-export */
import { SHIPMENT_OPTIONS } from 'shared/constants';

export const shipmentTypes = {
  [SHIPMENT_OPTIONS.HHG]: 'HHG',
  [SHIPMENT_OPTIONS.PPM]: 'PPM',
  [SHIPMENT_OPTIONS.NTS]: 'NTS',
  [SHIPMENT_OPTIONS.NTSR]: 'NTSR',
  [SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC]: 'HHG',
  [SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC]: 'HHG',
};

export const mtoShipmentTypes = {
  [SHIPMENT_OPTIONS.HHG]: 'Household goods',
  [SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC]: 'Household goods',
  [SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC]: 'Household goods',
  [SHIPMENT_OPTIONS.PPM]: 'Personally procured move',
  [SHIPMENT_OPTIONS.NTS]: 'Non-temp storage',
  [SHIPMENT_OPTIONS.NTSR]: 'Non-temp storage release',
};

export const shipmentStatuses = {
  DRAFT: 'DRAFT',
  SUBMITTED: 'SUBMITTED',
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
  CANCELLATION_REQUESTED: 'CANCELLATION_REQUESTED',
};
