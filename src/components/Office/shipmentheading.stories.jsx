import React from 'react';
import { storiesOf } from '@storybook/react';
import { object, text } from '@storybook/addon-knobs';

import ShipmentHeading from './ShipmentHeading';

storiesOf('TOO/TIO Components/ShipmentHeading', module).add('Shipment Heading', () => (
  <ShipmentHeading
    shipmentInfo={{
      shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
      originCity: text('ShipmentInfo.originCity', 'San Antonio'),
      originState: text('ShipmentInfo.originState', 'TX'),
      originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
      destinationAddress: object('MTOShipment.destinationAddress', {
        street_address_1: '123 Any Street',
        city: 'Tacoma',
        state: 'WA',
        postal_code: '98421',
      }),
      scheduledPickupDate: text('ShipmentInfo.scheduledPickupDate', '27 Mar 2020'),
    }}
  />
));
