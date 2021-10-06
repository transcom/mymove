import React from 'react';

import ShipmentSITExtensions from '../ShipmentSITExtensions/ShipmentSITExtensions';
import { SITStatusOrigin } from '../ShipmentSITExtensions/ShipmentSITExtensionsTestParams';

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
  <ShipmentSITExtensions
    {...{
      sitExtensions: [sitExtension],
      sitStatus: SITStatusOrigin,
      shipment: { sitDaysAllowance: 90 },
      hideSITExtensionAction: true,
    }}
  />
);

export const Basic = () => (
  <SubmitSITExtensionModal Submit={() => {}} onClose={() => {}} summarySITComponent={summarySITExtension} />
);
