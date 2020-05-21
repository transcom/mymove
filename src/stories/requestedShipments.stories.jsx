import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs } from '@storybook/addon-knobs';

import RequestedShipments from 'components/Office/RequestedShipments';

storiesOf('TOO/TIO Components|RequestedShipments', module)
  .addDecorator(withKnobs)
  .add('with no details', () => {
    return (
      <div style={{ padding: '20px' }}>
        <RequestedShipments />
      </div>
    );
  })
  .add('with details', () => {
    return (
      <div style={{ position: 'relative', padding: '20px' }}>
        <RequestedShipments>
          <>Inside requested shipments component.</>
        </RequestedShipments>
      </div>
    );
  });
