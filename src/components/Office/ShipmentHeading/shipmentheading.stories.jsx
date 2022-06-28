import React from 'react';
import { object, text, boolean } from '@storybook/addon-knobs';

import ShipmentHeading from './ShipmentHeading';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export const shipmentHeading = () => (
  <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
    <ShipmentHeading
      shipmentInfo={{
        shipmentID: text('ShipmentInfo.shipmentID', '1'),
        shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
        originCity: text('ShipmentInfo.originCity', 'San Antonio'),
        originState: text('ShipmentInfo.originState', 'TX'),
        originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
        destinationAddress: object('MTOShipment.destinationAddress', {
          streetAddress1: '123 Any Street',
          city: 'Tacoma',
          state: 'WA',
          postalCode: '98421',
        }),
        scheduledPickupDate: text('ShipmentInfo.scheduledPickupDate', '27 Mar 2020'),
      }}
    />
  </MockProviders>
);

export const shipmentHeadingDiversion = () => (
  <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
    <ShipmentHeading
      shipmentInfo={{
        shipmentID: text('ShipmentInfo.shipmentID', '1'),
        shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
        isDiversion: boolean('ShipmentInfo.isDiversion', true),
        originCity: text('ShipmentInfo.originCity', 'San Antonio'),
        originState: text('ShipmentInfo.originState', 'TX'),
        originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
        destinationAddress: object('MTOShipment.destinationAddress', {
          streetAddress1: '123 Any Street',
          city: 'Tacoma',
          state: 'WA',
          postalCode: '98421',
        }),
        scheduledPickupDate: text('ShipmentInfo.scheduledPickupDate', '27 Mar 2020'),
      }}
    />
  </MockProviders>
);

export const shipmentHeadingCancelled = () => (
  <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
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
          streetAddress1: '123 Any Street',
          city: 'Tacoma',
          state: 'WA',
          postalCode: '98421',
        }),
        scheduledPickupDate: text('ShipmentInfo.scheduledPickupDate', '27 Mar 2020'),
      }}
    />
  </MockProviders>
);

export default { title: 'Office Components/ShipmentHeading' };
