import React from 'react';
import { storiesOf } from '@storybook/react';

import ImportantShipmentDates from './ImportantShipmentDates';

storiesOf('TOO/TIO Components|ImportantShipmentDates', module).add('default', () => {
  return (
    <ImportantShipmentDates requestedPickupDate="Thursday, 26 Mar 2020" scheduledPickupDate="Friday, 27 Mar 2020" />
  );
});
