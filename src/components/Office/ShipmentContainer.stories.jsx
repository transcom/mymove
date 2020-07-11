import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs, text } from '@storybook/addon-knobs';

import ShipmentContainer from './ShipmentContainer';
import ShipmentHeading from './ShipmentHeading';

storiesOf('TOO/TIO Components|ShipmentContainer', module)
  .addDecorator(withKnobs)
  .add('Shipment Container', () => (
    <ShipmentContainer shipmentType={text('ShipmentContainer.shipmentType', 'HHG')}>
      <ShipmentHeading
        shipmentInfo={{
          shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
          originCity: text('ShipmentInfo.originCity', 'San Antonio'),
          originState: text('ShipmentInfo.originState', 'TX'),
          originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
          destinationCity: text('ShipmentInfo.destinationCity', 'Tacoma'),
          destinationState: text('ShipmentInfo.destinationState', 'WA'),
          destinationPostalCode: text('ShipmentInfo.destinationPostalCode', '98421'),
          scheduledPickupDate: text('ShipmentInfo.destinationPostalCode', '27 Mar 2020'),
        }}
      />
    </ShipmentContainer>
  ));
