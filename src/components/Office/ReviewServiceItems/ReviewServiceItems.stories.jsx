import React from 'react';

import ReviewServiceItems from './ReviewServiceItems';

import { SHIPMENT_OPTIONS } from 'shared/constants';

// Left Nav
export default {
  title: 'TOO/TIO Components|Review service items',
  component: ReviewServiceItems,
  decorators: [
    (storyFn) => (
      <div id="l-nav" style={{ padding: '20px', background: '#f0f0f0' }}>
        {storyFn()}
      </div>
    ),
  ],
};

const serviceItemCards = [
  {
    id: '1',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    serviceItemName: 'Domestic linehaul',
    amount: 6423,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
  {
    id: '2',
    shipmentType: null,
    serviceItemName: 'Counseling fee',
    amount: 2.8,
    createdAt: '2020-01-02T00:08:00.999Z',
  },
  {
    id: '3',
    shipmentType: SHIPMENT_OPTIONS.NTS,
    serviceItemName: 'Fuel surcharge',
    amount: 2.8,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
];

export const Basic = () => (
  <ReviewServiceItems shipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC} serviceItemCards={serviceItemCards} />
);
