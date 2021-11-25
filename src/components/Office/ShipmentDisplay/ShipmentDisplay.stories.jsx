import React from 'react';
import { object } from '@storybook/addon-knobs';

import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/Shipment Display',
  component: ShipmentDisplay,
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

const hhgInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  secondaryPickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
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
  customerRemarks: 'Ut enim ad minima veniam',
};

const ntsInfo = {
  heading: 'NTS',
  requestedPickupDate: '26 Mar 2020',
  shipmentId: 'testShipmentId394',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const ntsReleaseInfo = {
  heading: 'NTS-release',
  shipmentId: 'testShipmentId394',
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

const postalOnlyInfo = {
  heading: 'HHG',
  requestedPickupDate: '26 Mar 2020',
  shipmentId: 'testShipmentId394',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    postalCode: '98421',
  },
};

const diversionInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  isDiversion: true,
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const cancelledInfo = {
  heading: 'HHG',
  shipmentId: 'testShipmentId394',
  isDiversion: false,
  shipmentStatus: 'CANCELED',
  requestedPickupDate: '26 Mar 2020',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

export const HHGShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} shipmentType={SHIPMENT_OPTIONS.HHG} isSubmitted />
  </div>
);

export const HHGShipmentNoIcon = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={object('displayInfo', hhgInfo)}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
      showIcon={false}
    />
  </div>
);

export const HHGShipmentWithCounselorRemarks = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
      showIcon={false}
    />
  </div>
);

export const HHGShipmentEditable = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
      showIcon={false}
      editURL="/"
    />
  </div>
);

export const NTSShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={object('displayInfo', ntsInfo)} shipmentType={SHIPMENT_OPTIONS.NTS} isSubmitted />
  </div>
);

export const NTSReleaseShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={object('displayInfo', ntsReleaseInfo)}
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      isSubmitted
    />
  </div>
);

export const ApprovedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={object('displayInfo', hhgInfo)}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted={false}
    />
  </div>
);

export const PostalOnlyDestination = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={object('displayInfo', postalOnlyInfo)}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
    />
  </div>
);

export const DivertedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      shipmentId="1"
      displayInfo={object('displayInfo', diversionInfo)}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
    />
  </div>
);

export const CancelledShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      shipmentId="1"
      displayInfo={object('displayInfo', cancelledInfo)}
      shipmentType={SHIPMENT_OPTIONS.HHG}
      isSubmitted
    />
  </div>
);
