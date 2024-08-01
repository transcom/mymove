import React from 'react';
import { object, text } from '@storybook/addon-knobs';

import ShipmentInfoList from './ShipmentInfoList';

import { ADDRESS_UPDATE_STATUS } from 'constants/shipments';

export default {
  title: 'Office Components/HHG Shipment Info List',
  component: ShipmentInfoList,
};

const info = {
  requestedPickupDate: '2021-06-01',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  secondaryPickupAddress: {
    streetAddress1: '444 S 131st St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  tertiaryPickupAddress: {
    streetAddress1: '1909 Electric Ave',
    city: 'Missoula',
    state: 'MT',
    postalCode: '59801',
  },
  destinationAddress: {
    streetAddress1: '7 Q St',
    city: 'Austin',
    state: 'TX',
    postalCode: '78722',
  },
  deliveryAddressUpdate: {
    status: 'REQUESTED',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '17 8th St',
    city: 'Austin',
    state: 'TX',
    postalCode: '78751',
  },
  tertiaryDeliveryAddress: {
    streetAddress1: '14 9th St',
    city: 'Austin',
    state: 'TX',
    postalCode: '78752',
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

export const WithTertiaryAddresses = () => (
  <ShipmentInfoList
    shipment={{
      requestedPickupDate: text('requestedPickupDate', info.requestedPickupDate),
      pickupAddress: object('pickupAddress', info.pickupAddress),
      secondaryPickupAddress: object('secondaryPickupAddress', info.secondaryPickupAddress),
      tertiaryPickupAddress: object('tertiaryPickupAddress', info.tertiaryPickupAddress),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      secondaryDeliveryAddress: object('secondaryDeliveryAddress', info.secondaryDeliveryAddress),
      tertiaryDeliveryAddress: object('tertiaryDeliveryAddress', info.tertiaryDeliveryAddress),
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
      tertiaryPickupAddress: object('tertiaryPickupAddress', info.tertiaryPickupAddress),
      destinationAddress: object('destinationAddress', info.destinationAddress),
      secondaryDeliveryAddress: object('secondaryDeliveryAddress', info.secondaryDeliveryAddress),
      tertiaryDeliveryAddress: object('tertiaryDeliveryAddress', info.tertiaryDeliveryAddress),
      agents: [object('agent1', info.agents[0]), object('agent2', info.agents[1])],
      counselorRemarks: text('counselorRemarks', info.counselorRemarks),
      customerRemarks: text('customerRemarks', info.customerRemarks),
    }}
  />
);

export const WithDeliveryAddressUpdateRequested = () => (
  <ShipmentInfoList
    shipment={{
      requestedPickupDate: text('requestedPickupDate', info.requestedPickupDate),
      pickupAddress: object('pickupAddress', info.pickupAddress),
      deliveryAddressUpdate: object('deliveryAddressUpdate', info.deliveryAddressUpdate),
      destinationAddress: object('destinationAddress', info.destinationAddress),
    }}
    errorIfMissing={[
      {
        fieldName: 'destinationAddress',
        condition: (shipment) => shipment.deliveryAddressUpdate?.status === ADDRESS_UPDATE_STATUS.REQUESTED,
        optional: true,
      },
    ]}
  />
);
