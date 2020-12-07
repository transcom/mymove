import React from 'react';
import { withKnobs } from '@storybook/addon-knobs';

import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';

export default {
  title: 'TOO/TIO Components/ShipmentWeightDetails',
  decorators: [
    withKnobs,
    (storyFn) => (
      <div id="containers" style={{ padding: '20px' }}>
        {storyFn()}
      </div>
    ),
  ],
};

export const WithNoDetails = () => <ShipmentWeightDetails />;

export const WithDetails = () => <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} />;
