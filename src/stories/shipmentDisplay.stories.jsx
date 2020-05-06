import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs, object } from '@storybook/addon-knobs';
import { SHIPMENT_TYPE } from 'shared/constants';

import ShipmentDisplay from 'components/Office/ShipmentDisplay';
import RequestedShipments from 'components/Office/RequestedShipments';

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
  .add('with one shipment requested', () => {
    return (
      <div style={{ padding: '20px' }}>
        <ShipmentDisplay displayInfo={object('displayInfo', hhgInfo)} />
      </div>
    );
  })
  .add('with two shipment requested', () => {
    return (
      <div style={{ padding: '20px' }}>
        <RequestedShipments>
          <ShipmentDisplay displayInfo={hhgInfo} />
          <ShipmentDisplay shipmentType={SHIPMENT_TYPE.NTS} displayInfo={object('displayInfo', ntsInfo)} />
        </RequestedShipments>
      </div>
    );
  });
