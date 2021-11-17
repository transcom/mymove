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
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
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
    <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} isSubmitted />
  </div>
);

export const HHGShipmentNoIcon = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} isSubmitted showIcon={false} />
  </div>
);

export const HHGShipmentWithCounselorRemarks = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
      isSubmitted
      showIcon={false}
    />
  </div>
);

export const HHGShipmentEditable = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay
      displayInfo={{ ...hhgInfo, counselorRemarks: 'counselor approved' }}
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

export const ApprovedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} isSubmitted={false} />
  </div>
);

export const PostalOnlyDestination = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={object('displayInfo', postalOnlyInfo)} isSubmitted />
  </div>
);

export const DivertedShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay shipmentId="1" displayInfo={object('displayInfo', diversionInfo)} isSubmitted />
  </div>
);

export const CancelledShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay shipmentId="1" displayInfo={object('displayInfo', cancelledInfo)} isSubmitted />
  </div>
);
