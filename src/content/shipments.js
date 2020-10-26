/* eslint-disable import/prefer-default-export */
import { SHIPMENT_OPTIONS } from 'shared/constants';

export const shipmentTypes = {
  [SHIPMENT_OPTIONS.HHG]: 'HHG',
  [SHIPMENT_OPTIONS.PPM]: 'PPM',
  [SHIPMENT_OPTIONS.NTS]: 'NTS',
  [SHIPMENT_OPTIONS.NTSR]: 'NTS-R',
};

export const shipmentForm = {
  header: {
    [SHIPMENT_OPTIONS.HHG]: 'When and where can the movers pick up and deliver this shipment?',
    [SHIPMENT_OPTIONS.NTS]: 'Where and when should the movers pick up your things going into storage?',
    [SHIPMENT_OPTIONS.NTSR]: 'Where and when should the movers release your things from storage?',
  },
};
