import React from 'react';

import ReviewAccountingCodes from './ReviewAccountingCodes';

import { LOA_TYPE, SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Office Components/ReviewServiceItems/ReviewAccountingCodes',
  component: ReviewAccountingCodes,
};

const TACs = { HHG: '1234', NTS: '5678' };
const SACs = { HHG: 'AB12', NTS: 'CD34' };

const serviceItemsHHG = [
  {
    status: 'APPROVED',
    mtoShipmentTacType: LOA_TYPE.HHG,
    mtoShipmentID: '10',
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    amount: 23.45,
  },
  {
    status: 'APPROVED',
    mtoShipmentTacType: LOA_TYPE.HHG,
    mtoShipmentID: '10',
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    amount: 99.55,
  },
];

const serviceItemsNTSR = [
  {
    status: 'APPROVED',
    mtoShipmentTacType: LOA_TYPE.NTS,
    mtoShipmentSacType: LOA_TYPE.NTS,
    mtoShipmentID: '20',
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    amount: 559,
  },
  {
    status: 'APPROVED',
    mtoShipmentTacType: LOA_TYPE.NTS,
    mtoShipmentSacType: LOA_TYPE.NTS,
    mtoShipmentID: '20',
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    amount: 552.11,
  },
];

export const withOneShipment = () => <ReviewAccountingCodes TACs={TACs} SACs={SACs} cards={[...serviceItemsHHG]} />;

export const withMultipleShipments = () => (
  <ReviewAccountingCodes TACs={TACs} SACs={SACs} cards={[...serviceItemsHHG, ...serviceItemsNTSR]} />
);
