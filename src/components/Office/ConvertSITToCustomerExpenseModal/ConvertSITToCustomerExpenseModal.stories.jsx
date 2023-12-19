import React from 'react';
import MockDate from 'mockdate';
import { addons } from '@storybook/preview-api';

import ConvertSITToCustomerExpenseModal from './ConvertSITToCustomerExpenseModal';

// Based on sitStatus below. The date is 31 days after the entry date.
const mockedDate = '2023-04-19T00:00:00.000Z';
export default {
  title: 'Office Components/ConvertSITToCustomerExpenseModal',
  component: ConvertSITToCustomerExpenseModal,
  decorators: [
    (Story) => {
      MockDate.set(mockedDate);
      addons.getChannel().on('storyRendered', MockDate.reset);
      return (
        <div style={{ padding: '1em', backgroundColor: '#f9f9f9' }}>
          <Story />
        </div>
      );
    },
  ],
  parameters: {
    docs: {
      inlineStories: false,
    },
  },
};

const sitStatus = {
  totalDaysRemaining: 210,
  totalSITDaysUsed: 60,
  currentSIT: {
    location: 'DESTINATION',
    daysInSIT: 30,
    sitEntryDate: '2023-03-19T00:00:00.000Z',
  },
};

export const Basic = () => (
  <div className="officeApp">
    <ConvertSITToCustomerExpenseModal
      Submit={() => {}}
      onClose={() => {}}
      sitStatus={sitStatus}
      shipment={{ sitDaysAllowance: 270 }}
    />
  </div>
);
