import React from 'react';
import { render, screen } from '@testing-library/react';

import BillableWeightCard from './BillableWeightCard';

describe('BillableWeightCard', () => {
  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
    const shipments = [
      { id: '0001', shipmentType: 'HHG', billableWeightCap: '6,161', primeEstimatedWeight: '5,600' },
      { id: '0002', shipmentType: 'HHG', billableWeightCap: '3,200', reweigh: { id: '1234' } },
      { id: '0003', shipmentType: 'HHG', billableWeightCap: '3,400', primeEstimatedWeight: '5,000' },
    ];

    const defaultProps = {
      maxBillableWeight: '13,750',
      totalBillableWeight: '12,460',
      weightRequested: '12,260',
      weightAllowance: '8,000',
      shipments,
    };

    render(<BillableWeightCard {...defaultProps} />);

    // labels
    expect(screen.getByText('Maximum billable weight')).toBeInTheDocument();
    expect(screen.getByText('Weight requested')).toBeInTheDocument();
    expect(screen.getByText('Weight allowance')).toBeInTheDocument();
    expect(screen.getByText('Total billable weight')).toBeInTheDocument();

    // weights
    expect(screen.getByText(`${defaultProps.maxBillableWeight} lbs`)).toBeInTheDocument();
    expect(screen.getByText(`${defaultProps.totalBillableWeight} lbs`)).toBeInTheDocument();
    expect(screen.getByText(`${defaultProps.weightRequested} lbs`)).toBeInTheDocument();
    expect(screen.getByText(`${defaultProps.weightAllowance} lbs`)).toBeInTheDocument();

    // flags
    expect(screen.getByText('Over weight')).toBeInTheDocument();
    expect(screen.getByText('Missing weight')).toBeInTheDocument();

    // shipment weights
    expect(screen.getByText(`${shipments[0].billableWeightCap} lbs`)).toBeInTheDocument();
    expect(screen.getByText(`${shipments[1].billableWeightCap} lbs`)).toBeInTheDocument();
    expect(screen.getByText(`${shipments[2].billableWeightCap} lbs`)).toBeInTheDocument();
  });
});
