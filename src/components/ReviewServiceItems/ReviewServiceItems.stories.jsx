import React from 'react';

import ReviewServiceItems from './ReviewServiceItems';

import { SHIPMENT_TYPE } from 'shared/constants';

// Left Nav
export default {
  title: 'TOO/TIO Components|Review service items',
  component: ReviewServiceItems,
};

const serviceItemCards = [
  {
    id: '1',
    shipmentType: SHIPMENT_TYPE.HHG,
    serviceItemName: 'Domestic linehaul',
    amount: 6423,
  },
];

export const Basic = () => (
  <div id="l-nav" style={{ padding: '20px', background: '#f0f0f0' }}>
    <ReviewServiceItems shipmentType={SHIPMENT_TYPE.HHG_SHORTHAUL_DOMESTIC} serviceItemCards={serviceItemCards} />
  </div>
);
