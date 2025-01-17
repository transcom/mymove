import React from 'react';
import { render, screen } from '@testing-library/react';

import ExternalVendorWeightSummary from './ExternalVendorWeightSummary';

import { MockProviders } from 'testUtils';

describe('ExternalVendorWeightSummary component', () => {
  it('renders with one NTS shipment', () => {
    render(
      <MockProviders>
        <ExternalVendorWeightSummary shipments={[{ shipmentType: 'HHG_INTO_NTS' }]} />
      </MockProviders>,
    );

    expect(screen.getByText('1 other shipment:')).toBeInTheDocument();
  });

  it('renders with many NTS shipments', () => {
    const shipments = [{ shipmentType: 'HHG_INTO_NTS' }, { shipmentType: 'HHG_INTO_NTS' }];
    render(
      <MockProviders>
        <ExternalVendorWeightSummary shipments={shipments} />
      </MockProviders>,
    );

    expect(screen.getByText('2 other shipments:')).toBeInTheDocument();
  });

  it('renders with one NTSR shipment', () => {
    render(
      <MockProviders>
        <ExternalVendorWeightSummary shipments={[{ ntsRecordedWeight: 1500, shipmentType: 'HHG_OUTOF_NTS' }]} />
      </MockProviders>,
    );

    expect(screen.getByText('1 other shipment:')).toBeInTheDocument();
    expect(screen.getByText('1,500 lbs')).toBeInTheDocument();
  });

  it('renders with many NTSR shipments', () => {
    const shipments = [
      { ntsRecordedWeight: 1000, shipmentType: 'HHG_OUTOF_NTS' },
      { ntsRecordedWeight: 500, shipmentType: 'HHG_OUTOF_NTS' },
      { ntsRecordedWeight: 1500, shipmentType: 'HHG_OUTOF_NTS' },
    ];
    render(
      <MockProviders>
        <ExternalVendorWeightSummary shipments={shipments} />
      </MockProviders>,
    );

    expect(screen.getByText('3 other shipments:')).toBeInTheDocument();
    expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
  });

  it('renders with many NTSR and NTS shipments', () => {
    const shipments = [
      { ntsRecordedWeight: 1000, shipmentType: 'HHG_OUTOF_NTS' },
      { shipmentType: 'HHG_INTO_NTS' },
      { ntsRecordedWeight: 500, shipmentType: 'HHG_OUTOF_NTS' },
      { shipmentType: 'HHG_INTO_NTS' },
      { ntsRecordedWeight: 1500, shipmentType: 'HHG_OUTOF_NTS' },
    ];
    render(
      <MockProviders>
        <ExternalVendorWeightSummary shipments={shipments} />
      </MockProviders>,
    );

    expect(screen.getByText('5 other shipments:')).toBeInTheDocument();
    expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
  });
});
