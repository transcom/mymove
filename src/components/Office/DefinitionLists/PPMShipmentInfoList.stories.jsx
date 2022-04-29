import React from 'react';
// import { object, text } from '@storybook/addon-knobs';

import PPMShipmentInfoList from './PPMShipmentInfoList';

export default {
  title: 'Office Components/PPM Shipment Info List',
  component: PPMShipmentInfoList,
};

const ppmInfo = {
  expectedDepartureDate: '2021-06-01',
  pickupPostalCode: '08540',
  destinationPostalCode: '90120',
  sitExpected: false,
  counselorRemarks:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ' +
    'Morbi porta nibh nibh, ac malesuada tortor egestas.',
  customerRemarks: 'Ut enim ad minima veniam',
};

export const Basic = () => <PPMShipmentInfoList shipment={ppmInfo} />;
