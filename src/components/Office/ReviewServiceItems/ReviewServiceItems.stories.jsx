import React from 'react';
import { action } from '@storybook/addon-actions';

import ReviewServiceItems from './ReviewServiceItems';

import { SHIPMENT_OPTIONS, SERVICE_ITEM_STATUS } from 'shared/constants';

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
    shipmentId: '10',
    serviceItemName: 'Domestic linehaul',
    amount: 6423,
    status: SERVICE_ITEM_STATUS.SUBMITTED,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
  {
    id: '2',
    shipmentType: null,
    shipmentId: null,
    serviceItemName: 'Counseling fee',
    amount: 2.8,
    status: SERVICE_ITEM_STATUS.SUBMITTED,
    createdAt: '2020-01-02T00:08:00.999Z',
  },
  {
    id: '3',
    shipmentType: SHIPMENT_OPTIONS.NTS,
    shipmentId: '30',
    serviceItemName: 'Fuel surcharge',
    amount: 2.8,
    status: SERVICE_ITEM_STATUS.SUBMITTED,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
];

export const Basic = () => <ReviewServiceItems serviceItemCards={serviceItemCards} handleClose={action('clicked')} />;
