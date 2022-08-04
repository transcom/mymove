import React from 'react';
import { render, screen } from '@testing-library/react';

import SITCostDetails from './SITCostDetails';

const args = {
  cost: 12300,
  weight: 234,
  sitLocation: 'ORIGIN',
  location: '23456',
  departureDate: '2022-10-29',
  entryDate: '2022-08-06',
};

describe('components/Office/SITCostDetails', () => {
  it('renders correctly', () => {
    render(<SITCostDetails {...args} />);
    expect(screen.getByText(/Storage in transit \(SIT\)/)).toBeInTheDocument();
  });

  it('displays passed SIT details correctly', () => {
    render(<SITCostDetails {...args} />);
    expect(screen.getByText(/Government constructed cost: \$123/)).toBeInTheDocument();
    expect(screen.getByText(/Maximum reimbursement for storing 234 lbs of origin SIT at 23456/)).toBeInTheDocument();
  });

  it('correctly computes days between entry and departure', async () => {
    render(<SITCostDetails {...args} />);
    expect(screen.queryByText(/for 84 days./));
  });

  it('comma-separates thousands values', async () => {
    args.cost = 123400;
    args.weight = 23456;
    render(<SITCostDetails {...args} />);
    expect(screen.queryByText(/Government constructed cost: \$1,234/)).toBeInTheDocument();
    expect(screen.queryByText(/Maximum reimbursement for storing 23,456 lbs/)).toBeInTheDocument();
  });

  it('displays singular "day" when there is only one', async () => {
    args.entryDate = '2022-10-28';
    render(<SITCostDetails {...args} />);
    expect(screen.queryByText(/for 1 day./)).toBeInTheDocument();
  });
});
