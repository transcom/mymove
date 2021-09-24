import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import BillableWeightCard from './BillableWeightCard';

import { formatWeight } from 'shared/formatters';

describe('BillableWeightCard', () => {
  const defaultProps = {
    maxBillableWeight: 13750,
    totalBillableWeight: 12460,
    weightRequested: 12260,
    weightAllowance: 8000,
    onReviewWeights: jest.fn(),
  };

  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 2161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234', weight: 40 },
      },
      {
        id: '0002',
        shipmentType: 'HHG',
        calculatedBillableWeight: 3200,
        estimatedWeight: 5000,
        primeEstimatedWeight: 1000,
        reweigh: { id: '1234', weight: 300 },
      },
      {
        id: '0003',
        shipmentType: 'HHG',
        calculatedBillableWeight: 3400,
        estimatedWeight: 5000,
        primeEstimatedWeight: 200,
        reweigh: { id: '1234', weight: 500 },
      },
    ];

    render(<BillableWeightCard {...defaultProps} shipments={shipments} />);

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
    expect(screen.queryByText('Over weight')).not.toBeInTheDocument();
    expect(screen.queryByText('Missing weight')).not.toBeInTheDocument();

    // shipment weights
    expect(screen.getByText(formatWeight(shipments[0].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[1].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[2].calculatedBillableWeight))).toBeInTheDocument();
  });

  it('shows an over weight flag if a shipment is over weight', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234', weight: 40 },
      },
    ];

    render(<BillableWeightCard {...defaultProps} shipments={shipments} />);

    expect(screen.queryByText('Missing weight')).not.toBeInTheDocument();
    expect(screen.getByText('Over weight')).toBeInTheDocument();
  });

  it('shows a missing weight flag if a shipment is missing the reweigh weight', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 1161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234' },
      },
    ];

    render(<BillableWeightCard {...defaultProps} shipments={shipments} />);

    expect(screen.queryByText('Over weight')).not.toBeInTheDocument();
    expect(screen.getByText('Missing weight')).toBeInTheDocument();
  });

  it('shows a missing weight flag if a shipment is missing a prime estimated weight', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 4161,
        estimatedWeight: 5600,
        reweigh: { id: '1234', weight: 40 },
      },
    ];

    render(<BillableWeightCard {...defaultProps} shipments={shipments} />);

    expect(screen.queryByText('Over weight')).not.toBeInTheDocument();
    expect(screen.getByText('Missing weight')).toBeInTheDocument();
  });

  it('implements the review weights handler when the review weights button is clicked', async () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234', weight: 40 },
      },
    ];
    render(<BillableWeightCard {...defaultProps} shipments={shipments} />);

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });

    userEvent.click(reviewWeights);

    await waitFor(() => {
      expect(defaultProps.onReviewWeights).toHaveBeenCalled();
    });
  });
});
