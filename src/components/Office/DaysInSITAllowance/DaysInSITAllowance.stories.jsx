import React from 'react';

import DaysInSITAllowance from './DaysInSITAllowance';

export default {
  title: 'Office Components/DaysInSITAllowance',
  component: DaysInSITAllowance,
  argTypes: {
    previouslyBilledDays: {
      type: 'number',
      defaultValue: 30,
      required: false,
    },
    previouslyBilledEndDate: {
      type: 'string',
      defaultValue: '2021-06-08',
      required: false,
    },
    pendingSITDaysInvoiced: {
      type: 'number',
      defaultValue: 60,
    },
    pendingBilledEndDate: {
      type: 'string',
      defaultValue: '2021-08-08',
    },
    totalSITDaysAuthorized: {
      type: 'number',
      defaultValue: 120,
    },
    totalSITDaysRemaining: {
      type: 'number',
      defaultValue: 30,
    },
    totalSITEndDate: {
      type: 'string',
      defaultValue: '2021-09-08',
    },
  },
};

const defaultArgs = {
  previouslyBilledDays: 30,
  previouslyBilledEndDate: '2021-06-08',
  pendingSITDaysInvoiced: 60,
  pendingBilledEndDate: '2021-08-08',
  totalSITDaysAuthorized: 120,
  totalSITDaysRemaining: 30,
  totalSITEndDate: '2021-09-08',
};

const Template = (args) => <DaysInSITAllowance shipmentPaymentSITBalance={{ ...defaultArgs, ...args }} />;

export const PastPendingRemaining = Template.bind({});

export const NoPastBilledDays = Template.bind({});

NoPastBilledDays.args = {
  previouslyBilledDays: 0,
  previouslyBilledEndDate: undefined,
  totalSITDaysRemaining: 60,
  totalSITEndDate: '2021-10-07',
};
