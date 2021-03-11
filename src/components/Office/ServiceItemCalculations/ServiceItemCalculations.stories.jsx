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

const largeData = [
  {
    value: '85 cwt',
    label: 'Billable weight (cwt)',
    details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000'],
  },
  {
    value: '2,337',
    label: 'Mileage',
    details: ['Zip 322 to Zip 919'],
  },
  {
    value: '0.03',
    label: 'Baseline linehaul price',
    details: ['Domestic non-peak', 'Origin service area: 176', 'Pickup date: 24 Jan 2020'],
  },
  {
    value: '1.033',
    label: 'Price escalation factor',
    details: null,
  },
  {
    value: '$6.423',
    label: 'Total amount requested',
    details: [],
  },
];

// modify data to match high fidelity mock up
const smallData = [...largeData];
smallData[0].value = '85';

export const LargeTable = () => <ServiceItemCalculations calculations={largeData} />;

export const SmallTable = () => <ServiceItemCalculations calculations={smallData} tableSize="small" />;
