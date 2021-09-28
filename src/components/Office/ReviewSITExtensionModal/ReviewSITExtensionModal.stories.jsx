import React from 'react';

import ReviewSITExtensionModal from './ReviewSITExtensionModal';

export default {
  title: 'Office Components/ReviewSITEXtensionModal',
  component: ReviewSITExtensionModal,
};

const sitExtension = {
  requestedDays: 45,
  requestReason: 'AWAITING_COMPLETION_OF_RESIDENCE',
  contractorRemarks: 'The customer requested an extension',
  id: '123',
};

export const Basic = () => (
  <ReviewSITExtensionModal sitExtension={sitExtension} onSubmit={() => {}} onClose={() => {}} />
);
