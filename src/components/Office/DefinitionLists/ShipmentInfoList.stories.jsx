import React from 'react';
import { object, text } from '@storybook/addon-knobs';

import ShipmentInfoList from './ShipmentInfoList';

export default {
  title: 'Office Components/ShipmentInfoList',
  component: ShipmentInfoList,
};

const info = {
  requestedMoveDate: '2021-06-01',
  originAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  secondPickupAddress: {
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
  secondDestinationAddress: {
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
    requestedMoveDate={text('requestedMoveDate', info.requestedMoveDate)}
    originAddress={object('originAddress', info.originAddress)}
    destinationAddress={object('destinationAddress', info.destinationAddress)}
  />
);

export const WithSecondaryAddresses = () => (
  <ShipmentInfoList
    requestedMoveDate={text('requestedMoveDate', info.requestedMoveDate)}
    originAddress={object('originAddress', info.originAddress)}
    secondPickupAddress={object('secondPickupAddress', info.secondPickupAddress)}
    destinationAddress={object('destinationAddress', info.destinationAddress)}
    secondDestinationAddress={object('secondDestinationAddress', info.secondDestinationAddress)}
  />
);

export const WithAgents = () => (
  <ShipmentInfoList
    requestedMoveDate={text('requestedMoveDate', info.requestedMoveDate)}
    originAddress={object('originAddress', info.originAddress)}
    destinationAddress={object('destinationAddress', info.destinationAddress)}
    agents={[object('agent1', info.agents[0]), object('agent2', info.agents[1])]}
  />
);

export const WithRemarks = () => (
  <ShipmentInfoList
    requestedMoveDate={text('requestedMoveDate', info.requestedMoveDate)}
    originAddress={object('originAddress', info.originAddress)}
    destinationAddress={object('destinationAddress', info.destinationAddress)}
    counselorRemarks={text('counselorRemarks', info.counselorRemarks)}
    customerRemarks={text('customerRemarks', info.customerRemarks)}
  />
);

export const WithAllInfo = () => (
  <ShipmentInfoList
    requestedMoveDate={text('requestedMoveDate', info.requestedMoveDate)}
    originAddress={object('originAddress', info.originAddress)}
    secondPickupAddress={object('secondPickupAddress', info.secondPickupAddress)}
    destinationAddress={object('destinationAddress', info.destinationAddress)}
    secondDestinationAddress={object('secondDestinationAddress', info.secondDestinationAddress)}
    agents={[object('agent1', info.agents[0]), object('agent2', info.agents[1])]}
    counselorRemarks={text('counselorRemarks', info.counselorRemarks)}
    customerRemarks={text('customerRemarks', info.customerRemarks)}
  />
);
