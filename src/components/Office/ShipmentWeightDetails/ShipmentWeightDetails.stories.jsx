import React from 'react';

import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';

export default {
  title: 'Office Components/ShipmentWeightDetails',
  decorators: [
    (storyFn) => (
      <div id="containers" style={{ padding: '20px' }}>
        {storyFn()}
      </div>
    ),
  ],
};

export const WithNoDetails = () => <ShipmentWeightDetails />;

export const WithDetails = () => <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} />;
