import React from 'react';
import { action } from '@storybook/addon-actions';

import ReviewServiceItems from './ReviewServiceItems';

import { SHIPMENT_OPTIONS, SERVICE_ITEM_STATUS } from 'shared/constants';

export default {
  title: 'TOO/TIO Components|ReviewServiceItems',
  component: ReviewServiceItems,
  decorators: [
    (storyFn) => (
      <div style={{ margin: '10px', height: '80vh', display: 'flex', flexDirection: 'column' }}>{storyFn()}</div>
    ),
  ],
};

export const Basic = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        serviceItemName: 'Counseling services',
        amount: 1234.0,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);

export const BasicWithTwoItems = () => {
  return (
    <ReviewServiceItems
      serviceItemCards={[
        {
          id: '1',
          serviceItemName: 'Counseling services',
          amount: 1234.0,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
        {
          id: '2',
          serviceItemName: 'Move management',
          amount: 1234.0,
          createdAt: '2020-01-01T00:08:00.999Z',
        },
      ]}
      handleClose={action('clicked')}
    />
  );
};
// TODO - Skipping this story for now since this has animations and lokiAsync() can't be used in Storybook CSF format
BasicWithTwoItems.story = {
  parameters: {
    loki: { skip: true },
  },
};

export const HHG = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        shipmentId: '10',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        serviceItemName: 'Domestic linehaul',
        amount: 5678.05,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);

export const NonTemporaryStorage = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        shipmentId: '10',
        shipmentType: SHIPMENT_OPTIONS.NTS,
        serviceItemName: 'Domestic linehaul',
        amount: 6423.51,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);

export const MultipleShipmentsGroups = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        serviceItemName: 'Counseling services',
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
      {
        id: '3',
        shipmentId: '20',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        serviceItemName: 'Domestic linehaul',
        amount: 5678.05,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '4',
        shipmentId: '30',
        shipmentType: SHIPMENT_OPTIONS.NTS,
        serviceItemName: 'Domestic linehaul',
        amount: 6423.51,
        createdAt: '2020-01-01T00:07:30.999Z',
      },
      {
        id: '5',
        shipmentId: '30',
        shipmentType: SHIPMENT_OPTIONS.NTS,
        serviceItemName: 'Fuel Surcharge',
        amount: 100000000000000,
        createdAt: '2020-01-01T00:07:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);

export const WithStatusAndReason = () => (
  <ReviewServiceItems
    serviceItemCards={[
      {
        id: '1',
        serviceItemName: 'Counseling services',
        amount: 0.01,
        createdAt: '2020-01-01T00:09:00.999Z',
      },
      {
        id: '2',
        serviceItemName: 'Move management',
        amount: 1234.0,
        status: SERVICE_ITEM_STATUS.REJECTED,
        rejectionReason: 'Amount exceeds limit',
        createdAt: '2020-01-01T00:06:00.999Z',
      },
      {
        id: '3',
        shipmentId: '20',
        shipmentType: SHIPMENT_OPTIONS.HHG,
        serviceItemName: 'Domestic linehaul',
        amount: 5678.05,
        status: SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:08:00.999Z',
      },
      {
        id: '4',
        shipmentId: '30',
        shipmentType: SHIPMENT_OPTIONS.NTS,
        serviceItemName: 'Domestic linehaul',
        amount: 6423.51,
        status: SERVICE_ITEM_STATUS.APPROVED,
        createdAt: '2020-01-01T00:07:30.999Z',
      },
      {
        id: '5',
        shipmentId: '30',
        shipmentType: SHIPMENT_OPTIONS.NTS,
        serviceItemName: 'Fuel Surcharge',
        amount: 100000000000000,
        createdAt: '2020-01-01T00:07:00.999Z',
      },
    ]}
    handleClose={action('clicked')}
  />
);
