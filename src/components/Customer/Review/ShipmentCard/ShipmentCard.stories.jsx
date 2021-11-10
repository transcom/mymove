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
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
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
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
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

const secondaryDeliveryAddress = {
  secondaryDeliveryAddress: {
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
};

const secondaryPickupAddress = {
  secondaryPickupAddress: {
    streetAddress1: '812 S 129th Street',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
};

export default {
  title: 'Customer Components / ShipmentCard',
  decorators: [
    (Story) => (
      <div style={{ padding: 10 }}>
        <Story />
      </div>
    ),
  ],
};

export const HHGShipment = () => <HHGShipmentCard {...hhgDefaultProps} />;

export const HHGShipmentWithSecondaryDestinationAddress = () => (
  <HHGShipmentCard {...hhgDefaultProps} {...secondaryDeliveryAddress} />
);

export const HHGShipmentWithSecondaryPickupAddress = () => (
  <HHGShipmentCard {...hhgDefaultProps} {...secondaryPickupAddress} />
);

export const HHGShipmentWithSecondaryAddresses = () => (
  <HHGShipmentCard {...hhgDefaultProps} {...secondaryPickupAddress} {...secondaryDeliveryAddress} />
);

export const NTSShipment = () => <NTSShipmentCard {...ntsDefaultProps} />;

export const NTSShipmentWithSecondaryPickupAddress = () => (
  <NTSShipmentCard {...ntsDefaultProps} {...secondaryPickupAddress} />
);

export const NTSRShipment = () => <NTSRShipmentCard {...ntsrDefaultProps} />;

export const NTSRShipmentWithSecondaryDestinationAddress = () => (
  <NTSRShipmentCard {...ntsrDefaultProps} {...secondaryDeliveryAddress} />
);

export const PPMShipment = () => <PPMShipmentCard {...ppmDefaultProps} />;
