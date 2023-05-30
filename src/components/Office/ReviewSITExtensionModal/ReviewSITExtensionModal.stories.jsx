import React from 'react';
import MockDate from 'mockdate';
import { addons } from '@storybook/preview-api';
import { isHappoRun } from 'happo-plugin-storybook/register';

import ReviewSITExtensionModal from './ReviewSITExtensionModal';

// Based on sitStatus below. The date is 31 days after the entry date.
const mockedDate = '2023-09-22T00:00:00.000Z';
export default {
  title: 'Office Components/ReviewSITExtensionModal',
  component: ReviewSITExtensionModal,
  decorators: [
    (Story) => {
      if (isHappoRun()) {
        MockDate.set(mockedDate);
        addons.getChannel().on('storyRendered', MockDate.reset);
      }
      return <Story />;
    },
  ],
};

const sitExtension = {
  requestedDays: 45,
  requestReason: 'AWAITING_COMPLETION_OF_RESIDENCE',
  contractorRemarks: 'The customer requested an extension',
  status: 'PENDING',
  id: '123',
};

const sitStatus = {
  totalDaysRemaining: 30,
  sitEntryDate: new Date('22 Aug 2023'),
  totalSITDaysUsed: 15,
  daysInSIT: 15,
};

const shipment = {
  sitDaysAllowance: 45,
};

export const Basic = () => (
  <div className="officeApp">
    <ReviewSITExtensionModal
      sitExtension={sitExtension}
      onSubmit={() => {}}
      onClose={() => {}}
      shipment={shipment}
      sitStatus={sitStatus}
    />
  </div>
);
