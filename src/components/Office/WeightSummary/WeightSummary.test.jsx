import React from 'react';
import { render, screen } from '@testing-library/react';

import { formatWeight } from 'shared/formatters';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';

describe('WeightSummary', () => {
  it('renders without crashing', () => {
    const shipments = [
      { id: '0001', shipmentType: 'HHG', billableWeight: 6161, estimatedWeight: 5600 },
      {
        id: '0002',
        shipmentType: 'HHG',
        billableWeight: 3200,
        estimatedWeight: 5000,
        reweigh: { id: '1234' },
      },
      { id: '0003', shipmentType: 'HHG', billableWeight: 3400, estimatedWeight: 5000 },
    ];

    const defaultProps = {
      maxBillableWeight: 13750,
      totalBillableWeight: 12460,
      weightRequested: 12260,
      weightAllowance: 8000,
      shipments,
    };

    render(<WeightSummary {...defaultProps} />);
    // labels
    expect(screen.getByText('Max billable weight')).toBeInTheDocument();
    expect(screen.getByText('Weight requested')).toBeInTheDocument();
    expect(screen.getByText('Weight allowance')).toBeInTheDocument();
    expect(screen.getByText('Total billable weight')).toBeInTheDocument();

    // weights
    expect(screen.getByText(formatWeight(defaultProps.maxBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.totalBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightRequested))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();

    // shipment weights
    // expect(screen.getByText(formatWeight(shipments[0].billableWeight))).toBeInTheDocument();
    // expect(screen.getByText(formatWeight(shipments[1].billableWeight))).toBeInTheDocument();
    // expect(screen.getByText(formatWeight(shipments[2].billableWeight))).toBeInTheDocument();
  });
});
