import React from 'react';

import ShipmentSITDisplay from '../ShipmentSITDisplay/ShipmentSITDisplay';
import { SITStatusOrigin } from '../ShipmentSITDisplay/ShipmentSITDisplayTestParams';

import SubmitSITExtensionModal from './SubmitSITExtensionModal';

export default {
  title: 'Office Components/SubmitSITExtensionModal',
  component: SubmitSITExtensionModal,
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
    <SubmitSITExtensionModal Submit={() => {}} onClose={() => {}} summarySITComponent={summarySITExtension} />
  </div>
);
