import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PrimeSimulatorAvailableMoves from './AvailableMovesQueue';

describe('Prime component', () => {
  it('renders on the screen', () => {
    render(<PrimeSimulatorAvailableMoves />);

    expect(screen.queryByTestId('primedatefilter')).toBeInTheDocument();
  });

  it('setFilterByDate correctly parses the input', async () => {
    const input = screen.getByTestId('primedatefilter');
    const submit = screen.getByRole('button');

    await userEvent.type(input, '2023-13-29');
    await userEvent.click(submit);

    expect(screen.getAllByText('Enter a valid date.')).toBeInTheDocument();
  });
});
