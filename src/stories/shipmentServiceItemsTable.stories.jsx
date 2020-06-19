import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs } from '@storybook/addon-knobs';

import ShipmentServiceItemsTable from '../components/Office/ShipmentServiceItemsTable';

storiesOf('TOO/TIO Components|ShipmentServiceItemsTable', module)
  .addDecorator(withKnobs)
  .add('Shipment Service Items Table', () => <ShipmentServiceItemsTable shipmentType="hhg" />);
