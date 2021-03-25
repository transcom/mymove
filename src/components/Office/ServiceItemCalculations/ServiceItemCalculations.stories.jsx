import React from 'react';

import ServiceItemCalculations from './ServiceItemCalculations';
import testParams from './serviceItemTestParams';

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

export const LargeTableDLH = () => (
  <ServiceItemCalculations serviceItemParams={testParams.DomesticLongHaul} totalAmountRequested={642} itemCode="DLH" />
);

export const SmallTableDLH = () => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticLongHaul}
    totalAmountRequested={642}
    itemCode="DLH"
    tableSize="small"
  />
);

export const LargeTableDOP = () => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginPrice}
    totalAmountRequested={642}
    itemCode="DOP"
  />
);

export const SmallTableDOP = () => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginPrice}
    totalAmountRequested={642}
    itemCode="DOP"
    tableSize="small"
  />
);
