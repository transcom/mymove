import React from 'react';

import PPMSummaryList from './PPMSummaryList';

import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';
import { MockProviders } from 'testUtils';

export default {
  title: 'Components / PPMSummaryList',
  component: PPMSummaryList,
  argTypes: {
    onUploadClick: { action: 'upload button clicked' },
  },
};

const Template = (args) => (
  <MockProviders>
    <PPMSummaryList {...args} />
  </MockProviders>
);

export const Submitted = Template.bind({});
Submitted.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.SUBMITTED,
      ppmShipment: {
        id: '11',
        status: ppmShipmentStatuses.SUBMITTED,
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
        pickupAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Pickup Test City',
          state: 'NY',
          postalCode: '10001',
        },
        destinationAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Destination Test City',
          state: 'NY',
          postalCode: '11111',
        },
      },
    },
  ],
};

export const Approved = Template.bind({});
Approved.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '11',
        status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
        approvedAt: '2022-04-15T15:38:07.103Z',
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
        pickupAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Pickup Test City',
          state: 'NY',
          postalCode: '10001',
        },
        destinationAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Destination Test City',
          state: 'NY',
          postalCode: '11111',
        },
      },
    },
  ],
};

export const ApprovedMultiple = Template.bind({});
ApprovedMultiple.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '11',
        status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
        approvedAt: '2022-04-15T15:38:07.103Z',
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
        pickupAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Pickup Test City',
          state: 'NY',
          postalCode: '10001',
        },
        destinationAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Destination Test City',
          state: 'NY',
          postalCode: '11111',
        },
      },
    },
    {
      id: '2',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '2',
        status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
        approvedAt: '2022-04-20T15:38:07.103Z',
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
        pickupAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Pickup Test City',
          state: 'NY',
          postalCode: '10001',
        },
        destinationAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Destination Test City',
          state: 'NY',
          postalCode: '11111',
        },
      },
    },
  ],
};

export const PaymentSubmitted = Template.bind({});
PaymentSubmitted.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '11',
        status: ppmShipmentStatuses.NEEDS_CLOSEOUT,
        approvedAt: '2022-04-15T15:38:07.103Z',
        submittedAt: '2022-04-19T15:38:07.103Z',
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
        pickupAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Pickup Test City',
          state: 'NY',
          postalCode: '10001',
        },
        destinationAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Destination Test City',
          state: 'NY',
          postalCode: '11111',
        },
      },
    },
  ],
};

export const PaymentReviewed = Template.bind({});
PaymentReviewed.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '11',
        status: ppmShipmentStatuses.CLOSEOUT_COMPLETE,
        approvedAt: '2022-04-15T15:38:07.103Z',
        submittedAt: '2022-04-19T15:38:07.103Z',
        reviewedAt: '2022-04-23T15:38:07.103Z',
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
        pickupAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Pickup Test City',
          state: 'NY',
          postalCode: '10001',
        },
        destinationAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Destination Test City',
          state: 'NY',
          postalCode: '11111',
        },
      },
    },
  ],
};
