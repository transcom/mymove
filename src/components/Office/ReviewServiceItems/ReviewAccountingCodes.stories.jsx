import React from 'react';

import ReviewAccountingCodes from './ReviewAccountingCodes';

import { LOA_TYPE, SHIPMENT_OPTIONS, PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';

export default {
  title: 'Office Components/ReviewServiceItems/ReviewAccountingCodes',
  component: ReviewAccountingCodes,
};

const TACs = { HHG: '1234', NTS: '5678' };
const SACs = { HHG: 'AB12', NTS: 'CD34' };

const serviceItemsHHG = [
  {
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentTacType: LOA_TYPE.HHG,
    mtoShipmentID: '10',
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    amount: 23.45,
  },
  {
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentTacType: LOA_TYPE.HHG,
    mtoShipmentID: '10',
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    amount: 99.55,
  },
];

const serviceItemsNTSR = [
  {
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentTacType: LOA_TYPE.NTS,
    mtoShipmentSacType: LOA_TYPE.NTS,
    mtoShipmentID: '20',
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    amount: 559,
  },
  {
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentTacType: LOA_TYPE.NTS,
    mtoShipmentSacType: LOA_TYPE.NTS,
    mtoShipmentID: '20',
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    amount: 552.11,
  },
];

const moveLevelServices = [
  {
    amount: 44.33,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoServiceItemName: 'Move management',
  },
  {
    amount: 20.65,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoServiceItemName: 'Counseling',
  },
];

export const withOneShipment = () => <ReviewAccountingCodes TACs={TACs} SACs={SACs} cards={[...serviceItemsHHG]} />;

export const withMultipleShipments = () => (
  <ReviewAccountingCodes TACs={TACs} SACs={SACs} cards={[...serviceItemsHHG, ...serviceItemsNTSR]} />
);

export const withMultipleShipmentsAndMoveLevelServices = () => (
  <ReviewAccountingCodes
    TACs={TACs}
    SACs={SACs}
    cards={[...serviceItemsHHG, ...serviceItemsNTSR, ...moveLevelServices]}
  />
);
