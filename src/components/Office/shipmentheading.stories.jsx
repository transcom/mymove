import React from 'react';
import { object, text, boolean } from '@storybook/addon-knobs';

import ShipmentHeading from './ShipmentHeading';

export const shipmentHeading = () => (
  <ShipmentHeading
    shipmentInfo={{
      shipmentID: text('ShipmentInfo.shipmentID', '1'),
      shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
      originCity: text('ShipmentInfo.originCity', 'San Antonio'),
      originState: text('ShipmentInfo.originState', 'TX'),
      originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
      destinationAddress: object('MTOShipment.destinationAddress', {
        street_address_1: '123 Any Street',
        city: 'Tacoma',
        state: 'WA',
        postal_code: '98421',
      }),
      scheduledPickupDate: text('ShipmentInfo.scheduledPickupDate', '27 Mar 2020'),
    }}
  />
);

export const shipmentHeadingDiversion = () => (
  <ShipmentHeading
    shipmentInfo={{
      shipmentID: text('ShipmentInfo.shipmentID', '1'),
      shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
      isDiversion: boolean('ShipmentInfo.isDiversion', true),
      originCity: text('ShipmentInfo.originCity', 'San Antonio'),
      originState: text('ShipmentInfo.originState', 'TX'),
      originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
      destinationAddress: object('MTOShipment.destinationAddress', {
        street_address_1: '123 Any Street',
        city: 'Tacoma',
        state: 'WA',
        postal_code: '98421',
      }),
      scheduledPickupDate: text('ShipmentInfo.scheduledPickupDate', '27 Mar 2020'),
    }}
  />
);

export const shipmentHeadingCancelled = () => (
  <ShipmentHeading
    shipmentInfo={{
      shipmentID: text('ShipmentInfo.shipmentID', '1'),
      shipmentStatus: text('ShipmentInfo.shipmentStatus', 'CANCELED'),
      shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
      isDiversion: boolean('ShipmentInfo.isDiversion', false),
      originCity: text('ShipmentInfo.originCity', 'San Antonio'),
      originState: text('ShipmentInfo.originState', 'TX'),
      originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
      destinationAddress: object('MTOShipment.destinationAddress', {
        street_address_1: '123 Any Street',
        city: 'Tacoma',
        state: 'WA',
        postal_code: '98421',
      }),
      scheduledPickupDate: text('ShipmentInfo.scheduledPickupDate', '27 Mar 2020'),
    }}
  />
);

export default { title: 'Office Components/ShipmentHeading' };
