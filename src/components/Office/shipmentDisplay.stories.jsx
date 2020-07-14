import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs, object } from '@storybook/addon-knobs';

import ShipmentDisplay from 'components/Office/ShipmentDisplay';

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

storiesOf('TOO/TIO Components|ShipmentDisplay', module)
  .addDecorator(withKnobs)
  .add('HHG Shipment', () => {
    return (
      <div style={{ padding: '20px' }}>
        <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} isSubmitted />
      </div>
    );
  })
  .add('NTS Shipment', () => {
    return (
      <div style={{ padding: '20px' }}>
        <ShipmentDisplay displayInfo={object('displayInfo', ntsInfo)} shipmentType="NTS" isSubmitted />
      </div>
    );
  })
  .add('Approved Shipment', () => {
    return (
      <div style={{ padding: '20px' }}>
        <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} isSubmitted={false} />
      </div>
    );
  });
