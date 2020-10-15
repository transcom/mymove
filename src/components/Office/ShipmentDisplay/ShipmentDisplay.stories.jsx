import React from 'react';
import { withKnobs, object } from '@storybook/addon-knobs';

import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'TOO/TIO Components|Shipment Display',
  component: ShipmentDisplay,
  decorators: [withKnobs],
};

const hhgInfo = {
  heading: 'HHG',
  requestedMoveDate: '26 Mar 2020',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
};

const ntsInfo = {
  heading: 'NTS',
  requestedMoveDate: '26 Mar 2020',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
};

const postalOnlyInfo = {
  heading: 'HHG',
  requestedMoveDate: '26 Mar 2020',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    postal_code: '98421',
  },
};

export const HHGShipment = () => (
  <div style={{ padding: '20px' }}>
    <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} isSubmitted />
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
