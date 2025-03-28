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

export const WithMultipleNTSRShipments = () => (
  <ExternalVendorWeightSummary
    shipments={[
      {
        ntsRecordedWeight: 1000,
        shipmentType: 'HHG_OUTOF_NTS',
      },
      {
        ntsRecordedWeight: 2000,
        shipmentType: 'HHG_OUTOF_NTS',
      },
      {
        ntsRecordedWeight: 1500,
        shipmentType: 'HHG_OUTOF_NTS',
      },
    ]}
  />
);

export const WithMultipleNTSShipments = () => (
  // NTS shipments from external vendors don't have weights
  <ExternalVendorWeightSummary
    shipments={[
      {
        shipmentType: 'HHG_INTO_NTS',
      },
      {
        shipmentType: 'HHG_INTO_NTS',
      },
      {
        shipmentType: 'HHG_INTO_NTS',
      },
    ]}
  />
);

export const WithMultipleShipmentsOfBothTypes = () => (
  <ExternalVendorWeightSummary
    shipments={[
      {
        ntsRecordedWeight: 1000,
        shipmentType: 'HHG_OUTOF_NTS',
      },
      {
        ntsRecordedWeight: 2000,
        shipmentType: 'HHG_OUTOF_NTS',
      },
      {
        shipmentType: 'HHG_INTO_NTS',
      },
      {
        ntsRecordedWeight: 1500,
        shipmentType: 'HHG_OUTOF_NTS',
      },
      {
        shipmentType: 'HHG_INTO_NTS',
      },
    ]}
  />
);

export const WithOneNTSShipment = () => (
  <ExternalVendorWeightSummary
    shipments={[
      {
        NTS: 'HHG_INTO_NTS',
      },
    ]}
  />
);

export const WithOneNTSRShipment = () => (
  <ExternalVendorWeightSummary
    shipments={[
      {
        ntsRecordedWeight: 1500,
        shipmentType: 'HHG_OUTOF_NTS',
      },
    ]}
  />
);
