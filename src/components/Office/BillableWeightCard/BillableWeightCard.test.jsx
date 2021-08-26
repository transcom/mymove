import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import BillableWeightCard from './BillableWeightCard';

import { formatWeight } from 'shared/formatters';

describe('BillableWeightCard', () => {
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
    onReviewWeights: jest.fn(),
  };

  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
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

  it('implements the review weights handler when the review weights button is clicked', async () => {
    render(<BillableWeightCard {...defaultProps} />);

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });

    userEvent.click(reviewWeights);

    await waitFor(() => {
      expect(defaultProps.onReviewWeights).toHaveBeenCalled();
    });
  });
});
