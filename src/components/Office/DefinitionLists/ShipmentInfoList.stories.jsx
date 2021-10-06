import React from 'react';
import { object, text } from '@storybook/addon-knobs';

import ShipmentInfoList from './ShipmentInfoList';

export default {
  title: 'Office Components/Shipment Info List',
  component: ShipmentInfoList,
};

const info = {
  requestedPickupDate: '2021-06-01',
  pickupAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  secondaryPickupAddress: {
    street_address_1: '444 S 131st St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '7 Q St',
    city: 'Austin',
    state: 'TX',
    postal_code: '78722',
  },
  secondaryDeliveryAddress: {
    street_address_1: '17 8th St',
    city: 'Austin',
    state: 'TX',
    postal_code: '78751',
  },
  agents: [
    {
      agentType: 'RELEASING_AGENT',
      firstName: 'Quinn',
      lastName: 'Ocampo',
      phone: '999-999-9999',
      email: 'quinnocampo@myemail.com',
    },
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  counselorRemarks:
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam vulputate commodo erat. ' +
    'Morbi porta nibh nibh, ac malesuada tortor egestas.',
  customerRemarks: 'Ut enim ad minima veniam',
};

export const Basic = () => (
  <ShipmentInfoList
    shipment={{
      requestedPickupDate: text('requestedPickupDate', info.requestedPickupDate),
      pickupAddress: object('pickupAddress', info.pickupAddress),
      destinationAddress: object('destinationAddress', info.destinationAddress),
    }}
  />
);

export const WithSecondaryAddresses = () => (
  <ShipmentInfoList
    shipment={{
      requestedPickupDate: text('requestedPickupDate', info.requestedPickupDate),
      pickupAddress: object('pickupAddress', info.pickupAddress),
      secondaryPickupAddress: object('secondaryPickupAddress', info.secondaryPickupAddress),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      secondaryDeliveryAddress: object('secondaryDeliveryAddress', info.secondaryDeliveryAddress),
    }}
  />
);

export const WithAgents = () => (
  <ShipmentInfoList
    shipment={{
      requestedPickupDate: text('requestedPickupDate', info.requestedPickupDate),
      pickupAddress: object('pickupAddress', info.pickupAddress),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      agents: [object('agent1', info.agents[0]), object('agent2', info.agents[1])],
    }}
  />
);

export const WithRemarks = () => (
  <ShipmentInfoList
    shipment={{
      requestedPickupDate: text('requestedPickupDate', info.requestedPickupDate),
      pickupAddress: object('pickupAddress', info.pickupAddress),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      counselorRemarks: text('counselorRemarks', info.counselorRemarks),
      customerRemarks: text('customerRemarks', info.customerRemarks),
    }}
  />
);

export const WithAllInfo = () => (
  <ShipmentInfoList
    shipment={{
      requestedPickupDate: text('requestedPickupDate', info.requestedPickupDate),
      pickupAddress: object('pickupAddress', info.pickupAddress),
      secondaryPickupAddress: object('secondaryPickupAddress', info.secondaryPickupAddress),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      secondaryDeliveryAddress: object('secondaryDeliveryAddress', info.secondaryDeliveryAddress),
      agents: [object('agent1', info.agents[0]), object('agent2', info.agents[1])],
      counselorRemarks: text('counselorRemarks', info.counselorRemarks),
      customerRemarks: text('customerRemarks', info.customerRemarks),
    }}
  />
);
