import React from 'react';
// import { object, text } from '@storybook/addon-knobs';

import PPMShipmentInfoList from './PPMShipmentInfoList';

export default {
  title: 'Office Components/PPM Shipment Info List',
  component: PPMShipmentInfoList,
};

const ppmInfo = {
  ppmShipment: {
    actualMoveDate: null,
    hasRequestedAdvance: true,
    advanceAmountRequested: 598700,
    approvedAt: null,
    createdAt: '2022-04-29T21:48:21.581Z',
    deletedAt: null,
    destinationPostalCode: '30813',
    eTag: 'MjAyMi0wNC0yOVQyMTo0ODoyMS41ODE0MzFa',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: true,
    id: 'b6ec215c-2cef-45fe-8d4a-35f445cd4768',
    netWeight: null,
    pickupPostalCode: '90210',
    proGearWeight: 1987,
    reviewedAt: null,
    secondaryDestinationPostalCode: '30814',
    secondaryPickupPostalCode: '90211',
    shipmentId: 'b5c2d9a1-d1e6-485d-9678-8b62deb0d801',
    spouseProGearWeight: 498,
    status: 'SUBMITTED',
    submittedAt: '2022-04-29T21:48:21.573Z',
    updatedAt: '2022-04-29T21:48:21.581Z',
  },
};

export const Basic = () => (
  <PPMShipmentInfoList
    shipment={{
      ...ppmInfo,
      counselorRemarks: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ',
    }}
    warnIfMissing={['counselorRemarks']}
    isExpanded
  />
);
export const DefaultView = () => (
  <PPMShipmentInfoList
    shipment={{
      ...ppmInfo,
      counselorRemarks: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ',
    }}
    warnIfMissing={['counselorRemarks']}
  />
);
export const MissingInfo = () => <PPMShipmentInfoList shipment={ppmInfo} warnIfMissing={['counselorRemarks']} />;
