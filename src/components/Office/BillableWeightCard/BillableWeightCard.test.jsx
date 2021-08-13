import React from 'react';
import { render } from '@testing-library/react';

import BillableWeightCard from './BillableWeightCard';

describe('BillableWeightCard', () => {
  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
    const weights = {
      maxBillableWeight: '13,750',
      totalBillableWeight: '12,460',
      weightRequested: '12,260',
      weightAllowance: '8,000',
    };
    const { getByText } = render(<BillableWeightCard {...weights} />);

    // labels
    expect(getByText('Maximum billable weight')).toBeInTheDocument();
    expect(getByText('Weight requested')).toBeInTheDocument();
    expect(getByText('Weight allowance')).toBeInTheDocument();
    expect(getByText('Total billable weight')).toBeInTheDocument();

    // weights
    expect(getByText(`${weights.maxBillableWeight} lbs`)).toBeInTheDocument();
    expect(getByText(`${weights.totalBillableWeight} lbs`)).toBeInTheDocument();
    expect(getByText(`${weights.weightRequested} lbs`)).toBeInTheDocument();
    expect(getByText(`${weights.weightAllowance} lbs`)).toBeInTheDocument();
  });
});
