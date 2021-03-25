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

export const LargeTable = () => (
  <ServiceItemCalculations serviceItemParams={testParams.DomesticLongHaul} totalAmountRequested={642} itemCode="DLH" />
);

export const SmallTable = () => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticLongHaul}
    totalAmountRequested={642}
    itemCode="DLH"
    tableSize="small"
  />
);

export const LargeDOASIT = () => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginAdditionalSIT}
    totalAmountRequested={642}
    itemCode="DOASIT"
  />
);

export const SmallDOASIT = () => (
  <ServiceItemCalculations
    serviceItemParams={testParams.DomesticOriginAdditionalSIT}
    totalAmountRequested={642}
    itemCode="DOASIT"
    tableSize="small"
  />
);
