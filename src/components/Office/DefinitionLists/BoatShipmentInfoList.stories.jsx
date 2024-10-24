import React from 'react';
import { object, text, boolean } from '@storybook/addon-knobs';

import BoatShipmentInfoList from './BoatShipmentInfoList';

import { boatShipmentTypes } from 'constants/shipments';

export default {
  title: 'Office Components/Boat Shipment Info List',
  component: BoatShipmentInfoList,
};

const boatShipment = {
  boatShipment: {
    type: boatShipmentTypes.TOW_AWAY,
    year: 2020,
    make: 'Yamaha',
    model: '242X E-Series',
    lengthInInches: 276,
    widthInInches: 102,
    heightInInches: 120,
    hasTrailer: true,
    isRoadworthy: true,
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
  counselorRemarks: 'Please be cautious with the boat hull.',
  customerRemarks: 'Handle with care.',
};

export const Basic = () => (
  <BoatShipmentInfoList
    shipment={{
      boatShipment: {
        type: text('type', boatShipment.type),
        year: text('year', boatShipment.year),
        make: text('make', boatShipment.make),
        model: text('model', boatShipment.model),
        lengthInInches: text('lengthInInches', boatShipment.lengthInInches),
        widthInInches: text('widthInInches', boatShipment.widthInInches),
        heightInInches: text('heightInInches', boatShipment.heightInInches),
        hasTrailer: boolean('hasTrailer', boatShipment.hasTrailer),
        isRoadworthy: boolean('isRoadworthy', boatShipment.isRoadworthy),
      },
      pickupAddress: object('pickupAddress', boatShipment.pickupAddress),
      destinationAddress: object('destinationAddress', boatShipment.destinationAddress),
    }}
  />
);

export const DefaultView = () => (
  <BoatShipmentInfoList
    shipment={{
      ...boatShipment,
    }}
  />
);
