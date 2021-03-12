/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import HHGShipmentCard from './HHGShipmentCard';
import PPMShipmentCard from './PPMShipmentCard';
import NTSShipmentCard from './NTSShipmentCard';
import NTSRShipmentCard from './NTSRShipmentCard';

const hhgDefaultProps = {
  moveId: 'testMove123',
  shipmentNumber: 1,
  shipmentType: 'HHG',
  shipmentId: 'ABC123K',
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  pickupLocation: {
    street_address_1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postal_code: '111111',
  },
  releasingAgent: {
    firstName: 'Jo',
    lastName: 'Xi',
    phone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
  requestedDeliveryDate: new Date('03/01/2020').toISOString(),
  destinationZIP: '73523',
  receivingAgent: {
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    phone: '(999) 999-9999',
    email: 'dorothy.lagomarsino@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

const ntsDefaultProps = {
  moveId: 'testMove123',
  shipmentType: 'HHG_INTO_NTS_DOMESTIC',
  shipmentId: 'ABC123K',
  showEditBtn: true,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  pickupLocation: {
    street_address_1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postal_code: '111111',
  },
  releasingAgent: {
    firstName: 'Jo',
    lastName: 'Xi',
    phone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

const ntsrDefaultProps = {
  moveId: 'testMove123',
  shipmentNumber: 1,
  shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
  shipmentId: 'ABC123K',
  showEditBtn: true,
  requestedDeliveryDate: new Date('03/01/2020').toISOString(),
  receivingAgent: {
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    phone: '(999) 999-9999',
    email: 'dorothy.lagomarsino@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

const ppmDefaultProps = {
  moveId: 'testMove123',
  destinationZIP: '11111',
  estimatedWeight: '5,000',
  expectedDepartureDate: new Date('01/01/2020').toISOString(),
  shipmentId: 'ABC123K',
  sitDays: '24',
  originZIP: '00000',
};

export default {
  title: 'Customer Components / ShipmentCard',
};

export const HHGShipment = () => (
  <div style={{ padding: 10 }}>
    <HHGShipmentCard {...hhgDefaultProps} />
  </div>
);

export const NTSShipment = () => (
  <div style={{ padding: 10 }}>
    <NTSShipmentCard {...ntsDefaultProps} />
  </div>
);

export const NTSRShipment = () => (
  <div style={{ padding: 10 }}>
    <NTSRShipmentCard {...ntsrDefaultProps} />
  </div>
);

export const PPMShipment = () => (
  <div style={{ padding: 10 }}>
    <PPMShipmentCard {...ppmDefaultProps} />
  </div>
);
