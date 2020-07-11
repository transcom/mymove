import React from 'react';
import { action } from '@storybook/addon-actions';

import ReviewServiceItems from './ReviewServiceItems';

import { SHIPMENT_TYPE } from 'shared/constants';

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
    shipmentType: SHIPMENT_TYPE.HHG,
    serviceItemName: 'Domestic linehaul',
    amount: 6423,
    createdAt: '2020-01-01T00:08:00.999Z',
  },
];

export const Basic = () => (
  <ReviewServiceItems
    shipmentType={SHIPMENT_TYPE.HHG_SHORTHAUL_DOMESTIC}
    serviceItemCards={serviceItemCards}
    handleClose={action('clicked')}
  />
);
