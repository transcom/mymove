import React from 'react';

import ServiceItemCalculations from './ServiceItemCalculations';

export default {
  title: 'Office Components/ServiceItemCalculations',
  decorators: [
    (Story) => {
      return (
        <div style={{ padding: '20px' }}>
          <Story />
        </div>
      );
    },
  ],
};

export const Default = () => <ServiceItemCalculations />;
