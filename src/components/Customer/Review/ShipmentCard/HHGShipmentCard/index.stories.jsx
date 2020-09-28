/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import HHGShipmentCard from '.';

const defaultProps = {
  shipmentNumber: 1,
  shipmentId: '#ABC123K-001',
  requestedPickupDate: new Date().toISOString(),
  pickupLocation: {
    street_address_1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postal_code: '111111',
  },
  releasingAgent: {
    name: 'Jo Xi',
    telephone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
  requestedDeliveryDate: new Date().toISOString(),
  destinationZIP: '73523',
  receivingAgent: {
    name: 'Dorothy Lagomarsino',
    telephone: '(999) 999-9999',
    email: 'dorothy.lagomarsino@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

export default {
  title: 'Customer Components | HHGShipmentCard',
};

export const Basic = () => (
  <div style={{ padding: 10 }}>
    <HHGShipmentCard {...defaultProps} />
  </div>
);
