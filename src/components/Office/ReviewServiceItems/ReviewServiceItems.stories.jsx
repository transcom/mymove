import React from 'react';
import { action } from '@storybook/addon-actions';

import ReviewServiceItems from './ReviewServiceItems';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'TOO/TIO Components|ReviewServiceItems',
  component: ReviewServiceItems,
  decorators: [
    (storyFn) => (
      <div id="l-nav" style={{ padding: '20px', background: '#f0f0f0' }}>
        {storyFn()}
      </div>
    ),
  ],
};

export const Basic = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        serviceItemName: 'Domestic linehaul',
        amount: 1234.0,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);

export const HHG = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        serviceItemName: 'Domestic linehaul',
        amount: 5678.05,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);

export const NTS = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.NTS,
        serviceItemName: 'Domestic linehaul',
        amount: 6423.51,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);
