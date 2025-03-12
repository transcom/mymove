import React from 'react';

import PPMShipmentInfoList from './PPMShipmentInfoList';

import { PPM_TYPES } from 'shared/constants';

export default {
  title: 'Office Components/PPM Shipment Info List',
  component: PPMShipmentInfoList,
};

const ppmInfo = {
  ppmShipment: {
    ppmType: PPM_TYPES.INCENTIVE_BASED,
    actualMoveDate: null,
    hasRequestedAdvance: true,
    advanceAmountRequested: 598700,
    approvedAt: null,
    createdAt: '2022-04-29T21:48:21.581Z',
    eTag: 'MjAyMi0wNC0yOVQyMTo0ODoyMS41ODE0MzFa',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: true,
    id: 'b6ec215c-2cef-45fe-8d4a-35f445cd4768',
    proGearWeight: 1987,
    reviewedAt: null,
    shipmentId: 'b5c2d9a1-d1e6-485d-9678-8b62deb0d801',
    spouseProGearWeight: 498,
    status: 'SUBMITTED',
    submittedAt: '2022-04-29T21:48:21.573Z',
  },
};

export const Basic = () => (
  <PPMShipmentInfoList
    shipment={{
      ...ppmInfo,
      counselorRemarks: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ',
    }}
    warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
    isExpanded
  />
);

export const DefaultView = () => (
  <PPMShipmentInfoList
    shipment={{
      ...ppmInfo,
      counselorRemarks: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ',
    }}
    warnIfMissing={[{ fieldName: 'counselorRemarks' }]}
  />
);

export const Warning = () => (
  <PPMShipmentInfoList shipment={ppmInfo} warnIfMissing={[{ fieldName: 'counselorRemarks' }]} />
);

export const MissingInfo = () => (
  <PPMShipmentInfoList
    shipment={ppmInfo}
    errorIfMissing={[
      {
        fieldName: 'advanceStatus',
        condition: (shipment) => shipment?.ppmShipment?.hasRequestedAdvance === true,
      },
    ]}
  />
);
