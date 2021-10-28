import React from 'react';

import ShipmentSITDisplay from '../ShipmentSITDisplay/ShipmentSITDisplay';
import { SITStatusOrigin } from '../ShipmentSITDisplay/ShipmentSITDisplayTestParams';

import FinancialReviewModal from './FinancialReviewModal';

export default {
  title: 'Office Components/FinancialReviewModal',
  component: FinancialReviewModal,
};

const sitExtension = {
  requestedDays: 45,
  requestReason: 'AWAITING_COMPLETION_OF_RESIDENCE',
  contractorRemarks: 'The customer requested an extension',
  status: 'PENDING',
  id: '123',
};

const summarySITExtension = (
  <ShipmentSITDisplay
    {...{
      sitExtensions: [sitExtension],
      sitStatus: SITStatusOrigin,
      shipment: { sitDaysAllowance: 90 },
      hideSITExtensionAction: true,
    }}
  />
);

export const Basic = () => (
  <FinancialReviewModal Submit={() => {}} onClose={() => {}} summarySITComponent={summarySITExtension} />
);
