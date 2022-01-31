import React from 'react';
import { render, screen } from '@testing-library/react';

import ExternalVendorWeightSummary from './ExternalVendorWeightSummary';

import { MockProviders } from 'testUtils';

const shipments = [{ ntsRecordedWeight: 1000 }, { ntsRecordedWeight: 500 }, { ntsRecordedWeight: 1500 }];

describe('ExternalVendorWeightSummary component', () => {
  it('renders with one shipment', () => {
    render(
      <MockProviders>
        <ExternalVendorWeightSummary shipments={[shipments[0]]} />
      </MockProviders>,
    );

    expect(screen.getByText('1 other shipment:')).toBeInTheDocument();
    expect(screen.getByText('1,000 lbs')).toBeInTheDocument();
  });

  it('renders with one shipment', () => {
    render(
      <MockProviders>
        <ExternalVendorWeightSummary shipments={shipments} />
      </MockProviders>,
    );

    expect(screen.getByText('3 other shipments:')).toBeInTheDocument();
    expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
  });
});
