import React from 'react';
import { render, screen } from '@testing-library/react';

import SITCostDetails from './SITCostDetails';

import { LOCATION_TYPES } from 'types/sitStatusShape';

const args = {
  cost: 12300,
  weight: 234,
  sitLocation: LOCATION_TYPES.ORIGIN,
  originZip: '23456',
  destinationZip: '54321',
  departureDate: '2022-10-29',
  entryDate: '2022-08-06',
};

describe('components/Office/SITCostDetails', () => {
  it('renders correctly', () => {
    render(<SITCostDetails {...args} />);
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent(/Storage in transit \(SIT\)/);
  });

  it('displays passed SIT details correctly', () => {
    render(<SITCostDetails {...args} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(/Government constructed cost: \$123/);
    expect(screen.getByText(/234 lbs of origin SIT/)).toBeInTheDocument();
    expect(screen.getByText(/at 23456/)).toBeInTheDocument();
  });

  it('displays destination ZIP when SIT location is destination', () => {
    args.sitLocation = LOCATION_TYPES.DESTINATION;
    render(<SITCostDetails {...args} />);
    expect(screen.getByText(/at 54321/)).toBeInTheDocument();
  });

  it('correctly computes days between entry and departure', async () => {
    render(<SITCostDetails {...args} />);
    expect(screen.getByText(/for 84 days./)).toBeInTheDocument();
  });

  it('comma-separates thousands values', async () => {
    args.cost = 123400;
    args.weight = 23456;
    render(<SITCostDetails {...args} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(/Government constructed cost: \$1,234/);
    expect(screen.getByText(/23,456 lbs/)).toBeInTheDocument();
  });

  it('displays singular "day" when there is only one', async () => {
    args.entryDate = '2022-10-28';
    render(<SITCostDetails {...args} />);
    expect(screen.getByText(/for 1 day./)).toBeInTheDocument();
  });
});
