/* eslint-disable import/prefer-default-export */
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';

export const shipmentTypeLabels = {
  [SHIPMENT_OPTIONS.HHG]: 'HHG',
  [SHIPMENT_OPTIONS.PPM]: 'PPM',
  [SHIPMENT_OPTIONS.NTS]: 'NTS',
  [SHIPMENT_OPTIONS.NTSR]: 'NTS-release',
  [SHIPMENT_OPTIONS.BOAT]: 'Boat',
  [SHIPMENT_TYPES.BOAT_HAUL_AWAY]: 'Boat',
  [SHIPMENT_TYPES.BOAT_TOW_AWAY]: 'Boat',
};

export const shipmentForm = {
  header: {
    [SHIPMENT_OPTIONS.HHG]: 'Movers pack and transport this shipment',
    [SHIPMENT_OPTIONS.NTS]: 'Where and when should the movers pick up your personal property going into storage?',
    [SHIPMENT_OPTIONS.NTSR]: 'Where and when should the movers deliver your personal property from storage?',
  },
};

export const shipmentSectionLabels = {
  HHG: 'HHG shipment',
  HHG_INTO_NTS_DOMESTIC: 'NTS shipment',
  HHG_OUTOF_NTS_DOMESTIC: 'NTS-release shipment',
};
