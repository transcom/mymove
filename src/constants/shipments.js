/* eslint-disable import/prefer-default-export */
import { SHIPMENT_OPTIONS } from 'shared/constants';

export const shipmentTypes = {
  [SHIPMENT_OPTIONS.HHG]: 'HHG',
  [SHIPMENT_OPTIONS.PPM]: 'PPM',
  [SHIPMENT_OPTIONS.NTS]: 'NTS',
  [SHIPMENT_OPTIONS.NTSR]: 'NTS-release',
  [SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC]: 'HHG',
  [SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC]: 'HHG',
};

export const shipmentModificationTypes = {
  CANCELED: 'CANCELED',
  DIVERSION: 'DIVERSION',
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
  CANCELED: 'CANCELED',
  DIVERSION_REQUESTED: 'DIVERSION_REQUESTED',
};

export const ppmShipmentStatuses = {
  DRAFT: 'DRAFT',
  SUBMITTED: 'SUBMITTED',
  WAITING_ON_CUSTOMER: 'WAITING_ON_CUSTOMER',
  NEEDS_ADVANCE_APPROVAL: 'NEEDS_ADVANCE_APPROVAL',
  NEEDS_PAYMENT_APPROVAL: 'NEEDS_PAYMENT_APPROVAL',
  PAYMENT_APPROVED: 'PAYMENT_APPROVED',
};

export const shipmentDestinationTypes = {
  HOME_OF_RECORD: 'Home of record (HOR)',
  HOME_OF_SELECTION: 'Home of selection (HOS)',
  PLACE_ENTERED_ACTIVE_DUTY: 'Place entered active duty (PLEAD)',
  OTHER_THAN_AUTHORIZED: 'Other than authorized',
};

export const LONGHAUL_MIN_DISTANCE = 50;

export const PPM_MAX_ADVANCE_RATIO = 0.6;
