import React from 'react';
import { object, text } from '@storybook/addon-knobs';

import NTSRShipmentInfoList from './NTSRShipmentInfoList';

export default {
  title: 'Office Components/Shipment Info List',
  component: NTSRShipmentInfoList,
};

const info = {
  primeActualWeight: 2000,
  storageFacility: {
    address: {
      city: 'Anytown',
      country: 'USA',
      postalCode: '90210',
      state: 'OK',
      streetAddress1: '555 Main Ave',
      streetAddress2: 'Apartment 900',
    },
    facilityName: 'my storage',
    lotNumber: '2222',
  },
  serviceOrderNumber: '12341234',
  requestedDeliveryDate: '26 Mar 2020',
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  agents: [
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
  tacType: 'HHG',
  sacType: 'NTS',
};

export const NTSRBasic = () => (
  <NTSRShipmentInfoList
    shipment={{
      counselorRemarks: text('counselorRemarks', info.counselorRemarks),
      requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
      storageFacility: object('storageFacility', info.storageFacility),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      tacType: text('tacType', info.tacType),
      sacType: text('sacType', info.sacType),
      primeActualWeight: text('primeActualWeight', info.primeActualWeight),
      serviceOrderNumber: text('serviceOrderNumber', info.serviceOrderNumber),
    }}
  />
);

export const NTSRMissingInfo = () => (
  <NTSRShipmentInfoList
    shipment={{
      requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      tacType: text('tacType', info.tacType),
      sacType: text('sacType', info.sacType),
      primeActualWeight: text('primeActualWeight', info.primeActualWeight),
      serviceOrderNumber: text('serviceOrderNumber', info.serviceOrderNumber),
    }}
    errorIfMissing={['storageFacility']}
  />
);

export const WithAgents = () => (
  <NTSRShipmentInfoList
    shipment={{
      counselorRemarks: text('counselorRemarks', info.counselorRemarks),
      requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
      storageFacility: object('storageFacility', info.storageFacility),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      tacType: text('tacType', info.tacType),
      sacType: text('sacType', info.sacType),
      primeActualWeight: text('primeActualWeight', info.primeActualWeight),
      serviceOrderNumber: text('serviceOrderNumber', info.serviceOrderNumber),
      agents: [object('agent1', info.agents[0]), object('agent2', info.agents[1])],
    }}
  />
);

export const WithRemarks = () => (
  <NTSRShipmentInfoList
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
  <NTSRShipmentInfoList
    shipment={{
      requestedDeliveryDate: text('requestedDeliveryDate', info.requestedDeliveryDate),
      storageFacility: object('storageFacility', info.storageFacility),
      tacType: text('tacType', info.tacType),
      sacType: text('sacType', info.sacType),
      primeActualWeight: text('primeActualWeight', info.primeActualWeight),
      serviceOrderNumber: text('serviceOrderNumber', info.serviceOrderNumber),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      secondaryDeliveryAddress: object('secondaryDeliveryAddress', info.secondaryDeliveryAddress),
      agents: [object('agent1', info.agents[0]), object('agent2', info.agents[1])],
      counselorRemarks: text('counselorRemarks', info.counselorRemarks),
      customerRemarks: text('customerRemarks', info.customerRemarks),
    }}
  />
);
