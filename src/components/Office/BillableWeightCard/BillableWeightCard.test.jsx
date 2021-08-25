import React from 'react';
import { render, screen } from '@testing-library/react';

import BillableWeightCard from './BillableWeightCard';

import { formatWeight } from 'shared/formatters';

describe('BillableWeightCard', () => {
  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
    const shipments = [
      { id: '0001', shipmentType: 'HHG', billableWeightCap: '6,161', primeEstimatedWeight: '5,600' },
      { id: '0002', shipmentType: 'HHG', billableWeightCap: '3,200', reweigh: { id: '1234' } },
      { id: '0003', shipmentType: 'HHG', billableWeightCap: '3,400', primeEstimatedWeight: '5,000' },
    ];

    const defaultProps = {
      maxBillableWeight: 13750,
      totalBillableWeight: 12460,
      weightRequested: 12260,
      weightAllowance: 8000,
      shipments,
      reviewWeights: () => {},
    };

    render(<BillableWeightCard {...defaultProps} />);

    // labels
    expect(screen.getByText('Maximum billable weight')).toBeInTheDocument();
    expect(screen.getByText('Weight requested')).toBeInTheDocument();
    expect(screen.getByText('Weight allowance')).toBeInTheDocument();
    expect(screen.getByText('Total billable weight')).toBeInTheDocument();

    // weights
    expect(screen.getByText(formatWeight(defaultProps.maxBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.totalBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightRequested))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();

    // flags
    expect(screen.getByText('Over weight')).toBeInTheDocument();
    expect(screen.getByText('Missing weight')).toBeInTheDocument();

    // shipment weights
    expect(screen.getByText(formatWeight(shipments[0].billableWeightCap))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[1].billableWeightCap))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[2].billableWeightCap))).toBeInTheDocument();
  });
});
