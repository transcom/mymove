import React from 'react';

import PPMSummaryList from './PPMSummaryList';

import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';

export default {
  title: 'Components / PPMSummaryList',
  component: PPMSummaryList,
  argTypes: {
    onUploadClick: { action: 'upload button clicked' },
  },
};

const Template = (args) => <PPMSummaryList {...args} />;

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
        status: ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL,
        approvedAt: '2022-04-15T15:38:07.103Z',
        submittedAt: '2022-04-19T15:38:07.103Z',
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
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
        status: ppmShipmentStatuses.PAYMENT_APPROVED,
        approvedAt: '2022-04-15T15:38:07.103Z',
        submittedAt: '2022-04-19T15:38:07.103Z',
        reviewedAt: '2022-04-23T15:38:07.103Z',
        hasRequestedAdvance: true,
        advanceAmountRequested: 10000,
      },
    },
  ],
};
