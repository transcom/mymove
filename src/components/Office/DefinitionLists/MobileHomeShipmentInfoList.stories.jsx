import React from 'react';
import { object, text } from '@storybook/addon-knobs';

import MobileHomeShipmentInfoList from './MobileHomeShipmentInfoList';

export default {
  title: 'Office Components/Mobile Home Shipment Info List',
  component: MobileHomeShipmentInfoList,
};

const mobileHomeShipment = {
  mobileHomeShipment: {
    year: 2020,
    make: 'Yamaha',
    model: '242X E-Series',
    lengthInInches: 276,
    widthInInches: 102,
    heightInInches: 120,
  },
  pickupAddress: {
    streetAddress1: '123 Harbor Dr',
    city: 'Miami',
    state: 'FL',
    postalCode: '33101',
  },
  destinationAddress: {
    streetAddress1: '456 Marina Blvd',
    city: 'Key West',
    state: 'FL',
    postalCode: '33040',
  },
  secondaryPickupAddress: {
    streetAddress1: '789 Seaport Ln',
    city: 'Fort Lauderdale',
    state: 'FL',
    postalCode: '33316',
  },
  tertiaryPickupAddress: {
    streetAddress1: '101 Yacht Club Rd',
    city: 'Naples',
    state: 'FL',
    postalCode: '34102',
  },
  secondaryDeliveryAddress: {
    streetAddress1: '111 Ocean Dr',
    city: 'Palm Beach',
    state: 'FL',
    postalCode: '33480',
  },
  tertiaryDeliveryAddress: {
    streetAddress1: '222 Shoreline Dr',
    city: 'Clearwater',
    state: 'FL',
    postalCode: '33767',
  },
  mtoAgents: [
    {
      agentType: 'RELEASING_AGENT',
      firstName: 'John',
      lastName: 'Doe',
      phone: '123-456-7890',
      email: 'john.doe@example.com',
    },
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Jane',
      lastName: 'Doe',
      phone: '987-654-3210',
      email: 'jane.doe@example.com',
    },
  ],
  counselorRemarks: 'Please be cautious with the mobile home.',
  customerRemarks: 'Handle with care.',
};

export const Basic = () => (
  <MobileHomeShipmentInfoList
    shipment={{
      mobileHomeShipment: {
        year: text('year', mobileHomeShipment.year),
        make: text('make', mobileHomeShipment.make),
        model: text('model', mobileHomeShipment.model),
        lengthInInches: text('lengthInInches', mobileHomeShipment.lengthInInches),
        widthInInches: text('widthInInches', mobileHomeShipment.widthInInches),
        heightInInches: text('heightInInches', mobileHomeShipment.heightInInches),
      },
      pickupAddress: object('pickupAddress', mobileHomeShipment.pickupAddress),
      destinationAddress: object('destinationAddress', mobileHomeShipment.destinationAddress),
    }}
  />
);

export const DefaultView = () => (
  <MobileHomeShipmentInfoList
    shipment={{
      ...mobileHomeShipment,
    }}
  />
);
