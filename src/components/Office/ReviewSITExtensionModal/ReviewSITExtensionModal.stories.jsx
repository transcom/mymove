import React from 'react';

import ShipmentSITDisplay from '../ShipmentSITDisplay/ShipmentSITDisplay';

import ReviewSITExtensionModal from './ReviewSITExtensionModal';

import { SITStatusOrigin } from 'components/Office/ShipmentSITDisplay/ShipmentSITDisplayTestParams';

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

const summarySITExtension = (
  <div className="officeApp">
    <ShipmentSITDisplay
      {...{
        sitExtensions: [sitExtension],
        sitStatus: SITStatusOrigin,
        shipment: { sitDaysAllowance: 90 },
        hideSITExtensionAction: true,
      }}
    />
  </div>
);

export const Basic = () => (
  <div className="officeApp">
    <ReviewSITExtensionModal
      sitExtension={sitExtension}
      onSubmit={() => {}}
      onClose={() => {}}
      summarySITComponent={summarySITExtension}
    />
  </div>
);
