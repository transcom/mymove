import React from 'react';

import ReviewSITExtensionModal from './ReviewSITExtensionModal';

export default {
  title: 'Office Components/ReviewSITExtensionModal',
  component: ReviewSITExtensionModal,
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
