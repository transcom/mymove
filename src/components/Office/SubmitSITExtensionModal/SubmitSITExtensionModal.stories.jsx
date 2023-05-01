import React from 'react';

import SubmitSITExtensionModal from './SubmitSITExtensionModal';

export default {
  title: 'Office Components/SubmitSITExtensionModal',
  component: SubmitSITExtensionModal,
};

const sitStatus = {
  daysInSIT: 30,
  location: 'DESTINATION',
  sitEntryDate: '2023-03-19T00:00:00.000Z',
  totalDaysRemaining: 210,
  totalSITDaysUsed: 60,
};

export const Basic = () => (
  <div className="officeApp">
    <SubmitSITExtensionModal
      Submit={() => {}}
      onClose={() => {}}
      sitStatus={sitStatus}
      shipment={{ sitDaysAllowance: 270 }}
    />
  </div>
);
