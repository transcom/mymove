import React from 'react';

import ExternalVendorWeightSummary from './ExternalVendorWeightSummary';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/ExternalVendorWeightSummary',
  component: ExternalVendorWeightSummary,
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

export const WithMultipleShipments = () => (
  <ExternalVendorWeightSummary
    shipments={[
      {
        ntsRecordedWeight: 1000,
      },
      {
        ntsRecordedWeight: 2000,
      },
      {
        ntsRecordedWeight: 1500,
      },
    ]}
  />
);

export const WithOneShipment = () => (
  <ExternalVendorWeightSummary
    shipments={[
      {
        ntsRecordedWeight: 1000,
      },
    ]}
  />
);
