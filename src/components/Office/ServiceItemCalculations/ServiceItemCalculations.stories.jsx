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

const data = [
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
    label: ['Total amount requested'],
    details: [],
  },
];

export const Default = () => <ServiceItemCalculations calculations={data} />;
