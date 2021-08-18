import React from 'react';
import { render } from '@testing-library/react';

import BillableWeightCard from './BillableWeightCard';

describe('BillableWeightCard', () => {
  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
    const shipments = [
      { id: '0001', shipmentType: 'HHG', billableWeightCap: '5,600' },
      { id: '0002', shipmentType: 'HHG', billableWeightCap: '3,200', reweigh: { id: '1234' } },
      { id: '0003', shipmentType: 'HHG', billableWeightCap: '3,400' },
    ];
    const entitlements = [
      { id: '1234', shipmentId: '0001', authorizedWeight: '4,600' },
      { id: '12346', shipmentId: '0002', authorizedWeight: '4,600' },
      { id: '12347', shipmentId: '0003', authorizedWeight: '4,600' },
    ];

    const defaultProps = {
      maxBillableWeight: '13,750',
      totalBillableWeight: '12,460',
      weightRequested: '12,260',
      weightAllowance: '8,000',
      entitlements,
      shipments,
    };

    const { getByText } = render(<BillableWeightCard {...defaultProps} />);

    // labels
    expect(getByText('Maximum billable weight')).toBeInTheDocument();
    expect(getByText('Weight requested')).toBeInTheDocument();
    expect(getByText('Weight allowance')).toBeInTheDocument();
    expect(getByText('Total billable weight')).toBeInTheDocument();

    // weights
    expect(getByText(`${defaultProps.maxBillableWeight} lbs`)).toBeInTheDocument();
    expect(getByText(`${defaultProps.totalBillableWeight} lbs`)).toBeInTheDocument();
    expect(getByText(`${defaultProps.weightRequested} lbs`)).toBeInTheDocument();
    expect(getByText(`${defaultProps.weightAllowance} lbs`)).toBeInTheDocument();

    // flags
    expect(getByText('Over weight')).toBeInTheDocument();
    expect(getByText('Missing weight')).toBeInTheDocument();

    // shipment weights
    expect(getByText(`${shipments[0].billableWeightCap} lbs`)).toBeInTheDocument();
    expect(getByText(`${shipments[1].billableWeightCap} lbs`)).toBeInTheDocument();
    expect(getByText(`${shipments[2].billableWeightCap} lbs`)).toBeInTheDocument();
  });
});
