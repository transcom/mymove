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

export const SingleDisabled = Template.bind({});
SingleDisabled.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.SUBMITTED,
      ppmShipment: { id: '11', status: ppmShipmentStatuses.SUBMITTED, advanceRequested: true, advance: 10000 },
    },
  ],
};

export const SingleEnabled = Template.bind({});
SingleEnabled.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '11',
        status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
        approvedAt: '2022-04-15T15:38:07.103Z',
        advanceRequested: true,
        advance: 10000,
      },
    },
  ],
};

export const MultipleEnabled = Template.bind({});
MultipleEnabled.args = {
  shipments: [
    {
      id: '1',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '11',
        status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
        approvedAt: '2022-04-15T15:38:07.103Z',
        advanceRequested: true,
        advance: 10000,
      },
    },
    {
      id: '2',
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        id: '2',
        status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
        approvedAt: '2022-04-20T15:38:07.103Z',
        advanceRequested: true,
        advance: 10000,
      },
    },
  ],
};
