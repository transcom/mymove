import React from 'react';
import { render, screen } from '@testing-library/react';

import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';

describe('WeightDisplay', () => {
  it('renders without crashing', () => {
    render(<WeightDisplay heading="heading test" />);

    expect(screen.getByText('heading test')).toBeInTheDocument();
  });

  it('renders with weight value', () => {
    render(<WeightDisplay heading="heading test" value={1234} />);

    expect(screen.getByText('1,234 lbs')).toBeInTheDocument();
  });

  it('renders with edit button', () => {
    render(<WeightDisplay heading="heading test" value={1234} showEditBtn />);

    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('edit button is clicked', () => {
    const mockEditBtn = jest.fn();
    render(<WeightDisplay heading="heading test" value={1234} showEditBtn onEdit={mockEditBtn} />);
    screen.getByRole('button').click();

    expect(mockEditBtn).toHaveBeenCalledTimes(1);
  });
});
